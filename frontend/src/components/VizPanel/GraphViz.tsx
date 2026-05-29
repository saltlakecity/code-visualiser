import { useMemo } from 'react';
import ReactFlow, { Background, Controls, type Edge, type Node } from 'reactflow';
import 'reactflow/dist/style.css';
import type { VizModel } from '../../types';

export function GraphViz({ model }: { model: VizModel }) {
  const { nodes, edges } = useMemo(() => {
    const nodes: Node[] = (model.nodes ?? []).map((n) => ({
      id: n.id,
      position: { x: n.x ?? 0, y: n.y ?? 0 },
      data: { label: n.label },
      style: {
        background: '#16161d',
        border: '1px solid var(--color-accent)',
        color: '#e8e8ef',
        borderRadius: 8,
        padding: 8,
        fontSize: 12,
        fontFamily: 'JetBrains Mono, monospace',
        minWidth: 48,
        textAlign: 'center',
      },
    }));
    const edges: Edge[] = (model.edges ?? []).map((e, i) => ({
      id: `e-${i}`,
      source: e.from,
      target: e.to,
      animated: true,
      style: { stroke: 'var(--color-accent)' },
    }));
    return { nodes, edges };
  }, [model]);

  if (!nodes.length) return null;

  return (
    <div className="h-56 w-full p-2">
      <h3 className="mb-1 font-mono text-xs text-accent">{model.name}</h3>
      <div className="h-48 rounded border border-border">
        <ReactFlow nodes={nodes} edges={edges} fitView proOptions={{ hideAttribution: true }}>
          <Background color="#2a2a35" gap={16} />
          <Controls showInteractive={false} />
        </ReactFlow>
      </div>
    </div>
  );
}
