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

export interface VizModel {
  type: string;
  name: string;
  values?: string[];
  rows?: string[][];
  highlight?: number[];
  highlightNode?: string;
  rangeHighlight?: { from: number; to: number };
  pointers?: { name: string; index: number }[];
  nodes?: { id: string; label: string; x?: number; y?: number; active?: boolean }[];
  edges?: { from: string; to: string }[];
  entries?: { key: string; value: string }[];
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
}

export interface Trace {
  sessionId: string;
  frames: TraceFrame[];
  status: string;
  panicMsg?: string;
}

export interface StepResult {
  frame?: TraceFrame;
  trace?: Trace;
  finished: boolean;
  error?: string;
}

export interface RunOptions {
  maxSteps: number;
}
