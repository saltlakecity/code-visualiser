package models

// AppSettings are user preferences exposed to the UI.
type AppSettings struct {
	Language      string  `json:"language"`
	Theme         string  `json:"theme"`
	UIStyle       string  `json:"uiStyle"`
	AccentColor   string  `json:"accentColor"`
	CustomCodeDir string  `json:"customCodeDir"`
	DataDir       string  `json:"dataDir"`
	PlaybackSpeed float64 `json:"playbackSpeed"`
}
