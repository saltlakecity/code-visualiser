import { t } from '../../i18n';
import { useSettingsStore, currentLang } from '../../store/settingsStore';

type LegendItem = {
  kind: string;
  label: string;
  description: string;
  active?: boolean;
  edge?: boolean;
};

function LegendSwatch({ kind, active, edge }: { kind?: string; active?: boolean; edge?: boolean }) {
  if (edge) {
    return <span className="flow-legend-edge" aria-hidden />;
  }
  const k = kind || 'process';
  let cls = `flow-legend-swatch flow-node flow-node-${k}`;
  if (active) cls += ' flow-node-active';
  return <span className={cls} aria-hidden />;
}

export function FlowLegend({ showStepHints = false }: { showStepHints?: boolean }) {
  const { settings } = useSettingsStore();
  const tr = t(currentLang(settings));

  const items: LegendItem[] = [
    { kind: 'start', label: tr.flowLegendStart, description: tr.flowLegendStartDesc },
    { kind: 'end', label: tr.flowLegendEnd, description: tr.flowLegendEndDesc },
    { kind: 'process', label: tr.flowLegendProcess, description: tr.flowLegendProcessDesc },
    { kind: 'decision', label: tr.flowLegendDecision, description: tr.flowLegendDecisionDesc },
    { kind: 'loop', label: tr.flowLegendLoop, description: tr.flowLegendLoopDesc },
    { kind: 'return', label: tr.flowLegendReturn, description: tr.flowLegendReturnDesc },
  ];

  if (showStepHints) {
    items.push(
      { kind: 'process', label: tr.flowLegendActive, description: tr.flowLegendActiveDesc, active: true },
      { kind: 'process', label: tr.flowLegendEdge, description: tr.flowLegendEdgeDesc, edge: true },
    );
  }

  return (
    <div className="flow-legend">
      <div className="flow-legend-title">{tr.flowLegendTitle}</div>
      <ul className="flow-legend-list">
        {items.map((item) => (
          <li key={item.label + (item.active ? '-active' : '') + (item.edge ? '-edge' : '')} className="flow-legend-item">
            <LegendSwatch kind={item.kind} active={item.active} edge={item.edge} />
            <div className="flow-legend-text">
              <span className="flow-legend-label">{item.label}</span>
              <span className="flow-legend-desc">{item.description}</span>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
