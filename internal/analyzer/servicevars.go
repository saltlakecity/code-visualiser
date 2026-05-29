package analyzer

import (
	"strings"

	"visualisator/internal/models"
)

// IsServiceVar reports debugger/compiler-internal names that should be hidden in the UI.
func IsServiceVar(name string) bool {
	if name == "" {
		return true
	}
	if strings.HasPrefix(name, "~") {
		return true
	}
	if name == "(return)" {
		return true
	}
	return false
}

// FilterVariables removes service variables from a Delve variable list.
func FilterVariables(vars []models.Variable) []models.Variable {
	if len(vars) == 0 {
		return nil
	}
	out := make([]models.Variable, 0, len(vars))
	for _, v := range vars {
		if IsServiceVar(v.Name) {
			continue
		}
		if len(v.Children) > 0 {
			v.Children = FilterVariables(v.Children)
		}
		out = append(out, v)
	}
	return out
}

// FilterVarDeltas removes service variable rows from change lists.
func FilterVarDeltas(deltas []VarDelta) []VarDelta {
	if len(deltas) == 0 {
		return nil
	}
	out := make([]VarDelta, 0, len(deltas))
	for _, d := range deltas {
		if IsServiceVar(d.Name) {
			continue
		}
		out = append(out, d)
	}
	return out
}
