package workspace

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"

	"golang.org/x/tools/imports"
)

// FormatAndImports runs goimports and gofmt on Go source.
func FormatAndImports(filename string, source []byte) ([]byte, error) {
	processed, err := imports.Process(filename, source, &imports.Options{
		Comments:  true,
		TabIndent: true,
		TabWidth:  4,
	})
	if err != nil {
		return nil, fmt.Errorf("imports: %w", err)
	}
	formatted, err := format.Source(processed)
	if err != nil {
		return processed, nil
	}
	return formatted, nil
}

// PrepareUserSource ensures package, imports, and main for custom/snippets.
func PrepareUserSource(source string, mode string) (string, string, error) {
	src := strings.TrimSpace(source)
	if src == "" {
		return "", "", fmt.Errorf("empty source")
	}

	filename := "main.go"
	if mode == "function_wrapper" {
		filename = "solution.go"
	}

	if !strings.Contains(src, "package ") {
		src = "package main\n\n" + src
	}

	processed, err := FormatAndImports(filename, []byte(src))
	if err != nil {
		// retry after ensuring parseable
		processed = []byte(src)
	}

	final := string(processed)
	mainExtra := ""

	if !HasMainPackage(final) {
		mainExtra, err = generateStubMain(final)
		if err != nil {
			return "", "", err
		}
	}

	return final, mainExtra, nil
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
}`, nil
	}

	call := buildDemoCall(target)
	var buf bytes.Buffer
	buf.WriteString("package main\n\n")
	buf.WriteString("import \"fmt\"\n\n")
	buf.WriteString("func main() {\n")
	buf.WriteString("\tfmt.Println(" + call + ")\n")
	buf.WriteString("}\n")
	out, err := FormatAndImports("main.go", buf.Bytes())
	if err != nil {
		return buf.String(), nil
	}
	return string(out), nil
}

func buildDemoCall(fn *ast.FuncDecl) string {
	if fn.Type.Params == nil || len(fn.Type.Params.List) == 0 {
		return fn.Name.Name + "()"
	}
	var args []string
	for _, field := range fn.Type.Params.List {
		arg := demoValueForType(field.Type)
		args = append(args, arg)
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
