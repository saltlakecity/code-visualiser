import { useEffect, useRef } from 'react';
import { highlightGoLine } from '../utils/goHighlight';

interface Props {
  source: string;
  highlightLine: number;
  scrollKey?: number;
}

export function LineHighlightEditor({ source, highlightLine, scrollKey }: Props) {
  const containerRef = useRef<HTMLDivElement>(null);
  const lineRefs = useRef<Map<number, HTMLDivElement>>(new Map());
  const lastScrollKey = useRef<number | undefined>(undefined);
  const lines = source.split('\n');

  useEffect(() => {
    if (highlightLine <= 0) return;
    if (scrollKey !== undefined && scrollKey === lastScrollKey.current) return;
    lastScrollKey.current = scrollKey;
    const el = lineRefs.current.get(highlightLine);
    el?.scrollIntoView({ block: 'center', behavior: 'smooth' });
  }, [highlightLine, scrollKey]);

  return (
    <div ref={containerRef} className="go-editor min-h-0 flex-1 overflow-auto font-mono text-sm leading-relaxed">
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
            <pre
              className="min-w-0 flex-1 whitespace-pre-wrap break-all px-2 py-0 text-[var(--color-text)]"
              dangerouslySetInnerHTML={{
                __html: highlightGoLine(line) || '\u00A0',
              }}
            />
          </div>
        );
      })}
    </div>
  );
}
