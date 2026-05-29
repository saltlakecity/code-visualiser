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
  input: string;
  expected: string;
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
  highlight?: number[];
  nodes?: { id: string; label: string; x?: number; y?: number }[];
  edges?: { from: string; to: string }[];
}

export interface TraceFrame {
  index: number;
  file: string;
  line: number;
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
