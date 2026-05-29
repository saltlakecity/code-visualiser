package debugger

import (
	"sync"

	"visualisator/internal/models"
	"visualisator/internal/viz"
)

const maxFrames = 10000

// Recorder stores trace frames for replay.
type Recorder struct {
	mu                 sync.RWMutex
	frames             []models.TraceFrame
	status             string
	panicMsg           string
	failureFrameIndex  int
	failureEditorLine  int
	failureReason      string
	failureKind        string
}

// NewRecorder creates an empty recorder.
func NewRecorder() *Recorder {
	return &Recorder{status: "idle"}
}

// Reset clears recorded frames.
func (r *Recorder) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.frames = nil
	r.status = "running"
	r.panicMsg = ""
	r.failureFrameIndex = -1
	r.failureEditorLine = 0
	r.failureReason = ""
	r.failureKind = ""
}

// Append adds a frame if under limit (skips duplicate consecutive positions).
func (r *Recorder) Append(frame models.TraceFrame) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.frames) >= maxFrames {
		return
	}
	if n := len(r.frames); n > 0 {
		last := r.frames[n-1]
		if last.File == frame.File && last.Line == frame.Line && last.Function == frame.Function {
			// refresh variables on same line
			r.frames[n-1] = frame
			r.frames[n-1].Index = last.Index
			return
		}
	}
	frame.Index = len(r.frames)
	r.frames = append(r.frames, frame)
}

// SetFailure records failure context on the trace.
func (r *Recorder) SetFailure(frameIndex, editorLine int, reason, kind string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.failureFrameIndex = frameIndex
	r.failureEditorLine = editorLine
	r.failureReason = reason
	r.failureKind = kind
	if reason != "" {
		r.panicMsg = reason
	}
}

// MarkFailureFrame marks a frame as the failure point.
func (r *Recorder) MarkFailureFrame(index int, reason string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if index < 0 || index >= len(r.frames) {
		return
	}
	r.frames[index].IsFailure = true
	r.frames[index].FailureReason = reason
}

// SetStatus updates session status.
func (r *Recorder) SetStatus(status, panicMsg string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.status = status
	if panicMsg != "" {
		r.panicMsg = panicMsg
	}
}

// Trace returns a copy of the trace.
func (r *Recorder) Trace(sessionID string) models.Trace {
	r.mu.RLock()
	defer r.mu.RUnlock()
	frames := TrimTraceFrames(r.frames)
	failureIdx := r.failureFrameIndex
	if failureIdx >= 0 && len(frames) > 0 {
		// Adjust failure index after trim
		trimStart := 0
		if len(r.frames) > 4 {
			trimStart = 2
		}
		failureIdx = failureIdx - trimStart
		if failureIdx < 0 {
			failureIdx = 0
		}
		if failureIdx >= len(frames) {
			failureIdx = len(frames) - 1
		}
	}
	return models.Trace{
		SessionID:         sessionID,
		Frames:            frames,
		Status:            r.status,
		PanicMsg:          r.panicMsg,
		FailureFrameIndex: failureIdx,
		FailureEditorLine: r.failureEditorLine,
		FailureReason:     r.failureReason,
		FailureKind:       r.failureKind,
	}
}

// FrameAt returns frame by index.
func (r *Recorder) FrameAt(index int) *models.TraceFrame {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if index < 0 || index >= len(r.frames) {
		return nil
	}
	f := r.frames[index]
	return &f
}

// Count returns number of frames.
func (r *Recorder) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.frames)
}

// BuildFrame constructs a TraceFrame from debugger state.
func BuildFrame(file string, line, editorLine int, function string, vars []models.Variable, stack []models.StackFrame) models.TraceFrame {
	return BuildFrameWithViz(file, line, editorLine, function, vars, stack, nil)
}

// BuildFrameWithViz constructs a TraceFrame with extended visualization.
func BuildFrameWithViz(
	file string, line, editorLine int, function string,
	vars []models.Variable, stack []models.StackFrame,
	_ *VizBuildConfig,
) models.TraceFrame {
	pointers := viz.DetectPointers(vars)
	vizModels := viz.Extract(vars, pointers)

	if editorLine <= 0 {
		editorLine = line
	}
	return models.TraceFrame{
		File:       basename(file),
		Line:       line,
		EditorLine: editorLine,
		Function:   function,
		Variables:  vars,
		Stack:      stack,
		Viz:        vizModels,
	}
}

// VizBuildConfig is visualization context for a frame.
type VizBuildConfig struct {
	EditorSource string
	LoopIter     int
	Condition    string
}

func buildStepVizModel(cfg *VizBuildConfig, editorLine int, function string) *models.VizModel {
	if cfg == nil || cfg.EditorSource == "" {
		return nil
	}
	return viz.BuildStepViz(cfg.EditorSource, editorLine, function, cfg.LoopIter, cfg.Condition)
}

func upsertStepFlow(vizModels []models.VizModel, step models.VizModel) []models.VizModel {
	out := make([]models.VizModel, 0, len(vizModels)+1)
	replaced := false
	for _, m := range vizModels {
		if m.Type == "step_flow" {
			if !replaced {
				out = append(out, step)
				replaced = true
			}
			continue
		}
		out = append(out, m)
	}
	if !replaced {
		out = append(out, step)
	}
	return out
}

func basename(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[i+1:]
		}
	}
	return path
}
