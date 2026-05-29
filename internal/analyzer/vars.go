package analyzer

import (
	"strings"

	"visualisator/internal/models"
)

// ComputeVarDeltas compares flattened variable trees between two frames.
func ComputeVarDeltas(prev, curr []models.Variable) []VarDelta {
	prevMap := flatVars(prev, "")
	currMap := flatVars(curr, "")
	var out []VarDelta

	for name, cv := range currMap {
		if IsServiceVar(name) {
			continue
		}
		newVal := displayValue(cv)
		pv, ok := prevMap[name]
		if !ok {
			out = append(out, VarDelta{
				Name:     name,
				Type:     cv.Type,
				OldValue: "—",
				NewValue: newVal,
			})
			continue
		}
		oldVal := displayValue(pv)
		if oldVal != newVal || pv.Type != cv.Type {
			out = append(out, VarDelta{
				Name:     name,
				Type:     cv.Type,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
	}
	return out
}

func flatVars(vars []models.Variable, prefix string) map[string]models.Variable {
	m := make(map[string]models.Variable)
	for _, v := range vars {
		key := v.Name
		if prefix != "" {
			key = prefix + "." + v.Name
		}
		m[key] = v
		if len(v.Children) > 0 {
			for k, c := range flatVars(v.Children, key) {
				m[k] = c
			}
		}
	}
	return m
}

func displayValue(v models.Variable) string {
	if len(v.Children) > 0 {
		var items []string
		for i, c := range v.Children {
			if i >= 8 {
				break
			}
			val := strings.TrimSpace(c.Value)
			if val == "" {
				val = c.Name
			}
			items = append(items, val)
		}
		suffix := ""
		if len(v.Children) > 8 {
			suffix = ", …"
		}
		return "[" + strings.Join(items, ", ") + suffix + "]"
	}
	val := strings.TrimSpace(v.Value)
	if len(val) > 60 {
		return val[:57] + "…"
	}
	if val == "" {
		return "—"
	}
	return val
}

// VarLookup builds a flat name→display value map for expression evaluation.
func VarLookup(vars []models.Variable) map[string]string {
	m := make(map[string]string)
	for k, v := range flatVars(vars, "") {
		m[k] = displayValue(v)
	}
	return m
}
