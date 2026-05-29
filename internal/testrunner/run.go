package testrunner

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"visualisator/internal/models"
)

func runBinary(binaryPath string) (string, error) {
	cmd := exec.Command(binaryPath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return stdout.String(), fmt.Errorf("%s", msg)
	}
	return stdout.String(), nil
}

// RunTest executes a compiled binary and compares output.
func RunTest(binaryPath string, expected string) models.TestResult {
	actual, err := runBinary(binaryPath)
	if err != nil {
		return models.TestResult{
			Passed:   false,
			Actual:   strings.TrimSpace(actual),
			Expected: expected,
			Error:    err.Error(),
		}
	}
	passed := CompareOutput(actual, expected)
	return models.TestResult{
		Passed:   passed,
		Actual:   strings.TrimSpace(actual),
		Expected: expected,
	}
}
