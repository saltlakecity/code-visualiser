package app

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"visualisator/internal/models"
	"visualisator/internal/paths"
	"visualisator/internal/problems"
	"visualisator/internal/session"
	"visualisator/internal/settings"
)

// App is the Wails application struct.
type App struct {
	ctx      context.Context
	sessions *session.Manager
	sessMu   sync.Mutex
	problems *problems.Loader
	settings *settings.Store
}

// NewApp creates the application.
func NewApp() *App {
	return &App{}
}

func logf(format string, args ...any) {
	paths.AppendLog(format, args...)
}

// Startup is called on app start.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	logf("startup")

	pl, err := problems.NewLoader()
	if err != nil {
		logf("problems loader error: %v", err)
		runtime.LogError(ctx, fmt.Sprintf("problems loader: %v", err))
	} else {
		a.problems = pl
	}

	st, err := settings.NewStore()
	if err != nil {
		logf("settings error: %v", err)
		runtime.LogError(ctx, fmt.Sprintf("settings: %v", err))
	} else {
		a.settings = st
	}

}

func (a *App) ensureSessions() (*session.Manager, error) {
	a.sessMu.Lock()
	defer a.sessMu.Unlock()
	if a.sessions != nil {
		return a.sessions, nil
	}
	sm, err := session.NewManager()
	if err != nil {
		logf("session manager error: %v", err)
		return nil, err
	}
	a.sessions = sm
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
	return a.sessions, nil
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
	p, err := a.problems.Get(id)
	if err != nil {
		return p, err
	}
	NormalizeProblemTests(&p)
	return p, nil
}

// Compile builds user source. mainRunner is optional main func body for algorithm problems.
func (a *App) Compile(source string, mode string, mainRunner string) (models.CompileResult, error) {
	sm, err := a.ensureSessions()
	if err != nil {
		return models.CompileResult{}, err
	}
	res, err := sm.Compile(source, mode, mainRunner)
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
	sm, err := a.ensureSessions()
	if err != nil {
		return err
	}
	if err := sm.StartDebug(sessionID, opts); err != nil {
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
		sm, err := a.ensureSessions()
		if err != nil {
			return
		}
		trace, err := sm.GetTrace(sessionID)
		if err != nil {
			return
		}
		if len(trace.Frames) > 0 {
			last := trace.Frames[len(trace.Frames)-1]
			runtime.EventsEmit(a.ctx, "debug:frame", last)
		}
		if trace.Status == "finished" || trace.Status == "panic" || trace.Status == "error" {
			return
		}
		select {
		case <-a.ctx.Done():
			return
		default:
		}
		time.Sleep(150 * time.Millisecond)
	}
	sm, _ := a.ensureSessions()
	if sm == nil {
		return
	}
	trace, _ := sm.GetTrace(sessionID)
	if trace.Status == "paused" || trace.Status == "running" {
		trace.Status = "finished"
	}
	runtime.EventsEmit(a.ctx, "debug:finished", trace)
}

// GetTrace returns full recorded trace.
func (a *App) GetTrace(sessionID string) (models.Trace, error) {
	sm, err := a.ensureSessions()
	if err != nil {
		return models.Trace{}, err
	}
	return sm.GetTrace(sessionID)
}

// Step executes debugger command.
func (a *App) Step(sessionID string, command string) (models.StepResult, error) {
	sm, err := a.ensureSessions()
	if err != nil {
		return models.StepResult{}, err
	}
	res, err := sm.Step(sessionID, command)
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
	sm, err := a.ensureSessions()
	if err != nil {
		return nil, err
	}
	return sm.GetFrame(sessionID, index)
}

// Stop ends session.
func (a *App) Stop(sessionID string) {
	sm, err := a.ensureSessions()
	if err != nil {
		return
	}
	sm.Stop(sessionID)
}

// GetTempDir returns base temp path for display.
func (a *App) GetTempDir() string {
	return filepath.Join(paths.DataDir(), "sessions")
}

// HasAcceptedDisclaimer is a stub for first-run warning persistence.
func (a *App) HasAcceptedDisclaimer() bool {
	return true
}
