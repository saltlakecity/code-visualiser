import type { VizModel } from '../../types';

const POINTER_STYLES = [
  'border-accent text-accent',
  'border-sky-400 text-sky-400',
  'border-amber-400 text-amber-400',
];

export function MatrixViz({ model }: { model: VizModel }) {
  const rows = model.rows ?? [];
  const highlight = new Set(model.highlight ?? []);
  const pointers = model.pointers ?? [];

  return (
    <div className="p-3">
      <h3 className="mb-2 font-mono text-xs text-accent">{model.name}</h3>
      {pointers.length > 0 && (
        <div className="mb-2 flex flex-wrap gap-1">
          {pointers.map((p, idx) => (
            <span
              key={p.name}
              className={`rounded border px-1.5 py-0.5 text-[9px] font-semibold ${POINTER_STYLES[idx % POINTER_STYLES.length]}`}
            >
              {p.name}
            </span>
          ))}
        </div>
      )}
      <div className="inline-block overflow-x-auto">
        <table className="border-collapse font-mono text-sm">
          <tbody>
            {rows.map((row, ri) => (
              <tr key={ri}>
                {row.map((val, ci) => {
                  const flatIdx = ri * row.length + ci;
                  const active = highlight.has(flatIdx) || highlight.has(ci);
                  return (
                    <td
                      key={ci}
                      className={`h-10 min-w-[40px] border border-border px-2 text-center ${
                        active ? 'bg-accent/25 text-accent' : 'bg-bg'
                      }`}
                    >
                      {val}
                    </td>
                  );
                })}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export function MapViz({ model }: { model: VizModel }) {
  const entries = model.entries ?? [];
  return (
    <div className="p-3">
      <h3 className="mb-2 font-mono text-xs text-accent">{model.name}</h3>
      <table className="w-full border-collapse text-xs">
        <thead>
          <tr className="text-muted">
            <th className="border border-border px-2 py-1 text-left">Key</th>
            <th className="border border-border px-2 py-1 text-left">Value</th>
          </tr>
        </thead>
        <tbody>
          {entries.map((e, i) => (
            <tr key={i}>
              <td className="border border-border px-2 py-1 font-mono text-accent">{e.key}</td>
              <td className="border border-border px-2 py-1 font-mono">{e.value}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
