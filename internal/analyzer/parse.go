package analyzer

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// ParseSource builds a line-indexed map from editor source (user-visible code).
func ParseSource(source string) (*SourceMap, error) {
	src := strings.TrimSpace(source)
	if src == "" {
		return nil, fmt.Errorf("empty source")
	}

	editorOffset := 0
	parseSrc := src
	if !strings.Contains(src, "package ") {
		parseSrc = "package main\n\n" + src
		editorOffset = 2
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "solution.go", parseSrc, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	lines := sourceLines(src)
	sm := &SourceMap{
		Lines:  make(map[int]StatementInfo),
		Source: src,
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
		sm.Outline = FunctionOutline{Name: "main"}
		return sm, nil
	}

	sm.Outline = buildOutline(fset, target, editorOffset)
	indexFuncBody(fset, target.Body, sm.Lines, lines, editorOffset)
	fillLineGaps(sm.Lines, lines)

	return sm, nil
}

func sourceLines(src string) []string {
	return strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n")
}

func lineText(lines []string, line int) string {
	if line < 1 || line > len(lines) {
		return ""
	}
	return strings.TrimSpace(lines[line-1])
}

func indexFuncBody(fset *token.FileSet, body *ast.BlockStmt, out map[int]StatementInfo, lines []string, editorOffset int) {
	if body == nil {
		return
	}
	for _, stmt := range body.List {
		indexStmt(fset, stmt, out, lines, editorOffset)
	}
}

func indexStmt(fset *token.FileSet, stmt ast.Stmt, out map[int]StatementInfo, lines []string, editorOffset int) {
	if stmt == nil {
		return
	}
	line := fset.Position(stmt.Pos()).Line - editorOffset
	if line < 1 {
		return
	}
	info := stmtInfo(fset, stmt, lines, line)
	if info.Line > 0 {
		out[info.Line] = info
	}

	switch s := stmt.(type) {
	case *ast.BlockStmt:
		for _, inner := range s.List {
			indexStmt(fset, inner, out, lines, editorOffset)
		}
	case *ast.IfStmt:
		if s.Init != nil {
			indexStmt(fset, s.Init, out, lines, editorOffset)
		}
		if s.Body != nil {
			indexFuncBody(fset, s.Body, out, lines, editorOffset)
		}
		if s.Else != nil {
			indexElse(fset, s.Else, out, lines, editorOffset)
		}
	case *ast.ForStmt:
		if s.Body != nil {
			indexFuncBody(fset, s.Body, out, lines, editorOffset)
		}
	case *ast.RangeStmt:
		if s.Body != nil {
			indexFuncBody(fset, s.Body, out, lines, editorOffset)
		}
	case *ast.SwitchStmt:
		if s.Body != nil {
			for _, clause := range s.Body.List {
				cc, ok := clause.(*ast.CaseClause)
				if !ok {
					continue
				}
				for _, inner := range cc.Body {
					indexStmt(fset, inner, out, lines, editorOffset)
				}
			}
		}
	}
}

func indexElse(fset *token.FileSet, stmt ast.Stmt, out map[int]StatementInfo, lines []string, editorOffset int) {
	switch s := stmt.(type) {
	case *ast.BlockStmt:
		indexFuncBody(fset, s, out, lines, editorOffset)
	case *ast.IfStmt:
		indexStmt(fset, s, out, lines, editorOffset)
	default:
		indexStmt(fset, stmt, out, lines, editorOffset)
	}
}

func stmtInfo(fset *token.FileSet, stmt ast.Stmt, lines []string, line int) StatementInfo {
	text := lineText(lines, line)
	info := StatementInfo{
		Line: line,
		Text: text,
		Kind: "other",
	}

	switch s := stmt.(type) {
	case *ast.AssignStmt:
		info.Kind = "assign"
		for _, lhs := range s.Lhs {
			info.Writes = append(info.Writes, identNames(lhs)...)
		}
		for _, rhs := range s.Rhs {
			info.Reads = append(info.Reads, identNames(rhs)...)
		}
	case *ast.IncDecStmt:
		info.Kind = "inc"
		info.Writes = identNames(s.X)
	case *ast.IfStmt:
		info.Kind = "if"
		if s.Cond != nil {
			info.ConditionExpr = exprString(fset, s.Cond)
			info.Reads = identNames(s.Cond)
		}
	case *ast.ForStmt:
		info.Kind = "for"
		if s.Cond != nil {
			info.ConditionExpr = exprString(fset, s.Cond)
			info.Reads = identNames(s.Cond)
		}
		if s.Init != nil {
			if as, ok := s.Init.(*ast.AssignStmt); ok {
				for _, rhs := range as.Rhs {
					info.Reads = append(info.Reads, identNames(rhs)...)
				}
			}
		}
	case *ast.RangeStmt:
		info.Kind = "range"
		info.ConditionExpr = exprString(fset, s.X)
		info.Reads = identNames(s.X)
		if s.Key != nil {
			info.Writes = identNames(s.Key)
		}
		if s.Value != nil {
			info.Writes = append(info.Writes, identNames(s.Value)...)
		}
	case *ast.ReturnStmt:
		info.Kind = "return"
		for _, r := range s.Results {
			info.Reads = append(info.Reads, identNames(r)...)
		}
	case *ast.ExprStmt:
		if call, ok := s.X.(*ast.CallExpr); ok {
			info.Kind = "call"
			info.Reads = identNames(call)
		} else {
			info.Kind = "expr"
		}
	case *ast.SwitchStmt:
		info.Kind = "switch"
		if s.Tag != nil {
			info.ConditionExpr = exprString(fset, s.Tag)
			info.Reads = identNames(s.Tag)
		}
	case *ast.DeferStmt:
		info.Kind = "defer"
		info.Reads = identNames(s.Call)
	default:
		info.Kind = "other"
	}

	if info.Text == "" {
		var buf bytes.Buffer
		_ = ast.Fprint(&buf, fset, stmt, nil)
		info.Text = strings.TrimSpace(buf.String())
	}
	return info
}

func fillLineGaps(out map[int]StatementInfo, lines []string) {
	var last StatementInfo
	for i := 1; i <= len(lines); i++ {
		if info, ok := out[i]; ok {
			last = info
			continue
		}
		t := strings.TrimSpace(lines[i-1])
		if t == "" || t == "{" || t == "}" {
			continue
		}
		if last.Line > 0 {
			copy := last
			copy.Line = i
			copy.Text = t
			out[i] = copy
		}
	}
}

func buildOutline(fset *token.FileSet, fn *ast.FuncDecl, editorOffset int) FunctionOutline {
	o := FunctionOutline{Name: fn.Name.Name}
	if fn.Type.Params != nil {
		o.Params = formatFieldList(fset, fn.Type.Params)
	}
	if fn.Type.Results != nil {
		o.Returns = formatFieldList(fset, fn.Type.Results)
	}
	o.DemoCall = demoCallName(fn)

	el := func(line int) int {
		if l := line - editorOffset; l > 0 {
			return l
		}
		return line
	}

	var walk func(ast.Node)
	walk = func(n ast.Node) {
		if n == nil {
			return
		}
		switch node := n.(type) {
		case *ast.ForStmt:
			line := el(fset.Position(node.For).Line)
			text := "for"
			if node.Cond != nil {
				text = "for " + exprString(fset, node.Cond)
			}
			o.Loops = append(o.Loops, OutlinePoint{Line: line, Text: text})
		case *ast.RangeStmt:
			line := el(fset.Position(node.For).Line)
			o.Loops = append(o.Loops, OutlinePoint{Line: line, Text: "range " + exprString(fset, node.X)})
		case *ast.IfStmt:
			line := el(fset.Position(node.If).Line)
			cond := ""
			if node.Cond != nil {
				cond = exprString(fset, node.Cond)
			}
			o.Branches = append(o.Branches, OutlinePoint{Line: line, Text: "if " + cond})
		case *ast.ReturnStmt:
			line := el(fset.Position(node.Return).Line)
			o.ReturnsAt = append(o.ReturnsAt, OutlinePoint{Line: line, Text: "return"})
		}
		ast.Inspect(n, func(c ast.Node) bool {
			if c == n {
				return true
			}
			switch c.(type) {
			case *ast.FuncDecl, *ast.FuncLit:
				return false
			}
			walk(c)
			return false
		})
	}
	if fn.Body != nil {
		for _, stmt := range fn.Body.List {
			walk(stmt)
		}
	}
	return o
}

func demoCallName(fn *ast.FuncDecl) string {
	if fn.Type.Params == nil || len(fn.Type.Params.List) == 0 {
		return fn.Name.Name + "()"
	}
	var args []string
	for _, f := range fn.Type.Params.List {
		args = append(args, demoValueForType(f.Type))
	}
	return fn.Name.Name + "(" + strings.Join(args, ", ") + ")"
}

func demoValueForType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		switch t.Name {
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "byte", "rune":
			return "0"
		case "float32", "float64":
			return "0"
		case "string":
			return `""`
		case "bool":
			return "false"
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
	default:
		return "nil"
	}
}

