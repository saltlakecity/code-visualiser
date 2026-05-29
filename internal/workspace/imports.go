package workspace

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/tools/imports"
)

// PrepareUserSource ensures package and optional main stub (imports fixed per-directory later).
func PrepareUserSource(source string, mode string) (string, string, error) {
	src := strings.TrimSpace(source)
	if src == "" {
		return "", "", fmt.Errorf("empty source")
	}

	if !strings.Contains(src, "package ") {
		src = "package main\n\n" + src
	}

	mainExtra := ""
	if !HasMainPackage(src) {
		var err error
		mainExtra, err = generateStubMain(src)
		if err != nil {
			return "", "", err
		}
	}

	return src, mainExtra, nil
}

// ProcessDirImports runs goimports on all Go files in a module directory.
func ProcessDirImports(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		out, err := imports.Process(path, b, &imports.Options{
			Comments:  true,
			TabIndent: true,
			TabWidth:  4,
		})
		if err != nil {
			return fmt.Errorf("%s: %w", e.Name(), err)
		}
		if err := os.WriteFile(path, out, 0o644); err != nil {
			return err
		}
	}
	return nil
}

// ProcessDirImportsExec uses goimports CLI when available (better for multi-file packages).
func ProcessDirImportsExec(dir string) error {
	if _, err := exec.LookPath("goimports"); err == nil {
		cmd := exec.Command("goimports", "-w", ".")
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("goimports: %s: %w", out, err)
		}
		return nil
	}
	return ProcessDirImports(dir)
}

func generateStubMain(source string) (string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "solution.go", source, 0)
	if err != nil {
		return "", fmt.Errorf("parse for main stub: %w", err)
	}

	var target *ast.FuncDecl
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name == "main" || fn.Recv != nil {
			continue
		}
		target = fn
		break
	}
	if target == nil {
		return `package main

func main() {
}
`, nil
	}

	call := buildDemoCall(target)
	var buf bytes.Buffer
	buf.WriteString("package main\n\n")
	buf.WriteString("import \"fmt\"\n\n")
	buf.WriteString("func main() {\n")
	buf.WriteString("\tfmt.Println(" + call + ")\n")
	buf.WriteString("}\n")
	return buf.String(), nil
}

func buildDemoCall(fn *ast.FuncDecl) string {
	if fn.Type.Params == nil || len(fn.Type.Params.List) == 0 {
		return fn.Name.Name + "()"
	}
	var args []string
	for _, field := range fn.Type.Params.List {
		args = append(args, demoValueForType(field.Type))
	}
	return fn.Name.Name + "(" + strings.Join(args, ", ") + ")"
}

func demoValueForType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		switch t.Name {
		case "int", "int8", "int16", "int32", "int64":
			return "0"
		case "uint", "uint8", "uint16", "uint32", "uint64":
			return "0"
		case "float32", "float64":
			return "0"
		case "string":
			return `""`
		case "bool":
			return "false"
		case "byte", "rune":
			return "0"
		default:
			return "nil"
		}
	case *ast.StarExpr:
		return "nil"
	case *ast.ArrayType:
		if t.Len == nil {
			if ident, ok := t.Elt.(*ast.Ident); ok && ident.Name == "int" {
				return "[]int{1, 2, 3, 1}"
			}
			return "nil"
		}
		return "[1]int{0}"
	case *ast.MapType:
		return "nil"
	case *ast.InterfaceType:
		return "nil"
	case *ast.SelectorExpr:
		return "nil"
	default:
		return "nil"
	}
}
