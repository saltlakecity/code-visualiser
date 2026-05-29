import { lazy, Suspense, useEffect, useState } from 'react';
import { t } from '../../i18n';
import { useDebuggerStore } from '../../store/debuggerStore';
import { useSettingsStore, currentLang } from '../../store/settingsStore';
import type { VizModel } from '../../types';
import { ArrayViz } from './ArrayViz';
import { MapViz, MatrixViz } from './MatrixViz';
import { StackViz } from './StackViz';
import { StepFlowViz } from './FlowViz';
import { StringViz } from './StringViz';

const GraphViz = lazy(() => import('./GraphViz').then((m) => ({ default: m.GraphViz })));

type VizTab = 'data' | 'execution';

function isDataModel(m: VizModel) {
  return m.type !== 'step_flow' && m.type !== 'flowchart' && m.type !== 'scalars';
}

function renderDataModel(m: VizModel, dimmed: boolean) {
  const wrap = (child: React.ReactNode) => (
    <div key={m.name + m.type} className={dimmed ? 'opacity-40' : undefined}>
      {child}
    </div>
  );
  if (m.type === 'array') return wrap(<ArrayViz model={m} />);
  if (m.type === 'string') return wrap(<StringViz model={m} />);
  if (m.type === 'stack') return wrap(<StackViz model={m} />);
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
}

export function VisualizationPanel({
  models,
  hasTrace,
}: {
  models: VizModel[];
  hasTrace: boolean;
}) {
  const { settings } = useSettingsStore();
  const { highlightedVar } = useDebuggerStore();
  const tr = t(currentLang(settings));
  const [tab, setTab] = useState<VizTab>('data');

  const dataModels = models.filter(isDataModel);
  const hasData = dataModels.length > 0;
  const stepModel = models.find((m) => m.type === 'step_flow');

  useEffect(() => {
    if (!hasData && stepModel && hasTrace) {
      setTab('execution');
    } else if (hasData) {
      setTab('data');
    }
  }, [hasData, stepModel, hasTrace]);

  const sortedData = [...dataModels].sort((a, b) => {
    if (highlightedVar) {
      if (a.name === highlightedVar) return -1;
      if (b.name === highlightedVar) return 1;
    }
    return 0;
  });

  return (
    <div className="flex h-full min-h-0 flex-col">
      <div className="flex shrink-0 border-b border-border">
        <button
          type="button"
          onClick={() => setTab('data')}
          className={`flex-1 px-2 py-1.5 text-[10px] font-medium uppercase tracking-wide ${
            tab === 'data' ? 'border-b-2 border-accent text-accent' : 'text-muted'
          }`}
        >
          {tr.vizTabData}
        </button>
        <button
          type="button"
          onClick={() => setTab('execution')}
          className={`flex-1 px-2 py-1.5 text-[10px] font-medium uppercase tracking-wide ${
            tab === 'execution' ? 'border-b-2 border-accent text-accent' : 'text-muted'
          }`}
        >
          {tr.vizTabExecution}
        </button>
      </div>

      <div className="min-h-0 flex-1 overflow-y-auto">
        {tab === 'data' ? (
          hasData ? (
            <div className="p-1">{sortedData.map((m) => renderDataModel(m, !!(highlightedVar && m.name !== highlightedVar)))}</div>
          ) : (
            <div className="flex h-full items-center justify-center p-4 text-center text-xs text-muted">
              {hasTrace ? tr.vizEmptyData : tr.vizEmpty}
            </div>
          )
        ) : stepModel ? (
          <StepFlowViz model={stepModel} />
        ) : (
          <div className="flex h-full items-center justify-center p-4 text-center text-xs text-muted">
            {hasTrace ? tr.vizEmptyExecution : tr.vizEmpty}
          </div>
        )}
      </div>
    </div>
  );
}

export function VizPanel({ models }: { models: VizModel[] }) {
  return <VisualizationPanel models={models} hasTrace={models.length > 0} />;
}
