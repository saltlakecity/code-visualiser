import { create } from 'zustand';
import type { CompileError, Problem, ProblemMeta, Trace, TraceFrame } from '../types';

interface DebuggerState {
  problems: ProblemMeta[];
  currentProblem: Problem | null;
  customMode: boolean;
  source: string;
  sessionId: string;
  trace: Trace | null;
  frameIndex: number;
  compileErrors: CompileError[];
  status: string;
  panicMsg: string;
  isRunning: boolean;
  setProblems: (p: ProblemMeta[]) => void;
  setCurrentProblem: (p: Problem | null) => void;
  setCustomMode: (v: boolean) => void;
  setSource: (s: string) => void;
  setSessionId: (id: string) => void;
  setTrace: (t: Trace | null) => void;
  setFrameIndex: (i: number) => void;
  setCompileErrors: (e: CompileError[]) => void;
  setStatus: (s: string) => void;
  setPanicMsg: (m: string) => void;
  setIsRunning: (v: boolean) => void;
  currentFrame: () => TraceFrame | null;
}

export const useDebuggerStore = create<DebuggerState>((set, get) => ({
  problems: [],
  currentProblem: null,
  customMode: false,
  source: defaultCustomSource(),
  sessionId: '',
  trace: null,
  frameIndex: 0,
  compileErrors: [],
  status: 'idle',
  panicMsg: '',
  isRunning: false,
  setProblems: (problems) => set({ problems }),
  setCurrentProblem: (currentProblem) => set({ currentProblem }),
  setCustomMode: (customMode) => set({ customMode }),
  setSource: (source) => set({ source }),
  setSessionId: (sessionId) => set({ sessionId }),
  setTrace: (trace) => set({ trace, frameIndex: trace?.frames?.length ? trace.frames.length - 1 : 0 }),
  setFrameIndex: (frameIndex) => set({ frameIndex }),
  setCompileErrors: (compileErrors) => set({ compileErrors }),
  setStatus: (status) => set({ status }),
  setPanicMsg: (panicMsg) => set({ panicMsg }),
  setIsRunning: (isRunning) => set({ isRunning }),
  currentFrame: () => {
    const { trace, frameIndex } = get();
    if (!trace?.frames?.length) return null;
    return trace.frames[frameIndex] ?? null;
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
