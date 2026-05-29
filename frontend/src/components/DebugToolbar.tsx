import {
  Play,
  Square,
  StepForward,
  CornerDownLeft,
  ArrowDown,
  SkipForward,
  RotateCcw,
  Save,
  FolderOpen,
  Pause,
  GitBranch,
} from 'lucide-react';
import { t, translateStatus } from '../i18n';
import { useSettingsStore, currentLang } from '../store/settingsStore';

const SPEEDS = [0.25, 0.5, 1, 2, 4];

interface Props {
  status: string;
  isRunning: boolean;
  canStep: boolean;
  customMode: boolean;
  onCompileRun: () => void;
  onStop: () => void;
  onStep: (cmd: string) => void;
  onPlay: () => void;
  onPause?: () => void;
  isPlaying?: boolean;
  canPlay?: boolean;
  playbackSpeed?: number;
  onSpeedChange?: (speed: number) => void;
  onSaveCustom?: () => void;
  onLoadCustom?: () => void;
  onFlowchart?: () => void;
}

export function DebugToolbar({
  status,
  isRunning,
  canStep,
  customMode,
  onCompileRun,
  onStop,
  onStep,
  onPlay,
  onPause,
  isPlaying = false,
  canPlay = false,
  playbackSpeed = 1,
  onSpeedChange,
  onSaveCustom,
  onLoadCustom,
  onFlowchart,
}: Props) {
  const { settings } = useSettingsStore();
  const tr = t(currentLang(settings));
  const statusLabel = translateStatus(currentLang(settings), status);

  return (
    <div className="flex flex-wrap items-center gap-2 border-b border-border bg-panel px-3 py-2">
      <button
        onClick={onCompileRun}
        disabled={isRunning}
        className="flex items-center gap-1.5 rounded-md bg-accent px-3 py-1.5 text-xs font-semibold text-black disabled:opacity-50"
      >
        <Play className="h-3.5 w-3.5" />
        {tr.run}
      </button>
      <button
        onClick={onStop}
        className="flex items-center gap-1.5 rounded-md border border-border px-3 py-1.5 text-xs hover:bg-bg"
      >
        <Square className="h-3.5 w-3.5" />
        {tr.stop}
      </button>
      {onFlowchart && (
        <button
          onClick={onFlowchart}
          disabled={isRunning}
          className="flex items-center gap-1.5 rounded-md border border-border px-3 py-1.5 text-xs hover:bg-bg disabled:opacity-50"
        >
          <GitBranch className="h-3.5 w-3.5" />
          {tr.flowchartButton}
        </button>
      )}
      {customMode && (
        <>
          <button
            onClick={onSaveCustom}
            className="flex items-center gap-1 rounded-md border border-border px-2 py-1.5 text-xs hover:bg-bg"
          >
            <Save className="h-3.5 w-3.5" />
            {tr.custom.save}
          </button>
          <button
            onClick={onLoadCustom}
            className="flex items-center gap-1 rounded-md border border-border px-2 py-1.5 text-xs hover:bg-bg"
          >
            <FolderOpen className="h-3.5 w-3.5" />
            {tr.custom.load}
          </button>
        </>
      )}
      <div className="mx-1 h-5 w-px bg-border" />
      <ToolBtn icon={<StepForward className="h-3.5 w-3.5" />} label={tr.next} disabled={!canStep} onClick={() => onStep('next')} />
      <ToolBtn icon={<ArrowDown className="h-3.5 w-3.5" />} label={tr.stepIn} disabled={!canStep} onClick={() => onStep('step')} />
      <ToolBtn icon={<CornerDownLeft className="h-3.5 w-3.5" />} label={tr.stepOut} disabled={!canStep} onClick={() => onStep('stepOut')} />
      <ToolBtn icon={<SkipForward className="h-3.5 w-3.5" />} label={tr.continue} disabled={!canStep} onClick={() => onStep('continue')} />
      {isPlaying ? (
        <ToolBtn icon={<Pause className="h-3.5 w-3.5" />} label={tr.pause} disabled={!canPlay} onClick={() => onPause?.()} />
      ) : (
        <ToolBtn icon={<RotateCcw className="h-3.5 w-3.5" />} label={tr.playTrace} disabled={!canPlay} onClick={onPlay} />
      )}
      {canPlay && onSpeedChange && (
        <select
          value={playbackSpeed}
          onChange={(e) => onSpeedChange(Number(e.target.value))}
          className="rounded border border-border bg-bg px-2 py-1.5 text-xs"
          title={tr.playbackSpeed}
        >
          {SPEEDS.map((s) => (
            <option key={s} value={s}>
              {s}x
            </option>
          ))}
        </select>
      )}
      <span className="ml-auto rounded-full bg-bg px-2 py-0.5 text-[10px] uppercase tracking-wide text-muted">
        {statusLabel}
      </span>
    </div>
  );
}

function ToolBtn({
  icon,
  label,
  disabled,
  onClick,
}: {
  icon: React.ReactNode;
  label: string;
  disabled?: boolean;
  onClick: () => void;
}) {
  return (
    <button
      title={label}
      disabled={disabled}
      onClick={onClick}
      className="flex items-center gap-1 rounded-md border border-border px-2 py-1.5 text-xs hover:bg-bg disabled:opacity-40"
    >
      {icon}
      <span className="hidden lg:inline">{label}</span>
    </button>
  );
}
