import type { CompileError } from '../types';
import { LineHighlightEditor } from './LineHighlightEditor';

interface Props {
  source: string;
  onChange: (v: string) => void;
  highlightLine?: number;
  errors?: CompileError[];
  readOnly?: boolean;
}

export function CodeEditor({
  source,
  onChange,
  highlightLine = 0,
  errors = [],
  readOnly = false,
}: Props) {
  const lineErrs = errors.filter((e) => e.line > 0);

  return (
    <div className="relative flex h-full flex-col overflow-hidden rounded-lg border border-border bg-panel">
      {readOnly ? (
        <LineHighlightEditor source={source} highlightLine={highlightLine} />
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
