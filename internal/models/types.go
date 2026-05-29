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
	Name     string `json:"name,omitempty"`
	Input    string `json:"input"`
	Args     string `json:"args,omitempty"`
	Setup    string `json:"setup,omitempty"`
	Expected string `json:"expected"`
	Hidden   bool   `json:"hidden,omitempty"`
}

// TestResult is the outcome of running one test.
type TestResult struct {
	Index    int    `json:"index"`
	Passed   bool   `json:"passed"`
	Actual   string `json:"actual,omitempty"`
	Expected string `json:"expected"`
	Error    string `json:"error,omitempty"`
}

// AnalyzeResult is static analysis before debug run.
type AnalyzeResult struct {
	Outline SourceOutline `json:"outline"`
	Summary string        `json:"summary"`
}

// CompileResult from go build.
type CompileResult struct {
	Success          bool           `json:"success"`
	Session          string         `json:"sessionId"`
	Binary           string         `json:"binaryPath,omitempty"`
	CompiledSource   string         `json:"compiledSource,omitempty"`
	DebugSource      string         `json:"debugSource,omitempty"`
	DebugFile        string         `json:"debugFile,omitempty"`
	EditorLineOffset int            `json:"editorLineOffset"`
	Errors           []CompileError `json:"errors,omitempty"`
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
	MaxSteps int    `json:"maxSteps"`
	Lang     string `json:"lang,omitempty"` // "ru" or "en"
}

// Trace is recorded execution frames.
type Trace struct {
	SessionID           string       `json:"sessionId"`
	Frames              []TraceFrame `json:"frames"`
	Status              string       `json:"status"` // running, paused, finished, panic, error
	PanicMsg            string       `json:"panicMsg,omitempty"`
	FailureFrameIndex   int          `json:"failureFrameIndex,omitempty"`
	FailureEditorLine   int          `json:"failureEditorLine,omitempty"`
	FailureReason       string       `json:"failureReason,omitempty"`
	FailureKind         string       `json:"failureKind,omitempty"` // panic, runtime, error
}

// VarDelta is a variable change between execution steps.
type VarDelta struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	OldValue string `json:"oldValue"`
	NewValue string `json:"newValue"`
}

// FrameExplanation is a human-readable breakdown of one execution step.
type FrameExplanation struct {
	Summary   string     `json:"summary"`
	LineText  string     `json:"lineText"`
	StmtKind  string     `json:"stmtKind"`
	Condition string     `json:"condition,omitempty"`
	Changes   []VarDelta `json:"changes"`
	LoopIter  int        `json:"loopIter,omitempty"`
}

// SourceOutline is static function overview for the breakdown tab.
type SourceOutline struct {
	Name      string         `json:"name"`
	Params    string         `json:"params"`
	Returns   string         `json:"returns"`
	Loops     []OutlinePoint `json:"loops,omitempty"`
	Branches  []OutlinePoint `json:"branches,omitempty"`
	ReturnsAt []OutlinePoint `json:"returnsAt,omitempty"`
	DemoCall  string         `json:"demoCall,omitempty"`
}

// OutlinePoint marks control-flow at a line.
type OutlinePoint struct {
	Line int    `json:"line"`
	Text string `json:"text"`
}

// TraceFrame is one debugger snapshot.
type TraceFrame struct {
	Index         int               `json:"index"`
	File          string            `json:"file"`
	Line          int               `json:"line"`
	EditorLine    int               `json:"editorLine"`
	Function      string            `json:"function"`
	Variables     []Variable        `json:"variables"`
	Stack         []StackFrame      `json:"stack"`
	Viz           []VizModel        `json:"viz"`
	Explanation   *FrameExplanation `json:"explanation,omitempty"`
	IsFailure     bool              `json:"isFailure,omitempty"`
	FailureReason string            `json:"failureReason,omitempty"`
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

// VizPointer labels an index variable pointing into an array.
type VizPointer struct {
	Name  string `json:"name"`
	Index int    `json:"index"`
}

// VizScalar is a single scalar variable for card display.
type VizScalar struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// VizModel is visualization data for the UI.
type VizModel struct {
	Type           string       `json:"type"` // array, matrix, linked_list, tree, map, string, stack, scalars, step_flow, flowchart
	Name           string       `json:"name"`
	Values         []string     `json:"values,omitempty"`
	Rows           [][]string   `json:"rows,omitempty"`
	Highlight      []int        `json:"highlight,omitempty"`
	HighlightNode  string       `json:"highlightNode,omitempty"`
	RangeHighlight *VizRange    `json:"rangeHighlight,omitempty"`
	Pointers       []VizPointer `json:"pointers,omitempty"`
	Nodes          []VizNode    `json:"nodes,omitempty"`
	Edges          []VizEdge    `json:"edges,omitempty"`
	Entries        []VizEntry   `json:"entries,omitempty"`
	Scalars        []VizScalar  `json:"scalars,omitempty"`
	LoopIter       int          `json:"loopIter,omitempty"`
}

// FlowchartResult is a static control-flow diagram for a function.
type FlowchartResult struct {
	FunctionName string    `json:"functionName"`
	Summary      string    `json:"summary"`
	Nodes        []VizNode `json:"nodes"`
	Edges        []VizEdge `json:"edges"`
}

// VizRange highlights an inclusive index range in an array.
type VizRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// VizEntry is a map key-value pair for table display.
type VizEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// VizNode for graph visualizations.
type VizNode struct {
	ID     string  `json:"id"`
	Label  string  `json:"label"`
	Kind   string  `json:"kind,omitempty"` // start, process, decision, loop, return, end
	Line   int     `json:"line,omitempty"`
	X      float64 `json:"x,omitempty"`
	Y      float64 `json:"y,omitempty"`
	Active bool    `json:"active,omitempty"`
}

// VizEdge connects viz nodes.
type VizEdge struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Label    string `json:"label,omitempty"`
	Animated bool   `json:"animated,omitempty"`
}

// StepResult returned after a step command.
type StepResult struct {
	Frame   *TraceFrame `json:"frame,omitempty"`
	Trace   *Trace      `json:"trace,omitempty"`
	Finished bool       `json:"finished"`
	Error   string      `json:"error,omitempty"`
}
