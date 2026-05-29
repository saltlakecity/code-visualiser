package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"visualisator/internal/models"
	"visualisator/internal/problems"
	"visualisator/internal/session"
	"visualisator/internal/settings"
)

// App is the Wails application struct.
type App struct {
	ctx      context.Context
	sessions *session.Manager
	problems *problems.Loader
	settings *settings.Store
}

// NewApp creates the application.
func NewApp() *App {
	return &App{}
}

// Startup is called on app start.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	sm, err := session.NewManager()
	if err != nil {
		runtime.LogError(ctx, fmt.Sprintf("session manager: %v", err))
		return
	}
	a.sessions = sm
	pl, err := problems.NewLoader()
	if err != nil {
		runtime.LogError(ctx, fmt.Sprintf("problems loader: %v", err))
		return
	}
	a.problems = pl
	st, err := settings.NewStore()
	if err != nil {
		runtime.LogError(ctx, fmt.Sprintf("settings: %v", err))
	} else {
		a.settings = st
	}
	a.sessions.SetDebugCompleteHandler(func(sessionID string, trace models.Trace, err error) {
		if err != nil && trace.Status == "running" {
			trace.Status = "error"
			trace.PanicMsg = err.Error()
		}
		if trace.Status == "paused" || trace.Status == "running" {
			trace.Status = "finished"
		}
		if len(trace.Frames) > 0 {
			runtime.EventsEmit(a.ctx, "debug:frame", trace.Frames[len(trace.Frames)-1])
		}
		if trace.Status == "panic" {
			runtime.EventsEmit(a.ctx, "debug:panic", trace.PanicMsg)
		}
		runtime.EventsEmit(a.ctx, "debug:finished", trace)
	})
}

// ListProblems returns catalog metadata.
func (a *App) ListProblems() []models.ProblemMeta {
	if a.problems == nil {
		return nil
	}
	return a.problems.List()
}

// LoadProblem returns a problem by id.
func (a *App) LoadProblem(id string) (models.Problem, error) {
	if a.problems == nil {
		return models.Problem{}, fmt.Errorf("problems not loaded")
	}
	return a.problems.Get(id)
}

// Compile builds user source. mainRunner is optional main func body for algorithm problems.
func (a *App) Compile(source string, mode string, mainRunner string) (models.CompileResult, error) {
	if a.sessions == nil {
		return models.CompileResult{}, fmt.Errorf("app not ready")
	}
	res, err := a.sessions.Compile(source, mode, mainRunner)
	if err != nil {
		return res, err
	}
	if !res.Success {
		runtime.EventsEmit(a.ctx, "compile:error", res.Errors)
	}
	return res, nil
}

// RunDebug starts recording execution trace.
func (a *App) RunDebug(sessionID string, opts models.RunOptions) error {
	if a.sessions == nil {
		return fmt.Errorf("app not ready")
	}
	if err := a.sessions.StartDebug(sessionID, opts); err != nil {
		runtime.EventsEmit(a.ctx, "debug:error", err.Error())
		runtime.EventsEmit(a.ctx, "debug:finished", models.Trace{
			SessionID: sessionID,
			Status:    "error",
			PanicMsg:  err.Error(),
		})
		return err
	}
	go a.pollTrace(sessionID)
	return nil
}

func (a *App) pollTrace(sessionID string) {
	deadline := time.Now().Add(45 * time.Second)
	for time.Now().Before(deadline) {
		trace, err := a.sessions.GetTrace(sessionID)
		if err != nil {
			return
		}
		if len(trace.Frames) > 0 {
			last := trace.Frames[len(trace.Frames)-1]
			runtime.EventsEmit(a.ctx, "debug:frame", last)
		}
		if trace.Status == "finished" || trace.Status == "panic" || trace.Status == "error" {
			return // finished event emitted by handler
		}
		select {
		case <-a.ctx.Done():
			return
		default:
		}
		time.Sleep(150 * time.Millisecond)
	}
	trace, _ := a.sessions.GetTrace(sessionID)
	if trace.Status == "paused" || trace.Status == "running" {
		trace.Status = "finished"
	}
	runtime.EventsEmit(a.ctx, "debug:finished", trace)
}

// GetTrace returns full recorded trace.
func (a *App) GetTrace(sessionID string) (models.Trace, error) {
	return a.sessions.GetTrace(sessionID)
}

// Step executes debugger command.
func (a *App) Step(sessionID string, command string) (models.StepResult, error) {
	res, err := a.sessions.Step(sessionID, command)
	if err != nil {
		return res, err
	}
	if res.Frame != nil {
		runtime.EventsEmit(a.ctx, "debug:frame", *res.Frame)
	}
	if res.Finished {
		runtime.EventsEmit(a.ctx, "debug:finished", res.Trace)
	}
	return res, nil
}

// GetFrame returns trace frame by index.
func (a *App) GetFrame(sessionID string, index int) (*models.TraceFrame, error) {
	return a.sessions.GetFrame(sessionID, index)
}

// Stop ends session.
func (a *App) Stop(sessionID string) {
	a.sessions.Stop(sessionID)
}

// GetTempDir returns base temp path for display.
func (a *App) GetTempDir() string {
	return filepath.Join(os.TempDir(), "visualisator")
}

// HasAcceptedDisclaimer is a stub for first-run warning persistence.
func (a *App) HasAcceptedDisclaimer() bool {
	return true
}
