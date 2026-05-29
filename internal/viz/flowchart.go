package viz

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"visualisator/internal/models"
)

const maxLabelLen = 28

type cfgBuilder struct {
	fset         *token.FileSet
	editorOffset int
	endID        string
	nodes        []models.VizNode
	edges        []models.VizEdge
	counter      int
}

// BuildFlowGraph parses source and builds a control-flow flowchart.
func BuildFlowGraph(source string) (models.FlowchartResult, error) {
	src := strings.TrimSpace(source)
	if src == "" {
		return models.FlowchartResult{}, fmt.Errorf("empty source")
	}

	editorOffset := 0
	parseSrc := src
	if !strings.Contains(src, "package ") {
		parseSrc = "package main\n\n" + src
		editorOffset = 2
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "solution.go", parseSrc, 0)
	if err != nil {
		return models.FlowchartResult{}, fmt.Errorf("parse: %w", err)
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

	result := models.FlowchartResult{FunctionName: "main"}
	if target == nil || target.Body == nil {
		return result, nil
	}
	result.FunctionName = target.Name.Name

	b := &cfgBuilder{fset: fset, editorOffset: editorOffset}
	startLine := b.el(fset.Position(target.Body.Lbrace).Line + 1)
	endLine := b.el(fset.Position(target.Body.Rbrace).Line)
	if startLine < 1 {
		startLine = 1
	}
	if endLine < startLine {
		endLine = startLine
	}

	startID := b.addNode("start", "Start", startLine)
	b.endID = b.addNode("end", "End", endLine)

	entry, exits := b.buildSequence(target.Body.List)
	if entry != "" {
		b.connect(startID, entry, "")
	} else {
		b.connect(startID, b.endID, "")
	}
	for _, ex := range exits {
		if ex != "" {
			b.connect(ex, b.endID, "")
		}
	}

	startIDRef := startID
	ApplyFlowLayout(&b.nodes, b.edges, startIDRef, b.endID)

	result.Nodes = b.nodes
	result.Edges = b.edges
	result.Summary = fmt.Sprintf("%s — %d nodes", result.FunctionName, len(b.nodes))
	return result, nil
}

func (b *cfgBuilder) el(line int) int {
	if l := line - b.editorOffset; l > 0 {
		return l
	}
	return line
}

func (b *cfgBuilder) nextID(kind string) string {
	b.counter++
	return fmt.Sprintf("%s_%d", kind, b.counter)
}

func (b *cfgBuilder) addNode(kind, label string, line int) string {
	id := b.nextID(kind)
	label = shortenFlowLabel(label, maxLabelLen)
	if label == "" {
		label = flowKindLabel(kind)
	}
	b.nodes = append(b.nodes, models.VizNode{
		ID:    id,
		Kind:  kind,
		Label: label,
		Line:  line,
	})
	return id
}

func (b *cfgBuilder) connect(from, to, label string) {
	if from == "" || to == "" {
		return
	}
	b.edges = append(b.edges, models.VizEdge{From: from, To: to, Label: label})
}

func (b *cfgBuilder) buildSequence(stmts []ast.Stmt) (entry string, exits []string) {
	var first string
	var prev []string

	for _, stmt := range stmts {
		e, out := b.buildStmt(stmt)
		if e == "" {
			continue
		}
		if first == "" {
			first = e
		}
		for _, p := range prev {
			b.connect(p, e, "")
		}
		if out == nil {
			return first, nil
		}
		prev = out
	}
	return first, prev
}

func (b *cfgBuilder) buildBlock(body *ast.BlockStmt) (entry string, exits []string) {
	if body == nil || len(body.List) == 0 {
		return "", []string{""}
	}
	return b.buildSequence(body.List)
}

func (b *cfgBuilder) buildStmt(stmt ast.Stmt) (entry string, exits []string) {
	switch s := stmt.(type) {
	case *ast.IfStmt:
		return b.buildIf(s)
	case *ast.ForStmt:
		return b.buildFor(s)
	case *ast.RangeStmt:
		return b.buildRange(s)
	case *ast.SwitchStmt:
		return b.buildSwitch(s)
	case *ast.ReturnStmt:
		line := b.el(b.fset.Position(s.Return).Line)
		id := b.addNode("return", "return", line)
		b.connect(id, b.endID, "")
		return id, nil
	case *ast.BlockStmt:
		return b.buildBlock(s)
	default:
		label, rawLine := friendlyFlowLabel(b.fset, stmt)
		id := b.addNode("process", label, b.el(rawLine))
		return id, []string{id}
	}
}

func (b *cfgBuilder) buildIf(s *ast.IfStmt) (entry string, exits []string) {
	line := b.el(b.fset.Position(s.If).Line)
	cond := "if"
	if s.Cond != nil {
		cond = "if " + shortenFlowExpr(b.fset, s.Cond)
	}
	decID := b.addNode("decision", cond, line)

	thenEntry, thenExits := b.buildBlock(s.Body)
	if thenEntry == "" {
		thenEntry = b.addNode("process", "then", line)
		thenExits = []string{thenEntry}
	}
	b.connect(decID, thenEntry, "yes")

	var elseEntry string
	var elseExits []string
	if s.Else != nil {
		elseEntry, elseExits = b.buildElse(s.Else)
		if elseEntry == "" {
			elseEntry = b.addNode("process", "else", line)
			elseExits = []string{elseEntry}
		}
		b.connect(decID, elseEntry, "no")
	} else {
		// false branch skips then-body — mark decision as alternate exit
		elseExits = []string{decID + ":no"}
	}

	return decID, mergeExits(thenExits, elseExits, decID)
}

func (b *cfgBuilder) buildElse(stmt ast.Stmt) (entry string, exits []string) {
	switch s := stmt.(type) {
	case *ast.BlockStmt:
		return b.buildBlock(s)
	case *ast.IfStmt:
		return b.buildIf(s)
	default:
		return b.buildStmt(stmt)
	}
}

func mergeExits(a, bExits []string, decID string) []string {
	var out []string
	add := func(list []string) {
		for _, x := range list {
			if x == "" {
				continue
			}
			if strings.HasSuffix(x, ":no") {
				out = append(out, decID)
				continue
			}
			out = append(out, x)
		}
	}
	add(a)
	add(bExits)
	if len(out) == 0 {
		return nil
	}
	seen := make(map[string]bool)
	var uniq []string
	for _, x := range out {
		if !seen[x] {
			seen[x] = true
			uniq = append(uniq, x)
		}
	}
	return uniq
}

func (b *cfgBuilder) buildFor(s *ast.ForStmt) (entry string, exits []string) {
	line := b.el(b.fset.Position(s.For).Line)
	label := "for"
	if s.Cond != nil {
		label = "for " + shortenFlowExpr(b.fset, s.Cond)
	}
	loopID := b.addNode("loop", label, line)

	bodyEntry, bodyExits := b.buildBlock(s.Body)
	if bodyEntry == "" {
		bodyEntry = b.addNode("process", "body", line)
		bodyExits = []string{bodyEntry}
	}
	b.connect(loopID, bodyEntry, "enter")
	for _, ex := range bodyExits {
		if ex != "" {
			b.connect(ex, loopID, "next")
			b.edges[len(b.edges)-1].Animated = true
		}
	}

	exitID := b.addNode("process", "exit loop", line)
	b.connect(loopID, exitID, "exit")
	return loopID, []string{exitID}
}

func (b *cfgBuilder) buildRange(s *ast.RangeStmt) (entry string, exits []string) {
	line := b.el(b.fset.Position(s.For).Line)
	label := "range"
	if key := flowExprNames(s.Key); len(key) > 0 {
		label = "range " + strings.Join(key, ",")
	}
	loopID := b.addNode("loop", label, line)

	bodyEntry, bodyExits := b.buildBlock(s.Body)
	if bodyEntry == "" {
		bodyEntry = b.addNode("process", "body", line)
		bodyExits = []string{bodyEntry}
	}
	b.connect(loopID, bodyEntry, "enter")
	for _, ex := range bodyExits {
		if ex != "" {
			b.connect(ex, loopID, "next")
			b.edges[len(b.edges)-1].Animated = true
		}
	}

	exitID := b.addNode("process", "exit loop", line)
	b.connect(loopID, exitID, "exit")
	return loopID, []string{exitID}
}

func (b *cfgBuilder) buildSwitch(s *ast.SwitchStmt) (entry string, exits []string) {
	line := b.el(b.fset.Position(s.Switch).Line)
	tag := "switch"
	if s.Tag != nil {
		tag = "switch " + shortenFlowExpr(b.fset, s.Tag)
	}
	decID := b.addNode("decision", tag, line)

	if s.Body == nil || len(s.Body.List) == 0 {
		exitID := b.addNode("process", "exit switch", line)
		b.connect(decID, exitID, "")
		return decID, []string{exitID}
	}

	var allExits []string
	for _, clause := range s.Body.List {
		cc, ok := clause.(*ast.CaseClause)
		if !ok {
			continue
		}
		caseLabel := "default"
		if len(cc.List) > 0 {
			caseLabel = shortenFlowExpr(b.fset, cc.List[0])
		}
		caseID := b.addNode("process", caseLabel, line)
		b.connect(decID, caseID, caseLabel)

		if len(cc.Body) > 0 {
			bodyEntry, bodyExits := b.buildSequence(cc.Body)
			if bodyEntry != "" {
				b.connect(caseID, bodyEntry, "")
			}
			allExits = append(allExits, bodyExits...)
		} else {
			allExits = append(allExits, caseID)
		}
	}

	exitID := b.addNode("process", "exit switch", line)
	var fallExits []string
	for _, ex := range allExits {
		if ex != "" {
			fallExits = append(fallExits, ex)
		}
	}
	if len(fallExits) == 0 {
		b.connect(decID, exitID, "")
		return decID, []string{exitID}
	}
	for _, ex := range fallExits {
		b.connect(ex, exitID, "")
	}
	return decID, []string{exitID}
}

func shortenFlowLabel(s string, max int) string {
	s = strings.TrimSpace(strings.ReplaceAll(s, "\n", " "))
	s = strings.Join(strings.Fields(s), " ")
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func flowKindLabel(kind string) string {
	switch kind {
	case "start":
		return "Start"
	case "end":
		return "End"
	case "decision":
		return "?"
	case "loop":
		return "Loop"
	case "return":
		return "return"
	default:
		return "step"
	}
}

func friendlyFlowLabel(fset *token.FileSet, stmt ast.Stmt) (string, int) {
	line := fset.Position(stmt.Pos()).Line
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		var names []string
		for _, lhs := range s.Lhs {
			names = append(names, flowExprNames(lhs)...)
		}
		if len(names) > 0 {
			return ":=" + strings.Join(names, ","), line
		}
		return "assign", line
	case *ast.IncDecStmt:
		if n := flowExprNames(s.X); len(n) > 0 {
			return n[0] + s.Tok.String(), line
		}
		return s.Tok.String(), line
	case *ast.ExprStmt:
		if call, ok := s.X.(*ast.CallExpr); ok {
			if id, ok := call.Fun.(*ast.Ident); ok {
				return id.Name + "()", line
			}
		}
		return "call", line
	case *ast.BranchStmt:
		return s.Tok.String(), line
	default:
		return "step", line
	}
}

func flowExprNames(e ast.Expr) []string {
	if e == nil {
		return nil
	}
	switch x := e.(type) {
	case *ast.Ident:
		return []string{x.Name}
	case *ast.SelectorExpr:
		if id, ok := x.X.(*ast.Ident); ok {
			return []string{id.Name + "." + x.Sel.Name}
		}
	}
	return nil
}

func shortenFlowExpr(fset *token.FileSet, e ast.Expr) string {
	if e == nil {
		return ""
	}
	switch x := e.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.BasicLit:
		return x.Value
	case *ast.BinaryExpr:
		return shortenFlowExpr(fset, x.X) + " " + x.Op.String() + " " + shortenFlowExpr(fset, x.Y)
	case *ast.UnaryExpr:
		return x.Op.String() + shortenFlowExpr(fset, x.X)
	case *ast.ParenExpr:
		return "(" + shortenFlowExpr(fset, x.X) + ")"
	case *ast.IndexExpr:
		return shortenFlowExpr(fset, x.X) + "[…]"
	case *ast.SelectorExpr:
		return shortenFlowExpr(fset, x.X) + "." + x.Sel.Name
	case *ast.CallExpr:
		if id, ok := x.Fun.(*ast.Ident); ok {
			return id.Name + "()"
		}
		return "call()"
	default:
		return "expr"
	}
}
