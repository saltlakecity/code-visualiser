package session

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"visualisator/internal/analyzer"
	"visualisator/internal/debugger"
	"visualisator/internal/models"
	"visualisator/internal/viz"
	"visualisator/internal/workspace"
)

const defaultTimeout = 30 * time.Second

// DebugCompleteFunc is called when automated run finishes.
type DebugCompleteFunc func(sessionID string, trace models.Trace, err error)

// Manager coordinates compile and debug sessions.
type Manager struct {
	mu               sync.RWMutex
	builder          *workspace.Builder
	sessions         map[string]*entry
	onDebugComplete  DebugCompleteFunc
}

type entry struct {
	paths            *workspace.SessionPaths
	debug            *debugger.Session
	mode             string
	source           string
	editorSource     string
	debugSource      string
	debugFile        string
	editorLineOffset int
	sourceMap        *analyzer.SourceMap
	done             chan struct{}
}

// NewManager creates a session manager.
func NewManager() (*Manager, error) {
	base := workspace.BuildRoot()
	b, err := workspace.NewBuilder(base)
	if err != nil {
		return nil, err
	}
	return &Manager{
		builder:  b,
		sessions: make(map[string]*entry),
	}, nil
}

// SetDebugCompleteHandler registers callback after RunToEnd.
func (m *Manager) SetDebugCompleteHandler(fn DebugCompleteFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onDebugComplete = fn
}

// Compile builds user source. mainExtra is optional main.go body for problem runners.
func (m *Manager) Compile(source string, mode string, mainExtra string) (models.CompileResult, error) {
	if mode == "" || mode == "auto" {
		mode = workspace.DetectMode(source)
	}

	editorSource := source
	prepared, stubMain, err := workspace.PrepareUserSource(source, mode)
	if err != nil {
		return models.CompileResult{Success: false, Errors: []models.CompileError{{
			File: "main.go", Line: 1, Column: 1, Message: err.Error(),
		}}}, nil
	}
	source = prepared
	if mainExtra == "" && stubMain != "" {
		mainExtra = stubMain
	}
	if mode == "auto" {
		if workspace.HasMainPackage(source) {
			mode = "full_main"
		} else {
			mode = "function_wrapper"
		}
	}

	paths, err := m.builder.Prepare()
	if err != nil {
		return models.CompileResult{}, err
	}
	if _, err := m.builder.WriteSource(paths, source, mode); err != nil {
		m.builder.Cleanup(paths.ID)
		return models.CompileResult{Success: false, Errors: []models.CompileError{{
			File: "main.go", Line: 1, Column: 1, Message: err.Error(),
		}}}, nil
	}
	if mainExtra != "" {
		mainSrc := mainExtra
		if !strings.Contains(mainSrc, "package ") {
			mainSrc = "package main\n\n" + mainSrc
		}
		if err := m.builder.WriteMain(paths, mainSrc); err != nil {
			m.builder.Cleanup(paths.ID)
			return models.CompileResult{}, err
		}
	}
	if err := workspace.ProcessDirImportsExec(paths.Dir); err != nil {
		m.builder.Cleanup(paths.ID)
		return models.CompileResult{Success: false, Errors: []models.CompileError{{
			File: "main.go", Line: 1, Column: 1, Message: err.Error(),
		}}}, nil
	}
	result := m.builder.Compile(paths)
	if !result.Success {
		m.builder.Cleanup(paths.ID)
		return result, nil
	}
	debugFile := "main.go"
	if mode == "function_wrapper" {
		debugFile = "solution.go"
	}
	debugPath := filepath.Join(paths.Dir, debugFile)
	debugBytes, err := os.ReadFile(debugPath)
	if err != nil {
		m.builder.Cleanup(paths.ID)
		return models.CompileResult{}, fmt.Errorf("read debug file: %w", err)
	}
	debugSource := string(debugBytes)
	offset := workspace.EditorLineOffset(editorSource, debugSource)

	var sourceMap *analyzer.SourceMap
	if sm, err := analyzer.ParseSource(editorSource); err == nil {
		sourceMap = sm
	}

	m.mu.Lock()
	m.sessions[paths.ID] = &entry{
		paths:            paths,
		mode:             mode,
		source:           source,
		editorSource:     editorSource,
		debugSource:      debugSource,
		debugFile:        debugFile,
		editorLineOffset: offset,
		sourceMap:        sourceMap,
	}
	m.mu.Unlock()

	result.CompiledSource = debugSource
	result.DebugSource = debugSource
	result.DebugFile = debugFile
	result.EditorLineOffset = offset
	return result, nil
}

// CompileAndRun builds and compiles without starting debug session.
func (m *Manager) CompileAndRun(source string, mode string, mainExtra string) (models.CompileResult, error) {
	res, err := m.Compile(source, mode, mainExtra)
	if err != nil {
		return res, err
	}
	return res, nil
}

// CleanupSession removes a compiled session without debug.
func (m *Manager) CleanupSession(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if e, ok := m.sessions[sessionID]; ok {
		if e.debug != nil {
			e.debug.Stop()
		}
		m.builder.Cleanup(sessionID)
		delete(m.sessions, sessionID)
	}
}

