package viz

import (
	"strings"

	"visualisator/internal/analyzer"
	"visualisator/internal/models"
)

// BuildStepViz creates a step_flow model highlighting the current execution position.
func BuildStepViz(
	source string,
	editorLine int,
	functionName string,
	loopIter int,
	condition string,
) *models.VizModel {
	if strings.TrimSpace(source) == "" || editorLine <= 0 {
		return nil
	}

	flow, err := BuildFlowGraph(source)
	if err != nil || len(flow.Nodes) == 0 {
		return nil
	}

	activeID := ""
	for i := range flow.Nodes {
		n := &flow.Nodes[i]
		if n.Line == editorLine {
			n.Active = true
			activeID = n.ID
		}
	}

	// Fallback: match by statement kind on nearest line
	if activeID == "" {
		sm, err := analyzer.ParseSource(source)
		if err == nil {
			if info, ok := sm.StmtAt(editorLine); ok {
				for i := range flow.Nodes {
					n := &flow.Nodes[i]
					if n.Line == info.Line || strings.Contains(strings.ToLower(n.Label), info.Kind) {
						n.Active = true
						activeID = n.ID
						break
					}
				}
			}
		}
	}

	// Highlight edge leading to active node
	if activeID != "" {
		for i := range flow.Edges {
			if flow.Edges[i].To == activeID {
				flow.Edges[i].Animated = true
			}
		}
	}

	name := functionName
	if name == "" {
		name = flow.FunctionName
	}

	return &models.VizModel{
		Type:     "step_flow",
		Name:     name,
		Nodes:    flow.Nodes,
		Edges:    flow.Edges,
		LoopIter: loopIter,
	}
}
