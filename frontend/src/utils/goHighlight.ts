import Prism from 'prismjs';
import 'prismjs/components/prism-go';
import type { Language } from '../i18n';
import { formatFuncTip, formatTypeTip, lookupFunc, lookupType } from './goReference';
import { formatUserFuncTip, parseFunctions, type ParsedFunc } from './goParse';

export interface HighlightOptions {
  lang?: Language;
  source?: string;
  changedVars?: string[];
  activeVar?: string | null;
  funcMap?: Map<string, ParsedFunc>;
}

let cachedFuncMap: { source: string; map: Map<string, ParsedFunc> } | null = null;

function getFuncMap(source?: string): Map<string, ParsedFunc> {
  if (!source) return new Map();
  if (cachedFuncMap?.source === source) return cachedFuncMap.map;
  const map = parseFunctions(source);
  cachedFuncMap = { source, map };
  return map;
}

function escapeRegex(s: string): string {
  return s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

function escapeAttr(s: string): string {
  return s.replace(/&/g, '&amp;').replace(/"/g, '&quot;').replace(/</g, '&lt;');
}

function escapeHtml(s: string): string {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;');
}

function wrapTip(className: string, tip: string, inner: string): string {
  return `<span class="${className}" data-tip="${escapeAttr(tip)}">${inner}</span>`;
}

function annotateHtml(html: string, opts: HighlightOptions): string {
  const lang = opts.lang ?? 'en';
  const funcMap = opts.funcMap ?? getFuncMap(opts.source);
  const changed = new Set((opts.changedVars ?? []).map((v) => v.split('.')[0]));
  const active = opts.activeVar?.split('.')[0] ?? null;

  const parts = html.split(/(<[^>]+>)/g);
  return parts
    .map((part) => {
      if (part.startsWith('<')) return part;

      let text = part;

      // Built-in / user functions: longest match first
      const funcNames = [
        ...Array.from(funcMap.keys()).filter((k) => !k.includes('(')),
        ...['fmt.Println', 'fmt.Printf', 'fmt.Print', 'fmt.Sprintf', 'len', 'cap', 'make', 'append', 'copy', 'delete', 'new', 'panic', 'recover'],
      ];
      const seen = new Set<string>();
      const ordered = funcNames.filter((n) => {
        if (seen.has(n)) return false;
        seen.add(n);
        return true;
      }).sort((a, b) => b.length - a.length);

      for (const fn of ordered) {
        const re = new RegExp(`\\b(${escapeRegex(fn)})\\b`, 'g');
        text = text.replace(re, (match) => {
          const userFn = funcMap.get(fn);
          const stdFn = lookupFunc(fn, lang);
          const tip = userFn
            ? formatUserFuncTip(userFn)
            : stdFn
              ? formatFuncTip(stdFn, fn)
              : fn;
          return wrapTip('go-tip go-tip-func', tip, match);
        });
      }

      // Primitive / known types
      const typeWords = ['bool', 'string', 'byte', 'rune', 'int', 'int8', 'int16', 'int32', 'int64', 'uint', 'uint8', 'uint16', 'uint32', 'uint64', 'uintptr', 'float32', 'float64', 'complex64', 'complex128', 'any', 'error'];
      for (const tp of typeWords.sort((a, b) => b.length - a.length)) {
        const re = new RegExp(`\\b(${escapeRegex(tp)})\\b`, 'g');
        text = text.replace(re, (match) => {
          const info = lookupType(tp, lang);
          return wrapTip('go-tip go-tip-type', formatTypeTip(info, tp), match);
        });
      }

      // Changed / active variables
      const varNames = [...changed];
      if (active && !varNames.includes(active)) varNames.push(active);
      varNames.sort((a, b) => b.length - a.length);

      for (const v of varNames) {
        const re = new RegExp(`\\b(${escapeRegex(v)})\\b`, 'g');
        text = text.replace(re, (match) => {
          const cls =
            active === v
              ? 'go-var-active go-tip'
              : changed.has(v)
                ? 'go-var-changed go-tip'
                : 'go-tip';
          const tip =
            active === v
              ? lang === 'ru'
                ? `Активная переменная: ${v}`
                : `Active variable: ${v}`
              : lang === 'ru'
                ? `Изменилась: ${v}`
                : `Changed: ${v}`;
          return `<span class="${cls}" data-tip="${escapeAttr(tip)}">${match}</span>`;
        });
      }

      return text;
    })
    .join('');
}

export function highlightGoLine(line: string, opts: HighlightOptions = {}): string {
  const grammar = Prism.languages.go;
  if (!grammar) return escapeHtml(line);
  const html = Prism.highlight(line, grammar, 'go');
  if (!opts.changedVars?.length && !opts.activeVar && !opts.source) return html;
  return annotateHtml(html, opts);
}

export function highlightGoSource(source: string, opts: Omit<HighlightOptions, 'source'> = {}): string[] {
  return source.split('\n').map((line) => highlightGoLine(line, { ...opts, source }));
}
