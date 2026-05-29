import type { Monaco } from '@monaco-editor/react';
import type { IDisposable, languages } from 'monaco-editor';
import type { Language } from './i18n';
import { BUILTIN_FUNCS, GO_TYPES } from './utils/goBuiltins';
import { formatFuncTip, formatTypeTip, lookupFunc, lookupType } from './utils/goReference';
import { formatUserFuncTip, parseFunctions } from './utils/goParse';

let registered = false;
let disposables: IDisposable[] = [];
let monacoApi: Monaco | null = null;

let currentLang: Language = 'en';
let currentSource = '';

export function setMonacoGoContext(lang: Language, source: string): void {
  currentLang = lang;
  currentSource = source;
}

export function registerMonacoGoLanguage(monaco: Monaco): void {
  if (registered) return;
  registered = true;
  monacoApi = monaco;

  disposables.push(
    monaco.languages.registerCompletionItemProvider('go', {
      triggerCharacters: ['.', '('],
      provideCompletionItems(model, position) {
        const word = model.getWordUntilPosition(position);
        const range = {
          startLineNumber: position.lineNumber,
          endLineNumber: position.lineNumber,
          startColumn: word.startColumn,
          endColumn: word.endColumn,
        };

        const funcMap = parseFunctions(currentSource);
        const suggestions: languages.CompletionItem[] = [];

        for (const fn of BUILTIN_FUNCS) {
          const info = lookupFunc(fn, currentLang);
          suggestions.push({
            label: fn,
            kind: monaco.languages.CompletionItemKind.Function,
            detail: info?.signature,
            documentation: info?.description,
            insertText: fn,
            range,
          });
        }

        const seen = new Set<string>();
        for (const key of funcMap.keys()) {
          if (key.includes('(') || seen.has(key)) continue;
          seen.add(key);
          const fn = funcMap.get(key)!;
          suggestions.push({
            label: key,
            kind: monaco.languages.CompletionItemKind.Function,
            detail: formatUserFuncTip(fn),
            insertText: key,
            range,
          });
        }

        for (const tp of GO_TYPES) {
          const info = lookupType(tp, currentLang);
          suggestions.push({
            label: tp,
            kind: monaco.languages.CompletionItemKind.TypeParameter,
            detail: info?.title,
            documentation: info?.description,
            insertText: tp,
            range,
          });
        }

        return { suggestions };
      },
    }),
  );

  disposables.push(
    monaco.languages.registerHoverProvider('go', {
      provideHover(model, position) {
        const word = model.getWordAtPosition(position);
        if (!word) return null;

        const name = word.word;
        const funcMap = parseFunctions(currentSource);
        const range = {
          startLineNumber: position.lineNumber,
          startColumn: word.startColumn,
          endLineNumber: position.lineNumber,
          endColumn: word.endColumn,
        };

        const userFn = funcMap.get(name);
        if (userFn) {
          return {
            range,
            contents: [{ value: `**${name}**\n\n${formatUserFuncTip(userFn)}` }],
          };
        }

        const stdFn = lookupFunc(name, currentLang);
        if (stdFn) {
          return {
            range,
            contents: [{ value: formatFuncTip(stdFn, name) }],
          };
        }

        const typeInfo = lookupType(name, currentLang);
        if (typeInfo) {
          return {
            range,
            contents: [{ value: formatTypeTip(typeInfo, name) }],
          };
        }

        return null;
      },
    }),
  );
}

export function disposeMonacoGoLanguage(): void {
  for (const d of disposables) d.dispose();
  disposables = [];
  registered = false;
}
