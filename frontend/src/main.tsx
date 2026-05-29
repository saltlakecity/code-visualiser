import { createRoot } from 'react-dom/client';
import App from './App';
import './index.css';
import 'monaco-editor/esm/vs/editor/standalone/browser/standalone-tokens.css';
import './monacoSetup';

const rootEl = document.getElementById('root');
if (rootEl) {
  createRoot(rootEl).render(<App />);
}
