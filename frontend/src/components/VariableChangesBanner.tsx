import { useMemo } from 'react';
import type { Language } from '../i18n';
import { t } from '../i18n';
import type { VarChange } from '../utils/varChanges';

interface Props {
  changes: VarChange[];
  lang: Language;
}

export function VariableChangesBanner({ changes, lang }: Props) {
  const tr = t(lang);

  const visible = useMemo(() => changes.slice(0, 6), [changes]);

  if (!visible.length) return null;

  return (
    <div className="shrink-0 border-t border-border bg-panel/80 px-3 py-1.5">
      <div className="mb-1 text-[10px] font-semibold uppercase tracking-wide text-muted">
        {tr.varChanges}
      </div>
      <div className="flex flex-wrap gap-1.5">
        {visible.map((c) => (
          <div
            key={c.name}
            className="flex max-w-full items-center gap-1 rounded border border-amber-500/30 bg-amber-500/10 px-2 py-0.5 font-mono text-[11px]"
            title={`${c.name}: ${c.type}`}
          >
            <span className="font-semibold text-accent">{c.name}</span>
            <span className="truncate text-muted">{c.oldValue}</span>
            <span className="text-amber-400">→</span>
            <span className="truncate text-success">{c.newValue}</span>
          </div>
        ))}
        {changes.length > visible.length && (
          <span className="self-center text-[10px] text-muted">+{changes.length - visible.length}</span>
        )}
      </div>
    </div>
  );
}
