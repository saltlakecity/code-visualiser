import type { CompileError } from '../types';
import { LineHighlightEditor } from './LineHighlightEditor';

interface Props {
  source: string;
  onChange: (v: string) => void;
  highlightLine?: number;
  scrollKey?: number;
  errors?: CompileError[];
  readOnly?: boolean;
  onEditAgain?: () => void;
  editAgainLabel?: string;
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
}: Props) {
  const lineErrs = errors.filter((e) => e.line > 0);

  return (
    <div className="relative flex h-full flex-col overflow-hidden rounded-lg border border-border bg-panel">
      {readOnly ? (
        <>
          {onEditAgain && (
            <div className="flex items-center justify-between border-b border-border bg-bg/50 px-3 py-1.5">
              <span className="text-[10px] uppercase tracking-wide text-muted">Replay mode</span>
              <button
                onClick={onEditAgain}
                className="rounded border border-border px-2 py-0.5 text-xs hover:bg-bg"
              >
                {editAgainLabel ?? 'Edit again'}
              </button>
            </div>
          )}
          <LineHighlightEditor
            source={source}
            highlightLine={highlightLine}
            scrollKey={scrollKey}
          />
        </>
      ) : (
        <textarea
          className="min-h-0 flex-1 resize-none bg-transparent p-3 font-mono text-sm leading-relaxed text-[var(--color-text)] outline-none"
          value={source}
          spellCheck={false}
          onChange={(e) => onChange(e.target.value)}
        />
      )}
      {lineErrs.length > 0 && (
        <div className="max-h-24 overflow-y-auto border-t border-border bg-error/10 p-2 text-xs text-error">
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
