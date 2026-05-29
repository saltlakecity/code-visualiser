package debugger

import (
	"testing"

	"visualisator/internal/models"
)

func TestTrimTraceFrames(t *testing.T) {
	frames := make([]models.TraceFrame, 8)
	for i := range frames {
		frames[i] = models.TraceFrame{Index: i, Function: "f"}
	}
	out := TrimTraceFrames(frames)
	if len(out) != 4 {
		t.Fatalf("len=%d want 4", len(out))
	}
	if out[0].Index != 0 || out[3].Index != 3 {
		t.Fatalf("reindex failed: %+v", out)
	}
}

func TestTrimTraceFrames_short(t *testing.T) {
	frames := []models.TraceFrame{{Index: 0}, {Index: 1}, {Index: 2}}
	out := TrimTraceFrames(frames)
	if len(out) != 3 {
		t.Fatalf("short trace should be unchanged, len=%d", len(out))
	}
}
