import type { VizModel } from '../../types';

const POINTER_STYLES = [
  'bg-accent/30 text-accent border-accent',
  'bg-sky-500/25 text-sky-400 border-sky-400',
  'bg-amber-500/25 text-amber-400 border-amber-400',
  'bg-emerald-500/25 text-emerald-400 border-emerald-400',
  'bg-violet-500/25 text-violet-400 border-violet-400',
];

export function ArrayViz({ model }: { model: VizModel }) {
  const values = model.values ?? [];
  const highlight = new Set(model.highlight ?? []);
  const pointers = model.pointers ?? [];

  const pointersByIndex = new Map<number, { name: string; styleIdx: number }[]>();
  pointers.forEach((p, idx) => {
    const list = pointersByIndex.get(p.index) ?? [];
    list.push({ name: p.name, styleIdx: idx % POINTER_STYLES.length });
    pointersByIndex.set(p.index, list);
  });

  return (
    <div className="p-3">
      <h3 className="mb-2 font-mono text-xs text-accent">{model.name}</h3>
      <div className="flex flex-wrap items-end gap-2">
        {values.map((val, i) => {
          const cellPointers = pointersByIndex.get(i) ?? [];
          const active = highlight.has(i);
          return (
            <div key={i} className="flex flex-col items-center gap-1">
              <div className="flex min-h-[18px] flex-wrap justify-center gap-0.5">
                {cellPointers.map((p) => (
                  <span
                    key={p.name}
                    className={`rounded border px-1 text-[9px] font-semibold ${POINTER_STYLES[p.styleIdx]}`}
                  >
                    ↓ {p.name}
                  </span>
                ))}
              </div>
              <div
                className={`flex h-12 min-w-[40px] items-center justify-center rounded-md border px-2 font-mono text-sm ${
                  active
                    ? 'border-accent bg-accent/25 text-accent shadow-[0_0_0_1px_var(--color-accent)]'
                    : 'border-border bg-bg'
                }`}
              >
                {val}
              </div>
              <span className={`text-[10px] ${active ? 'font-semibold text-accent' : 'text-muted'}`}>
                {i}
              </span>
            </div>
          );
        })}
      </div>
    </div>
  );
}
