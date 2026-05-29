import Prism from 'prismjs';
import 'prismjs/components/prism-go';

export function highlightGoLine(line: string): string {
  const grammar = Prism.languages.go;
  if (!grammar) return escapeHtml(line);
  return Prism.highlight(line, grammar, 'go');
}

function escapeHtml(s: string): string {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;');
}
