import { t } from '../i18n';
import { useSettingsStore, currentLang } from '../store/settingsStore';
import type { StackFrame } from '../types';

export function CallStack({ stack }: { stack: StackFrame[] }) {
  const { settings } = useSettingsStore();
  const tr = t(currentLang(settings));

  if (!stack.length) {
    return <p className="p-3 text-xs text-muted">{tr.callStackEmpty}</p>;
  }
  return (
    <div className="p-2">
      {stack.map((f, i) => (
        <div
          key={`${f.function}-${i}`}
          className={`mb-1 rounded-md border px-2 py-1.5 font-mono text-xs ${
            i === 0 ? 'border-accent/40 bg-accent/10' : 'border-border bg-bg'
          }`}
        >
          <div className="font-medium text-accent">{f.function || '(unknown)'}</div>
          <div className="text-muted">
            {f.file}:{f.line}
          </div>
        </div>
      ))}
    </div>
  );
}
