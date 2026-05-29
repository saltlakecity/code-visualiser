package debugger

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"visualisator/internal/models"
)

func TestRunToEndTwoSum(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "solution.go"), `package main

func twoSum(nums []int, target int) []int {
	m := make(map[int]int)
	for i, n := range nums {
		if j, ok := m[target-n]; ok {
			return []int{j, i}
		}
		m[n] = i
	}
	return nil
}
`)
	writeFile(t, filepath.Join(dir, "main.go"), `package main

import "fmt"

func main() {
	fmt.Println(twoSum([]int{2, 7, 11, 15}, 9))
}
`)
	cmd := exec.Command("go", "mod", "init", "usercode")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("mod init: %s %v", out, err)
	}
	bin := filepath.Join(dir, "user.exe")
	if runtime.GOOS != "windows" {
		bin = filepath.Join(dir, "user")
	}
	build := exec.Command("go", "build", "-gcflags=all=-N -l", "-o", bin, ".")
	build.Dir = dir
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %s %v", out, err)
	}

	s, err := NewSession(bin, dir, 10*time.Second)
	if err != nil {
		t.Fatalf("new session: %v", err)
	}
	defer s.Stop()

	done := make(chan struct{})
	var trace *models.Trace
	go func() {
		tr, err := s.RunToEnd(200, 0)
		if err != nil {
			t.Logf("run error: %v", err)
		}
		trace = tr
		close(done)
	}()

	select {
	case <-done:
		t.Logf("status=%s frames=%d", trace.Status, len(trace.Frames))
		for i, f := range trace.Frames {
			t.Logf("  frame %d: %s:%d %s vars=%d viz=%d", i, f.File, f.Line, f.Function, len(f.Variables), len(f.Viz))
		}
		if len(trace.Frames) == 0 {
			t.Fatal("no frames recorded")
		}
		if len(trace.Frames) < 3 {
			t.Fatalf("expected multiple frames, got %d", len(trace.Frames))
		}
		var inTwoSum bool
		var hasMapVar bool
		for _, f := range trace.Frames {
			if f.File == "solution.go" && f.Line > 0 {
				inTwoSum = true
			}
			for _, v := range f.Variables {
				if len(v.Type) >= 3 && v.Type[:3] == "map" {
					hasMapVar = true
					break
				}
			}
		}
		if !inTwoSum {
			t.Fatal("expected frames in solution.go")
		}
		if !hasMapVar {
			t.Fatal("expected map variable in trace (twoSum uses map)")
		}
	case <-time.After(15 * time.Second):
		s.Stop()
		t.Fatal("timeout waiting for RunToEnd")
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
