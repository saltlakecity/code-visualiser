package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"visualisator/internal/models"
	"visualisator/internal/paths"
	"visualisator/internal/settings"
)

func toModel(s settings.AppSettings) models.AppSettings {
	return models.AppSettings{
		Language:      s.Language,
		Theme:         s.Theme,
		AccentColor:   s.AccentColor,
		CustomCodeDir: s.CustomCodeDir,
		DataDir:       s.DataDir,
		PlaybackSpeed: s.PlaybackSpeed,
	}
}

func fromModel(m models.AppSettings) settings.AppSettings {
	return settings.AppSettings{
		Language:      m.Language,
		Theme:         m.Theme,
		AccentColor:   m.AccentColor,
		CustomCodeDir: m.CustomCodeDir,
		DataDir:       m.DataDir,
		PlaybackSpeed: m.PlaybackSpeed,
	}
}

// GetDataDir returns the application data directory path.
func (a *App) GetDataDir() string {
	return paths.DataDir()
}

// GetSettings returns persisted app settings.
func (a *App) GetSettings() models.AppSettings {
	if a.settings == nil {
		return toModel(settings.DefaultSettings())
	}
	return toModel(a.settings.Get())
}

// SaveSettings persists app settings.
func (a *App) SaveSettings(s models.AppSettings) error {
	if a.settings == nil {
		return fmt.Errorf("settings not available")
	}
	if s.CustomCodeDir != "" {
		if err := os.MkdirAll(s.CustomCodeDir, 0o755); err != nil {
			return err
		}
	}
	return a.settings.Save(fromModel(s))
}

// SaveCustomCode writes custom code to the configured directory.
func (a *App) SaveCustomCode(source string, filename string) error {
	if a.settings == nil {
		return fmt.Errorf("settings not available")
	}
	dir := a.settings.Get().CustomCodeDir
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	name := filename
	if name == "" {
		name = "custom.go"
	}
	if !strings.HasSuffix(name, ".go") {
		name += ".go"
	}
	path := filepath.Join(dir, filepath.Base(name))
	return os.WriteFile(path, []byte(source), 0o644)
}

// LoadCustomCode reads custom code from default or named file.
func (a *App) LoadCustomCode(filename string) (string, error) {
	if a.settings == nil {
		return "", fmt.Errorf("settings not available")
	}
	dir := a.settings.Get().CustomCodeDir
	var path string
	if filename == "" {
		path = a.settings.CustomCodePath()
	} else {
		path = filepath.Join(dir, filepath.Base(filename))
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(b), nil
}

// ListCustomFiles lists .go files in custom code directory.
func (a *App) ListCustomFiles() ([]string, error) {
	if a.settings == nil {
		return nil, fmt.Errorf("settings not available")
	}
	dir := a.settings.Get().CustomCodeDir
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".go") {
			files = append(files, e.Name())
		}
	}
	return files, nil
}

// PickCustomCodeDir opens native folder picker (returns chosen path).
func (a *App) PickCustomCodeDir() (string, error) {
	dir, err := PickDirectoryDialog(a.ctx)
	if err != nil {
		return "", err
	}
	return dir, nil
}
