package app

import "visualisator/internal/paths"

// LogClient writes a message from the frontend into app.log.
func (a *App) LogClient(message string) {
	paths.AppendLog("client: %s", message)
}
