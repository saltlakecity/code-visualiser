import { en, type Locale } from './en';
import { ru } from './ru';

export type Language = 'en' | 'ru';

const dict: Record<Language, Locale> = { en, ru };

export function t(lang: Language): Locale {
  return dict[lang] ?? en;
}

export function translateStatus(lang: Language, status: string): string {
  const map = dict[lang]?.status ?? en.status;
  const normalized = status.toLowerCase().replace(/\s+/g, '');
  const aliases: Record<string, keyof typeof map> = {
    idle: 'idle',
    compiling: 'compiling',
    debugging: 'debugging',
    paused: 'paused',
    finished: 'finished',
    panic: 'panic',
    error: 'error',
    stopped: 'stopped',
    compileerror: 'compileError',
  };
  const key = aliases[normalized];
  if (key && map[key]) return map[key];
  return status;
}
