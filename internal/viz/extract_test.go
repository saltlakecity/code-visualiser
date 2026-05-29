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
	pointers := DetectPointers(vars)
	out := Extract(vars, pointers)
	if len(out) != 1 {
		t.Fatalf("expected 1 viz model, got %d", len(out))
	}
	if out[0].Type != "array" || len(out[0].Values) != 4 {
		t.Fatalf("unexpected array viz: %+v", out[0])
	}
	if len(out[0].Pointers) != 1 || out[0].Pointers[0].Name != "i" || out[0].Pointers[0].Index != 1 {
		t.Fatalf("unexpected pointers: %+v", out[0].Pointers)
	}
}

func TestDetectPointers(t *testing.T) {
	vars := []models.Variable{
		{Name: "i", Type: "int", Value: "1"},
		{Name: "j", Type: "int", Value: "3"},
		{Name: "target", Type: "int", Value: "9"},
	}
	pointers := DetectPointers(vars)
	if len(pointers) != 2 {
		t.Fatalf("expected 2 pointers, got %d: %+v", len(pointers), pointers)
	}
	found := map[string]int{}
	for _, p := range pointers {
		found[p.Name] = p.Index
	}
	if found["i"] != 1 || found["j"] != 3 {
		t.Fatalf("unexpected pointer indices: %+v", found)
	}
}

func TestExtractSliceFromValueString(t *testing.T) {
	v := models.Variable{
		Name:  "nums",
		Type:  "[]int",
		Value: "[]int len: 4, cap: 4, [0]:2, [1]:7, [2]:11, [3]:15",
	}
	vals := parseSliceValues(v)
	if len(vals) != 4 {
		t.Fatalf("expected 4 values, got %d: %v", len(vals), vals)
	}
	if vals[0] != "2" || vals[3] != "15" {
		t.Fatalf("unexpected values: %v", vals)
	}
}

func TestDetectRange(t *testing.T) {
	vars := []models.Variable{
		{Name: "left", Value: "1"},
		{Name: "right", Value: "5"},
	}
	r := DetectRange(vars)
	if r == nil || r.From != 1 || r.To != 5 {
		t.Fatalf("unexpected range: %+v", r)
	}
}

func TestExtractLinkedListPointer(t *testing.T) {
	vars := []models.Variable{
		{
			Name:  "head",
			Type:  "*main.ListNode",
			Value: "(*main.ListNode)(0xc0000b4000)",
			Children: []models.Variable{
				{
					Name:  "main.ListNode",
					Type:  "main.ListNode",
					Value: "main.ListNode {Val: 1, Next: *main.ListNode}",
					Children: []models.Variable{
						{Name: "Val", Type: "int", Value: "1"},
						{
							Name:  "Next",
							Type:  "*main.ListNode",
							Value: "(*main.ListNode)(0xc0000b4020)",
							Children: []models.Variable{
								{
									Name: "main.ListNode",
									Children: []models.Variable{
										{Name: "Val", Value: "2"},
										{Name: "Next", Value: "nil"},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	out := Extract(vars, nil)
	if len(out) != 1 {
		t.Fatalf("expected 1 viz model, got %d", len(out))
	}
	if out[0].Type != "linked_list" {
		t.Fatalf("expected linked_list, got %s", out[0].Type)
	}
	if len(out[0].Nodes) < 2 {
		t.Fatalf("expected at least 2 nodes, got %d", len(out[0].Nodes))
	}
}
