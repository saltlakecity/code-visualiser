import type { Monaco } from '@monaco-editor/react';

function cssVar(name: string, fallback: string): string {
  const v = getComputedStyle(document.documentElement).getPropertyValue(name).trim();
  return v || fallback;
}

function hexToRgb(hex: string): { r: number; g: number; b: number } | null {
  const h = hex.replace('#', '');
  if (h.length !== 6) return null;
  return {
    r: parseInt(h.slice(0, 2), 16),
    g: parseInt(h.slice(2, 4), 16),
    b: parseInt(h.slice(4, 6), 16),
  };
}

/** Monaco themes require #RRGGBB or #RRGGBBAA — rgba() is ignored and falls back to defaults (often red guides). */
function colorAlpha(hex: string, alpha: number): string {
  const normalized = hex.startsWith('#') ? hex : `#${hex}`;
  const rgb = hexToRgb(normalized);
  if (!rgb) return '#ffa11620';
  const a = Math.round(Math.min(1, Math.max(0, alpha)) * 255)
    .toString(16)
    .padStart(2, '0');
  return `#${rgb.r.toString(16).padStart(2, '0')}${rgb.g.toString(16).padStart(2, '0')}${rgb.b.toString(16).padStart(2, '0')}${a}`;
}

export function defineMonacoThemes(monaco: Monaco): void {
  const accent = cssVar('--color-accent', '#ffa116');
  const bg = cssVar('--color-bg', '#0f0f13');
  const panel = cssVar('--color-panel', '#16161d');
  const border = cssVar('--color-border', '#2a2a35');
  const text = cssVar('--color-text', '#e8e8ef');
  const muted = cssVar('--color-muted', '#8b8b9e');
  const base = {
    base: 'vs-dark' as const,
    inherit: true,
    rules: [
      { token: 'comment', foreground: '6a737d' },
      { token: 'comment.doc', foreground: '6a737d' },
      { token: 'keyword', foreground: 'ff7b72' },
      { token: 'keyword.go', foreground: 'ff7b72' },
      { token: 'string', foreground: 'a5d6ff' },
      { token: 'string.go', foreground: 'a5d6ff' },
      { token: 'string.escape', foreground: 'a5d6ff' },
      { token: 'number', foreground: '79c0ff' },
      { token: 'number.go', foreground: '79c0ff' },
      { token: 'number.float', foreground: '79c0ff' },
      { token: 'number.hex', foreground: '79c0ff' },
      { token: 'type', foreground: '79c0ff' },
      { token: 'type.identifier', foreground: '79c0ff' },
      { token: 'delimiter', foreground: '8b8b9e' },
      { token: 'delimiter.go', foreground: '8b8b9e' },
      { token: 'operator', foreground: 'e8e8ef' },
      { token: 'identifier', foreground: 'e8e8ef' },
      { token: 'identifier.go', foreground: 'e8e8ef' },
    ],
    colors: {
      'editor.background': panel,
      'editor.foreground': text,
      'editor.lineHighlightBackground': colorAlpha(accent, 0.12),
      'editor.lineHighlightBorder': colorAlpha(accent, 0.45),
      'editorGutter.background': bg,
      'editorLineNumber.foreground': muted,
      'editorLineNumber.activeForeground': accent,
      'editor.selectionBackground': colorAlpha(accent, 0.25),
      'editor.inactiveSelectionBackground': colorAlpha(accent, 0.12),
      'editorCursor.foreground': accent,
      'editorWidget.background': panel,
      'editorWidget.border': border,
      'editorSuggestWidget.background': panel,
      'editorSuggestWidget.border': border,
      'editorHoverWidget.background': panel,
      'editorHoverWidget.border': border,
      'scrollbarSlider.background': colorAlpha(accent, 0.2),
      'minimap.background': bg,
      'editorBracketPairGuide.activeBackground1': colorAlpha(accent, 0.35),
      'editorBracketPairGuide.activeBackground2': colorAlpha(accent, 0.25),
      'editorBracketPairGuide.activeBackground3': colorAlpha(accent, 0.2),
      'editorBracketPairGuide.activeBackground4': colorAlpha(muted, 0.35),
      'editorBracketPairGuide.activeBackground5': colorAlpha(muted, 0.25),
      'editorBracketPairGuide.activeBackground6': colorAlpha(muted, 0.2),
      'editorBracketHighlight.foreground1': colorAlpha(accent, 0.9),
      'editorBracketHighlight.foreground2': colorAlpha(accent, 0.7),
      'editorBracketHighlight.foreground3': colorAlpha(accent, 0.5),
    },
  };

  monaco.editor.defineTheme('visualisator-dark', base);

  monaco.editor.defineTheme('visualisator-light', {
    ...base,
    base: 'vs',
    rules: [
      { token: 'comment', foreground: '6e7781' },
      { token: 'keyword', foreground: 'cf222e' },
      { token: 'string', foreground: '0a3069' },
      { token: 'number', foreground: '0550ae' },
      { token: 'type', foreground: '0550ae' },
      { token: 'function', foreground: '8250df' },
    ],
    colors: {
      ...base.colors,
      'editor.background': panel,
      'editor.foreground': text,
      'editorGutter.background': bg,
      'editorLineNumber.foreground': muted,
      'editorLineNumber.activeForeground': accent,
      'editor.lineHighlightBackground': colorAlpha(accent, 0.1),
      'editor.lineHighlightBorder': colorAlpha(accent, 0.4),
      'editor.selectionBackground': colorAlpha(accent, 0.2),
    },
  });

}

export function monacoThemeId(theme: string): string {
  return theme === 'light' ? 'visualisator-light' : 'visualisator-dark';
}
