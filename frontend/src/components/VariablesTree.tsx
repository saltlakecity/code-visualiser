import { useMemo, useState } from 'react';
import { ChevronDown, ChevronRight } from 'lucide-react';
import { t } from '../i18n';
import { useDebuggerStore } from '../store/debuggerStore';
import { useSettingsStore, currentLang } from '../store/settingsStore';
import type { Variable } from '../types';

interface Props {
  variables: Variable[];
  prevVariables?: Variable[];
  onSelectVar?: (name: string) => void;
}

export function VariablesTree({ variables, prevVariables = [], onSelectVar }: Props) {
  const { settings } = useSettingsStore();
  const tr = t(currentLang(settings));
  const { highlightedVar } = useDebuggerStore();

  const prevMap = useMemo(() => buildVarMap(prevVariables), [prevVariables]);
  const changed = useMemo(
    () => new Set(findChanged(variables, prevMap)),
    [variables, prevMap],
  );

  if (!variables.length) {
    return <p className="p-3 text-xs text-muted">{tr.noVariables}</p>;
  }

  return (
    <div className="overflow-y-auto p-2 text-xs">
      {variables.map((v) => (
        <VarNode
          key={v.name}
          v={v}
          depth={0}
          changed={changed.has(v.name)}
          highlighted={highlightedVar === v.name}
          onSelect={onSelectVar}
        />
      ))}
    </div>
  );
}

function VarNode({
  v,
  depth,
  changed,
  highlighted,
  onSelect,
}: {
  v: Variable;
  depth: number;
  changed: boolean;
  highlighted: boolean;
  onSelect?: (name: string) => void;
}) {
  const [open, setOpen] = useState(depth < 2);
  const hasChildren = (v.children?.length ?? 0) > 0;
  const preview = formatPreview(v);

  return (
    <div style={{ paddingLeft: depth * 10 }}>
      <div
        className={`mb-1 rounded px-1 py-1 ${
          highlighted ? 'ring-1 ring-accent' : changed ? 'bg-amber-500/15' : ''
        } ${onSelect ? 'cursor-pointer hover:bg-bg' : ''}`}
        onClick={() => {
          if (hasChildren) setOpen(!open);
          onSelect?.(v.name);
        }}
      >
        <div className="flex flex-wrap items-center gap-1">
          {hasChildren ? (
            open ? (
              <ChevronDown className="h-3 w-3 shrink-0" />
            ) : (
              <ChevronRight className="h-3 w-3 shrink-0" />
            )
          ) : (
            <span className="w-3" />
          )}
          <span className="font-semibold text-accent">{v.name}</span>
          <span className="rounded bg-bg px-1 py-0.5 font-mono text-[10px] text-muted">{v.type}</span>
          {changed && <span className="text-[10px] text-amber-400">Δ</span>}
        </div>
        <div className="ml-4 truncate font-mono text-success">{preview}</div>
      </div>
      {open &&
        v.children?.map((c) => (
          <VarNode
            key={`${v.name}.${c.name}`}
            v={c}
            depth={depth + 1}
            changed={false}
            highlighted={false}
            onSelect={onSelect}
          />
        ))}
    </div>
  );
}

function formatPreview(v: Variable): string {
  if (v.children?.length) {
    const slicePreview = formatSliceChildren(v);
    if (slicePreview) return slicePreview;
    const mapPreview = formatMapChildren(v);
    if (mapPreview) return mapPreview;
  }
  const val = v.value.trim();
  if (val.length > 80) return val.slice(0, 77) + '…';
  return val || '—';
}

function formatSliceChildren(v: Variable): string | null {
  const typeLower = v.type.toLowerCase();
  if (!typeLower.startsWith('[]') && !typeLower.includes('slice')) return null;
  const children = v.children ?? [];
  const items = children.map((c) => formatChildVal(c)).filter(Boolean);
  if (!items.length) return null;
  return `[${items.join(', ')}]`;
}

function formatMapChildren(v: Variable): string | null {
  if (!v.type.startsWith('map[')) return null;
  const children = v.children ?? [];
  const pairs = children.slice(0, 6).map((c) => `${c.name}:${c.value}`);
  const suffix = children.length > 6 ? ', …' : '';
  return `{${pairs.join(', ')}${suffix}}`;
}

function formatChildVal(c: Variable): string {
  let val = c.value.trim();
  if (val.includes(':')) {
    val = val.split(':').pop()?.trim() ?? val;
  }
  return val;
}

function buildVarMap(vars: Variable[]): Map<string, string> {
  const m = new Map<string, string>();
  for (const v of vars) {
    m.set(v.name, serializeVar(v));
  }
  return m;
}

function serializeVar(v: Variable): string {
  return `${v.type}|${v.value}|${(v.children ?? []).map(serializeVar).join(';')}`;
}

function findChanged(vars: Variable[], prev: Map<string, string>): string[] {
  const out: string[] = [];
  for (const v of vars) {
    const cur = serializeVar(v);
    const prevVal = prev.get(v.name);
    if (prevVal === undefined || prevVal !== cur) {
      out.push(v.name);
    }
  }
  return out;
}
