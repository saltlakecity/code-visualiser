package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"visualisator/internal/app"
	"visualisator/internal/paths"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	defer func() {
		if r := recover(); r != nil {
			paths.AppendLog("panic: %v\n%s", r, debug.Stack())
		}
	}()

	paths.AppendLog("main: start exe=%s", mustExe())

	application := app.NewApp()

	err := wails.Run(&options.App{
		Title:  "Go Code Visualizer",
		Width:  1280,
		Height: 800,
		MinWidth:  900,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 15, G: 15, B: 19, A: 255},
		OnStartup: application.Startup,
		OnDomReady: func(ctx context.Context) {
			paths.AppendLog("lifecycle: dom ready")
		},
		OnShutdown: func(ctx context.Context) {
			paths.AppendLog("lifecycle: shutdown")
		},
		Bind: []interface{}{
			application,
		},
	})

	if err != nil {
		paths.AppendLog("wails error: %v", err)
		fmt.Fprintf(os.Stderr, "visualisator: %v\nlog: %s\n", err, paths.LogPath())
		os.Exit(1)
	}
	paths.AppendLog("main: exit ok")
}

func mustExe() string {
	p, err := os.Executable()
	if err != nil {
		return "?"
	}
	return p
}
