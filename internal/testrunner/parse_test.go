package testrunner

import (
	"testing"

	"visualisator/internal/models"
)

func TestParseInput(t *testing.T) {
	args, err := ParseInput("nums=[2,7,11,15], target=9")
	if err != nil {
		t.Fatal(err)
	}
	if args != "[]int{2, 7, 11, 15}, 9" {
		t.Fatalf("unexpected args: %q", args)
	}
}

func TestResolveArgsPrefersExplicit(t *testing.T) {
	tc := models.TestCase{Args: "[]int{1}, 2", Input: "nums=[3], target=4"}
	args, err := ResolveArgs(tc)
	if err != nil || args != "[]int{1}, 2" {
		t.Fatalf("got %q err=%v", args, err)
	}
}

func TestCompareOutput(t *testing.T) {
	if !CompareOutput("[0 1]\n", "[0, 1]") {
		t.Fatal("expected match")
	}
	if CompareOutput("1", "2") {
		t.Fatal("expected mismatch")
	}
}

func TestBuildMain(t *testing.T) {
	main, err := BuildMain("twoSum", models.TestCase{
		Input:    "nums=[1,2], target=3",
		Expected: "[0 1]",
	}, false)
	if err != nil {
		t.Fatal(err)
	}
	if !containsAll(main, "func main()", "twoSum(", "[]int{1, 2}") {
		t.Fatalf("unexpected main: %s", main)
	}
}

func containsAll(s string, parts ...string) bool {
	for _, p := range parts {
		if !contains(s, p) {
			return false
		}
	}
	return true
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