func formatFieldList(fset *token.FileSet, fl *ast.FieldList) string {
	if fl == nil {
		return ""
	}
	var parts []string
	for _, f := range fl.List {
		names := make([]string, len(f.Names))
		for i, n := range f.Names {
			names[i] = n.Name
		}
		t := exprString(fset, f.Type)
		if len(names) == 0 {
			parts = append(parts, t)
		} else {
			parts = append(parts, strings.Join(names, ", ")+" "+t)
		}
	}
	return strings.Join(parts, ", ")
}

func exprString(fset *token.FileSet, e ast.Expr) string {
	if e == nil {
		return ""
	}
	var buf bytes.Buffer
	_ = ast.Fprint(&buf, fset, e, nil)
	return strings.TrimSpace(buf.String())
}

func identNames(e ast.Expr) []string {
	var names []string
	ast.Inspect(e, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok && id.Name != "_" {
			names = append(names, id.Name)
		}
		return true
	})
	return names
}

// StmtAt returns statement info for an editor line (1-based).
func (sm *SourceMap) StmtAt(line int) (StatementInfo, bool) {
	if sm == nil || sm.Lines == nil {
		return StatementInfo{}, false
	}
	if info, ok := sm.Lines[line]; ok {
		return info, true
	}
	best := 0
	var bestInfo StatementInfo
	for l, info := range sm.Lines {
		if l <= line && l > best {
			best = l
			bestInfo = info
		}
	}
	if best > 0 {
		bestInfo.Line = line
		bestInfo.Text = lineText(sourceLines(sm.Source), line)
		return bestInfo, true
	}
	return StatementInfo{}, false
}
