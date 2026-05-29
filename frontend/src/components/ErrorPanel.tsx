import { t } from '../i18n';
import { useSettingsStore, currentLang } from '../store/settingsStore';

export function ErrorPanel({ message }: { message: string }) {
  const { settings } = useSettingsStore();
  const tr = t(currentLang(settings));

  if (!message) return null;
  return (
    <div className="max-h-32 overflow-y-auto border-t border-error/50 bg-error/10 p-3">
      <h3 className="mb-1 text-xs font-semibold text-error">{tr.runtimeError}</h3>
      <pre className="whitespace-pre-wrap font-mono text-xs text-error/90">{message}</pre>
    </div>
  );
}
