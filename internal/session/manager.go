package session

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"visualisator/internal/debugger"
	"visualisator/internal/models"
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
	paths           *workspace.SessionPaths
	debug           *debugger.Session
	mode            string
	source          string
	compiledSource  string
	debugFile       string
	editorLineOffset int
	done            chan struct{}
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
	offset := workspace.EditorLineOffset(editorSource, prepared)

	m.mu.Lock()
	m.sessions[paths.ID] = &entry{
		paths:            paths,
		mode:             mode,
		source:           source,
		compiledSource:   prepared,
		debugFile:        debugFile,
		editorLineOffset: offset,
	}
	m.mu.Unlock()

	result.CompiledSource = prepared
	result.DebugFile = debugFile
	result.EditorLineOffset = offset
	return result, nil
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
	dbg.StartTimeout()
	e.debug = dbg

	maxSteps := opts.MaxSteps
	if maxSteps <= 0 {
		maxSteps = 300
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