// StartDebug launches delve on a compiled session.
func (m *Manager) StartDebug(sessionID string, opts models.RunOptions) error {
	m.mu.Lock()
	e, ok := m.sessions[sessionID]
	m.mu.Unlock()
	if !ok {
		return fmt.Errorf("session %q not found", sessionID)
	}
	if e.debug != nil {
		e.debug.Stop()
	}
	dbg, err := debugger.NewSession(e.paths.Binary, e.paths.Dir, defaultTimeout)
	if err != nil {
		return err
	}
	lang := analyzer.LangEN
	if opts.Lang == "ru" {
		lang = analyzer.LangRU
	}
	dbg.SetExplainConfig(&analyzer.ExplainConfig{
		SourceMap:    e.sourceMap,
		EditorSource: e.editorSource,
		Narrate:      analyzer.NewNarrateContext(),
		Lang:         lang,
	})
	dbg.StartTimeout()
	e.debug = dbg

	maxSteps := opts.MaxSteps
	if maxSteps <= 0 {
		maxSteps = 800
	}
	done := make(chan struct{})
	e.done = done
	go func() {
		defer close(done)
		_, runErr := dbg.RunToEnd(maxSteps, e.editorLineOffset)
		trace := dbg.Recorder().Trace(sessionID)
		m.mu.Lock()
		cb := m.onDebugComplete
		m.mu.Unlock()
		if cb != nil {
			cb(sessionID, trace, runErr)
		}
	}()
	return nil
}

// Step runs a debugger command.
func (m *Manager) Step(sessionID, command string) (models.StepResult, error) {
	e := m.getEntry(sessionID)
	if e == nil || e.debug == nil {
		return models.StepResult{}, fmt.Errorf("no active debug session")
	}
	res, err := e.debug.Step(command)
	if err != nil {
		return models.StepResult{Error: err.Error()}, err
	}
	if res == nil {
		trace := e.debug.Recorder().Trace(sessionID)
		return models.StepResult{Trace: &trace}, nil
	}
	return *res, nil
}

// GetTrace returns recorded trace.
func (m *Manager) GetTrace(sessionID string) (models.Trace, error) {
	e := m.getEntry(sessionID)
	if e == nil {
		return models.Trace{}, fmt.Errorf("session %q not found", sessionID)
	}
	if e.debug != nil {
		return e.debug.Recorder().Trace(sessionID), nil
	}
	return models.Trace{SessionID: sessionID, Status: "idle"}, nil
}

// GetFrame returns a frame by index.
func (m *Manager) GetFrame(sessionID string, index int) (*models.TraceFrame, error) {
	e := m.getEntry(sessionID)
	if e == nil || e.debug == nil {
		return nil, fmt.Errorf("session not found")
	}
	f := e.debug.Recorder().FrameAt(index)
	if f == nil {
		return nil, fmt.Errorf("frame %d not found", index)
	}
	return f, nil
}

// BuildFlowchart returns a static flowchart for editor source.
func (m *Manager) BuildFlowchart(source string, lang string) (models.FlowchartResult, error) {
	sm, err := analyzer.ParseSource(source)
	if err != nil {
		return models.FlowchartResult{}, err
	}
	result, err := viz.BuildFlowGraph(source)
	if err != nil {
		return models.FlowchartResult{}, err
	}
	sig := sm.Outline.Name + "(" + sm.Outline.Params + ")"
	if sm.Outline.Returns != "" {
		sig += " " + sm.Outline.Returns
	}
	if lang == "ru" {
		result.Summary = fmt.Sprintf("%s — %d узлов", sig, len(result.Nodes))
	} else {
		result.Summary = fmt.Sprintf("%s — %d nodes", sig, len(result.Nodes))
	}
	return result, nil
}

// AnalyzeSource returns static outline for editor source (before run).
func (m *Manager) AnalyzeSource(source string) (models.SourceOutline, string, error) {
	sm, err := analyzer.ParseSource(source)
	if err != nil {
		return models.SourceOutline{}, "", err
	}
	lang := analyzer.LangRU
	return analyzer.ToSourceOutline(sm.Outline), analyzer.OutlineSummary(lang, sm.Outline), nil
}

// GetSourceOutline returns outline stored for a compiled session.
func (m *Manager) GetSourceOutline(sessionID string) (models.SourceOutline, bool) {
	e := m.getEntry(sessionID)
	if e == nil || e.sourceMap == nil {
		return models.SourceOutline{}, false
	}
	return analyzer.ToSourceOutline(e.sourceMap.Outline), true
}

// Stop ends debug and cleans up.
func (m *Manager) Stop(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.sessions[sessionID]
	if !ok {
		return
	}
	if e.debug != nil {
		e.debug.Stop()
	}
	m.builder.Cleanup(sessionID)
	delete(m.sessions, sessionID)
}

func (m *Manager) getEntry(sessionID string) *entry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions[sessionID]
}
