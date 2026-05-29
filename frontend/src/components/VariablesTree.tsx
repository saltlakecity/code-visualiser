import { useState } from 'react';
import { ChevronDown, ChevronRight } from 'lucide-react';
import { t } from '../i18n';
import { useSettingsStore, currentLang } from '../store/settingsStore';
import type { Variable } from '../types';

export function VariablesTree({ variables }: { variables: Variable[] }) {
  const { settings } = useSettingsStore();
  const tr = t(currentLang(settings));

  if (!variables.length) {
    return <p className="p-3 text-xs text-muted">{tr.noVariables}</p>;
  }
  return (
    <div className="overflow-y-auto p-2 font-mono text-xs">
      {variables.map((v) => (
        <VarNode key={v.name} v={v} depth={0} />
      ))}
    </div>
  );
}

function VarNode({ v, depth }: { v: Variable; depth: number }) {
  const [open, setOpen] = useState(depth < 2);
  const hasChildren = (v.children?.length ?? 0) > 0;

  return (
    <div style={{ paddingLeft: depth * 12 }}>
      <div
        className="flex cursor-pointer items-start gap-1 rounded py-0.5 hover:bg-bg"
        onClick={() => hasChildren && setOpen(!open)}
      >
        {hasChildren ? (
          open ? <ChevronDown className="mt-0.5 h-3 w-3 shrink-0" /> : <ChevronRight className="mt-0.5 h-3 w-3 shrink-0" />
        ) : (
          <span className="w-3" />
        )}
        <span className="text-accent">{v.name}</span>
        <span className="text-muted">{v.type}</span>
        <span className="ml-1 truncate text-success">{v.value}</span>
      </div>
      {open &&
        v.children?.map((c) => <VarNode key={`${v.name}.${c.name}`} v={c} depth={depth + 1} />)}
    </div>
  );
}
