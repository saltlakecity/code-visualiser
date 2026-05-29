import type { VizModel } from '../../types';

export function StackViz({ model }: { model: VizModel }) {
  const values = [...(model.values ?? [])].reverse();
  const pointers = model.pointers ?? [];

  return (
    <div className="p-3">
      <h3 className="mb-2 font-mono text-xs text-accent">{model.name} (stack)</h3>
      <div className="flex flex-col items-center gap-1">
        <div className="mb-1 text-[10px] uppercase tracking-wide text-muted">top ↓</div>
        {values.length === 0 ? (
          <div className="rounded border border-dashed border-border px-4 py-3 text-xs text-muted">empty</div>
        ) : (
          values.map((val, revIdx) => {
            const i = values.length - 1 - revIdx;
            const isTop = revIdx === 0;
            return (
              <div
                key={i}
                className={`flex w-full max-w-[120px] items-center justify-center rounded-md border px-3 py-2 font-mono text-sm ${
                  isTop
                    ? 'border-accent bg-accent/20 text-accent'
                    : 'border-border bg-bg text-[var(--color-text)]'
                }`}
              >
                {val}
              </div>
            );
          })
        )}
        {pointers.length > 0 && (
          <div className="mt-2 text-[10px] text-muted">
            {pointers.map((p) => `${p.name} → [${p.index}]`).join(', ')}
          </div>
        )}
      </div>
    </div>
  );
}
