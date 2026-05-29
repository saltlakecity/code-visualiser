import { useState } from 'react';
import { t, type Language } from '../i18n';
import { useSettingsStore, currentLang } from '../store/settingsStore';

const KEY = 'visualisator-disclaimer-accepted';

export function useDisclaimer() {
  const [open, setOpen] = useState(() => !localStorage.getItem(KEY));
  const accept = () => {
    localStorage.setItem(KEY, '1');
    setOpen(false);
  };
  return { open, accept };
}

export function DisclaimerModal({ onAccept }: { onAccept: () => void }) {
  const { settings } = useSettingsStore();
  const tr = t(currentLang(settings));

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 p-4">
      <div className="max-w-md rounded-lg border border-border bg-panel p-6 shadow-xl">
        <h2 className="text-lg font-semibold text-accent">{tr.disclaimer.title}</h2>
        <p className="mt-3 text-sm leading-relaxed text-muted">{tr.disclaimer.body}</p>
        <button
          onClick={onAccept}
          className="mt-5 w-full rounded-md bg-accent py-2 text-sm font-semibold text-black"
        >
          {tr.disclaimer.accept}
        </button>
      </div>
    </div>
  );
}
