import { create } from 'zustand';
import type { CompileError, Problem, ProblemMeta, TestCase, TestResult, Trace, TraceFrame } from '../types';

interface DebuggerState {
  problems: ProblemMeta[];
  currentProblem: Problem | null;
  customMode: boolean;
  source: string;
  debugSource: string;
  debugFile: string;
  editorLineOffset: number;
  sessionId: string;
  trace: Trace | null;
  frameIndex: number;
  compileErrors: CompileError[];
  status: string;
  panicMsg: string;
  isRunning: boolean;
  customTests: TestCase[];
  selectedTestIndex: number;
  testResults: TestResult[];
  highlightedVar: string | null;
  setProblems: (p: ProblemMeta[]) => void;
  setCurrentProblem: (p: Problem | null) => void;
  setCustomMode: (v: boolean) => void;
  setSource: (s: string) => void;
  setDebugSource: (s: string, file?: string, offset?: number) => void;
  setSessionId: (id: string) => void;
  setTrace: (t: Trace | null) => void;
  setFrameIndex: (i: number) => void;
  setCompileErrors: (e: CompileError[]) => void;
  setStatus: (s: string) => void;
  setPanicMsg: (m: string) => void;
  setIsRunning: (v: boolean) => void;
  setCustomTests: (t: TestCase[]) => void;
  setSelectedTestIndex: (i: number) => void;
  setTestResults: (r: TestResult[]) => void;
  setHighlightedVar: (name: string | null) => void;
  resetTrace: () => void;
  currentFrame: () => TraceFrame | null;
  prevFrame: () => TraceFrame | null;
  activeTests: () => TestCase[];
}

export const useDebuggerStore = create<DebuggerState>((set, get) => ({
  problems: [],
  currentProblem: null,
  customMode: false,
  source: defaultCustomSource(),
  debugSource: '',
  debugFile: '',
  editorLineOffset: 0,
  sessionId: '',
  trace: null,
  frameIndex: 0,
  compileErrors: [],
  status: 'idle',
  panicMsg: '',
  isRunning: false,
  customTests: defaultCustomTests(),
  selectedTestIndex: 0,
  testResults: [],
  highlightedVar: null,
  setProblems: (problems) => set({ problems }),
  setCurrentProblem: (currentProblem) =>
    set({ currentProblem, selectedTestIndex: 0, testResults: [] }),
  setCustomMode: (customMode) => set({ customMode }),
  setSource: (source) => set({ source }),
  setDebugSource: (debugSource, debugFile = '', offset = 0) =>
    set((s) => ({
      debugSource,
      debugFile: debugFile || s.debugFile,
      editorLineOffset: offset,
    })),
  setSessionId: (sessionId) => set({ sessionId }),
  setTrace: (trace) =>
    set({ trace, frameIndex: trace?.frames?.length ? trace.frames.length - 1 : 0 }),
  setFrameIndex: (frameIndex) => set({ frameIndex }),
  setCompileErrors: (compileErrors) => set({ compileErrors }),
  setStatus: (status) => set({ status }),
  setPanicMsg: (panicMsg) => set({ panicMsg }),
  setIsRunning: (isRunning) => set({ isRunning }),
  setCustomTests: (customTests) => set({ customTests }),
  setSelectedTestIndex: (selectedTestIndex) => set({ selectedTestIndex }),
  setTestResults: (testResults) => set({ testResults }),
  setHighlightedVar: (highlightedVar) => set({ highlightedVar }),
  resetTrace: () =>
    set({
      trace: null,
      frameIndex: 0,
      sessionId: '',
      debugSource: '',
      debugFile: '',
      editorLineOffset: 0,
      status: 'idle',
      panicMsg: '',
      isRunning: false,
    }),
  currentFrame: () => {
    const { trace, frameIndex } = get();
    if (!trace?.frames?.length) return null;
    return trace.frames[frameIndex] ?? null;
  },
  prevFrame: () => {
    const { trace, frameIndex } = get();
    if (!trace?.frames?.length || frameIndex <= 0) return null;
    return trace.frames[frameIndex - 1] ?? null;
  },
  activeTests: () => {
    const { customMode, currentProblem, customTests } = get();
    if (customMode) return customTests;
    return currentProblem?.tests ?? [];
  },
}));

function defaultCustomSource(): string {
  return `package main

import "fmt"

func main() {
\tnums := []int{2, 7, 11, 15}
\ttarget := 9
\tfor i := 0; i < len(nums); i++ {
\t\tfor j := i + 1; j < len(nums); j++ {
\t\t\tif nums[i]+nums[j] == target {
\t\t\t\tfmt.Println(i, j)
\t\t\t\treturn
\t\t\t}
\t\t}
\t}
}
`;
}

function defaultCustomTests(): TestCase[] {
  return [
    {
      name: 'Case 1',
      input: 'nums=[2,7,11,15], target=9',
      expected: '0 1',
    },
  ];
}
