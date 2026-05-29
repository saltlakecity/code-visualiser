package debugger

import (
	"testing"
)

func TestClassifyExit_success(t *testing.T) {
	status, msg, kind, failed := classifyExit(nil)
	if failed || status != "finished" || msg != "" || kind != "" {
		t.Fatalf("got status=%q msg=%q kind=%q failed=%v", status, msg, kind, failed)
	}

	status, msg, kind, failed = classifyExit(errString("Process 27452 has exited with status 0"))
	if failed || status != "finished" {
		t.Fatalf("status 0: status=%q failed=%v msg=%q", status, failed, msg)
	}
}

func TestClassifyExit_panic(t *testing.T) {
	status, _, kind, failed := classifyExit(errString("panic: index out of range"))
	if !failed || status != "panic" || kind != "panic" {
		t.Fatalf("panic: status=%q kind=%q failed=%v", status, kind, failed)
	}
}

func TestClassifyExit_nonZero(t *testing.T) {
	status, _, _, failed := classifyExit(errString("Process 1 has exited with status 2"))
	if !failed || status != "error" {
		t.Fatalf("status 2: status=%q failed=%v", status, failed)
	}
}

type errString string

func (e errString) Error() string { return string(e) }

func TestClassifyExit_hasExited(t *testing.T) {
	_, _, _, failed := classifyExit(errString("target process has exited"))
	if failed {
		t.Fatal("expected benign exit")
	}
}
