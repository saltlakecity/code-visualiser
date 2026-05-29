import { t } from '../i18n';
import { useSettingsStore, currentLang } from '../store/settingsStore';
import type { TestCase } from '../types';

export function TestPanel({ tests }: { tests: TestCase[] }) {
  const { settings } = useSettingsStore();
  const tr = t(currentLang(settings));

  if (!tests.length) return null;

  return (
    <div className="border-t border-border bg-panel p-3">
      <h3 className="mb-2 text-xs font-semibold uppercase tracking-wide text-muted">{tr.testCases}</h3>
      <div className="space-y-2">
        {tests.map((tc, i) => (
          <div key={i} className="rounded-md border border-border bg-bg p-2 font-mono text-xs">
            <div>
              <span className="text-muted">{tr.input}: </span>
              {tc.input}
            </div>
            <div>
              <span className="text-muted">{tr.expected}: </span>
              <span className="text-success">{tc.expected}</span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
