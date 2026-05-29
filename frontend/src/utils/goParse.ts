export interface ParsedFunc {
  name: string;
  receiver?: string;
  params: string;
  returns: string;
  line: number;
}

const FUNC_RE =
  /(?:^|\n)\s*(?:func\s+\(([^)]*)\)\s+)?(\w+)\s*\(([^)]*)\)\s*(?:\(([^)]*)\)|(\S+))?\s*\{/g;

export function parseFunctions(source: string): Map<string, ParsedFunc> {
  const map = new Map<string, ParsedFunc>();
  let m: RegExpExecArray | null;
  while ((m = FUNC_RE.exec(source)) !== null) {
    const receiver = m[1]?.trim();
    const name = m[2];
    const params = m[3]?.trim() ?? '';
    const returns = (m[4] ?? m[5] ?? '').trim();
    const line = source.slice(0, m.index).split('\n').length;
    const key = receiver ? `(${receiver}).${name}` : name;
    map.set(key, { name, receiver, params, returns, line });
    map.set(name, { name, receiver, params, returns, line });
  }
  return map;
}

export function formatUserFuncTip(fn: ParsedFunc): string {
  const recv = fn.receiver ? `(${fn.receiver}) ` : '';
  const ret = fn.returns ? ` → ${fn.returns}` : '';
  return `${recv}${fn.name}(${fn.params})${ret}`;
}
