package viz

import (
	"regexp"
	"strconv"
	"strings"

	"visualisator/internal/models"
)

var (
	intRe       = regexp.MustCompile(`^-?\d+$`)
	sliceElemRe = regexp.MustCompile(`\[(\d+)\]:\s*([^,\]]+)`)
)

// DetectPointers finds index variables and their array indices.
func DetectPointers(vars []models.Variable) []models.VizPointer {
	var out []models.VizPointer
	seen := make(map[string]bool)
	add := func(name string) {
		nameLower := strings.ToLower(name)
		for _, v := range vars {
			if strings.ToLower(v.Name) != nameLower {
				continue
			}
			val := strings.TrimSpace(v.Value)
			if !intRe.MatchString(val) {
				continue
			}
			idx, err := strconv.Atoi(val)
			if err != nil || idx < 0 {
				continue
			}
			key := nameLower + ":" + val
			if seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, models.VizPointer{Name: v.Name, Index: idx})
		}
	}
	for _, hn := range HighlightVars {
		add(hn)
	}
	add("k")
	return out
}

// DetectHighlights returns unique indices from DetectPointers (backward compatible).
func DetectHighlights(vars []models.Variable) []int {
	pointers := DetectPointers(vars)
	seen := make(map[int]bool)
	var highlights []int
	for _, p := range pointers {
		if !seen[p.Index] {
			seen[p.Index] = true
			highlights = append(highlights, p.Index)
		}
	}
	return highlights
}

// Extract builds visualization models from variables.
func Extract(vars []models.Variable, pointers []models.VizPointer) []models.VizModel {
	highlight := make([]int, 0, len(pointers))
	seenIdx := make(map[int]bool)
	for _, p := range pointers {
		if !seenIdx[p.Index] {
			seenIdx[p.Index] = true
			highlight = append(highlight, p.Index)
		}
	}
	var out []models.VizModel
	seen := make(map[string]bool)
	for _, v := range vars {
		if seen[v.Name] {
			continue
		}
		if vm := extractOne(v, highlight, pointers); vm != nil {
			seen[v.Name] = true
			out = append(out, *vm)
		}
	}
	return out
}

func unwrapVariable(v models.Variable) models.Variable {
	typeLower := strings.ToLower(strings.TrimSpace(v.Type))
	if !strings.HasPrefix(typeLower, "*") {
		return v
	}
	if len(v.Children) == 0 {
		return v
	}
	inner := v.Children[0]
	inner.Name = v.Name
	if inner.Type == "" {
		inner.Type = strings.TrimPrefix(typeLower, "*")
	}
	return inner
}

func extractOne(v models.Variable, highlight []int, pointers []models.VizPointer) *models.VizModel {
	v = unwrapVariable(v)
	typeLower := strings.ToLower(strings.TrimSpace(v.Type))
	valLower := strings.ToLower(v.Value)

	isSlice := strings.HasPrefix(typeLower, "[]") ||
		strings.HasPrefix(typeLower, "[") ||
		strings.Contains(valLower, "len:") ||
		strings.Contains(typeLower, "slice")

	if isSlice {
		vals := parseSliceValues(v)
		if len(vals) > 0 {
			return &models.VizModel{
				Type:      "array",
				Name:      v.Name,
				Values:    vals,
				Highlight: highlight,
				Pointers:  pointers,
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
		vals := make([]string, 0, len(v.Children))
		for _, c := range v.Children {
			vals = append(vals, formatChildValue(c))
		}
		if len(vals) > 0 {
			return vals
		}
	}
	inner := strings.TrimSpace(v.Value)
	if matches := sliceElemRe.FindAllStringSubmatch(inner, -1); len(matches) > 0 {
		maxIdx := -1
		for _, m := range matches {
			if i, err := strconv.Atoi(m[1]); err == nil && i > maxIdx {
				maxIdx = i
			}
		}
		if maxIdx >= 0 {
			vals := make([]string, maxIdx+1)
			for i := range vals {
				vals[i] = "?"
			}
			for _, m := range matches {
				i, _ := strconv.Atoi(m[1])
				if i >= 0 && i < len(vals) {
					vals[i] = strings.TrimSpace(m[2])
				}
			}
			return vals
		}
	}
	for _, prefix := range []string{"]int", "[]string", "[]"} {
		if idx := strings.Index(inner, prefix); idx >= 0 {
			inner = inner[idx+len(prefix):]
			break
		}
	}
	inner = strings.TrimPrefix(inner, "len:")
	if i := strings.Index(inner, "len:"); i >= 0 {
		inner = inner[:i]
	}
	inner = strings.Trim(inner, "[]")
	if inner == "" {
		return nil
	}
	parts := strings.Split(inner, ",")
	var vals []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			vals = append(vals, p)
		}
	}
	return vals
}

func formatChildValue(c models.Variable) string {
	val := strings.TrimSpace(c.Value)
	if val == "" && len(c.Children) > 0 {
		val = strings.TrimSpace(c.Children[0].Value)
	}
	if idx := strings.LastIndex(val, ":"); idx >= 0 {
		val = strings.TrimSpace(val[idx+1:])
	}
	return val
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
