import { Code2, Search, Settings } from 'lucide-react';
import { useState } from 'react';
import { t } from '../i18n';
import { useSettingsStore, currentLang } from '../store/settingsStore';
import type { ProblemMeta } from '../types';

interface Props {
  problems: ProblemMeta[];
  selectedId: string | null;
  customMode: boolean;
  onSelect: (id: string) => void;
  onCustom: () => void;
}

export function ProblemSidebar({ problems, selectedId, customMode, onSelect, onCustom }: Props) {
  const [query, setQuery] = useState('');
  const { settings, setSettingsOpen } = useSettingsStore();
  const tr = t(currentLang(settings));

  const filtered = problems.filter(
    (p) =>
      p.title.toLowerCase().includes(query.toLowerCase()) ||
      p.tags.some((tag) => tag.toLowerCase().includes(query.toLowerCase())),
  );

  return (
    <aside className="flex h-full min-h-0 w-60 flex-col overflow-hidden border-r border-border bg-panel">
      <div className="flex items-start justify-between border-b border-border p-3">
        <div>
          <h1 className="text-sm font-semibold text-accent">{tr.appTitle}</h1>
          <p className="mt-1 text-xs text-muted">{tr.appSubtitle}</p>
        </div>
        <button
          onClick={() => setSettingsOpen(true)}
          className="rounded p-1 text-muted hover:bg-bg hover:text-accent"
          title={tr.settings.title}
        >
          <Settings className="h-4 w-4" />
        </button>
      </div>
      <div className="p-2">
        <div className="relative">
          <Search className="absolute left-2 top-2.5 h-3.5 w-3.5 text-muted" />
          <input
            className="w-full rounded-md border border-border bg-bg py-2 pl-8 pr-2 text-xs outline-none focus:border-accent"
            placeholder={tr.searchProblems}
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
        </div>
      </div>
      <button
        onClick={onCustom}
        className={`mx-2 mb-2 flex items-center gap-2 rounded-md px-3 py-2 text-left text-xs transition ${
          customMode ? 'bg-accent/20 text-accent' : 'hover:bg-bg text-muted'
        }`}
      >
        <Code2 className="h-4 w-4" />
        {tr.customCode}
      </button>
      <nav className="flex-1 overflow-y-auto px-2 pb-3">
        {filtered.map((p) => (
          <button
            key={p.id}
            onClick={() => onSelect(p.id)}
            className={`mb-1 w-full rounded-md px-3 py-2 text-left transition ${
              selectedId === p.id && !customMode
                ? 'bg-accent/15 text-accent'
                : 'text-sm hover:bg-bg'
            }`}
          >
            <div className="font-medium">{p.title}</div>
            <div className="mt-1 flex flex-wrap gap-1">
              {p.tags.map((tag) => (
                <span key={tag} className="rounded bg-bg px-1.5 py-0.5 text-[10px] text-muted">
                  {tag}
                </span>
              ))}
            </div>
          </button>
        ))}
      </nav>
    </aside>
  );
}
