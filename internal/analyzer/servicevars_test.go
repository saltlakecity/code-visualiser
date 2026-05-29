package analyzer

import (
	"testing"

	"visualisator/internal/models"
)

func TestFilterVariables(t *testing.T) {
	vars := []models.Variable{
		{Name: "i", Value: "0"},
		{Name: "~r0", Value: "1"},
		{Name: "nums", Value: "[1,2]"},
	}
	out := FilterVariables(vars)
	if len(out) != 2 {
		t.Fatalf("got %d vars", len(out))
	}
}

func TestFilterVarDeltas(t *testing.T) {
	d := FilterVarDeltas([]VarDelta{
		{Name: "i", OldValue: "0", NewValue: "1"},
		{Name: "~r0", OldValue: "—", NewValue: "0"},
	})
	if len(d) != 1 || d[0].Name != "i" {
		t.Fatalf("%+v", d)
	}
}
