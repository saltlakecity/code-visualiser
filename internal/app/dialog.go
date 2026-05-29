package app

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// PickDirectoryDialog opens OS folder picker.
func PickDirectoryDialog(ctx context.Context) (string, error) {
	return runtime.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{
		Title: "Select folder for custom code",
	})
}
