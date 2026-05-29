package problems

import (
	"embed"
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strings"

	"visualisator/internal/models"
	"visualisator/internal/testrunner"
)

//go:embed data/*.json
var dataFS embed.FS

// Loader reads embedded problem definitions.
type Loader struct {
	problems map[string]models.Problem
}

// NewLoader parses all embedded problems.
func NewLoader() (*Loader, error) {
	entries, err := dataFS.ReadDir("data")
	if err != nil {
		return nil, err
	}
	m := make(map[string]models.Problem)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		b, err := dataFS.ReadFile(path.Join("data", e.Name()))
		if err != nil {
			return nil, err
		}
		var p models.Problem
		if err := json.Unmarshal(b, &p); err != nil {
			return nil, fmt.Errorf("%s: %w", e.Name(), err)
		}
		m[p.ID] = p
	}
	return &Loader{problems: m}, nil
}

// List returns problem metadata sorted by title.
func (l *Loader) List() []models.ProblemMeta {
	out := make([]models.ProblemMeta, 0, len(l.problems))
	for _, p := range l.problems {
		out = append(out, models.ProblemMeta{ID: p.ID, Title: p.Title, Tags: p.Tags})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Title < out[j].Title })
	return out
}

// Get returns a problem by id.
func (l *Loader) Get(id string) (models.Problem, error) {
	p, ok := l.problems[id]
	if !ok {
		return models.Problem{}, fmt.Errorf("problem %q not found", id)
	}
	NormalizeProblem(&p)
	return p, nil
}

// NormalizeProblem fills test metadata.
func NormalizeProblem(p *models.Problem) {
	for i := range p.Tests {
		if p.Tests[i].Args == "" && p.Tests[i].Input != "" {
			if args, err := testrunner.ResolveArgs(p.Tests[i]); err == nil {
				p.Tests[i].Args = args
			}
		}
		if p.Tests[i].Name == "" {
			p.Tests[i].Name = fmt.Sprintf("Case %d", i+1)
		}
	}
}

// BuildMain generates a main.go for problem template execution.
func (l *Loader) BuildMain(problemID string) (string, error) {
	p, err := l.Get(problemID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`package main

import "fmt"

%s

func main() {
	fmt.Println("ready")
}
`, p.Template), nil
}
