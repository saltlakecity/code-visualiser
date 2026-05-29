package workspace

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

var deniedImports = []string{
	"os/exec",
	"syscall",
	"unsafe",
	"plugin",
	"net/http",
	"net",
}

// ScanImports rejects dangerous import paths in user source.
func ScanImports(source string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "user.go", source, parser.ImportsOnly)
	if err != nil {
		// full parse may fail for partial snippets; try parsing as file with package
		f, err = parser.ParseFile(fset, "user.go", source, 0)
		if err != nil {
			return nil
		}
	}
	for _, imp := range f.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		for _, denied := range deniedImports {
			if path == denied || strings.HasPrefix(path, denied+"/") {
				return fmt.Errorf("import %q is not allowed in sandbox mode", path)
			}
		}
	}
	return nil
}

// ValidateGoMod blocks replace directives.
func ValidateGoMod(content string) error {
	for _, line := range strings.Split(content, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "replace ") {
			return fmt.Errorf("replace directives are not allowed in go.mod")
		}
	}
	return nil
}

// HasMainPackage checks if source declares package main with func main.
func HasMainPackage(source string) bool {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "main.go", source, 0)
	if err != nil {
		return false
	}
	if f.Name.Name != "main" {
		return false
	}
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "main" {
			continue
		}
		return true
	}
	return false
}
