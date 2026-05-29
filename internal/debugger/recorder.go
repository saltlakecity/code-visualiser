package debugger

import (
	"sync"

	"visualisator/internal/models"
	"visualisator/internal/viz"
)

const maxFrames = 10000

// Recorder stores trace frames for replay.
type Recorder struct {
	mu      sync.RWMutex
	frames  []models.TraceFrame
	status  string
	panicMsg string
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
}

// Append adds a frame if under limit.
func (r *Recorder) Append(frame models.TraceFrame) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.frames) >= maxFrames {
		return
	}
	frame.Index = len(r.frames)
	r.frames = append(r.frames, frame)
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
	frames := make([]models.TraceFrame, len(r.frames))
	copy(frames, r.frames)
	return models.Trace{
		SessionID: sessionID,
		Frames:    frames,
		Status:    r.status,
		PanicMsg:  r.panicMsg,
	}
}

// FrameAt returns frame by index.
func (r *Recorder) FrameAt(index int) *models.TraceFrame {
	r.mu.RLock()
	defer r.mu.Unlock()
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
func BuildFrame(file string, line int, function string, vars []models.Variable, stack []models.StackFrame) models.TraceFrame {
	highlight := viz.DetectHighlights(vars)
	vizModels := viz.Extract(vars, highlight)
	return models.TraceFrame{
		File:      basename(file),
		Line:      line,
		Function:  function,
		Variables: vars,
		Stack:     stack,
		Viz:       vizModels,
	}
}

func basename(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[i+1:]
		}
	}
	return path
}
