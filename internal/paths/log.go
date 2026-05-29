package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AppendLog writes a line to data/app.log (best-effort).
func AppendLog(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	_ = os.MkdirAll(DataDir(), 0o755)
	f, err := os.OpenFile(LogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	_, _ = fmt.Fprintf(f, "%s %s\n", time.Now().Format(time.RFC3339), msg)
	_ = f.Close()
}

// WebViewDataDir is the WebView2 user-data folder (must be writable).
func WebViewDataDir() string {
	dir := filepath.Join(DataDir(), "webview2")
	_ = os.MkdirAll(dir, 0o755)
	return dir
}
