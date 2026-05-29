# Go Code Visualizer

Desktop application for step-by-step visualization of Go algorithms. Built with **Wails v2**, **Delve**, **React**, and **Monaco Editor**.

Inspired by NeetCode-style visualizers: watch variables, call stack, current line, and data structure views (arrays, linked lists, trees).

## Features

- **Russian / English UI** — switch language in Settings
- **Settings** — theme (dark/light), accent color, folder for saved custom code
- **Auto-imports** — `goimports` adds standard library imports (`sort`, `fmt`, etc.) automatically
- **Save custom code** — persist your snippets to a folder on disk
- **Problem catalog** — 8 built-in algorithm problems (Two Sum, Binary Search, Linked List, Trees, DP, and more)
- **Custom code mode** — paste a full `main` program or just algorithm functions
- **Delve debugger** — full Go execution via embedded Delve (variables, pointers, structs)
- **Trace replay** — timeline slider to step backward/forward through recorded frames
- **Visualizations** — arrays with index highlights, linked lists and trees via React Flow
- **Error highlighting** — compile errors and runtime panics shown on the correct line

## Requirements

- Windows 10/11 (primary target)
- [Go 1.22+](https://go.dev/dl/) — required to compile user code at runtime
- Optional: `goimports` on PATH (otherwise the app uses bundled `golang.org/x/tools/imports`)
- [WebView2](https://developer.microsoft.com/microsoft-edge/webview2/) — usually preinstalled on Windows 11
- [Node.js 18+](https://nodejs.org/) — for frontend development only

## Quick Start

### Run from source

```bash
# Install Wails CLI (once)
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Development (hot reload)
cd c:\work\visualisator
wails dev
```

### Run release build

```bash
wails build
.\build\bin\visualisator.exe
```

## Usage

1. Pick a problem from the sidebar or choose **Custom Code**
2. Click **Run** — code compiles in a temp workspace (`%TEMP%\visualisator\`)
3. Use **Next / Step In / Step Out / Continue** while debugging, or scrub the **timeline** after the run
4. Inspect **Variables**, **Call Stack**, and **Visualization** panels on the right

## Security Notice

User code runs locally with your user permissions inside a subprocess. A 30-second timeout and process kill are applied per session. Imports such as `os/exec`, `net`, and `syscall` are blocked in sandbox scan mode.

## Project Structure

```
visualisator/
├── main.go                 # Wails entry
├── internal/
│   ├── app/                # Wails API bindings
│   ├── debugger/           # Delve integration + trace recorder
│   ├── workspace/          # Temp dirs, go build, import scan
│   ├── viz/                # Variable → visualization models
│   ├── problems/           # Embedded problem catalog
│   └── session/            # Session manager
├── frontend/               # React + Tailwind + Monaco UI
└── problems/data/          # Problem JSON definitions
```

## License

MIT
