import { useEffect, useState } from 'react';
import { PickCustomCodeDir } from '../../wailsjs/go/app/App';
import { t, type Language } from '../i18n';
import { models } from '../../wailsjs/go/models';
import { useSettingsStore, type AppSettings } from '../store/settingsStore';

export function SettingsModal() {
  const { settings, settingsOpen, setSettingsOpen, saveSettings } = useSettingsStore();
  const [draft, setDraft] = useState<AppSettings>(settings);
  const lang = (draft.language === 'en' ? 'en' : 'ru') as Language;
  const tr = t(lang);

  useEffect(() => {
    if (settingsOpen) setDraft(settings);
  }, [settingsOpen, settings]);

  if (!settingsOpen) return null;

  const onSave = async () => {
    await saveSettings(draft);
    setSettingsOpen(false);
  };

  const onBrowse = async () => {
    try {
      const dir = await PickCustomCodeDir();
      if (dir) setDraft({ ...draft, customCodeDir: dir });
    } catch {
      /* cancelled */
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
      <div className="w-full max-w-md rounded-lg border border-border bg-panel p-5 shadow-xl">
        <h2 className="text-lg font-semibold text-accent">{tr.settings.title}</h2>

        <label className="mt-4 block text-xs text-muted">{tr.settings.language}</label>
        <select
          className="mt-1 w-full rounded-md border border-border bg-bg px-3 py-2 text-sm"
          value={draft.language}
          onChange={(e) => setDraft({ ...draft, language: e.target.value })}
        >
          <option value="ru">{tr.settings.langRu}</option>
          <option value="en">{tr.settings.langEn}</option>
        </select>

        <label className="mt-3 block text-xs text-muted">{tr.settings.theme}</label>
        <select
          className="mt-1 w-full rounded-md border border-border bg-bg px-3 py-2 text-sm"
          value={draft.theme}
          onChange={(e) => setDraft({ ...draft, theme: e.target.value })}
        >
          <option value="dark">{tr.settings.dark}</option>
          <option value="light">{tr.settings.light}</option>
        </select>

        <label className="mt-3 block text-xs text-muted">{tr.settings.accentColor}</label>
        <div className="mt-1 flex gap-2">
          <input
            type="color"
            value={draft.accentColor}
            onChange={(e) => setDraft({ ...draft, accentColor: e.target.value })}
            className="h-10 w-14 cursor-pointer rounded border border-border bg-bg"
          />
          <input
            type="text"
            value={draft.accentColor}
            onChange={(e) => setDraft({ ...draft, accentColor: e.target.value })}
            className="flex-1 rounded-md border border-border bg-bg px-3 py-2 font-mono text-sm"
          />
        </div>

        <label className="mt-3 block text-xs text-muted">{tr.settings.customCodeDir}</label>
        <div className="mt-1 flex gap-2">
          <input
            type="text"
            value={draft.customCodeDir}
            onChange={(e) => setDraft({ ...draft, customCodeDir: e.target.value })}
            className="flex-1 rounded-md border border-border bg-bg px-3 py-2 text-xs"
          />
          <button
            onClick={onBrowse}
            className="rounded-md border border-border px-3 py-2 text-xs hover:bg-bg"
          >
            {tr.settings.browse}
          </button>
        </div>

        <div className="mt-6 flex justify-end gap-2">
          <button
            onClick={() => setSettingsOpen(false)}
            className="rounded-md border border-border px-4 py-2 text-sm hover:bg-bg"
          >
            {tr.settings.cancel}
          </button>
          <button
            onClick={onSave}
            className="rounded-md bg-accent px-4 py-2 text-sm font-semibold text-black"
          >
            {tr.settings.save}
          </button>
        </div>
      </div>
    </div>
  );
}
