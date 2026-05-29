package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"visualisator/internal/paths"
)

// AppSettings persisted user preferences.
type AppSettings struct {
	Language      string  `json:"language"`
	Theme         string  `json:"theme"`
	UIStyle       string  `json:"uiStyle"`
	AccentColor   string  `json:"accentColor"`
	CustomCodeDir string  `json:"customCodeDir"`
	DataDir       string  `json:"dataDir"`
	PlaybackSpeed float64 `json:"playbackSpeed"`
}

// DefaultSettings returns factory defaults (portable: next to exe).
func DefaultSettings() AppSettings {
	dataDir := paths.DataDir()
	return AppSettings{
		Language:      "ru",
		Theme:         "dark",
		UIStyle:       "default",
		AccentColor:   "#ffa116",
		CustomCodeDir: paths.DefaultCustomCodeDir(),
		DataDir:       dataDir,
		PlaybackSpeed: 1,
	}
}

// Store loads and saves settings JSON.
type Store struct {
	mu   sync.RWMutex
	path string
	data AppSettings
}

// NewStore creates settings store beside the executable.
func NewStore() (*Store, error) {
	s := &Store{
		path: paths.SettingsPath(),
		data: DefaultSettings(),
	}
	_ = s.Load()
	return s, nil
}

// Load reads settings from disk.
func (s *Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return s.saveLocked()
		}
		return err
	}
	var data AppSettings
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	if data.Language == "" {
		data.Language = "ru"
	}
	if data.Theme == "" {
		data.Theme = "dark"
	}
	if data.UIStyle == "" {
		data.UIStyle = "default"
	}
	if data.AccentColor == "" {
		data.AccentColor = "#ffa116"
	}
	if data.DataDir == "" {
		data.DataDir = paths.DataDir()
	}
	if data.CustomCodeDir == "" {
		data.CustomCodeDir = filepath.Join(data.DataDir, "custom")
	}
	if data.PlaybackSpeed <= 0 {
		data.PlaybackSpeed = 1
	}
	_ = os.MkdirAll(data.DataDir, 0o755)
	_ = os.MkdirAll(data.CustomCodeDir, 0o755)
	s.data = data
	return nil
}

// Get returns current settings.
func (s *Store) Get() AppSettings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data
}

// Save persists settings.
func (s *Store) Save(data AppSettings) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if data.DataDir == "" {
		data.DataDir = paths.DataDir()
	}
	if data.CustomCodeDir == "" {
		data.CustomCodeDir = filepath.Join(data.DataDir, "custom")
	}
	s.data = data
	_ = os.MkdirAll(s.data.DataDir, 0o755)
	_ = os.MkdirAll(s.data.CustomCodeDir, 0o755)
	return s.saveLocked()
}

func (s *Store) saveLocked() error {
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o644)
}

// CustomCodePath returns path for default custom code file.
func (s *Store) CustomCodePath() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return filepath.Join(s.data.CustomCodeDir, "custom.go")
}
