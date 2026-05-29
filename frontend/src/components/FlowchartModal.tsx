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
import { t } from '../i18n';
import { useSettingsStore, currentLang } from '../store/settingsStore';
import type { FlowchartResult } from '../types';
import { FlowLegend } from './VizPanel/FlowLegend';
import { flowNodeTypes, toFlowEdges, toFlowNodes } from './VizPanel/FlowNodes';

interface Props {
  open: boolean;
  data: FlowchartResult | null;
  loading: boolean;
  onClose: () => void;
  onLineClick?: (line: number) => void;
}

function FlowchartCanvas({
  data,
  onLineClick,
}: {
  data: FlowchartResult;
  onLineClick?: (line: number) => void;
}) {
  const initialNodes = useMemo(() => toFlowNodes(data.nodes ?? []), [data]);
  const initialEdges = useMemo(() => toFlowEdges(data.edges ?? []), [data]);

  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
  const { fitView } = useReactFlow();

  useEffect(() => {
    setNodes(initialNodes);
    setEdges(initialEdges);
    requestAnimationFrame(() => {
      fitView({ padding: 0.35, maxZoom: 1.2 });
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
      minZoom={0.15}
      maxZoom={1.5}
      proOptions={{ hideAttribution: true }}
      nodesDraggable
      nodesConnectable={false}
      elementsSelectable
      panOnDrag={[1, 2]}
      onNodeClick={(_, node) => {
        const line = (node.data as { line?: number }).line;
        if (line && line > 0) onLineClick?.(line);
      }}
    >
      <Background color="#2a2a35" gap={16} />
      <Controls showInteractive={false} />
    </ReactFlow>
  );
}

export function FlowchartModal({ open, data, loading, onClose, onLineClick }: Props) {
  const { settings } = useSettingsStore();
  const tr = t(currentLang(settings));

  if (!open) return null;

  const hasNodes = (data?.nodes?.length ?? 0) > 0;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
      <div className="flex max-h-[90vh] w-full max-w-5xl flex-col overflow-hidden rounded-lg border border-border bg-panel shadow-xl">
        <div className="flex shrink-0 items-center justify-between border-b border-border px-4 py-3">
          <div>
            <h2 className="text-sm font-semibold text-accent">
              {tr.flowchartTitle}
              {data?.functionName ? `: ${data.functionName}` : ''}
            </h2>
            {data?.summary && <p className="mt-0.5 font-mono text-xs text-muted">{data.summary}</p>}
          </div>
          <button
            type="button"
            onClick={onClose}
            className="rounded border border-border px-3 py-1 text-xs hover:bg-bg"
          >
            {tr.flowchartClose}
          </button>
        </div>
        <div className="flex min-h-0 flex-1 flex-col">
          <div className="h-[60vh] w-full shrink-0">
            {loading ? (
              <div className="flex h-full items-center justify-center text-xs text-muted">{tr.flowchartLoading}</div>
            ) : hasNodes && data ? (
              <ReactFlowProvider>
                <FlowchartCanvas data={data} onLineClick={onLineClick} />
              </ReactFlowProvider>
            ) : (
              <div className="flex h-full items-center justify-center text-xs text-muted">{tr.flowchartEmpty}</div>
            )}
          </div>
          {hasNodes && !loading ? <FlowLegend /> : null}
        </div>
      </div>
    </div>
  );
}
