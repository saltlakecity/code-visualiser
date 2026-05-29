package viz

import (
	"strings"

	"visualisator/internal/analyzer"
	"visualisator/internal/models"
)

// BuildScalars creates a scalars viz model from user variables.
func BuildScalars(vars []models.Variable) *models.VizModel {
	var scalars []models.VizScalar
	for _, v := range vars {
		if analyzer.IsServiceVar(v.Name) {
			continue
		}
		typeLower := strings.ToLower(strings.TrimSpace(v.Type))
		if isComplexType(typeLower) {
			continue
		}
		val := strings.TrimSpace(v.Value)
		if val == "" {
			continue
		}
		scalars = append(scalars, models.VizScalar{
			Name:  v.Name,
			Type:  v.Type,
			Value: val,
		})
	}
	if len(scalars) == 0 {
		return nil
	}
	return &models.VizModel{
		Type:    "scalars",
		Name:    "variables",
		Scalars: scalars,
	}
}

func isComplexType(typeLower string) bool {
	if strings.HasPrefix(typeLower, "[]") ||
		strings.HasPrefix(typeLower, "map[") ||
		strings.HasPrefix(typeLower, "[") ||
		strings.HasPrefix(typeLower, "*") ||
		strings.HasPrefix(typeLower, "struct") ||
		strings.Contains(typeLower, "interface") {
		return true
	}
	return false
}

func HasDataViz(models []models.VizModel) bool {
	for _, m := range models {
		switch m.Type {
		case "array", "matrix", "map", "linked_list", "tree", "string", "stack":
			return true
		}
	}
	return false
}
