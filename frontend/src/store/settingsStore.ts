import { create } from 'zustand';
import { GetSettings, SaveSettings } from '../../wailsjs/go/app/App';
import { models } from '../../wailsjs/go/models';
import type { Language } from '../i18n';

export type UIStyle = 'default' | 'minimal' | 'premium';
export type AppSettings = models.AppSettings & { uiStyle?: string };

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
  uiStyle: 'default',
  accentColor: '#ffa116',
  customCodeDir: '',
  playbackSpeed: 1,
}) as AppSettings;

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
          uiStyle: s.uiStyle || 'default',
          accentColor: s.accentColor || '#ffa116',
          customCodeDir: s.customCodeDir || '',
          dataDir: s.dataDir || '',
          playbackSpeed: s.playbackSpeed || 1,
        }) as AppSettings,
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
    root.dataset.uiStyle = (settings.uiStyle as UIStyle) || 'default';

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

    const uiStyle = (settings.uiStyle as UIStyle) || 'default';
    if (uiStyle === 'minimal') {
      root.style.setProperty('--editor-glow', 'none');
      root.style.setProperty('--panel-blur', 'none');
      root.style.setProperty('--radius-base', '2px');
    } else if (uiStyle === 'premium') {
      root.style.setProperty('--editor-glow', '0 0 24px color-mix(in srgb, var(--color-accent) 12%, transparent)');
      root.style.setProperty('--panel-blur', 'blur(8px)');
      root.style.setProperty('--radius-base', '10px');
    } else {
      root.style.setProperty('--editor-glow', 'none');
      root.style.setProperty('--panel-blur', 'none');
      root.style.setProperty('--radius-base', '8px');
    }
  },
}));

export function currentLang(settings: AppSettings): Language {
  return settings.language === 'en' ? 'en' : 'ru';
}
