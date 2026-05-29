package analyzer

// Lang is UI language for narration.
type Lang string

const (
	LangEN Lang = "en"
	LangRU Lang = "ru"
)

// StatementInfo describes an executable statement at a source line.
type StatementInfo struct {
	Kind          string   `json:"kind"` // assign, inc, if, for, range, return, call, switch, defer, other
	Text          string   `json:"text"`
	ConditionExpr string   `json:"conditionExpr,omitempty"`
	Reads         []string `json:"reads,omitempty"`
	Writes        []string `json:"writes,omitempty"`
	Line          int      `json:"line"`
}

// FunctionOutline is a static overview of the target function.
type FunctionOutline struct {
	Name       string         `json:"name"`
	Params     string         `json:"params"`
	Returns    string         `json:"returns"`
	Loops      []OutlinePoint `json:"loops,omitempty"`
	Branches   []OutlinePoint `json:"branches,omitempty"`
	ReturnsAt  []OutlinePoint `json:"returnsAt,omitempty"`
	DemoCall   string         `json:"demoCall,omitempty"`
}

// OutlinePoint marks a control-flow feature at a line.
type OutlinePoint struct {
	Line int    `json:"line"`
	Text string `json:"text"`
}

// SourceMap maps editor line numbers to statement info.
type SourceMap struct {
	Lines   map[int]StatementInfo `json:"lines"`
	Outline FunctionOutline       `json:"outline"`
	Source  string                `json:"-"`
}

// VarDelta is a variable change between steps.
type VarDelta struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	OldValue string `json:"oldValue"`
	NewValue string `json:"newValue"`
}

// FrameExplanation is human-readable step breakdown for the UI.
type FrameExplanation struct {
	Summary   string     `json:"summary"`
	LineText  string     `json:"lineText"`
	StmtKind  string     `json:"stmtKind"`
	Condition string     `json:"condition,omitempty"`
	Changes   []VarDelta `json:"changes"`
	LoopIter  int        `json:"loopIter,omitempty"`
}
