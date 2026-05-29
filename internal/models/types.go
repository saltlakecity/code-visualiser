package models

// ProblemMeta is a catalog entry.
type ProblemMeta struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

// Problem is a full problem definition.
type Problem struct {
	ID       string     `json:"id"`
	Title    string     `json:"title"`
	Tags     []string   `json:"tags"`
	Template string     `json:"template"`
	Entry    string     `json:"entry"`
	Runner   string     `json:"runner"`
	Tests    []TestCase `json:"tests"`
}

// TestCase is input/expected pair.
type TestCase struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

// CompileResult from go build.
type CompileResult struct {
	Success bool           `json:"success"`
	Session string         `json:"sessionId"`
	Binary  string         `json:"binaryPath,omitempty"`
	Errors  []CompileError `json:"errors,omitempty"`
}

// CompileError is a compiler diagnostic.
type CompileError struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Message string `json:"message"`
}

// RunOptions for debug session.
type RunOptions struct {
	MaxSteps int `json:"maxSteps"`
}

// Trace is recorded execution frames.
type Trace struct {
	SessionID string       `json:"sessionId"`
	Frames    []TraceFrame `json:"frames"`
	Status    string       `json:"status"` // running, paused, finished, panic, error
	PanicMsg  string       `json:"panicMsg,omitempty"`
}

// TraceFrame is one debugger snapshot.
type TraceFrame struct {
	Index     int            `json:"index"`
	File      string         `json:"file"`
	Line      int            `json:"line"`
	Function  string         `json:"function"`
	Variables []Variable     `json:"variables"`
	Stack     []StackFrame   `json:"stack"`
	Viz       []VizModel     `json:"viz"`
}

// Variable is a debug variable tree node.
type Variable struct {
	Name     string     `json:"name"`
	Type     string     `json:"type"`
	Value    string     `json:"value"`
	Children []Variable `json:"children,omitempty"`
}

// StackFrame is a call stack entry.
type StackFrame struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
}

// VizModel is visualization data for the UI.
type VizModel struct {
	Type      string   `json:"type"` // array, linked_list, tree, map
	Name      string   `json:"name"`
	Values    []string `json:"values,omitempty"`
	Highlight []int    `json:"highlight,omitempty"`
	Nodes     []VizNode `json:"nodes,omitempty"`
	Edges     []VizEdge `json:"edges,omitempty"`
}

// VizNode for graph visualizations.
type VizNode struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	X     float64 `json:"x,omitempty"`
	Y     float64 `json:"y,omitempty"`
}

// VizEdge connects viz nodes.
type VizEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// StepResult returned after a step command.
type StepResult struct {
	Frame   *TraceFrame `json:"frame,omitempty"`
	Trace   *Trace      `json:"trace,omitempty"`
	Finished bool       `json:"finished"`
	Error   string      `json:"error,omitempty"`
}
