package analyzer

import (
	"strings"
	"testing"
)

const twoSumSrc = `func twoSum(nums []int, target int) []int {
	m := make(map[int]int)
	for i, v := range nums {
		if need, ok := m[target-v]; ok {
			return []int{need, i}
		}
		m[v] = i
	}
	return nil
}`

func TestParseSource_twoSum(t *testing.T) {
	sm, err := ParseSource(twoSumSrc)
	if err != nil {
		t.Fatal(err)
	}
	if sm.Outline.Name != "twoSum" {
		t.Fatalf("name=%q", sm.Outline.Name)
	}
	if !strings.Contains(sm.Outline.Params, "nums") {
		t.Fatalf("params=%q", sm.Outline.Params)
	}
	if len(sm.Lines) < 3 {
		t.Fatalf("expected multiple line entries, got %d", len(sm.Lines))
	}
	info, ok := sm.StmtAt(2)
	if !ok {
		t.Fatal("no stmt at line 2")
	}
	if info.Kind != "assign" {
		t.Fatalf("line 2 kind=%q", info.Kind)
	}
	var loopLine int
	var loopKind string
	for _, info := range sm.Lines {
		if info.Kind == "for" || info.Kind == "range" {
			loopLine = info.Line
			loopKind = info.Kind
			break
		}
	}
	if loopLine == 0 {
		t.Fatal("no loop stmt")
	}
	if info, ok := sm.StmtAt(loopLine); !ok || info.Kind != loopKind {
		t.Fatalf("loop at %d: %+v ok=%v", loopLine, info, ok)
	}
}

func TestStmtAt_fallback(t *testing.T) {
	sm, err := ParseSource(twoSumSrc)
	if err != nil {
		t.Fatal(err)
	}
	info, ok := sm.StmtAt(4)
	if !ok {
		t.Fatal("expected fallback stmt")
	}
	if info.Text == "" {
		t.Fatal("empty text")
	}
}
