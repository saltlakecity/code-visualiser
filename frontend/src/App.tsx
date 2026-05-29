import { lazy, Suspense, useCallback, useEffect, useRef, useState } from 'react';
import { EventsOn } from '../wailsjs/runtime/runtime';
import {
  Compile,
  GetTrace,
  ListProblems,
  LoadProblem,
  LoadCustomCode,
  RunDebug,
  SaveCustomCode,
  Step,
  Stop,
} from '../wailsjs/go/app/App';
import { useDebuggerStore } from './store/debuggerStore';
import { useSettingsStore } from './store/settingsStore';
import { ProblemSidebar } from './components/ProblemSidebar';
import { CodeEditor } from './components/CodeEditor';
import { DebugToolbar } from './components/DebugToolbar';
import { Timeline } from './components/Timeline';
import { VariablesTree } from './components/VariablesTree';
import { CallStack } from './components/CallStack';
const VizPanel = lazy(() => import('./components/VizPanel').then((m) => ({ default: m.VizPanel })));
import { TestPanel } from './components/TestPanel';
import { ErrorPanel } from './components/ErrorPanel';
import { DisclaimerModal, useDisclaimer } from './components/DisclaimerModal';
import { SettingsModal } from './components/SettingsModal';
import { t, translateStatus } from './i18n';
import { currentLang } from './store/settingsStore';
import type { Trace } from './types';

type RightTab = 'vars' | 'stack';

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
    setProblems,
    setCurrentProblem,
    setCustomMode,
    setSource,
    setSessionId,
    setTrace,
    setFrameIndex,
    setCompileErrors,
    setStatus,
    setPanicMsg,
    setIsRunning,
    currentFrame,
  } = useDebuggerStore();

  const { settings, loadSettings, toast, setToast } = useSettingsStore();
  const tr = t(currentLang(settings));

  const [rightTab, setRightTab] = useState<RightTab>('vars');
  const playRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const sessionIdRef = useRef(sessionId);
  sessionIdRef.current = sessionId;
  const disclaimer = useDisclaimer();

  useEffect(() => {
    loadSettings();
  }, [loadSettings]);

  const selectProblem = useCallback(
    async (id: string) => {
      setCustomMode(false);
      const p = await LoadProblem(id);
      setCurrentProblem(p);
      setSource(p.template);
      setCompileErrors([]);
      setPanicMsg('');
      setTrace(null);
      setStatus('idle');
    },
    [setCustomMode, setCurrentProblem, setSource, setCompileErrors, setPanicMsg, setTrace, setStatus],
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

  const handleCustom = async () => {
    setCustomMode(true);
    setCurrentProblem(null);
    setCompileErrors([]);
    setPanicMsg('');
    setTrace(null);
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

  const handleCompileRun = async () => {
    setCompileErrors([]);
    setPanicMsg('');
    setTrace(null);
    setIsRunning(true);
    setStatus('compiling');

    const mode = customMode ? 'auto' : 'function_wrapper';
    const runner = customMode ? '' : currentProblem?.runner ?? '';
    const res = await Compile(source, mode, runner);
    if (!res.success) {
      setCompileErrors(res.errors ?? []);
      setStatus('compile error');
      setIsRunning(false);
      return;
    }

    if (res.compiledSource?.trim()) {
      setSource(res.compiledSource);
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
      // Fallback: poll trace in case events were missed
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
    if (sessionId) Stop(sessionId);
    if (playRef.current) clearInterval(playRef.current);
    setIsRunning(false);
    setStatus('stopped');
  };

  const handleStep = async (cmd: string) => {
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

  const handlePlay = () => {
    if (!trace?.frames?.length) return;
    if (playRef.current) clearInterval(playRef.current);
    let i = frameIndex;
    playRef.current = setInterval(() => {
      i += 1;
      if (i >= trace.frames.length) {
        if (playRef.current) clearInterval(playRef.current);
        return;
      }
      setFrameIndex(i);
    }, 500);
  };

  const frame = currentFrame();
  const canStep =
    !!sessionId && !isRunning && (status === 'paused' || status === 'debugging');
  const canPlay = (trace?.frames?.length ?? 0) > 0;

  return (
    <div className="flex h-screen flex-col bg-bg text-sm">
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

        <main className="flex min-w-0 flex-1 flex-col">
          <DebugToolbar
            status={translateStatus(currentLang(settings), status)}
            isRunning={isRunning}
            canStep={canStep}
            customMode={customMode}
            onCompileRun={handleCompileRun}
            onStop={handleStop}
            onStep={handleStep}
            canPlay={canPlay}
            onPlay={handlePlay}
            onSaveCustom={handleSaveCustom}
            onLoadCustom={handleLoadCustom}
          />
          <div className="min-h-0 flex-1 p-3">
            <CodeEditor
              source={source}
              onChange={setSource}
              currentLine={frame?.line}
              editorLine={frame?.editorLine ?? frame?.line}
              currentFile={frame?.file}
              errors={compileErrors}
            />
          </div>
          {currentProblem && <TestPanel tests={currentProblem.tests} />}
          <ErrorPanel message={panicMsg} />
          <Timeline
            frameIndex={frameIndex}
            total={trace?.frames?.length ?? 0}
            onChange={setFrameIndex}
          />
        </main>

        <aside className="flex w-96 flex-col border-l border-border bg-panel">
          <div className="flex border-b border-border">
            <TabButton active={rightTab === 'vars'} onClick={() => setRightTab('vars')}>
              {tr.variables}
            </TabButton>
            <TabButton active={rightTab === 'stack'} onClick={() => setRightTab('stack')}>
              {tr.callStack}
            </TabButton>
          </div>
          <div className="h-48 overflow-hidden border-b border-border">
            {rightTab === 'vars' ? (
              <VariablesTree variables={frame?.variables ?? []} />
            ) : (
              <CallStack stack={frame?.stack ?? []} />
            )}
          </div>
          <div className="min-h-0 flex-1">
            <div className="border-b border-border px-3 py-2 text-xs font-semibold uppercase text-muted">
              {tr.visualization}
            </div>
            {(frame?.viz?.length ?? 0) > 0 ? (
              <Suspense fallback={<div className="p-4 text-xs text-muted">…</div>}>
                <VizPanel models={frame?.viz ?? []} />
              </Suspense>
            ) : (
              <div className="flex h-full items-center justify-center p-4 text-xs text-muted">{tr.vizEmpty}</div>
            )}
          </div>
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
