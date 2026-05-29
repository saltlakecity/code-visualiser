package debugger

import "visualisator/internal/models"

const (
	trimHeadFrames = 2
	trimTailFrames = 2
)

// TrimTraceFrames drops the first and last service frames (main entry/exit) when enough frames exist.
func TrimTraceFrames(frames []models.TraceFrame) []models.TraceFrame {
	n := len(frames)
	minKeep := trimHeadFrames + trimTailFrames + 1
	if n < minKeep {
		return frames
	}
	trimmed := make([]models.TraceFrame, n-trimHeadFrames-trimTailFrames)
	copy(trimmed, frames[trimHeadFrames:n-trimTailFrames])
	for i := range trimmed {
		trimmed[i].Index = i
	}
	return trimmed
}
