import { lazy, Suspense, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { EventsOn } from '../wailsjs/runtime/runtime';
import {
  BuildRunnerFromTest,
  Compile,
  GetTrace,
  ListProblems,
  LoadCustomTests,
  LoadProblem,
  LoadCustomCode,
  RunDebug,
  SaveCustomCode,
  Step,
  Stop,
} from '../wailsjs/go/app/App';
import { models } from '../wailsjs/go/models';
import { useDebuggerStore } from './store/debuggerStore';
import { useSettingsStore } from './store/settingsStore';
import { ProblemSidebar } from './components/ProblemSidebar';
import { CodeEditor } from './components/CodeEditor';
import { DebugToolbar } from './components/DebugToolbar';
import { Timeline } from './components/Timeline';
import { VariablesTree } from './components/VariablesTree';
import { CallStack } from './components/CallStack';
import { ResizableRightPanel } from './components/ResizableRightPanel';
const VizPanel = lazy(() => import('./components/VizPanel').then((m) => ({ default: m.VizPanel })));
import { TestPanel } from './components/TestPanel';
import { ErrorPanel } from './components/ErrorPanel';
import { DisclaimerModal, useDisclaimer } from './components/DisclaimerModal';
import { SettingsModal } from './components/SettingsModal';
import { t, translateStatus } from './i18n';
import { currentLang } from './store/settingsStore';
import type { TestResult, Trace } from './types';
import { VariableChangesBanner } from './components/VariableChangesBanner';
import { computeVarChanges, changedVarNames } from './utils/varChanges';
import { hasSliceLikeVariables } from './utils/vizHelpers';

type RightTab = 'vars' | 'stack';

const BASE_PLAY_MS = 500;

function App() {
  const {
    problems,
    currentProblem,
    customMode,
    source,
    sessionId,
    trace,
    frameIndex,
    compileErrors,
    status,
    panicMsg,
    isRunning,
    selectedTestIndex,
    setProblems,
    setCurrentProblem,
    setCustomMode,
    setSource,
    setDebugSource,
    setSessionId,
    setTrace,
    setFrameIndex,
    setCompileErrors,
    setStatus,
    setPanicMsg,
    setIsRunning,
    setCustomTests,
    setTestResults,
    setHighlightedVar,
    highlightedVar,
    resetTrace,
    currentFrame,
    prevFrame,
    activeTests,
  } = useDebuggerStore();

  const { settings, loadSettings, saveSettings, toast, setToast } = useSettingsStore();
  const tr = t(currentLang(settings));

  const [rightTab, setRightTab] = useState<RightTab>('vars');
  const [isPlaying, setIsPlaying] = useState(false);
  const [varsHeight, setVarsHeight] = useState(220);
  const playRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const frameIndexRef = useRef(frameIndex);
  const traceRef = useRef(trace);
  const sessionIdRef = useRef(sessionId);
  const disclaimer = useDisclaimer();

  frameIndexRef.current = frameIndex;
  traceRef.current = trace;
  sessionIdRef.current = sessionId;

  useEffect(() => {
    loadSettings();
    LoadCustomTests()
      .then((tests) => {
        if (tests?.length) setCustomTests(tests.map((t) => ({ ...t })));
      })
      .catch(console.error);
  }, [loadSettings, setCustomTests]);

  const stopPlayback = useCallback(() => {
    if (playRef.current) {
      clearInterval(playRef.current);
      playRef.current = null;
    }
    setIsPlaying(false);
  }, []);

  const selectProblem = useCallback(
    async (id: string) => {
      stopPlayback();
      setCustomMode(false);
      const p = await LoadProblem(id);
      setCurrentProblem(p);
      setSource(p.template);
      setCompileErrors([]);
      setPanicMsg('');
      setTrace(null);
      setDebugSource('', '', 0);
      setStatus('idle');
    },
    [setCustomMode, setCurrentProblem, setSource, setCompileErrors, setPanicMsg, setTrace, setDebugSource, setStatus, stopPlayback],
  );

  useEffect(() => {
    ListProblems()
      .then((list) => {
        setProblems(list);
        if (list.length > 0 && !customMode && !currentProblem) {
          selectProblem(list[0].id);
        }
      })
      .catch(console.error);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    const offFrame = EventsOn('debug:frame', () => {
      const sid = sessionIdRef.current;
      if (!sid) return;
      GetTrace(sid).then((trc) => {
        setTrace(trc);
        if (trc.frames?.length) setFrameIndex(trc.frames.length - 1);
      });
    });
    const offFinished = EventsOn('debug:finished', (trc: Trace) => {
      setTrace(trc);
      setStatus(trc.status === 'paused' ? 'finished' : trc.status);
      if (trc.panicMsg) setPanicMsg(trc.panicMsg);
      setIsRunning(false);
      if (trc.frames?.length) setFrameIndex(trc.frames.length - 1);
    });
    const offPanic = EventsOn('debug:panic', (msg: string) => setPanicMsg(msg));
    const offCompile = EventsOn('compile:error', (errs: unknown) => setCompileErrors(errs as never));
    const offDebugErr = EventsOn('debug:error', (msg: string) => {
      setPanicMsg(msg);
      setIsRunning(false);
      setStatus('error');
    });
    return () => {
      offFrame();
      offFinished();
      offPanic();
      offCompile();
      offDebugErr();
    };
  }, [sessionId, setTrace, setFrameIndex, setStatus, setPanicMsg, setCompileErrors, setIsRunning]);

  const startPlayback = useCallback(() => {
    const trc = traceRef.current;
    if (!trc?.frames?.length) return;
    if (playRef.current) clearInterval(playRef.current);
    setIsPlaying(true);
    const speed = settings.playbackSpeed || 1;
    const delay = BASE_PLAY_MS / speed;
    playRef.current = setInterval(() => {
      const current = traceRef.current;
      if (!current?.frames?.length) {
        stopPlayback();
        return;
      }
      const i = frameIndexRef.current + 1;
      if (i >= current.frames.length) {
        stopPlayback();
        return;
      }
      setFrameIndex(i);
    }, delay);
  }, [settings.playbackSpeed, setFrameIndex, stopPlayback]);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      const tag = (e.target as HTMLElement)?.tagName;
      if (tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT') return;
      if (e.code === 'Space') {
        e.preventDefault();
        if (isPlaying) stopPlayback();
        else if (trace?.frames?.length) startPlayback();
      } else if (e.code === 'ArrowLeft' && trace?.frames?.length) {
        stopPlayback();
        setFrameIndex(Math.max(0, frameIndexRef.current - 1));
      } else if (e.code === 'ArrowRight' && trace?.frames?.length) {
        stopPlayback();
        setFrameIndex(Math.min(trace.frames.length - 1, frameIndexRef.current + 1));
      }
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [isPlaying, trace, setFrameIndex, stopPlayback, startPlayback]);

  const resolveRunner = async (): Promise<string> => {
    const tests = activeTests();
    const test = tests[selectedTestIndex] ?? tests[0];
    if (customMode) {
      if (!test) return '';
      const runner = await BuildRunnerFromTest('', new models.TestCase(test), source);
      return runner ?? '';
    }
    if (test && currentProblem?.entry) {
      const runner = await BuildRunnerFromTest(
        currentProblem.entry,
        new models.TestCase(test),
        source,
      );
      if (runner) return runner;
    }
    return currentProblem?.runner ?? '';
  };

  const handleCompileRun = async () => {
    stopPlayback();
    setCompileErrors([]);
    setPanicMsg('');
    setTrace(null);
    setDebugSource('', '', 0);
    setIsRunning(true);
    setStatus('compiling');

    const mode = customMode ? 'auto' : 'function_wrapper';
    const runner = await resolveRunner();
    const res = await Compile(source, mode, runner);
    if (!res.success) {
      setCompileErrors(res.errors ?? []);
      setStatus('compile error');
      setIsRunning(false);
      return;
    }

    if (res.debugSource?.trim()) {
      setDebugSource(res.debugSource, res.debugFile ?? '', res.editorLineOffset ?? 0);
    }

    const sid = res.sessionId;
    if (!sid) {
      setPanicMsg('Compile succeeded but session id is missing');
      setIsRunning(false);
      setStatus('error');
      return;
    }
    setSessionId(sid);
    setStatus('debugging');
    try {
      await RunDebug(sid, { maxSteps: 400 });
      setTimeout(async () => {
        try {
          const trc = await GetTrace(sid);
          if (trc.frames?.length) {
            setTrace(trc);
            setFrameIndex(trc.frames.length - 1);
          }
          if (trc.status === 'finished' || trc.status === 'panic' || trc.status === 'error') {
            setStatus(trc.status);
            setIsRunning(false);
          } else if (trc.frames?.length && trc.status === 'paused') {
            setStatus('finished');
            setIsRunning(false);
          }
        } catch {
          /* ignore */
        }
      }, 2000);
    } catch (e) {
      setPanicMsg(String(e));
      setIsRunning(false);
      setStatus('error');
    }
  };

  const handleStop = () => {
    stopPlayback();
    if (sessionId) Stop(sessionId);
    setIsRunning(false);
    setStatus('stopped');
  };

  const handleStep = async (cmd: string) => {
    stopPlayback();
    if (!sessionId) return;
    const res = await Step(sessionId, cmd);
    if (res.trace) setTrace(res.trace);
    if (res.trace?.frames?.length) setFrameIndex(res.trace.frames.length - 1);
    if (res.finished) {
      setIsRunning(false);
      setStatus('finished');
    }
    if (res.error) setPanicMsg(res.error);
  };

  useEffect(() => {
    if (!isPlaying) return;
    stopPlayback();
    startPlayback();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [settings.playbackSpeed, isPlaying]);

  const handleTimelineChange = (index: number) => {
    stopPlayback();
    setFrameIndex(index);
  };

  const handleSpeedChange = async (speed: number) => {
    const next = { ...settings, playbackSpeed: speed };
    await saveSettings(next);
  };

  const handleCustom = async () => {
    stopPlayback();
    setCustomMode(true);
    setCurrentProblem(null);
    setCompileErrors([]);
    setPanicMsg('');
    setTrace(null);
    setDebugSource('', '', 0);
    setStatus('idle');
    try {
      const saved = await LoadCustomCode('');
      if (saved?.trim()) setSource(saved);
    } catch {
      /* use default */
    }
  };

  const handleSaveCustom = async () => {
    try {
      await SaveCustomCode(source, 'custom.go');
      setToast(tr.custom.saved);
      setTimeout(() => setToast(null), 2500);
    } catch {
      setToast(tr.custom.saveFailed);
      setTimeout(() => setToast(null), 2500);
    }
  };

  const handleLoadCustom = async () => {
    try {
      const saved = await LoadCustomCode('');
      if (saved?.trim()) {
        setSource(saved);
        setToast(tr.custom.loaded);
      }
      setTimeout(() => setToast(null), 2500);
    } catch {
      setToast(tr.custom.saveFailed);
      setTimeout(() => setToast(null), 2500);
    }
  };

  const frame = currentFrame();
  const prev = prevFrame();
  const highlightLine = frame?.editorLine ?? frame?.line ?? 0;
  const lang = currentLang(settings);
  const varChanges = useMemo(
    () => computeVarChanges(prev?.variables ?? [], frame?.variables ?? []),
    [prev?.variables, frame?.variables],
  );
  const changedVars = useMemo(() => changedVarNames(varChanges), [varChanges]);
  const canStep = !!sessionId && !isRunning && (status === 'paused' || status === 'debugging');
  const canPlay = (trace?.frames?.length ?? 0) > 0;

  return (
    <div className="flex h-full flex-col overflow-hidden bg-bg text-sm">
      {disclaimer.open && <DisclaimerModal onAccept={disclaimer.accept} />}
      <SettingsModal />
      {toast && (
        <div className="fixed bottom-4 right-4 z-50 rounded-md bg-accent px-4 py-2 text-xs font-semibold text-black shadow-lg">
          {toast}
        </div>
      )}
      <div className="flex min-h-0 flex-1">
        <ProblemSidebar
          problems={problems}
          selectedId={currentProblem?.id ?? null}
          customMode={customMode}
          onSelect={selectProblem}
          onCustom={handleCustom}
        />

        <main className="flex min-h-0 min-w-0 flex-1 flex-col overflow-hidden">
          <DebugToolbar
            status={translateStatus(currentLang(settings), status)}
            isRunning={isRunning}
            canStep={canStep}
            customMode={customMode}
            onCompileRun={handleCompileRun}
            onStop={handleStop}
            onStep={handleStep}
            canPlay={canPlay}
            onPlay={startPlayback}
            onPause={stopPlayback}
            isPlaying={isPlaying}
            playbackSpeed={settings.playbackSpeed || 1}
            onSpeedChange={handleSpeedChange}
            onSaveCustom={handleSaveCustom}
            onLoadCustom={handleLoadCustom}
          />
          <div className="flex min-h-0 flex-1 flex-col overflow-hidden p-3 pb-0">
            <CodeEditor
              source={source}
              onChange={setSource}
              highlightLine={highlightLine}
              scrollKey={frameIndex}
              readOnly={!!trace}
              errors={compileErrors}
              lang={lang}
              changedVars={changedVars}
              activeVar={highlightedVar}
              onEditAgain={() => {
                stopPlayback();
                resetTrace();
              }}
              editAgainLabel={tr.editAgain}
            />
          </div>
          {trace && <VariableChangesBanner changes={varChanges} lang={lang} />}
          <div className="shrink-0 overflow-hidden">
            <TestPanel onRunSelected={handleCompileRun} onRunAll={(results: TestResult[]) => setTestResults(results)} />
            <ErrorPanel message={panicMsg} />
            <Timeline
              frameIndex={frameIndex}
              total={trace?.frames?.length ?? 0}
              onChange={handleTimelineChange}
              functionName={frame?.function}
              line={frame?.line}
              editorLine={frame?.editorLine}
            />
          </div>
        </main>

        <aside className="flex w-96 shrink-0 flex-col overflow-hidden border-l border-border bg-panel">
          <div className="flex border-b border-border">
            <TabButton active={rightTab === 'vars'} onClick={() => setRightTab('vars')}>
              {tr.variables}
            </TabButton>
            <TabButton active={rightTab === 'stack'} onClick={() => setRightTab('stack')}>
              {tr.callStack}
            </TabButton>
          </div>
          <ResizableRightPanel
            varsHeight={varsHeight}
            onResize={setVarsHeight}
            label={tr.visualization}
            varsPanel={
              rightTab === 'vars' ? (
                <VariablesTree
                  variables={frame?.variables ?? []}
                  prevVariables={prev?.variables ?? []}
                  onSelectVar={setHighlightedVar}
                />
              ) : (
                <CallStack stack={frame?.stack ?? []} />
              )
            }
            vizPanel={
              (frame?.viz?.length ?? 0) > 0 ? (
                <Suspense fallback={<div className="p-4 text-xs text-muted">…</div>}>
                  <VizPanel models={frame?.viz ?? []} />
                </Suspense>
              ) : (
                <div className="flex h-full items-center justify-center p-4 text-center text-xs text-muted">
                  {trace?.frames?.length
                    ? hasSliceLikeVariables(frame?.variables ?? [])
                      ? tr.vizUnrecognized
                      : tr.vizEmpty
                    : tr.vizEmpty}
                </div>
              )
            }
          />
        </aside>
      </div>
    </div>
  );
}

function TabButton({
  active,
  onClick,
  children,
}: {
  active: boolean;
  onClick: () => void;
  children: React.ReactNode;
}) {
  return (
    <button
      onClick={onClick}
      className={`flex-1 px-3 py-2 text-xs font-medium ${
        active ? 'border-b-2 border-accent text-accent' : 'text-muted hover:text-[var(--color-text)]'
      }`}
    >
      {children}
    </button>
  );
}

export default App;
