package viz

import (
	"testing"

	"visualisator/internal/models"
)

func TestExtractSlice(t *testing.T) {
	vars := []models.Variable{
		{
			Name:  "nums",
			Type:  "[]int",
			Value: "[]int len: 4, cap: 4",
			Children: []models.Variable{
				{Name: "[0]", Value: "1"},
				{Name: "[1]", Value: "2"},
				{Name: "[2]", Value: "3"},
				{Name: "[3]", Value: "1"},
			},
		},
		{Name: "i", Type: "int", Value: "1"},
	}
	out := Extract(vars, DetectHighlights(vars))
	if len(out) != 1 {
		t.Fatalf("expected 1 viz model, got %d", len(out))
	}
	if out[0].Type != "array" || len(out[0].Values) != 4 {
		t.Fatalf("unexpected array viz: %+v", out[0])
	}
}
