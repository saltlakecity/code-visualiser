package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// AppSettings persisted user preferences.
type AppSettings struct {
	Language      string `json:"language"`
	Theme         string `json:"theme"`
	AccentColor   string `json:"accentColor"`
	CustomCodeDir string `json:"customCodeDir"`
}

// DefaultSettings returns factory defaults.
func DefaultSettings() AppSettings {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, "Documents", "visualisator")
	return AppSettings{
		Language:      "ru",
		Theme:         "dark",
		AccentColor:   "#ffa116",
		CustomCodeDir: dir,
	}
}

// Store loads and saves settings JSON.
type Store struct {
	mu   sync.RWMutex
	path string
	data AppSettings
}

// NewStore creates settings store at config path.
func NewStore() (*Store, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir, _ = os.UserHomeDir()
		dir = filepath.Join(dir, ".visualisator")
	} else {
		dir = filepath.Join(dir, "visualisator")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	s := &Store{
		path: filepath.Join(dir, "settings.json"),
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
	if data.AccentColor == "" {
		data.AccentColor = "#ffa116"
	}
	if data.CustomCodeDir == "" {
		data.CustomCodeDir = DefaultSettings().CustomCodeDir
	}
	s.data = data
	_ = os.MkdirAll(s.data.CustomCodeDir, 0o755)
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
	s.data = data
	if s.data.CustomCodeDir != "" {
		_ = os.MkdirAll(s.data.CustomCodeDir, 0o755)
	}
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
