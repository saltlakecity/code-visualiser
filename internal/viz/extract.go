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
	ptrVarNames = []string{"curr", "current", "head", "node", "prev", "next", "root"}
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

// DetectRange returns left/right range for binary search visualization.
func DetectRange(vars []models.Variable) *models.VizRange {
	left, right := -1, -1
	for _, v := range vars {
		name := strings.ToLower(v.Name)
		val := strings.TrimSpace(v.Value)
		if !intRe.MatchString(val) {
			continue
		}
		idx, _ := strconv.Atoi(val)
		switch name {
		case "left", "low", "start":
			left = idx
		case "right", "high", "end":
			right = idx
		}
	}
	if left >= 0 && right >= 0 && left <= right {
		return &models.VizRange{From: left, To: right}
	}
	return nil
}

// DetectActiveNode finds pointer variable address for graph highlight.
func DetectActiveNode(vars []models.Variable) string {
	for _, name := range ptrVarNames {
		for _, v := range vars {
			if strings.EqualFold(v.Name, name) {
				if strings.Contains(v.Value, "nil") {
					continue
				}
				if strings.HasPrefix(v.Type, "*") || strings.Contains(v.Value, "0x") {
					return v.Value
				}
			}
		}
	}
	return ""
}

// Extract builds visualization models from variables.
func Extract(vars []models.Variable, pointers []models.VizPointer) []models.VizModel {
	globalRange := DetectRange(vars)
	activeAddr := DetectActiveNode(vars)
	var out []models.VizModel
	seen := make(map[string]bool)
	for _, v := range vars {
		if seen[v.Name] {
			continue
		}
		if vm := extractOne(v, pointers, globalRange, activeAddr); vm != nil {
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

func extractOne(v models.Variable, pointers []models.VizPointer, globalRange *models.VizRange, activeAddr string) *models.VizModel {
	v = unwrapVariable(v)
	typeLower := strings.ToLower(strings.TrimSpace(v.Type))
	valLower := strings.ToLower(v.Value)

	if isMatrixType(v) {
		rows := parseMatrixValues(v)
		if len(rows) > 0 {
			filtered := filterPointersForArray(pointers, len(rows)*maxRowLen(rows))
			hl := pointerIndices(filtered)
			return &models.VizModel{
				Type:           "matrix",
				Name:           v.Name,
				Rows:           rows,
				Highlight:      hl,
				Pointers:       filtered,
				RangeHighlight: globalRange,
			}
		}
	}

	isSlice := strings.HasPrefix(typeLower, "[]") ||
		strings.HasPrefix(typeLower, "[") ||
		strings.Contains(valLower, "len:") ||
		strings.Contains(typeLower, "slice")

	if isSlice && !isMatrixType(v) {
		vals := parseSliceValues(v)
		if len(vals) > 0 {
			filtered := filterPointersForArray(pointers, len(vals))
			hl := pointerIndices(filtered)
			return &models.VizModel{
				Type:           "array",
				Name:           v.Name,
				Values:         vals,
				Highlight:      hl,
				Pointers:       filtered,
				RangeHighlight: globalRange,
			}
		}
	}

	if strings.HasPrefix(typeLower, "map[") {
		return extractMap(v)
	}

	if isLinkedListType(v) {
		return extractLinkedList(v, activeAddr)
	}

	if isTreeType(v) {
		return extractTree(v, activeAddr)
	}

	return nil
}

func maxRowLen(rows [][]string) int {
	max := 0
	for _, r := range rows {
		if len(r) > max {
			max = len(r)
		}
	}
	return max
}

func filterPointersForArray(pointers []models.VizPointer, length int) []models.VizPointer {
	if length <= 0 {
		return nil
	}
	var out []models.VizPointer
	for _, p := range pointers {
		if p.Index >= 0 && p.Index < length {
			out = append(out, p)
		}
	}
	return out
}

func pointerIndices(pointers []models.VizPointer) []int {
	seen := make(map[int]bool)
	var hl []int
	for _, p := range pointers {
		if !seen[p.Index] {
			seen[p.Index] = true
			hl = append(hl, p.Index)
		}
	}
	return hl
}

func isMatrixType(v models.Variable) bool {
	typeLower := strings.ToLower(strings.TrimSpace(v.Type))
	return strings.HasPrefix(typeLower, "[][]")
}

func parseMatrixValues(v models.Variable) [][]string {
	if len(v.Children) == 0 {
		return nil
	}
	var rows [][]string
	for _, rowVar := range v.Children {
		vals := parseSliceValues(rowVar)
		if len(vals) > 0 {
			rows = append(rows, vals)
		}
	}
	return rows
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
	var entries []models.VizEntry
	for _, c := range v.Children {
		key := c.Name
		val := c.Value
		if idx := strings.Index(val, ":"); idx >= 0 {
			parts := strings.SplitN(val, ":", 2)
			if key == "" || strings.HasPrefix(key, "[") {
				key = strings.TrimSpace(parts[0])
			}
			val = strings.TrimSpace(parts[1])
		}
		entries = append(entries, models.VizEntry{Key: key, Value: val})
	}
	if len(entries) == 0 {
		return nil
	}
	return &models.VizModel{Type: "map", Name: v.Name, Entries: entries}
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

func extractLinkedList(v models.Variable, activeAddr string) *models.VizModel {
	var nodes []models.VizNode
	var edges []models.VizEdge
	addrMap := make(map[string]string)
	walkList(v, &nodes, &edges, addrMap, 0, 20)
	if len(nodes) == 0 {
		return nil
	}
	highlightNode := matchActiveNode(addrMap, activeAddr)
	for i := range nodes {
		if nodes[i].ID == highlightNode {
			nodes[i].Active = true
		}
	}
	return &models.VizModel{Type: "linked_list", Name: v.Name, Nodes: nodes, Edges: edges, HighlightNode: highlightNode}
}

func walkList(v models.Variable, nodes *[]models.VizNode, edges *[]models.VizEdge, addrMap map[string]string, depth, max int) {
	if depth >= max {
		return
	}
	id := strconv.Itoa(len(*nodes))
	label := v.Value
	addr := extractAddr(v.Value)
	if addr != "" {
		addrMap[addr] = id
	}
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
			walkList(c.Children[0], nodes, edges, addrMap, depth+1, max)
		}
	}
}

func extractAddr(val string) string {
	if idx := strings.Index(val, "0x"); idx >= 0 {
		end := idx + 2
		for end < len(val) && ((val[end] >= '0' && val[end] <= '9') || (val[end] >= 'a' && val[end] <= 'f') || (val[end] >= 'A' && val[end] <= 'F')) {
			end++
		}
		return val[idx:end]
	}
	return ""
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

func extractTree(v models.Variable, activeAddr string) *models.VizModel {
	var nodes []models.VizNode
	var edges []models.VizEdge
	addrMap := make(map[string]string)
	width := treeWidth(v, 4)
	if width < 120 {
		width = 120
	}
	layoutTree(v, &nodes, &edges, addrMap, "0", 0, 0, width, 4)
	if len(nodes) == 0 {
		return nil
	}
	highlightNode := matchActiveNode(addrMap, activeAddr)
	for i := range nodes {
		if nodes[i].ID == highlightNode {
			nodes[i].Active = true
		}
	}
	return &models.VizModel{Type: "tree", Name: v.Name, Nodes: nodes, Edges: edges, HighlightNode: highlightNode}
}

func matchActiveNode(addrMap map[string]string, activeAddr string) string {
	if activeAddr == "" {
		return ""
	}
	for addr, id := range addrMap {
		if strings.Contains(activeAddr, addr) || addr == activeAddr || strings.Contains(addr, activeAddr) {
			return id
		}
	}
	return ""
}

func treeWidth(v models.Variable, maxDepth int) float64 {
	if maxDepth <= 0 {
		return 1
	}
	w := 1.0
	for i := range v.Children {
		c := v.Children[i]
		if c.Name != "Left" && c.Name != "Right" {
			continue
		}
		if strings.Contains(c.Value, "nil") || len(c.Children) == 0 {
			continue
		}
		childW := treeWidth(c.Children[0], maxDepth-1)
		if childW > w {
			w = childW
		}
	}
	return w * 2
}

func layoutTree(v models.Variable, nodes *[]models.VizNode, edges *[]models.VizEdge, addrMap map[string]string, id string, depth int, xOffset float64, span float64, maxDepth int) {
	if depth >= maxDepth {
		return
	}
	label := v.Value
	addr := extractAddr(v.Value)
	if addr != "" {
		addrMap[addr] = id
	}
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
	half := span / 4
	if left != nil && !strings.Contains(left.Value, "nil") && len(left.Children) > 0 {
		lid := id + "L"
		*edges = append(*edges, models.VizEdge{From: id, To: lid})
		layoutTree(left.Children[0], nodes, edges, addrMap, lid, depth+1, x-half, half, maxDepth)
	}
	if right != nil && !strings.Contains(right.Value, "nil") && len(right.Children) > 0 {
		rid := id + "R"
		*edges = append(*edges, models.VizEdge{From: id, To: rid})
		layoutTree(right.Children[0], nodes, edges, addrMap, rid, depth+1, x+half, half, maxDepth)
	}
}
