package viz

import "visualisator/internal/models"

const (
	layoutGapY = 120.0
	layoutGapX = 240.0
)

// ApplyFlowLayout assigns X/Y positions using rank + branch column layout.
func ApplyFlowLayout(nodes *[]models.VizNode, edges []models.VizEdge, startID, endID string) {
	if len(*nodes) == 0 {
		return
	}

	idToIdx := make(map[string]int)
	for i, n := range *nodes {
		idToIdx[n.ID] = i
	}

	outEdges := make(map[string][]models.VizEdge)
	inEdges := make(map[string][]models.VizEdge)
	for _, e := range edges {
		if e.Label == "next" {
			continue
		}
		outEdges[e.From] = append(outEdges[e.From], e)
		inEdges[e.To] = append(inEdges[e.To], e)
	}

	rank := make(map[string]int)
	rank[startID] = 0
	changed := true
	for changed {
		changed = false
		for _, e := range edges {
			if e.Label == "next" {
				continue
			}
			if r, ok := rank[e.From]; ok {
				nr := r + 1
				if cur, exists := rank[e.To]; !exists || nr > cur {
					rank[e.To] = nr
					changed = true
				}
			}
		}
	}

	maxRank := 0
	for _, r := range rank {
		if r > maxRank {
			maxRank = r
		}
	}
	if endID != "" {
		rank[endID] = maxRank + 1
	}

	col := make(map[string]float64)
	col[startID] = 0
	assignBranchColumns(startID, outEdges, inEdges, col, rank)

	layerNodes := make(map[int][]string)
	for id, r := range rank {
		layerNodes[r] = append(layerNodes[r], id)
	}

	for r := 0; r <= maxRank+1; r++ {
		ids := layerNodes[r]
		if len(ids) == 0 {
			continue
		}
		for i, id := range ids {
			idx, ok := idToIdx[id]
			if !ok {
				continue
			}
			c := col[id]
			if c == 0 && len(ids) > 1 {
				c = (float64(i) - float64(len(ids)-1)/2) * 0.5
			}
			(*nodes)[idx].Y = float64(r) * layoutGapY
			(*nodes)[idx].X = c * layoutGapX
		}
	}
}

func assignBranchColumns(
	nodeID string,
	outEdges map[string][]models.VizEdge,
	inEdges map[string][]models.VizEdge,
	col map[string]float64,
	rank map[string]int,
) {
	edges := outEdges[nodeID]
	if len(edges) == 0 {
		return
	}

	parentCol := col[nodeID]
	isDecision := false
	for _, e := range edges {
		if e.Label == "yes" || e.Label == "no" {
			isDecision = true
			break
		}
	}

	if isDecision {
		for _, e := range edges {
			switch e.Label {
			case "yes":
				col[e.To] = parentCol - 1
			case "no":
				col[e.To] = parentCol + 1
			default:
				if _, ok := col[e.To]; !ok {
					col[e.To] = parentCol
				}
			}
			assignBranchColumns(e.To, outEdges, inEdges, col, rank)
		}
		return
	}

	for i, e := range edges {
		if _, ok := col[e.To]; !ok {
			offset := float64(i) - float64(len(edges)-1)/2
			col[e.To] = parentCol + offset*0.3
		}
		assignBranchColumns(e.To, outEdges, inEdges, col, rank)
	}
}
