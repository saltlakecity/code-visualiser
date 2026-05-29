import { useEffect, useRef, useState } from 'react';
import { AnalyzeSource } from '../../wailsjs/go/app/App';
import type { Language } from '../i18n';
import { t } from '../i18n';
import type { AnalyzeResult, Trace, TraceFrame } from '../types';

interface Props {
  trace: Trace | null;
  frameIndex: number;
  source: string;
  lang: Language;
  onSelectStep: (index: number) => void;
}

function stepLabel(frame: TraceFrame, index: number, tr: ReturnType<typeof t>): string {
  const exp = frame.explanation;
  if (exp?.summary) {
    const s = exp.summary;
    return s.length > 80 ? s.slice(0, 77) + '…' : s;
  }
  const line = frame.editorLine ?? frame.line;
  return `${tr.breakdownStep} ${index + 1} · ${frame.function || '?'} : ${line}`;
}

export function BreakdownPanel({ trace, frameIndex, source, lang, onSelectStep }: Props) {
  const tr = t(lang);
  const [outline, setOutline] = useState<AnalyzeResult | null>(null);
  const activeRef = useRef<HTMLButtonElement>(null);

  const frames = trace?.frames ?? [];
  const frame = frames[frameIndex];
  const exp = frame?.explanation;

  useEffect(() => {
    if (!source.trim()) {
      setOutline(null);
      return;
    }
    let cancelled = false;
    AnalyzeSource(source, lang)
      .then((res) => {
        if (!cancelled) setOutline(res);
      })
      .catch(() => {
        if (!cancelled) setOutline(null);
      });
    return () => {
      cancelled = true;
    };
  }, [source, lang]);

  useEffect(() => {
    activeRef.current?.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
  }, [frameIndex]);

  if (!frames.length) {
    return (
      <div className="flex flex-col p-3 text-xs">
        <p className="text-muted">{tr.breakdownEmpty}</p>
        {outline && (
          <div className="mt-4 rounded border border-border bg-bg/50 p-3">
            <div className="mb-2 font-semibold text-accent">{tr.breakdownOutline}</div>
            <p className="text-[var(--color-text)]">{outline.summary}</p>
            {outline.outline?.loops && outline.outline.loops.length > 0 && (
              <ul className="mt-2 list-inside list-disc text-muted">
                {outline.outline.loops.map((l) => (
                  <li key={`loop-${l.line}`}>
                    {l.line}: {l.text}
                  </li>
                ))}
              </ul>
            )}
          </div>
        )}
        <p className="mt-4 text-[10px] text-muted">{tr.breakdownLimits}</p>
      </div>
    );
  }

  return (
    <div className="flex flex-col text-xs">
      <div className="shrink-0 border-b border-border p-3">
        <div className="mb-1 flex flex-wrap items-center gap-2 text-[10px] uppercase tracking-wide text-muted">
          <span>
            {tr.breakdownStep} {frameIndex + 1} / {frames.length}
          </span>
          {frame?.function && <span className="font-mono normal-case">{frame.function}</span>}
          {(frame?.editorLine ?? frame?.line) ? (
            <span>
              {tr.breakdownLine} {frame?.editorLine ?? frame?.line}
            </span>
          ) : null}
          {exp?.loopIter ? (
            <span className="text-accent">
              {tr.breakdownLoop} #{exp.loopIter}
            </span>
          ) : null}
        </div>

        {exp?.lineText && (
          <pre className="mb-2 overflow-x-auto rounded bg-bg/60 px-2 py-1 font-mono text-[11px] text-muted">
            {exp.lineText}
          </pre>
        )}

        <p className="text-sm leading-relaxed text-[var(--color-text)]">{exp?.summary ?? '—'}</p>

        {exp?.condition && (
          <div className="mt-2 rounded border border-amber-500/30 bg-amber-500/10 px-2 py-1.5 font-mono text-[11px]">
            <span className="font-semibold text-amber-400">{tr.breakdownCondition}: </span>
            {exp.condition}
          </div>
        )}

        {exp?.changes && exp.changes.length > 0 && (
          <div className="mt-3">
            <div className="mb-1 font-semibold text-muted">{tr.breakdownChanges}</div>
            <table className="w-full border-collapse font-mono text-[11px]">
              <thead>
                <tr className="text-left text-[10px] text-muted">
                  <th className="pb-1 pr-2">{tr.breakdownColName}</th>
                  <th className="pb-1 pr-2">{tr.breakdownColWas}</th>
                  <th className="pb-1">{tr.breakdownColNow}</th>
                </tr>
              </thead>
              <tbody>
                {exp.changes.map((c) => (
                  <tr key={c.name} className="border-t border-border/50">
                    <td className="py-0.5 pr-2 font-semibold text-accent">{c.name}</td>
                    <td className="max-w-[80px] truncate py-0.5 pr-2 text-muted">{c.oldValue}</td>
                    <td className="max-w-[80px] truncate py-0.5 text-success">{c.newValue}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      <div className="flex flex-col">
        <div className="shrink-0 border-b border-border px-3 py-1.5 text-[10px] font-semibold uppercase tracking-wide text-muted">
          {tr.breakdownHistory}
        </div>
        <div className="p-2">
          {frames.map((f, i) => {
            const active = i === frameIndex;
            return (
              <button
                key={f.index ?? i}
                ref={active ? activeRef : undefined}
                type="button"
                onClick={() => onSelectStep(i)}
                className={`mb-1 w-full rounded px-2 py-1.5 text-left transition-colors ${
                  active
                    ? 'border border-accent/40 bg-accent/15 text-[var(--color-text)]'
                    : 'border border-transparent hover:bg-bg/50 text-muted'
                }`}
              >
                <span className="font-mono text-[10px] text-accent">{i + 1}.</span>{' '}
                <span className={active ? 'text-[var(--color-text)]' : ''}>{stepLabel(f, i, tr)}</span>
              </button>
            );
          })}
        </div>
        {frames.length >= 400 && (
          <p className="shrink-0 border-t border-border px-3 py-1.5 text-[10px] text-amber-500/90">
            {tr.breakdownLimits}
          </p>
        )}
      </div>
    </div>
  );
}

// exported for App changedVars wiring
export function explanationChangedVars(frame: TraceFrame | undefined): string[] {
  if (!frame?.explanation?.changes?.length) return [];
  return frame.explanation.changes.map((c) => c.name.split('.')[0]).filter(Boolean);
}
