import type { Variable } from '../types';

export interface VarChange {
  name: string;
  type: string;
  oldValue: string;
  newValue: string;
}

/** Delve/compiler-internal names (e.g. ~r0 = return slot). */
export function isServiceVar(name: string): boolean {
  return name.startsWith('~') || name === '(return)';
}

function flatVars(vars: Variable[], prefix = ''): Map<string, Variable> {
  const m = new Map<string, Variable>();
  for (const v of vars) {
    const key = prefix ? `${prefix}.${v.name}` : v.name;
    m.set(key, v);
    if (v.children?.length) {
      for (const [k, c] of flatVars(v.children, key)) m.set(k, c);
    }
  }
  return m;
}

function displayValue(v: Variable): string {
  if (v.children?.length) {
    const items = v.children.slice(0, 8).map((c) => c.value.trim() || c.name);
    const suffix = v.children.length > 8 ? ', …' : '';
    return `[${items.join(', ')}${suffix}]`;
  }
  const val = v.value.trim();
  if (val.length > 60) return val.slice(0, 57) + '…';
  return val || '—';
}

export function computeVarChanges(prev: Variable[], curr: Variable[]): VarChange[] {
  const prevMap = flatVars(prev);
  const currMap = flatVars(curr);
  const out: VarChange[] = [];

  for (const [name, cv] of currMap) {
    if (isServiceVar(name)) continue;
    const pv = prevMap.get(name);
    const newVal = displayValue(cv);
    if (!pv) {
      out.push({ name, type: cv.type, oldValue: '—', newValue: newVal });
      continue;
    }
    const oldVal = displayValue(pv);
    if (oldVal !== newVal || pv.type !== cv.type) {
      out.push({ name, type: cv.type, oldValue: oldVal, newValue: newVal });
    }
  }
  return out;
}

export function changedVarNames(changes: VarChange[]): string[] {
  return changes.map((c) => c.name.split('.')[0]).filter(Boolean);
}
