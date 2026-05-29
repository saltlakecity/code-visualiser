export interface ProblemMeta {
  id: string;
  title: string;
  tags: string[];
}

export interface Problem extends ProblemMeta {
  template: string;
  entry: string;
  runner?: string;
  tests: TestCase[];
}

export interface TestCase {
  name?: string;
  input: string;
  args?: string;
  setup?: string;
  expected: string;
  hidden?: boolean;
}

export interface TestResult {
  index: number;
  passed: boolean;
  actual?: string;
  expected: string;
  error?: string;
}

export interface CompileError {
  file: string;
  line: number;
  column: number;
  message: string;
}

export interface CompileResult {
  success: boolean;
  sessionId: string;
  binaryPath?: string;
  compiledSource?: string;
  debugSource?: string;
  debugFile?: string;
  editorLineOffset?: number;
  errors?: CompileError[];
}

export interface Variable {
  name: string;
  type: string;
  value: string;
  children?: Variable[];
}

export interface StackFrame {
  file: string;
  line: number;
  function: string;
}

export interface VizScalar {
  name: string;
  type: string;
  value: string;
}

export interface VizNode {
  id: string;
  label: string;
  kind?: string;
  line?: number;
  x?: number;
  y?: number;
  active?: boolean;
}

export interface VizEdge {
  from: string;
  to: string;
  label?: string;
  animated?: boolean;
}

export interface VizModel {
  type: string;
  name: string;
  values?: string[];
  rows?: string[][];
  highlight?: number[];
  highlightNode?: string;
  rangeHighlight?: { from: number; to: number };
  pointers?: { name: string; index: number }[];
  nodes?: VizNode[];
  edges?: VizEdge[];
  entries?: { key: string; value: string }[];
  scalars?: VizScalar[];
  loopIter?: number;
}

export interface FlowchartResult {
  functionName: string;
  summary: string;
  nodes: VizNode[];
  edges: VizEdge[];
}

export interface VarDelta {
  name: string;
  type: string;
  oldValue: string;
  newValue: string;
}

export interface FrameExplanation {
  summary: string;
  lineText: string;
  stmtKind: string;
  condition?: string;
  changes: VarDelta[];
  loopIter?: number;
}

export interface OutlinePoint {
  line: number;
  text: string;
}

export interface SourceOutline {
  name: string;
  params: string;
  returns: string;
  loops?: OutlinePoint[];
  branches?: OutlinePoint[];
  returnsAt?: OutlinePoint[];
  demoCall?: string;
}

export interface AnalyzeResult {
  outline: SourceOutline;
  summary: string;
}

export interface TraceFrame {
  index: number;
  file: string;
  line: number;
  editorLine?: number;
  function: string;
  variables: Variable[];
  stack: StackFrame[];
  viz: VizModel[];
  explanation?: FrameExplanation;
  isFailure?: boolean;
  failureReason?: string;
}

export interface Trace {
  sessionId: string;
  frames: TraceFrame[];
  status: string;
  panicMsg?: string;
  failureFrameIndex?: number;
  failureEditorLine?: number;
  failureReason?: string;
  failureKind?: string;
}

export interface StepResult {
  frame?: TraceFrame;
  trace?: Trace;
  finished: boolean;
  error?: string;
}

export interface RunOptions {
  maxSteps: number;
  lang?: string;
}
