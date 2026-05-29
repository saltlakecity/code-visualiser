import { loader } from '@monaco-editor/react';
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api';
import 'monaco-editor/esm/vs/basic-languages/go/go.contribution';
// Preload tokenizer so syntax highlighting is ready on first paint
import 'monaco-editor/esm/vs/basic-languages/go/go.js';
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker';

self.MonacoEnvironment = {
  getWorker() {
    return new editorWorker();
  },
};

loader.config({ monaco });

export const monacoReady = loader.init();
