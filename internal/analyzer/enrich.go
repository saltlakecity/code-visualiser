package analyzer

import (
	"visualisator/internal/models"
)

// ExplainConfig holds per-session data for frame narration.
type ExplainConfig struct {
	SourceMap *SourceMap
	Narrate   *NarrateContext
	Lang      Lang
}

// EnrichFrame attaches explanation to a trace frame using previous frame variables.
// rawCurr may include service vars (e.g. ~r0) for return-value narration only; it is not shown in Changes.
func EnrichFrame(cfg *ExplainConfig, tf *models.TraceFrame, prevVars, rawCurr []models.Variable) {
	if cfg == nil || tf == nil {
		return
	}
	line := tf.EditorLine
	if line <= 0 {
		line = tf.Line
	}
	if rawCurr == nil {
		rawCurr = tf.Variables
	}
	exp := NarrateStep(cfg.Lang, cfg.SourceMap, cfg.Narrate, line, tf.Function, prevVars, tf.Variables, rawCurr, tf.Index)
	if exp == nil {
		return
	}
	tf.Explanation = &models.FrameExplanation{
		Summary:   exp.Summary,
		LineText:  exp.LineText,
		StmtKind:  exp.StmtKind,
		Condition: exp.Condition,
		LoopIter:  exp.LoopIter,
		Changes:   toModelDeltas(FilterVarDeltas(exp.Changes)),
	}
}

func toModelDeltas(d []VarDelta) []models.VarDelta {
	out := make([]models.VarDelta, len(d))
	for i, v := range d {
		out[i] = models.VarDelta{
			Name:     v.Name,
			Type:     v.Type,
			OldValue: v.OldValue,
			NewValue: v.NewValue,
		}
	}
	return out
}

// ToSourceOutline converts analyzer outline to API model.
func ToSourceOutline(o FunctionOutline) models.SourceOutline {
	out := models.SourceOutline{
		Name:     o.Name,
		Params:   o.Params,
		Returns:  o.Returns,
		DemoCall: o.DemoCall,
	}
	for _, p := range o.Loops {
		out.Loops = append(out.Loops, models.OutlinePoint{Line: p.Line, Text: p.Text})
	}
	for _, p := range o.Branches {
		out.Branches = append(out.Branches, models.OutlinePoint{Line: p.Line, Text: p.Text})
	}
	for _, p := range o.ReturnsAt {
		out.ReturnsAt = append(out.ReturnsAt, models.OutlinePoint{Line: p.Line, Text: p.Text})
	}
	return out
}
