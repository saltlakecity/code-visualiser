# Go Code Visualizer (visualisator)

Desktop app for step-through Go code visualization (Wails + Delve + React).

## Requirements

- Go 1.22+
- [Node.js](https://nodejs.org/) 18+
- [Wails CLI v2](https://wails.io/docs/gettingstarted/installation): `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- [WebView2 Runtime](https://developer.microsoft.com/microsoft-edge/webview2/) (usually preinstalled on Windows 11)
- `go`, `dlv` on `PATH` for compile/debug

## Build and run (Windows)

**Always use Wails** — do not use plain `go build` for the desktop app.

```powershell
cd c:\work\visualisator
wails build
.\build\bin\visualisator.exe
```

The release binary is a **GUI app** (`-H windowsgui` in `wails.json`). You should **not** see a console window flash. If you do, rebuild with `wails build` (not `go run` / not `go build`).

### Dev mode (hot reload)

```powershell
wails dev
```

## Troubleshooting startup

Log file: `data\app.log` next to `visualisator.exe` (created on first run).

**Healthy startup:**

```
main: start exe=...
lifecycle: dom ready
main: exit ok
```

(`main: exit ok` appears only after you close the window.)

If the log shows `dom ready` but there is no window:

1. Check the taskbar / **Alt+Tab** for «Go Code Visualizer».
2. Delete WebView cache: `Remove-Item -Recurse -Force .\data\webview2`
3. Rebuild: `wails build`
4. Do **not** use `wails build -windowsconsole` unless you need a console for debugging.

Debug build (DevTools + console):

```powershell
wails build -debug
.\build\bin\visualisator.exe
```

## Portable data

```
visualisator.exe
data/
  settings.json
  custom/
  sessions/
  app.log
```

## Usage

1. Pick a problem or **Custom Code**
2. **Run** — compile + debug trace
3. Step / timeline / variables / visualization (graphs load when trace has structures)

Code editor is a lightweight text area (Monaco removed to keep WebView2 stable in production builds).

## Windows installer

```powershell
wails build -nsis
```

## Application icon

Place `build/appicon.png`, then rebuild. See Wails docs for `wails generate icons`.
