package debugger

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-delve/delve/pkg/proc"
	"github.com/go-delve/delve/service/api"
	"github.com/go-delve/delve/service/debugger"

	"visualisator/internal/models"
)

var loadCfg = proc.LoadConfig{
	FollowPointers:     true,
	MaxVariableRecurse: 5,
	MaxStringLen:       256,
	MaxArrayValues:     64,
	MaxStructFields:    -1,
}

// Session wraps a Delve debugger instance.
type Session struct {
	dbg      *debugger.Debugger
	recorder *Recorder
	workDir  string
	binary   string
	timeout  time.Duration
	done     chan struct{}
}

// NewSession launches Delve on the compiled binary.
func NewSession(binaryPath, workDir string, timeout time.Duration) (*Session, error) {
	cfg := &debugger.Config{
		WorkingDir: workDir,
		Backend:    "default",
	}
	dbg, err := debugger.New(cfg, []string{binaryPath})
	if err != nil {
		return nil, fmt.Errorf("debugger init: %w", err)
	}

	s := &Session{
		dbg:      dbg,
		recorder: NewRecorder(),
		workDir:  workDir,
		binary:   binaryPath,
		timeout:  timeout,
		done:     make(chan struct{}),
	}
	s.recorder.Reset()

	state, err := s.cmd(&api.DebuggerCommand{Name: api.Halt})
	if err != nil {
		// process may already be stopped at entry
		state, err = s.dbg.State(false)
		if err != nil {
			s.Stop()
			return nil, err
		}
	}
	if state != nil && !state.Exited && !state.Running {
		_, _ = s.captureState(state)
	}
	return s, nil
}

func (s *Session) cmd(command *api.DebuggerCommand) (*api.DebuggerState, error) {
	return s.dbg.Command(command, nil, nil, nil)
}

// Recorder returns the trace recorder.
func (s *Session) Recorder() *Recorder {
	return s.recorder
}

// Step executes a debugger command and records the frame.
func (s *Session) Step(command string) (*models.StepResult, error) {
	name, err := commandName(command)
	if err != nil {
		return nil, err
	}
	state, err := s.cmd(&api.DebuggerCommand{Name: name})
	if err != nil {
		return &models.StepResult{Error: err.Error()}, err
	}
	return s.processState(state)
}

func commandName(command string) (string, error) {
	switch strings.ToLower(command) {
	case "next":
		return api.Next, nil
	case "step", "stepin":
		return api.Step, nil
	case "stepout":
		return api.StepOut, nil
	case "continue", "cont":
		return api.Continue, nil
	case "halt":
		return api.Halt, nil
	default:
		return "", fmt.Errorf("unknown command %q", command)
	}
}

// RunToEnd steps through the program recording trace frames.
func (s *Session) RunToEnd(maxSteps int) (*models.Trace, error) {
	if maxSteps <= 0 {
		maxSteps = 500
	}

	for i := 0; i < maxSteps; i++ {
		select {
		case <-s.done:
			return nil, fmt.Errorf("session stopped")
		default:
		}

		state, err := s.dbg.State(false)
		if err != nil {
			break
		}
		if state.Exited {
			break
		}

		if !state.Running && state.CurrentThread != nil {
			_, _ = s.captureState(state)
		}

		nextState, err := s.cmd(&api.DebuggerCommand{Name: api.Next})
		if err != nil {
			break
		}
		if nextState == nil {
			break
		}
		if nextState.Exited {
			break
		}
		if !nextState.Running {
			if _, err := s.processState(nextState); err != nil {
				break
			}
		}
	}

	s.recorder.SetStatus("finished", "")
	trace := s.recorder.Trace("")
	return &trace, nil
}

