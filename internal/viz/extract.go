package viz

import (
	"regexp"
	"strconv"
	"strings"

	"visualisator/internal/models"
)

var intRe = regexp.MustCompile(`^-?\d+$`)

// DetectHighlights finds index variables to highlight in arrays.
func DetectHighlights(vars []models.Variable) []int {
	var highlights []int
	for _, v := range vars {
		nameLower := strings.ToLower(v.Name)
		for _, hn := range HighlightVars {
			if nameLower == hn && intRe.MatchString(strings.TrimSpace(v.Value)) {
				if idx, err := strconv.Atoi(strings.TrimSpace(v.Value)); err == nil {
					highlights = append(highlights, idx)
				}
			}
		}
	}
	return highlights
}

// Extract builds visualization models from variables.
func Extract(vars []models.Variable, highlight []int) []models.VizModel {
	var out []models.VizModel
	for _, v := range vars {
		if vm := extractOne(v, highlight); vm != nil {
			out = append(out, *vm)
		}
	}
	return out
}

func extractOne(v models.Variable, highlight []int) *models.VizModel {
	typeLower := strings.ToLower(v.Type)

	if strings.HasPrefix(typeLower, "[]") || strings.HasPrefix(typeLower, "[") {
		vals := parseSliceValues(v)
		if len(vals) > 0 {
			return &models.VizModel{
				Type:      "array",
				Name:      v.Name,
				Values:    vals,
				Highlight: highlight,
			}
		}
	}

	if strings.HasPrefix(typeLower, "map[") {
		return extractMap(v)
	}

	if isLinkedListType(v) {
		return extractLinkedList(v)
	}

	if isTreeType(v) {
		return extractTree(v)
	}

	return nil
}

func parseSliceValues(v models.Variable) []string {
	if len(v.Children) > 0 {
		var vals []string
		for _, c := range v.Children {
			vals = append(vals, strings.TrimSpace(c.Value))
		}
		return vals
	}
	inner := strings.Trim(v.Value, "[]")
	if inner == "" {
		return nil
	}
	parts := strings.Split(inner, ",")
	var vals []string
	for _, p := range parts {
		vals = append(vals, strings.TrimSpace(p))
	}
	return vals
}

func extractMap(v models.Variable) *models.VizModel {
	var nodes []models.VizNode
	for i, c := range v.Children {
		nodes = append(nodes, models.VizNode{
			ID:    strconv.Itoa(i),
			Label: c.Name + ": " + c.Value,
		})
	}
	if len(nodes) == 0 {
		return nil
	}
	return &models.VizModel{Type: "map", Name: v.Name, Nodes: nodes}
}

func isLinkedListType(v models.Variable) bool {
	hasVal, hasNext := false, false
	for _, c := range v.Children {
		switch c.Name {
		case "Val", "Value":
			hasVal = true
		case "Next":
			hasNext = true
		}
	}
	return hasVal && hasNext
}

func extractLinkedList(v models.Variable) *models.VizModel {
	var nodes []models.VizNode
	var edges []models.VizEdge
	walkList(v, &nodes, &edges, 0, 20)
	if len(nodes) == 0 {
		return nil
	}
	return &models.VizModel{Type: "linked_list", Name: v.Name, Nodes: nodes, Edges: edges}
}

func walkList(v models.Variable, nodes *[]models.VizNode, edges *[]models.VizEdge, depth, max int) {
	if depth >= max {
		return
	}
	id := strconv.Itoa(len(*nodes))
	label := v.Value
	for _, c := range v.Children {
		if c.Name == "Val" || c.Name == "Value" {
			label = c.Value
		}
	}
	*nodes = append(*nodes, models.VizNode{ID: id, Label: label, X: float64(depth) * 120, Y: 40})
	for _, c := range v.Children {
		if c.Name != "Next" {
			continue
		}
		if strings.Contains(c.Value, "nil") {
			return
		}
		if len(c.Children) > 0 {
			nextID := strconv.Itoa(len(*nodes))
			*edges = append(*edges, models.VizEdge{From: id, To: nextID})
			walkList(c.Children[0], nodes, edges, depth+1, max)
		}
	}
}

func isTreeType(v models.Variable) bool {
	hasVal, hasLeft, hasRight := false, false, false
	for _, c := range v.Children {
		switch c.Name {
		case "Val", "Value":
			hasVal = true
		case "Left":
			hasLeft = true
		case "Right":
			hasRight = true
		}
	}
	return hasVal && hasLeft && hasRight
}

func extractTree(v models.Variable) *models.VizModel {
	var nodes []models.VizNode
	var edges []models.VizEdge
	layoutTree(v, &nodes, &edges, "0", 0, 0, 4)
	if len(nodes) == 0 {
		return nil
	}
	return &models.VizModel{Type: "tree", Name: v.Name, Nodes: nodes, Edges: edges}
}

func layoutTree(v models.Variable, nodes *[]models.VizNode, edges *[]models.VizEdge, id string, depth int, xOffset float64, maxDepth int) {
	if depth >= maxDepth {
		return
	}
	label := v.Value
	var left, right *models.Variable
	for i := range v.Children {
		c := v.Children[i]
		switch c.Name {
		case "Val", "Value":
			label = c.Value
		case "Left":
			left = &v.Children[i]
		case "Right":
			right = &v.Children[i]
		}
	}
	x := xOffset
	y := float64(depth) * 80
	*nodes = append(*nodes, models.VizNode{ID: id, Label: label, X: x, Y: y})
	span := 120.0 / float64(depth+1)
	if left != nil && !strings.Contains(left.Value, "nil") {
		lid := id + "L"
		*edges = append(*edges, models.VizEdge{From: id, To: lid})
		if len(left.Children) > 0 {
			layoutTree(left.Children[0], nodes, edges, lid, depth+1, x-span, maxDepth)
		}
	}
	if right != nil && !strings.Contains(right.Value, "nil") {
		rid := id + "R"
		*edges = append(*edges, models.VizEdge{From: id, To: rid})
		if len(right.Children) > 0 {
			layoutTree(right.Children[0], nodes, edges, rid, depth+1, x+span, maxDepth)
		}
	}
}
