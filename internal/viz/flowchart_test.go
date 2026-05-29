package viz

import (
	"strings"
	"testing"

	"visualisator/internal/models"
)

func TestBuildFlowGraph_simple(t *testing.T) {
	src := `func sum(nums []int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}`
	result, err := BuildFlowGraph(src)
	if err != nil {
		t.Fatal(err)
	}
	if result.FunctionName != "sum" {
		t.Fatalf("function name = %q", result.FunctionName)
	}
	if len(result.Nodes) < 3 {
		t.Fatalf("expected at least 3 nodes, got %d", len(result.Nodes))
	}
	kinds := map[string]bool{}
	for _, n := range result.Nodes {
		kinds[n.Kind] = true
	}
	if !kinds["start"] || !kinds["end"] {
		t.Fatalf("missing start/end nodes: %+v", kinds)
	}
}

func TestBuildFlowGraph_binarySearch(t *testing.T) {
	src := `func search(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			return mid
		}
		if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}`
	result, err := BuildFlowGraph(src)
	if err != nil {
		t.Fatal(err)
	}
	if result.FunctionName != "search" {
		t.Fatalf("name = %q", result.FunctionName)
	}

	decisions := 0
	returns := 0
	var endID string
	for _, n := range result.Nodes {
		if n.Kind == "decision" {
			decisions++
		}
		if n.Kind == "return" {
			returns++
		}
		if n.Kind == "end" {
			endID = n.ID
		}
	}
	if decisions < 2 {
		t.Fatalf("expected >=2 decision nodes, got %d", decisions)
	}
	if returns < 2 {
		t.Fatalf("expected >=2 return nodes, got %d", returns)
	}
	if endID == "" {
		t.Fatal("missing end node")
	}

	// return inside if must connect to End, not to code after loop
	returnToEnd := false
	returnToAfterLoop := false
	var afterLoopID string
	for _, n := range result.Nodes {
		if n.Kind == "process" && strings.Contains(n.Label, "exit loop") {
			afterLoopID = n.ID
		}
	}
	for _, e := range result.Edges {
		for _, n := range result.Nodes {
			if n.Kind == "return" && n.Label == "return" && e.From == n.ID && e.To == endID {
				returnToEnd = true
			}
			if n.Kind == "return" && afterLoopID != "" && e.From == n.ID && e.To == afterLoopID {
				returnToAfterLoop = true
			}
		}
	}
	if !returnToEnd {
		t.Fatal("expected return edge to End")
	}
	if returnToAfterLoop {
		t.Fatal("return should not connect to after-loop node")
	}

	// branches should have different X
	var xs []float64
	for _, n := range result.Nodes {
		if n.Kind == "decision" {
			xs = append(xs, n.X)
		}
	}
	if len(xs) >= 2 && xs[0] == xs[1] {
		t.Log("warning: decision nodes share X", xs)
	}
}

func TestBuildFlowGraph_summaryNoAST(t *testing.T) {
	src := `func search(nums []int, target int) int { return -1 }`
	result, err := BuildFlowGraph(src)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(result.Summary, "*ast.") {
		t.Fatalf("summary contains AST dump: %q", result.Summary)
	}
}

func TestBuildStepViz_highlightsLine(t *testing.T) {
	src := `func f(n int) int {
	if n <= 2 {
		return n
	}
	return f(n-1) + f(n-2)
}`
	step := BuildStepViz(src, 2, "f", 0, "")
	if step == nil {
		t.Fatal("expected step viz")
	}
	if step.Type != "step_flow" {
		t.Fatalf("type = %q", step.Type)
	}
	active := false
	for _, n := range step.Nodes {
		if n.Active {
			active = true
		}
	}
	if !active {
		t.Fatal("no active node")
	}
}

func TestBuildScalars(t *testing.T) {
	vars := []models.Variable{
		{Name: "a", Type: "int", Value: "1"},
		{Name: "b", Type: "int", Value: "2"},
		{Name: "nums", Type: "[]int", Value: "len:3"},
	}
	m := BuildScalars(vars)
	if m == nil {
		t.Fatal("expected scalars")
	}
	if len(m.Scalars) != 2 {
		t.Fatalf("scalars count = %d", len(m.Scalars))
	}
}

func TestExtractString(t *testing.T) {
	v := models.Variable{Name: "s", Type: "string", Value: `"abc"`}
	vm := extractString(v, nil)
	if vm == nil || len(vm.Values) != 3 {
		t.Fatalf("string viz = %+v", vm)
	}
}

func TestHasDataViz(t *testing.T) {
	if HasDataViz([]models.VizModel{{Type: "scalars"}}) {
		t.Fatal("scalars should not count as data viz")
	}
	if !HasDataViz([]models.VizModel{{Type: "array"}}) {
		t.Fatal("array should count")
	}
}

func TestBuildFlowGraph_empty(t *testing.T) {
	_, err := BuildFlowGraph("   ")
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestExtractStackName(t *testing.T) {
	v := models.Variable{
		Name:  "stack",
		Type:  "[]rune",
		Value: "len:2",
		Children: []models.Variable{
			{Name: "[0]", Value: "40"},
			{Name: "[1]", Value: "41"},
		},
	}
	vm := extractOne(v, nil, nil, "")
	if vm == nil {
		t.Fatal("expected stack viz")
	}
	if vm.Type != "stack" {
		t.Fatalf("type = %q, want stack", vm.Type)
	}
}
