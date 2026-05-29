import { useEffect, useRef } from 'react';

interface Props {
  source: string;
  highlightLine: number;
}

export function LineHighlightEditor({ source, highlightLine }: Props) {
  const containerRef = useRef<HTMLDivElement>(null);
  const lineRefs = useRef<Map<number, HTMLDivElement>>(new Map());
  const lines = source.split('\n');

  useEffect(() => {
    if (highlightLine <= 0) return;
    const el = lineRefs.current.get(highlightLine);
    el?.scrollIntoView({ block: 'center', behavior: 'smooth' });
  }, [highlightLine, source]);

  return (
    <div ref={containerRef} className="min-h-0 flex-1 overflow-auto font-mono text-sm leading-relaxed">
      {lines.map((line, i) => {
        const lineNum = i + 1;
        const active = lineNum === highlightLine;
        return (
          <div
            key={lineNum}
            ref={(el) => {
              if (el) lineRefs.current.set(lineNum, el);
              else lineRefs.current.delete(lineNum);
            }}
            className={`flex ${
              active
                ? 'border-l-2 border-accent bg-accent/20'
                : 'border-l-2 border-transparent'
            }`}
          >
            <span
              className={`w-10 shrink-0 select-none pr-2 text-right text-xs ${
                active ? 'font-semibold text-accent' : 'text-muted'
              }`}
            >
              {lineNum}
            </span>
            <pre className="min-w-0 flex-1 whitespace-pre-wrap break-all px-2 py-0 text-[var(--color-text)]">
              {line || '\u00A0'}
            </pre>
          </div>
        );
      })}
    </div>
  );
}
