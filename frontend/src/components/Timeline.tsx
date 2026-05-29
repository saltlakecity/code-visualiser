import { ChevronLeft, ChevronRight } from 'lucide-react';

interface Props {
  frameIndex: number;
  total: number;
  onChange: (index: number) => void;
  functionName?: string;
  line?: number;
  editorLine?: number;
}

export function Timeline({
  frameIndex,
  total,
  onChange,
  functionName,
  line,
  editorLine,
}: Props) {
  const max = Math.max(0, total - 1);

  return (
    <div className="border-t border-border bg-panel px-4 py-2">
      <div className="flex items-center gap-3">
        <button
          disabled={frameIndex <= 0}
          onClick={() => onChange(Math.max(0, frameIndex - 1))}
          className="rounded border border-border p-1 disabled:opacity-40"
        >
          <ChevronLeft className="h-4 w-4" />
        </button>
        <input
          type="range"
          min={0}
          max={max}
          value={frameIndex}
          disabled={total === 0}
          onChange={(e) => onChange(Number(e.target.value))}
          className="flex-1 accent-accent"
        />
        <button
          disabled={frameIndex >= max}
          onClick={() => onChange(Math.min(max, frameIndex + 1))}
          className="rounded border border-border p-1 disabled:opacity-40"
        >
          <ChevronRight className="h-4 w-4" />
        </button>
        <span className="min-w-[80px] text-right font-mono text-xs text-muted">
          {total === 0 ? '0 / 0' : `${frameIndex + 1} / ${total}`}
        </span>
      </div>
      {functionName && (
        <div className="mt-1 truncate font-mono text-[10px] text-muted">
          {functionName}
          {(editorLine ?? line) ? ` : ${editorLine ?? line}` : ''}
        </div>
      )}
    </div>
  );
}
