import { useEffect, useMemo } from 'react';
import ReactFlow, {
  Background,
  Controls,
  ReactFlowProvider,
  useEdgesState,
  useNodesState,
  useReactFlow,
} from 'reactflow';
import 'reactflow/dist/style.css';
import type { VizModel } from '../../types';
import { FlowLegend } from './FlowLegend';
import { flowNodeTypes, toFlowEdges, toFlowNodes } from './FlowNodes';

function StepFlowCanvas({ model }: { model: VizModel }) {
  const initialNodes = useMemo(
    () =>
      toFlowNodes(
        (model.nodes ?? []).map((n) => ({
          id: n.id,
          label: n.label,
          kind: n.kind,
          line: n.line,
          x: n.x,
          y: n.y,
          active: n.active,
        })),
      ),
    [model],
  );
  const initialEdges = useMemo(() => toFlowEdges(model.edges ?? []), [model]);

  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
  const { fitView } = useReactFlow();

  useEffect(() => {
    setNodes(initialNodes);
    setEdges(initialEdges);
    requestAnimationFrame(() => {
      fitView({ padding: 0.3, maxZoom: 1.2 });
    });
  }, [initialNodes, initialEdges, setNodes, setEdges, fitView]);

  if (!nodes.length) return null;

  return (
    <ReactFlow
      nodes={nodes}
      edges={edges}
      nodeTypes={flowNodeTypes}
      onNodesChange={onNodesChange}
      onEdgesChange={onEdgesChange}
      minZoom={0.2}
      maxZoom={1.5}
      proOptions={{ hideAttribution: true }}
      nodesDraggable
      nodesConnectable={false}
      elementsSelectable
      panOnDrag={[1, 2]}
    >
      <Background color="#2a2a35" gap={16} />
      <Controls showInteractive={false} />
    </ReactFlow>
  );
}

export function StepFlowViz({ model }: { model: VizModel }) {
  return (
    <div className="flex h-full min-h-[280px] w-full flex-col p-2">
      <h3 className="mb-1 font-mono text-xs text-accent">{model.name}</h3>
      {model.loopIter ? (
        <div className="mb-1 text-[10px] text-muted">Loop #{model.loopIter}</div>
      ) : null}
      <div className="flex min-h-[240px] flex-1 flex-col overflow-hidden rounded border border-border">
        <div className="min-h-0 flex-1">
          <ReactFlowProvider>
            <StepFlowCanvas model={model} />
          </ReactFlowProvider>
        </div>
        <FlowLegend showStepHints />
      </div>
    </div>
  );
}
