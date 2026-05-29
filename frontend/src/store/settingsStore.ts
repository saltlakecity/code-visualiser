import { create } from 'zustand';
import { GetSettings, SaveSettings } from '../../wailsjs/go/app/App';
import { models } from '../../wailsjs/go/models';
import type { Language } from '../i18n';

export type AppSettings = models.AppSettings;

interface SettingsState {
  settings: AppSettings;
  settingsOpen: boolean;
  loaded: boolean;
  toast: string | null;
  loadSettings: () => Promise<void>;
  saveSettings: (s: AppSettings) => Promise<void>;
  setSettingsOpen: (v: boolean) => void;
  applyTheme: () => void;
  setToast: (msg: string | null) => void;
}

const defaults = new models.AppSettings({
  language: 'ru',
  theme: 'dark',
  accentColor: '#ffa116',
  customCodeDir: '',
});

export const useSettingsStore = create<SettingsState>((set, get) => ({
  settings: defaults,
  settingsOpen: false,
  loaded: false,
  toast: null,
  loadSettings: async () => {
    try {
      const s = await GetSettings();
      set({
        settings: new models.AppSettings({
          language: s.language || 'ru',
          theme: s.theme || 'dark',
          accentColor: s.accentColor || '#ffa116',
          customCodeDir: s.customCodeDir || '',
          dataDir: s.dataDir || '',
        }),
        loaded: true,
      });
      get().applyTheme();
    } catch {
      set({ loaded: true });
      get().applyTheme();
    }
  },
  saveSettings: async (s) => {
    await SaveSettings(s);
    set({ settings: s });
    get().applyTheme();
  },
  setSettingsOpen: (settingsOpen) => set({ settingsOpen }),
  setToast: (toast) => set({ toast }),
  applyTheme: () => {
    const { settings } = get();
    const root = document.documentElement;
    const accent = settings.accentColor || '#ffa116';
    root.style.setProperty('--color-accent', accent);
    root.dataset.theme = settings.theme || 'dark';
    if (settings.theme === 'light') {
      root.style.setProperty('--color-bg', '#f0f0f5');
      root.style.setProperty('--color-panel', '#ffffff');
      root.style.setProperty('--color-border', '#d8d8e0');
      root.style.setProperty('--color-text', '#1a1a22');
      root.style.setProperty('--color-muted', '#5c5c6e');
    } else {
      root.style.setProperty('--color-bg', '#0f0f13');
      root.style.setProperty('--color-panel', '#16161d');
      root.style.setProperty('--color-border', '#2a2a35');
      root.style.setProperty('--color-text', '#e8e8ef');
      root.style.setProperty('--color-muted', '#8b8b9e');
    }
  },
}));

export function currentLang(settings: AppSettings): Language {
  return settings.language === 'en' ? 'en' : 'ru';
}
