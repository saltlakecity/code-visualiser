import {
  useCallback,
  useEffect,
  useLayoutEffect,
  useRef,
  useState,
  type CSSProperties,
  type ReactNode,
} from 'react';
import { createPortal } from 'react-dom';

const VIEWPORT_MARGIN = 8;
const GAP = 6;
const MAX_WIDTH = 320;

interface Props {
  tip: string;
  children: ReactNode;
  className?: string;
}

function computePosition(anchor: DOMRect, tip: DOMRect): CSSProperties {
  const vw = window.innerWidth;
  const vh = window.innerHeight;

  let top = anchor.top - tip.height - GAP;
  let left = anchor.left;

  if (top < VIEWPORT_MARGIN) {
    top = anchor.bottom + GAP;
  }
  if (top + tip.height > vh - VIEWPORT_MARGIN) {
    top = Math.max(VIEWPORT_MARGIN, vh - VIEWPORT_MARGIN - tip.height);
  }

  if (left + tip.width > vw - VIEWPORT_MARGIN) {
    left = vw - VIEWPORT_MARGIN - tip.width;
  }
  if (left < VIEWPORT_MARGIN) {
    left = VIEWPORT_MARGIN;
  }

  return {
    position: 'fixed',
    top,
    left,
    zIndex: 10000,
    maxWidth: MAX_WIDTH,
  };
}

export function AdaptiveTooltip({ tip, children, className = '' }: Props) {
  const [visible, setVisible] = useState(false);
  const [ready, setReady] = useState(false);
  const [pos, setPos] = useState<CSSProperties>({});
  const anchorRef = useRef<HTMLSpanElement>(null);
  const popupRef = useRef<HTMLDivElement>(null);

  const reposition = useCallback(() => {
    const anchor = anchorRef.current;
    const popup = popupRef.current;
    if (!anchor || !popup || !tip) return;
    setPos(computePosition(anchor.getBoundingClientRect(), popup.getBoundingClientRect()));
    setReady(true);
  }, [tip]);

  useLayoutEffect(() => {
    if (!visible) {
      setReady(false);
      return;
    }
    reposition();
  }, [visible, tip, reposition]);

  useEffect(() => {
    if (!visible) return;
    const onMove = () => reposition();
    window.addEventListener('scroll', onMove, true);
    window.addEventListener('resize', onMove);
    return () => {
      window.removeEventListener('scroll', onMove, true);
      window.removeEventListener('resize', onMove);
    };
  }, [visible, reposition]);

  const show = () => tip && setVisible(true);
  const hide = () => setVisible(false);

  return (
    <>
      <span
        ref={anchorRef}
        className={className}
        onMouseEnter={show}
        onMouseLeave={hide}
        onFocus={show}
        onBlur={hide}
      >
        {children}
      </span>
      {visible &&
        tip &&
        createPortal(
          <div
            ref={popupRef}
            role="tooltip"
            className="go-tip-popup pointer-events-none rounded-md border border-border bg-panel px-2 py-1.5 font-sans text-[11px] leading-snug text-[var(--color-text)] shadow-lg"
            style={{
              ...pos,
              visibility: ready ? 'visible' : 'hidden',
              whiteSpace: 'pre-wrap',
            }}
          >
            {tip}
          </div>,
          document.body,
        )}
    </>
  );
}
