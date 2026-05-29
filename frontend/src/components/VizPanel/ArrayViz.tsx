import type { VizModel } from '../../types';

export function ArrayViz({ model }: { model: VizModel }) {
  const values = model.values ?? [];
  const highlight = new Set(model.highlight ?? []);

  return (
    <div className="p-3">
      <h3 className="mb-2 font-mono text-xs text-accent">{model.name}</h3>
      <div className="flex flex-wrap items-end gap-2">
        {values.map((val, i) => (
          <div key={i} className="flex flex-col items-center gap-1">
            <div
              className={`flex h-12 min-w-[40px] items-center justify-center rounded-md border px-2 font-mono text-sm ${
                highlight.has(i)
                  ? 'border-accent bg-accent/25 text-accent'
                  : 'border-border bg-bg'
              }`}
            >
              {val}
            </div>
            <span className="text-[10px] text-muted">{i}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
