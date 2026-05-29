import Editor, { type Monaco } from '@monaco-editor/react';
import type * as MonacoEditor from 'monaco-editor';
import { useCallback, useEffect, useMemo, useRef } from 'react';
import type { Language } from '../i18n';
import { registerMonacoGoLanguage, setMonacoGoContext } from '../monacoGoLanguage';
import { defineMonacoThemes, monacoThemeId } from '../monacoTheme';
import { useSettingsStore } from '../store/settingsStore';

interface Props {
  source: string;
  onChange?: (v: string) => void;
  readOnly?: boolean;
  highlightLine?: number;
  scrollKey?: number;
  errorLines?: Set<number>;
  lang?: Language;
  changedVars?: string[];
  activeVar?: string | null;
}

function escapeRegex(s: string): string {
  return s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

function varDecorations(
  monaco: Monaco,
  model: MonacoEditor.editor.ITextModel,
  changedVars: string[],
  activeVar: string | null,
): MonacoEditor.editor.IModelDeltaDecoration[] {
  const changed = new Set(changedVars.map((v) => v.split('.')[0]));
  const active = activeVar?.split('.')[0] ?? null;
  const names = [...changed];
  if (active && !names.includes(active)) names.push(active);
  if (names.length === 0) return [];

  names.sort((a, b) => b.length - a.length);
  const decos: MonacoEditor.editor.IModelDeltaDecoration[] = [];
  const lineCount = model.getLineCount();

  for (let line = 1; line <= lineCount; line++) {
    const text = model.getLineContent(line);
    for (const v of names) {
      const re = new RegExp(`\\b${escapeRegex(v)}\\b`, 'g');
      let m: RegExpExecArray | null;
      while ((m = re.exec(text)) !== null) {
        const start = m.index + 1;
        const end = start + v.length;
        const isActive = active === v;
        decos.push({
          range: new monaco.Range(line, start, line, end),
          options: {
            inlineClassName: isActive ? 'go-var-active' : changed.has(v) ? 'go-var-changed' : undefined,
            stickiness: monaco.editor.TrackedRangeStickiness.NeverGrowsWhenTypingAtEdges,
          },
        });
      }
    }
  }
  return decos;
}

export function MonacoGoEditor({
  source,
  onChange,
  readOnly = false,
  highlightLine = 0,
  scrollKey,
  errorLines,
  lang = 'en',
  changedVars = [],
  activeVar = null,
}: Props) {
  const theme = useSettingsStore((s) => s.settings.theme || 'dark');
  const editorRef = useRef<MonacoEditor.editor.IStandaloneCodeEditor | null>(null);
  const monacoRef = useRef<Monaco | null>(null);
  const decoIdsRef = useRef<string[]>([]);
  const varDecoIdsRef = useRef<string[]>([]);
  const lastScrollKey = useRef<number | undefined>(undefined);

  useEffect(() => {
    setMonacoGoContext(lang, source);
  }, [lang, source]);

  const applyDecorations = useCallback(() => {
    const editor = editorRef.current;
    const monaco = monacoRef.current;
    const model = editor?.getModel();
    if (!editor || !monaco || !model) return;

    const decos: MonacoEditor.editor.IModelDeltaDecoration[] = [];

    if (highlightLine > 0 && highlightLine <= model.getLineCount()) {
      decos.push({
        range: new monaco.Range(highlightLine, 1, highlightLine, 1),
        options: {
          isWholeLine: true,
          className: 'monaco-debug-current-line',
          linesDecorationsClassName: 'monaco-debug-gutter-bar',
        },
      });
    }

    if (!readOnly && errorLines?.size) {
      for (const line of errorLines) {
        if (line <= 0 || line > model.getLineCount()) continue;
        decos.push({
          range: new monaco.Range(line, 1, line, 1),
          options: {
            isWholeLine: true,
            className: 'monaco-compile-error-line',
            linesDecorationsClassName: 'monaco-compile-error-gutter',
          },
        });
      }
    }

    decoIdsRef.current = editor.deltaDecorations(decoIdsRef.current, decos);

    varDecoIdsRef.current = editor.deltaDecorations(
      varDecoIdsRef.current,
      varDecorations(monaco, model, changedVars, activeVar),
    );
  }, [highlightLine, changedVars, activeVar, readOnly, errorLines]);

  const clearErrorMarkers = useCallback(() => {
    const monaco = monacoRef.current;
    const model = editorRef.current?.getModel();
    if (monaco && model) monaco.editor.setModelMarkers(model, 'compile', []);
  }, []);

  useEffect(() => {
    applyDecorations();
  }, [applyDecorations, source]);

  useEffect(() => {
    clearErrorMarkers();
  }, [clearErrorMarkers, source, readOnly]);

  useEffect(() => {
    const editor = editorRef.current;
    const monaco = monacoRef.current;
    if (!editor || !monaco || highlightLine <= 0) return;
    if (scrollKey !== undefined && scrollKey === lastScrollKey.current) return;
    lastScrollKey.current = scrollKey;
    editor.revealLineInCenter(highlightLine);
  }, [highlightLine, scrollKey]);

  useEffect(() => {
    const monaco = monacoRef.current;
    const editor = editorRef.current;
    if (!monaco || !editor) return;
    defineMonacoThemes(monaco);
    monaco.editor.setTheme(monacoThemeId(theme));
  }, [theme]);

  const handleMount = useCallback(
    (editor: MonacoEditor.editor.IStandaloneCodeEditor, monaco: Monaco) => {
      editorRef.current = editor;
      monacoRef.current = monaco;
      defineMonacoThemes(monaco);
      registerMonacoGoLanguage(monaco);
      monaco.editor.setTheme(monacoThemeId(theme));
      const model = editor.getModel();
      if (model) monaco.editor.setModelLanguage(model, 'go');
      applyDecorations();
      clearErrorMarkers();
    },
    [theme, applyDecorations, clearErrorMarkers],
  );

  const options = useMemo<MonacoEditor.editor.IStandaloneEditorConstructionOptions>(
    () => ({
      readOnly,
      language: 'go',
      fontFamily: 'Consolas, "Courier New", monospace',
      fontSize: 14,
      lineHeight: 22,
      tabSize: 4,
      insertSpaces: false,
      minimap: { enabled: false },
      scrollBeyondLastLine: false,
      automaticLayout: true,
      fixedOverflowWidgets: true,
      bracketPairColorization: { enabled: true, independentColorPoolPerBracketType: false },
      guides: {
        bracketPairs: false,
        bracketPairsHorizontal: false,
        highlightActiveBracketPair: true,
        indentation: true,
      },
      folding: true,
      lineNumbers: 'on',
      renderLineHighlight: readOnly ? 'none' : 'line',
      unicodeHighlight: {
        ambiguousCharacters: false,
        invisibleCharacters: false,
        nonBasicASCII: false,
      },
      renderValidationDecorations: 'off',
      glyphMargin: false,
      padding: { top: 8, bottom: 8 },
      wordWrap: 'off',
      scrollbar: { verticalScrollbarSize: 10, horizontalScrollbarSize: 10 },
    }),
    [readOnly, highlightLine],
  );

  return (
    <div
      className="monaco-go-editor h-full min-h-0 flex-1 overflow-hidden"
      style={{ boxShadow: 'var(--editor-glow)' }}
    >
      <Editor
        height="100%"
        language="go"
        value={source}
        theme={monacoThemeId(theme)}
        options={options}
        loading={<div className="flex h-full items-center justify-center text-muted text-xs">Loading editor…</div>}
        beforeMount={(monaco) => {
          defineMonacoThemes(monaco);
          registerMonacoGoLanguage(monaco);
        }}
        onMount={handleMount}
        onChange={(v) => {
          if (!readOnly && onChange && v !== undefined) onChange(v);
        }}
      />
    </div>
  );
}
