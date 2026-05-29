package paths

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	exeDirOnce sync.Once
	exeDir     string
	dataDirOnce sync.Once
	dataDir    string
)

// ExeDir returns the directory containing the running executable.
func ExeDir() string {
	exeDirOnce.Do(func() {
		exe, err := os.Executable()
		if err != nil {
			exeDir = "."
			return
		}
		// Resolve symlinks (Windows installers, Go builds).
		if resolved, err := filepath.EvalSymlinks(exe); err == nil {
			exe = resolved
		}
		exeDir = filepath.Dir(exe)
	})
	return exeDir
}

// DataDir is the writable application data root.
func DataDir() string {
	dataDirOnce.Do(func() {
		candidates := []string{
			filepath.Join(ExeDir(), "data"),
		}
		if cfg, err := os.UserConfigDir(); err == nil {
			candidates = append(candidates, filepath.Join(cfg, "visualisator"))
		}
		candidates = append(candidates, filepath.Join(os.TempDir(), "visualisator-data"))

		for _, dir := range candidates {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				continue
			}
			probe := filepath.Join(dir, ".write_test")
			if err := os.WriteFile(probe, []byte("ok"), 0o644); err != nil {
				continue
			}
			_ = os.Remove(probe)
			dataDir = dir
			return
		}
		dataDir = filepath.Join(ExeDir(), "data")
		_ = os.MkdirAll(dataDir, 0o755)
	})
	return dataDir
}

// SettingsPath is settings.json inside the data directory.
func SettingsPath() string {
	return filepath.Join(DataDir(), "settings.json")
}

// DefaultCustomCodeDir is where custom Go snippets are stored by default.
func DefaultCustomCodeDir() string {
	dir := filepath.Join(DataDir(), "custom")
	_ = os.MkdirAll(dir, 0o755)
	return dir
}

// TempBuildDir is for compiled debug sessions.
func TempBuildDir() string {
	dir := filepath.Join(DataDir(), "sessions")
	_ = os.MkdirAll(dir, 0o755)
	return dir
}

// LogPath is startup/diagnostic log file path.
func LogPath() string {
	return filepath.Join(DataDir(), "app.log")
}
