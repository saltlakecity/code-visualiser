package debugger

import (
	"regexp"
	"strings"
)

var exitStatusRe = regexp.MustCompile(`(?i)exited with status\s+(\d+)`)

// classifyExit determines whether a debugger exit represents a failure.
func classifyExit(err error) (status, msg, kind string, failed bool) {
	if err == nil {
		return "finished", "", "", false
	}
	errStr := strings.TrimSpace(err.Error())
	if errStr == "" {
		return "finished", "", "", false
	}
	lower := strings.ToLower(errStr)

	if strings.Contains(lower, "panic") {
		return "panic", errStr, "panic", true
	}

	if m := exitStatusRe.FindStringSubmatch(errStr); len(m) == 2 {
		if m[1] == "0" {
			return "finished", "", "", false
		}
		return "error", errStr, "runtime", true
	}

	// Benign detach/exit wording from Delve on normal program completion.
	if strings.Contains(lower, "has exited") {
		return "finished", "", "", false
	}

	return "error", errStr, "runtime", true
}
