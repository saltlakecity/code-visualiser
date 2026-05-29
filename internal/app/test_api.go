package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"visualisator/internal/models"
	"visualisator/internal/paths"
	"visualisator/internal/testrunner"
	"visualisator/internal/workspace"
)

func customTestsPath() string {
	return filepath.Join(paths.DataDir(), "custom", "tests.json")
}

// LoadCustomTests reads user-defined test cases.
func (a *App) LoadCustomTests() ([]models.TestCase, error) {
	path := customTestsPath()
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var tests []models.TestCase
	if err := json.Unmarshal(b, &tests); err != nil {
		return nil, err
	}
	return tests, nil
}

// SaveCustomTests persists user-defined test cases.
func (a *App) SaveCustomTests(tests []models.TestCase) error {
	dir := filepath.Dir(customTestsPath())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(tests, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(customTestsPath(), b, 0o644)
}

// BuildRunnerFromTest generates main.go source for a test case.
func (a *App) BuildRunnerFromTest(entryFn string, testCase models.TestCase, source string) (string, error) {
	if workspace.HasMainPackage(source) {
		return "", nil
	}
	if entryFn == "" {
		entryFn = workspace.FirstFuncName(source)
	}
	return testrunner.BuildMain(entryFn, testCase, false)
}

// RunAllTests compiles and runs each test without debug trace.
func (a *App) RunAllTests(source string, mode string, entryFn string, tests []models.TestCase) ([]models.TestResult, error) {
	sm, err := a.ensureSessions()
	if err != nil {
		return nil, err
	}
	results := make([]models.TestResult, 0, len(tests))
	for i, tc := range tests {
		runner, err := testrunner.BuildMain(entryFn, tc, workspace.HasMainPackage(source))
		if err != nil {
			results = append(results, models.TestResult{
				Index:    i,
				Passed:   false,
				Expected: tc.Expected,
				Error:    err.Error(),
			})
			continue
		}
		res, err := sm.CompileAndRun(source, mode, runner)
		if err != nil || !res.Success {
			msg := compileErrMsg(err, res.Errors)
			results = append(results, models.TestResult{
				Index:    i,
				Passed:   false,
				Expected: tc.Expected,
				Error:    msg,
			})
			if res.Session != "" {
				sm.CleanupSession(res.Session)
			}
			continue
		}
		tr := testrunner.RunTest(res.Binary, tc.Expected)
		tr.Index = i
		results = append(results, tr)
		sm.CleanupSession(res.Session)
	}
	return results, nil
}

func compileErrMsg(err error, compileErrs []models.CompileError) string {
	if err != nil {
		return err.Error()
	}
	if len(compileErrs) > 0 {
		return compileErrs[0].Message
	}
	return "compile failed"
}

// ResolveTestArgs converts test input to Go args string.
func (a *App) ResolveTestArgs(testCase models.TestCase) (string, error) {
	return testrunner.ResolveArgs(testCase)
}

// NormalizeProblemTests fills missing Args from Input for embedded problems.
func NormalizeProblemTests(p *models.Problem) {
	for i := range p.Tests {
		if p.Tests[i].Args == "" && p.Tests[i].Input != "" {
			if args, err := testrunner.ResolveArgs(p.Tests[i]); err == nil {
				p.Tests[i].Args = args
			}
		}
		if p.Tests[i].Name == "" {
			p.Tests[i].Name = fmt.Sprintf("Case %d", i+1)
		}
	}
}
