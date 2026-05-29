import { useCallback, useEffect, useMemo, useRef } from 'react';
import type { Language } from '../i18n';
import { highlightGoLine } from '../utils/goHighlight';
import { parseFunctions } from '../utils/goParse';

interface Props {
  source: string;
  highlightLine: number;
  scrollKey?: number;
  lang?: Language;
  changedVars?: string[];
  activeVar?: string | null;
}

const LINE_HEIGHT = 22;

export function LineHighlightEditor({
  source,
  highlightLine,
  scrollKey,
  lang = 'en',
  changedVars = [],
  activeVar = null,
}: Props) {
  const scrollRef = useRef<HTMLDivElement>(null);
  const lineRefs = useRef<Map<number, HTMLDivElement>>(new Map());
  const lastScrollKey = useRef<number | undefined>(undefined);
  const lines = useMemo(() => source.split('\n'), [source]);
  const funcMap = useMemo(() => parseFunctions(source), [source]);

  useEffect(() => {
    if (highlightLine <= 0) return;
    if (scrollKey !== undefined && scrollKey === lastScrollKey.current) return;
    lastScrollKey.current = scrollKey;
    const el = lineRefs.current.get(highlightLine);
    el?.scrollIntoView({ block: 'center', behavior: 'smooth' });
  }, [highlightLine, scrollKey]);

  return (
    <div ref={scrollRef} className="go-editor min-h-0 flex-1 overflow-auto font-mono text-sm leading-[22px]">
      {lines.map((line, i) => {
        const lineNum = i + 1;
        const active = lineNum === highlightLine;
        const html = highlightGoLine(line, {
          lang,
          source,
          changedVars,
          activeVar,
          funcMap,
        });

        return (
          <div
            key={lineNum}
            ref={(el) => {
              if (el) lineRefs.current.set(lineNum, el);
              else lineRefs.current.delete(lineNum);
            }}
            className={`flex min-h-[22px] ${
              active
                ? 'bg-accent/15 shadow-[inset_3px_0_0_0_var(--color-accent)]'
                : 'hover:bg-bg/30'
            }`}
            style={{ minHeight: LINE_HEIGHT }}
          >
            <span
              className={`w-12 shrink-0 select-none border-r border-border/50 pr-3 text-right text-xs leading-[22px] ${
                active ? 'bg-accent/10 font-bold text-accent' : 'text-muted'
              }`}
            >
              {lineNum}
            </span>
            <pre
              className="min-w-0 flex-1 whitespace-pre px-3 py-0 text-[var(--color-text)]"
              dangerouslySetInnerHTML={{
                __html: html || '\u00A0',
              }}
            />
          </div>
        );
      })}
    </div>
  );
}

interface EditProps {
  source: string;
  onChange: (v: string) => void;
  errorLines?: Set<number>;
}

export function GoEditEditor({ source, onChange, errorLines }: EditProps) {
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const gutterRef = useRef<HTMLDivElement>(null);
  const lines = useMemo(() => source.split('\n'), [source]);
  const lineCount = Math.max(lines.length, 1);

  const syncScroll = useCallback(() => {
    const ta = textareaRef.current;
    const gutter = gutterRef.current;
    if (ta && gutter) gutter.scrollTop = ta.scrollTop;
  }, []);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Tab') {
      e.preventDefault();
      const ta = e.currentTarget;
      const start = ta.selectionStart;
      const end = ta.selectionEnd;
      const next = source.slice(0, start) + '\t' + source.slice(end);
      onChange(next);
      requestAnimationFrame(() => {
        ta.selectionStart = ta.selectionEnd = start + 1;
      });
    }
  };

  return (
    <div className="flex min-h-0 flex-1 overflow-hidden">
      <div
        ref={gutterRef}
        className="shrink-0 overflow-hidden border-r border-border/50 bg-bg/40 select-none"
        aria-hidden
      >
        <div className="py-3 pr-3 pl-2 font-mono text-xs leading-[22px] text-muted">
          {Array.from({ length: lineCount }, (_, i) => {
            const n = i + 1;
            const err = errorLines?.has(n);
            return (
              <div
                key={n}
                className={`h-[22px] text-right ${err ? 'font-bold text-error' : ''}`}
              >
                {n}
              </div>
            );
          })}
        </div>
      </div>
      <textarea
        ref={textareaRef}
        className="min-h-0 min-w-0 flex-1 resize-none overflow-auto bg-transparent px-3 py-3 font-mono text-sm leading-[22px] text-[var(--color-text)] outline-none"
        value={source}
        spellCheck={false}
        onChange={(e) => onChange(e.target.value)}
        onScroll={syncScroll}
        onKeyDown={handleKeyDown}
      />
    </div>
  );
}
