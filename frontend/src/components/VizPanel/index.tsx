import { lazy, Suspense } from 'react';
import { t } from '../../i18n';
import { useDebuggerStore } from '../../store/debuggerStore';
import { useSettingsStore, currentLang } from '../../store/settingsStore';
import type { VizModel } from '../../types';
import { ArrayViz } from './ArrayViz';
import { MapViz, MatrixViz } from './MatrixViz';

const GraphViz = lazy(() => import('./GraphViz').then((m) => ({ default: m.GraphViz })));

export function VizPanel({ models }: { models: VizModel[] }) {
  const { settings } = useSettingsStore();
  const { highlightedVar } = useDebuggerStore();
  const tr = t(currentLang(settings));

  if (!models.length) {
    return (
      <div className="flex h-full items-center justify-center p-4 text-xs text-muted">
        {tr.vizEmpty}
      </div>
    );
  }

  const sorted = [...models].sort((a, b) => {
    if (highlightedVar) {
      if (a.name === highlightedVar) return -1;
      if (b.name === highlightedVar) return 1;
    }
    return 0;
  });

  return (
    <div className="p-2">
      {sorted.map((m) => {
        const dimmed = highlightedVar && m.name !== highlightedVar;
        const wrap = (child: React.ReactNode) => (
          <div key={m.name} className={dimmed ? 'opacity-40' : undefined}>
            {child}
          </div>
        );
        if (m.type === 'array') return wrap(<ArrayViz model={m} />);
        if (m.type === 'matrix') return wrap(<MatrixViz model={m} />);
        if (m.type === 'map') return wrap(<MapViz model={m} />);
        if (m.type === 'linked_list' || m.type === 'tree') {
          return wrap(
            <Suspense fallback={<div className="p-2 text-xs text-muted">…</div>}>
              <GraphViz model={m} />
            </Suspense>,
          );
        }
        return null;
      })}
    </div>
  );
}
