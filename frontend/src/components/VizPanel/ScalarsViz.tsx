import type { VizModel } from '../../types';

export function ScalarsViz({ model }: { model: VizModel }) {
  const scalars = model.scalars ?? [];
  if (!scalars.length) return null;

  return (
    <div className="p-3">
      <h3 className="mb-2 font-mono text-xs text-accent">{model.name}</h3>
      <div className="grid grid-cols-2 gap-2">
        {scalars.map((s) => (
          <div
            key={s.name}
            className="rounded-md border border-border bg-bg/60 px-2 py-1.5"
          >
            <div className="font-mono text-[10px] font-semibold text-accent">{s.name}</div>
            <div className="truncate font-mono text-xs text-[var(--color-text)]">{s.value}</div>
            <div className="truncate text-[9px] text-muted">{s.type}</div>
          </div>
        ))}
      </div>
    </div>
  );
}
