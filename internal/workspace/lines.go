package workspace

import "strings"

// EditorLineOffset returns how many lines were prepended before the user's editor text.
func EditorLineOffset(editorSource, preparedSource string) int {
	editor := strings.TrimSpace(editorSource)
	prepared := strings.TrimSpace(preparedSource)
	if editor == "" || prepared == "" || editor == prepared {
		return 0
	}

	if idx := strings.Index(prepared, editor); idx >= 0 {
		return strings.Count(prepared[:idx], "\n")
	}

	// Match first non-empty line of editor in prepared file
	editorLines := strings.Split(editor, "\n")
	first := ""
	for _, l := range editorLines {
		first = strings.TrimSpace(l)
		if first != "" {
			break
		}
	}
	if first == "" {
		return 0
	}
	preparedLines := strings.Split(prepared, "\n")
	for i, l := range preparedLines {
		if strings.TrimSpace(l) == first {
			return i
		}
	}
	return 0
}
