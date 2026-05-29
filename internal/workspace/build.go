package workspace

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"

	"visualisator/internal/models"
	"visualisator/internal/paths"
)

// BuildRoot returns directory for temporary compile sessions.
func BuildRoot() string {
	return paths.TempBuildDir()
}

const moduleName = "usercode"

// Builder compiles user Go sources in isolated temp dirs.
type Builder struct {
	baseDir string
}

// NewBuilder creates a workspace builder under baseDir.
func NewBuilder(baseDir string) (*Builder, error) {
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, err
	}
	return &Builder{baseDir: baseDir}, nil
}

// SessionPaths holds paths for a compile session.
type SessionPaths struct {
	ID      string
	Dir     string
	Binary  string
	MainFile string
}

// Prepare creates a new session directory and go module.
func (b *Builder) Prepare() (*SessionPaths, error) {
	id := uuid.New().String()
	dir := filepath.Join(b.baseDir, id)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("go mod init: %s: %w", out, err)
	}
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	return &SessionPaths{
		ID:     id,
		Dir:    dir,
		Binary: filepath.Join(dir, "user"+ext),
	}, nil
}

// WriteSource writes user code files for the session.
func (b *Builder) WriteSource(paths *SessionPaths, source string, mode string) (string, error) {
	if err := ScanImports(source); err != nil {
		return "", err
	}

	switch mode {
	case "full_main":
		mainPath := filepath.Join(paths.Dir, "main.go")
		if err := os.WriteFile(mainPath, []byte(source), 0o644); err != nil {
			return "", err
		}
		paths.MainFile = "main.go"
		return mainPath, nil
	case "function_wrapper":
		solution := source
		if !strings.Contains(solution, "package ") {
			solution = "package main\n\n" + solution
		}
		solutionPath := filepath.Join(paths.Dir, "solution.go")
		if err := os.WriteFile(solutionPath, []byte(solution), 0o644); err != nil {
			return "", err
		}
		paths.MainFile = "solution.go"
		return solutionPath, nil
	default:
		return "", fmt.Errorf("unknown mode %q", mode)
	}
}

// WriteMain writes an explicit main.go for problem runner code.
func (b *Builder) WriteMain(paths *SessionPaths, mainSource string) error {
	if mainSource == "" {
		return nil
	}
	return os.WriteFile(filepath.Join(paths.Dir, "main.go"), []byte(mainSource), 0o644)
}

// Compile builds the session binary with debug symbols.
func (b *Builder) Compile(paths *SessionPaths) models.CompileResult {
	result := models.CompileResult{Session: paths.ID}
	cmd := exec.Command("go", "build",
		"-gcflags=all=-N -l",
		"-o", paths.Binary,
		".",
	)
	cmd.Dir = paths.Dir
	cmd.Env = append(os.Environ(),
		"CGO_ENABLED=0",
		"GOFLAGS=-trimpath",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		result.Success = false
		result.Errors = ParseCompileErrors(string(out))
		if len(result.Errors) == 0 {
			result.Errors = []models.CompileError{{
				File:    "main.go",
				Line:    1,
				Column:  1,
				Message: strings.TrimSpace(string(out)),
			}}
		}
		return result
	}
	result.Success = true
	result.Binary = paths.Binary
	return result
}

// Cleanup removes a session directory.
func (b *Builder) Cleanup(sessionID string) {
	dir := filepath.Join(b.baseDir, sessionID)
	_ = os.RemoveAll(dir)
}

// DetectMode returns compile mode from source.
func DetectMode(source string) string {
	if HasMainPackage(source) {
		return "full_main"
	}
	return "function_wrapper"
}
