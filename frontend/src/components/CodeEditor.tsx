import Editor, { type OnMount } from '@monaco-editor/react';
import type { editor } from 'monaco-editor';
import { useEffect, useRef } from 'react';
import type { CompileError } from '../types';

interface Props {
  source: string;
  onChange: (v: string) => void;
  currentLine?: number;
  currentFile?: string;
  errors?: CompileError[];
  readOnly?: boolean;
}

export function CodeEditor({
  source,
  onChange,
  currentLine,
  currentFile = 'main.go',
  errors = [],
  readOnly = false,
}: Props) {
  const editorRef = useRef<editor.IStandaloneCodeEditor | null>(null);
  const decoRef = useRef<editor.IEditorDecorationsCollection | null>(null);

  const onMount: OnMount = (ed) => {
    editorRef.current = ed;
    decoRef.current = ed.createDecorationsCollection([]);
  };

  useEffect(() => {
    const ed = editorRef.current;
    if (!ed || !decoRef.current) return;

    const decorations: editor.IModelDeltaDecoration[] = [];

    if (currentLine && currentLine > 0) {
      decorations.push({
        range: { startLineNumber: currentLine, startColumn: 1, endLineNumber: currentLine, endColumn: 1 },
        options: {
          isWholeLine: true,
          className: 'current-line-highlight',
          glyphMarginClassName: 'current-line-glyph',
        },
      });
      ed.revealLineInCenter(currentLine);
    }

    for (const err of errors) {
      if (err.file && !err.file.endsWith(currentFile) && err.file !== 'main.go' && err.file !== 'solution.go') {
        continue;
      }
      decorations.push({
        range: {
          startLineNumber: err.line,
          startColumn: err.column || 1,
          endLineNumber: err.line,
          endColumn: err.column + 1 || 2,
        },
        options: {
          isWholeLine: true,
          className: 'error-line-highlight',
          glyphMarginClassName: 'error-line-glyph',
          hoverMessage: { value: err.message },
        },
      });
    }

    decoRef.current.set(decorations);
  }, [currentLine, currentFile, errors]);

  return (
    <div className="h-full w-full overflow-hidden rounded-lg border border-border">
      <Editor
        height="100%"
        defaultLanguage="go"
        theme="vs-dark"
        value={source}
        onChange={(v) => onChange(v ?? '')}
        onMount={onMount}
        options={{
          readOnly,
          fontFamily: 'JetBrains Mono, Consolas, monospace',
          fontSize: 14,
          minimap: { enabled: false },
          scrollBeyondLastLine: false,
          glyphMargin: true,
          lineNumbers: 'on',
          renderLineHighlight: 'none',
          padding: { top: 12 },
        }}
      />
      <style>{`
        .current-line-highlight { background: rgba(255, 161, 22, 0.15) !important; }
        .current-line-glyph { background: #ffa116; width: 4px !important; margin-left: 3px; }
        .error-line-highlight { background: rgba(255, 71, 71, 0.2) !important; }
        .error-line-glyph { background: #ff4747; width: 4px !important; margin-left: 3px; }
      `}</style>
    </div>
  );
}
