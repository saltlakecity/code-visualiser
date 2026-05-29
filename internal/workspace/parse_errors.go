package workspace

import (
	"regexp"
	"strconv"
	"strings"

	"visualisator/internal/models"
)

var compileErrRe = regexp.MustCompile(`(?m)^(.+\.go):(\d+):(\d+):\s*(.+)$`)

// ParseCompileErrors extracts line/column from go build output.
func ParseCompileErrors(output string) []models.CompileError {
	var errs []models.CompileError
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		m := compileErrRe.FindStringSubmatch(strings.TrimSpace(line))
		if len(m) != 5 {
			continue
		}
		ln, _ := strconv.Atoi(m[2])
		col, _ := strconv.Atoi(m[3])
		file := m[1]
		if idx := strings.LastIndex(file, "/"); idx >= 0 {
			file = file[idx+1:]
		}
		if idx := strings.LastIndex(file, "\\"); idx >= 0 {
			file = file[idx+1:]
		}
		errs = append(errs, models.CompileError{
			File:    file,
			Line:    ln,
			Column:  col,
			Message: m[4],
		})
	}
	if len(errs) == 0 && strings.TrimSpace(output) != "" {
		errs = append(errs, models.CompileError{
			File:    "main.go",
			Line:    1,
			Column:  1,
			Message: strings.TrimSpace(output),
		})
	}
	return errs
}
