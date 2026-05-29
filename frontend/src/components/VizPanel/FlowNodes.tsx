import { memo } from 'react';
import { Handle, Position, type NodeProps } from 'reactflow';

function baseClass(kind: string, active?: boolean) {
  const k = kind || 'process';
  let cls = 'flow-node flow-node-' + k;
  if (active) cls += ' flow-node-active';
  return cls;
}

function FlowNodeInner({ data }: NodeProps<{ label: string; kind?: string; active?: boolean }>) {
  const kind = data.kind || 'process';
  const label = data.label || '';

  return (
    <div className={baseClass(kind, data.active)}>
      <Handle type="target" position={Position.Top} className="!bg-accent !w-2 !h-2" />
      <div className="flow-node-label">{label}</div>
      <Handle type="source" position={Position.Bottom} className="!bg-accent !w-2 !h-2" />
    </div>
  );
}

export const FlowNode = memo(FlowNodeInner);

export const flowNodeTypes = {
  flow: FlowNode,
};

export function toFlowNodes(
  nodes: { id: string; label: string; kind?: string; line?: number; x?: number; y?: number; active?: boolean }[],
) {
  return nodes.map((n) => ({
    id: n.id,
    type: 'flow',
    position: { x: n.x ?? 0, y: n.y ?? 0 },
    data: { label: n.label, kind: n.kind, active: n.active, line: n.line },
  }));
}

export function toFlowEdges(
  edges: { from: string; to: string; label?: string; animated?: boolean }[],
) {
  return edges.map((e, i) => ({
    id: `e-${i}`,
    source: e.from,
    target: e.to,
    label: e.label,
    animated: e.animated ?? false,
    style: { stroke: e.animated ? 'var(--color-accent)' : '#55556a' },
    labelStyle: { fill: '#a0a0b0', fontSize: 10 },
  }));
}
