import { useMemo } from 'react';
import type { Language } from '../i18n';
import type { CompileError } from '../types';
import { MonacoGoEditor } from './MonacoGoEditor';

interface Props {
  source: string;
  onChange: (v: string) => void;
  highlightLine?: number;
  scrollKey?: number;
  errors?: CompileError[];
  readOnly?: boolean;
  onEditAgain?: () => void;
  editAgainLabel?: string;
  lang?: Language;
  changedVars?: string[];
  activeVar?: string | null;
}

export function CodeEditor({
  source,
  onChange,
  highlightLine = 0,
  scrollKey,
  errors = [],
  readOnly = false,
  onEditAgain,
  editAgainLabel,
  lang = 'en',
  changedVars = [],
  activeVar = null,
}: Props) {
  const lineErrs = errors.filter((e) => e.line > 0);
  const errorLines = useMemo(() => new Set(lineErrs.map((e) => e.line)), [lineErrs]);

  return (
    <div className="relative flex h-full min-h-0 flex-col overflow-hidden rounded-lg border border-border bg-panel">
      {readOnly && onEditAgain && (
        <div className="flex shrink-0 items-center justify-between border-b border-border bg-bg/50 px-3 py-1.5">
          <span className="text-[10px] uppercase tracking-wide text-muted">Replay mode</span>
          <button
            onClick={onEditAgain}
            className="rounded border border-border px-2 py-0.5 text-xs hover:bg-bg"
          >
            {editAgainLabel ?? 'Edit again'}
          </button>
        </div>
      )}
      <MonacoGoEditor
        source={source}
        onChange={readOnly ? undefined : onChange}
        readOnly={readOnly}
        highlightLine={highlightLine}
        scrollKey={scrollKey}
        errorLines={errorLines}
        lang={lang}
        changedVars={changedVars}
        activeVar={activeVar}
      />
      {lineErrs.length > 0 && (
        <div className="max-h-24 shrink-0 overflow-y-auto border-t border-border bg-error/10 p-2 text-xs text-error">
          {lineErrs.map((e, i) => (
            <div key={`${e.line}-${i}`}>
              {e.line}: {e.message}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
