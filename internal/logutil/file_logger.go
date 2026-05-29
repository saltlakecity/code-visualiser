package logutil

import (
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2/pkg/logger"

	"visualisator/internal/paths"
)

// FileLogger writes Wails log output to data/app.log and stderr for errors.
type FileLogger struct{}

func (FileLogger) Print(message string)   { paths.AppendLog("print: %s", message) }
func (FileLogger) Trace(message string)   { paths.AppendLog("trace: %s", message) }
func (FileLogger) Debug(message string)   { paths.AppendLog("debug: %s", message) }
func (FileLogger) Info(message string)    { paths.AppendLog("info: %s", message) }
func (FileLogger) Warning(message string) { paths.AppendLog("warn: %s", message) }

func (FileLogger) Error(message string) {
	paths.AppendLog("error: %s", message)
	fmt.Fprintln(os.Stderr, "ERROR:", message)
}

func (FileLogger) Fatal(message string) {
	paths.AppendLog("fatal: %s", message)
	fmt.Fprintln(os.Stderr, "FATAL:", message)
}

var _ logger.Logger = FileLogger{}
