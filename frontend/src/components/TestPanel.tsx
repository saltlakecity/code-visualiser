import { useCallback, useEffect, useState } from 'react';
import { Check, Copy, Plus, Trash2, Play, ListChecks } from 'lucide-react';
import {
  BuildRunnerFromTest,
  LoadCustomTests,
  RunAllTests,
  SaveCustomTests,
} from '../../wailsjs/go/app/App';
import { models } from '../../wailsjs/go/models';
import { t } from '../i18n';
import { useDebuggerStore } from '../store/debuggerStore';
import { useSettingsStore, currentLang } from '../store/settingsStore';
import type { TestCase, TestResult } from '../types';

interface Props {
  onRunSelected: () => void;
  onRunAll: (results: TestResult[]) => void;
}

export function TestPanel({ onRunSelected, onRunAll }: Props) {
  const { settings } = useSettingsStore();
  const tr = t(currentLang(settings));
  const {
    customMode,
    currentProblem,
    source,
    customTests,
    selectedTestIndex,
    testResults,
    setCustomTests,
    setSelectedTestIndex,
    setTestResults,
    activeTests,
  } = useDebuggerStore();

  const tests = activeTests();
  const [editing, setEditing] = useState<TestCase | null>(null);

  useEffect(() => {
    if (customMode) {
      LoadCustomTests()
        .then((loaded) => {
          if (loaded?.length) {
            setCustomTests(loaded.map((t) => ({ ...t })));
          }
        })
        .catch(console.error);
    }
  }, [customMode, setCustomTests]);

  const persistCustom = useCallback(
    async (next: TestCase[]) => {
      if (!customMode) return;
      try {
        await SaveCustomTests(next.map((t) => new models.TestCase(t)));
      } catch (e) {
        console.error(e);
      }
    },
    [customMode],
  );

  const updateTests = (next: TestCase[]) => {
    if (customMode) setCustomTests(next);
    persistCustom(next);
  };

  const addTest = () => {
    const next = [
      ...tests,
      {
        name: `Case ${tests.length + 1}`,
        input: '',
        expected: '',
      },
    ];
    updateTests(next);
    setSelectedTestIndex(next.length - 1);
    setEditing(next[next.length - 1]);
  };

  const deleteTest = (idx: number) => {
    const next = tests.filter((_, i) => i !== idx);
    updateTests(next);
    setSelectedTestIndex(Math.max(0, Math.min(selectedTestIndex, next.length - 1)));
    setEditing(null);
  };

  const duplicateTest = (idx: number) => {
    const copy = { ...tests[idx], name: `${tests[idx].name ?? 'Case'} copy` };
    const next = [...tests.slice(0, idx + 1), copy, ...tests.slice(idx + 1)];
    updateTests(next);
    setSelectedTestIndex(idx + 1);
  };

  const saveEdit = () => {
    if (!editing) return;
    const next = tests.map((t, i) => (i === selectedTestIndex ? { ...editing } : t));
    updateTests(next);
    setEditing(null);
  };

  const handleRunAll = async () => {
    const entry = currentProblem?.entry ?? '';
    try {
      const results = await RunAllTests(source, customMode ? 'auto' : 'function_wrapper', entry, tests);
      setTestResults(results.map((r) => ({ ...r })));
      onRunAll(results.map((r) => ({ ...r })));
    } catch (e) {
      console.error(e);
    }
  };

  if (!customMode && !currentProblem) return null;

  const selected = tests[selectedTestIndex];
  const edit = editing ?? selected;

  return (
    <div className="max-h-48 shrink-0 overflow-y-auto border-t border-border bg-panel p-3">
      <div className="mb-2 flex flex-wrap items-center gap-2">
        <h3 className="text-xs font-semibold uppercase tracking-wide text-muted">{tr.testCases}</h3>
        <div className="ml-auto flex gap-1">
          {customMode && (
            <button
              onClick={addTest}
              className="flex items-center gap-1 rounded border border-border px-2 py-1 text-xs hover:bg-bg"
            >
              <Plus className="h-3 w-3" />
              {tr.tests.add}
            </button>
          )}
          <button
            onClick={onRunSelected}
            className="flex items-center gap-1 rounded border border-border px-2 py-1 text-xs hover:bg-bg"
          >
            <Play className="h-3 w-3" />
            {tr.tests.runSelected}
          </button>
          <button
            onClick={handleRunAll}
            disabled={tests.length === 0}
            className="flex items-center gap-1 rounded border border-border px-2 py-1 text-xs hover:bg-bg disabled:opacity-40"
          >
            <ListChecks className="h-3 w-3" />
            {tr.tests.runAll}
          </button>
        </div>
      </div>

      {tests.length === 0 ? (
        <p className="text-xs text-muted">{tr.tests.empty}</p>
      ) : (
        <div className="flex gap-3">
          <div className="flex max-h-40 min-w-[140px] flex-col gap-1 overflow-y-auto">
            {tests.map((tc, i) => {
              const result = testResults.find((r) => r.index === i);
              const active = i === selectedTestIndex;
              return (
                <button
                  key={i}
                  onClick={() => {
                    setSelectedTestIndex(i);
                    setEditing(null);
                  }}
                  className={`flex items-center gap-1 rounded border px-2 py-1 text-left text-xs ${
                    active ? 'border-accent bg-accent/10' : 'border-border hover:bg-bg'
                  }`}
                >
                  {result && (
                    <span className={result.passed ? 'text-success' : 'text-error'}>
                      {result.passed ? '✓' : '✗'}
                    </span>
                  )}
                  <span className="truncate">{tc.name ?? `Case ${i + 1}`}</span>
                </button>
              );
            })}
          </div>

          {edit && (
            <div className="min-w-0 flex-1 space-y-2">
              {customMode ? (
                <>
                  <input
                    className="w-full rounded border border-border bg-bg px-2 py-1 font-mono text-xs"
                    placeholder={tr.tests.name}
                    value={edit.name ?? ''}
                    onChange={(e) => setEditing({ ...edit, name: e.target.value })}
                  />
                  <input
                    className="w-full rounded border border-border bg-bg px-2 py-1 font-mono text-xs"
                    placeholder={tr.tests.inputHint}
                    value={edit.input}
                    onChange={(e) => setEditing({ ...edit, input: e.target.value })}
                  />
                  <input
                    className="w-full rounded border border-border bg-bg px-2 py-1 font-mono text-xs"
                    placeholder={tr.tests.argsHint}
                    value={edit.args ?? ''}
                    onChange={(e) => setEditing({ ...edit, args: e.target.value })}
                  />
                  <input
                    className="w-full rounded border border-border bg-bg px-2 py-1 font-mono text-xs"
                    placeholder={tr.expected}
                    value={edit.expected}
                    onChange={(e) => setEditing({ ...edit, expected: e.target.value })}
                  />
                  <div className="flex gap-1">
                    <button
                      onClick={saveEdit}
                      className="flex items-center gap-1 rounded border border-border px-2 py-1 text-xs hover:bg-bg"
                    >
                      <Check className="h-3 w-3" />
                      {tr.tests.save}
                    </button>
                    <button
                      onClick={() => duplicateTest(selectedTestIndex)}
                      className="flex items-center gap-1 rounded border border-border px-2 py-1 text-xs hover:bg-bg"
                    >
                      <Copy className="h-3 w-3" />
                      {tr.tests.duplicate}
                    </button>
                    <button
                      onClick={() => deleteTest(selectedTestIndex)}
                      className="flex items-center gap-1 rounded border border-border px-2 py-1 text-xs text-error hover:bg-bg"
                    >
                      <Trash2 className="h-3 w-3" />
                      {tr.tests.delete}
                    </button>
                  </div>
                </>
              ) : (
                <div className="rounded-md border border-border bg-bg p-2 font-mono text-xs">
                  <div>
                    <span className="text-muted">{tr.input}: </span>
                    {edit.input}
                  </div>
                  {edit.args && (
                    <div>
                      <span className="text-muted">{tr.tests.args}: </span>
                      {edit.args}
                    </div>
                  )}
                  <div>
                    <span className="text-muted">{tr.expected}: </span>
                    <span className="text-success">{edit.expected}</span>
                  </div>
                </div>
              )}

              {testResults
                .filter((r) => r.index === selectedTestIndex)
                .map((r, i) => (
                  <div
                    key={i}
                    className={`rounded border p-2 text-xs ${
                      r.passed ? 'border-success/40 bg-success/10' : 'border-error/40 bg-error/10'
                    }`}
                  >
                    {r.passed ? tr.tests.passed : tr.tests.failed}
                    {r.actual !== undefined && (
                      <div className="mt-1 font-mono">
                        {tr.tests.actual}: {r.actual || '(empty)'}
                      </div>
                    )}
                    {r.error && <div className="mt-1 text-error">{r.error}</div>}
                  </div>
                ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export async function buildRunnerForTest(
  source: string,
  entry: string,
  test: TestCase,
): Promise<string> {
  try {
    const runner = await BuildRunnerFromTest(entry, new models.TestCase(test), source);
    return runner ?? '';
  } catch {
    return '';
  }
}
