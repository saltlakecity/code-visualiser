import { useCallback, useEffect, useRef } from 'react';

interface Props {
  varsHeight: number;
  onResize: (height: number) => void;
  minHeight?: number;
  maxHeight?: number;
  varsPanel: React.ReactNode;
  vizPanel: React.ReactNode;
  label: string;
}

export function ResizableRightPanel({
  varsHeight,
  onResize,
  minHeight = 160,
  maxHeight = 480,
  varsPanel,
  vizPanel,
  label,
}: Props) {
  const dragging = useRef(false);
  const startY = useRef(0);
  const startH = useRef(0);

  const onMouseDown = (e: React.MouseEvent) => {
    dragging.current = true;
    startY.current = e.clientY;
    startH.current = varsHeight;
    e.preventDefault();
  };

  const onMouseMove = useCallback(
    (e: MouseEvent) => {
      if (!dragging.current) return;
      const delta = e.clientY - startY.current;
      const next = Math.min(maxHeight, Math.max(minHeight, startH.current + delta));
      onResize(next);
    },
    [maxHeight, minHeight, onResize],
  );

  const onMouseUp = useCallback(() => {
    dragging.current = false;
  }, []);

  useEffect(() => {
    window.addEventListener('mousemove', onMouseMove);
    window.addEventListener('mouseup', onMouseUp);
    return () => {
      window.removeEventListener('mousemove', onMouseMove);
      window.removeEventListener('mouseup', onMouseUp);
    };
  }, [onMouseMove, onMouseUp]);

  return (
    <div className="flex min-h-0 flex-1 flex-col overflow-hidden">
      <div
        style={{ height: varsHeight }}
        className="flex min-h-0 shrink-0 flex-col overflow-hidden border-b border-border"
      >
        <div className="panel-scroll">{varsPanel}</div>
      </div>
      <div
        onMouseDown={onMouseDown}
        className="flex h-1.5 shrink-0 cursor-row-resize items-center justify-center border-b border-border bg-bg hover:bg-accent/20"
      >
        <div className="h-0.5 w-8 rounded bg-border" />
      </div>
      <div className="flex min-h-0 flex-1 flex-col overflow-hidden">
        <div className="shrink-0 border-b border-border px-3 py-2 text-xs font-semibold uppercase text-muted">
          {label}
        </div>
        <div className="panel-scroll">{vizPanel}</div>
      </div>
    </div>
  );
}
