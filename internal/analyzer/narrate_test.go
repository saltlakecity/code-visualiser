package analyzer

import (
	"strings"
	"testing"

	"visualisator/internal/models"
)

func TestEvalCondition_simple(t *testing.T) {
	vars := map[string]string{"i": "0", "n": "3"}
	ok, detail, worked := EvalCondition("i < n", vars)
	if !worked {
		t.Fatal("eval failed", detail)
	}
	if !ok {
		t.Fatal("expected true")
	}
	if detail == "" {
		t.Fatal("empty detail")
	}
}

func TestEvalCondition_equality(t *testing.T) {
	vars := map[string]string{"x": "5", "y": "5"}
	ok, _, worked := EvalCondition("x == y", vars)
	if !worked || !ok {
		t.Fatalf("got ok=%v worked=%v", ok, worked)
	}
}

func TestComputeVarDeltas(t *testing.T) {
	prev := []models.Variable{{Name: "i", Type: "int", Value: "0"}}
	curr := []models.Variable{{Name: "i", Type: "int", Value: "1"}}
	d := ComputeVarDeltas(prev, curr)
	if len(d) != 1 || d[0].NewValue != "1" {
		t.Fatalf("%+v", d)
	}
}

func TestNarrateStep_assignRU(t *testing.T) {
	sm, err := ParseSource("func f() {\n\tx := 1\n}")
	if err != nil {
		t.Fatal(err)
	}
	curr := []models.Variable{{Name: "x", Type: "int", Value: "1"}}
	exp := NarrateStep(LangRU, sm, NewNarrateContext(), 2, "main.f", nil, curr, nil, 1)
	if exp == nil || !strings.Contains(exp.Summary, "Присваиваем") {
		t.Fatalf("%+v", exp)
	}
}

func TestNarrateStep_enter(t *testing.T) {
	sm, _ := ParseSource("func add(a int, b int) int {\n\treturn a + b\n}")
	curr := []models.Variable{
		{Name: "a", Type: "int", Value: "2"},
		{Name: "b", Type: "int", Value: "3"},
	}
	exp := NarrateStep(LangEN, sm, NewNarrateContext(), 1, "main.add", nil, curr, nil, 0)
	if !strings.Contains(exp.Summary, "Enter add") {
		t.Fatalf("%q", exp.Summary)
	}
}

func TestFormatConditionResult(t *testing.T) {
	s := FormatConditionResult("i < n", true, "0 < 3", LangRU)
	if !strings.Contains(s, "истина") {
		t.Fatal(s)
	}
}

func TestExplainFailure_index(t *testing.T) {
	s := ExplainFailure(LangEN, nil, 5, "runtime error: index out of range [3] with length 2")
	if !strings.Contains(s, "line 5") || !strings.Contains(s, "index out of bounds") {
		t.Fatal(s)
	}
}
