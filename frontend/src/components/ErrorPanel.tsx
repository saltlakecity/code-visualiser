import { t } from '../i18n';
import { useSettingsStore, currentLang } from '../store/settingsStore';
import type { Trace } from '../types';

interface Props {
  message: string;
  trace?: Trace | null;
  onGoToFailure?: () => void;
}

export function ErrorPanel({ message, trace, onGoToFailure }: Props) {
  const { settings } = useSettingsStore();
  const tr = t(currentLang(settings));

  const displayMsg = trace?.failureReason || message;
  const isRuntimeFailure = trace?.status === 'panic' || trace?.status === 'error';
  if (!displayMsg || (trace && !isRuntimeFailure && !message)) return null;

  const failureIdx = trace?.failureFrameIndex;
  const failureLine = trace?.failureEditorLine;

  return (
    <div className="max-h-32 overflow-y-auto border-t border-error/50 bg-error/10 p-3">
      <h3 className="mb-1 text-xs font-semibold text-error">{tr.runtimeError}</h3>
      {failureIdx != null && failureIdx >= 0 && (
        <p className="mb-1 text-[11px] text-error/80">
          {tr.failureAtStep} {failureIdx + 1}
          {failureLine ? ` · ${tr.failureLine} ${failureLine}` : ''}
        </p>
      )}
      <pre className="whitespace-pre-wrap font-mono text-xs text-error/90">{displayMsg}</pre>
      {onGoToFailure && failureIdx != null && failureIdx >= 0 && (
        <button
          type="button"
          onClick={onGoToFailure}
          className="mt-2 rounded border border-error/40 px-2 py-0.5 text-[10px] text-error hover:bg-error/10"
        >
          {tr.goToFailureStep}
        </button>
      )}
    </div>
  );
}