func (s *Session) processState(state *api.DebuggerState) (*models.StepResult, error) {
	if state == nil {
		return &models.StepResult{Finished: true}, nil
	}

	if state.Exited {
		exitStatus := "finished"
		panicMsg := ""
		if state.Err != nil {
			errStr := state.Err.Error()
			if strings.Contains(strings.ToLower(errStr), "panic") {
				exitStatus = "panic"
				panicMsg = errStr
			} else {
				exitStatus = "error"
				panicMsg = errStr
			}
		}
		s.recorder.SetStatus(exitStatus, panicMsg)
		trace := s.recorder.Trace("")
		return &models.StepResult{Finished: true, Trace: &trace}, nil
	}

	if state.Running {
		return &models.StepResult{Finished: false}, nil
	}

	frame, err := s.captureState(state)
	if err != nil {
		return nil, err
	}
	trace := s.recorder.Trace("")
	return &models.StepResult{Frame: frame, Trace: &trace, Finished: false}, nil
}

func (s *Session) captureState(state *api.DebuggerState) (*models.TraceFrame, error) {
	thread := state.CurrentThread
	if thread == nil {
		return nil, nil
	}

	// Skip runtime frames without source
	if thread.File == "" && thread.Line == 0 {
		return nil, nil
	}

	goid := int64(-1)
	if state.SelectedGoroutine != nil {
		goid = state.SelectedGoroutine.ID
	} else if thread.GoroutineID > 0 {
		goid = thread.GoroutineID
	}

	var stack []models.StackFrame
	rawStack, err := s.dbg.Stacktrace(goid, 32, 0)
	if err == nil {
		apiFrames, convErr := s.dbg.ConvertStacktrace(rawStack, &loadCfg)
		if convErr == nil {
			for _, sf := range apiFrames {
				fn := ""
				if sf.Function != nil {
					fn = sf.Function.Name()
				}
				stack = append(stack, models.StackFrame{
					File:     basename(sf.File),
					Line:     sf.Line,
					Function: fn,
				})
			}
		}
	}

	var vars []models.Variable
	if goid >= 0 {
		procVars, err := s.dbg.LocalVariables(goid, 0, 0, loadCfg)
		if err == nil {
			for _, v := range api.ConvertVars(procVars) {
				vars = append(vars, convertVariable(v))
			}
		}
	}

	fn := ""
	if thread.Function != nil {
		fn = thread.Function.Name()
	}

	file := basename(thread.File)
	if file == "" && len(stack) > 0 {
		file = stack[0].File
	}
	line := thread.Line
	if line == 0 && len(stack) > 0 {
		line = stack[0].Line
	}

	// Prefer recording frames in user source files
	if !isUserSource(file) && len(stack) > 0 {
		for _, sf := range stack {
			if isUserSource(sf.File) {
				file = sf.File
				line = sf.Line
				fn = sf.Function
				break
			}
		}
	}

	if file == "" || line == 0 {
		return nil, nil
	}

	tf := BuildFrame(file, line, fn, vars, stack)
	s.recorder.Append(tf)
	s.recorder.SetStatus("paused", "")
	return &tf, nil
}

func isUserSource(file string) bool {
	base := strings.ToLower(filepath.Base(file))
	return base == "main.go" || base == "solution.go"
}

func convertVariable(v api.Variable) models.Variable {
	mv := models.Variable{
		Name:  v.Name,
		Type:  v.Type,
		Value: v.Value,
	}
	for _, child := range v.Children {
		mv.Children = append(mv.Children, convertVariable(child))
	}
	return mv
}

// Stop halts and detaches the debugger.
func (s *Session) Stop() {
	select {
	case <-s.done:
	default:
		close(s.done)
	}
	if s.dbg != nil {
		_ = s.dbg.Detach(true)
		s.dbg = nil
	}
}

// StartTimeout kills the session after duration.
func (s *Session) StartTimeout() {
	if s.timeout <= 0 {
		return
	}
	go func() {
		select {
		case <-time.After(s.timeout):
			s.Stop()
		case <-s.done:
		}
	}()
}
