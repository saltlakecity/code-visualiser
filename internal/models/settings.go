package models

// AppSettings are user preferences exposed to the UI.
type AppSettings struct {
	Language      string `json:"language"`
	Theme         string `json:"theme"`
	AccentColor   string `json:"accentColor"`
	CustomCodeDir string `json:"customCodeDir"`
}
