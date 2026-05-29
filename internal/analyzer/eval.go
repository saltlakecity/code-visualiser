package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

// EvalCondition evaluates a boolean Go expression using variable values from the debugger.
func EvalCondition(expr string, vars map[string]string) (result bool, detail string, ok bool) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return false, "", false
	}

	node, err := parser.ParseExpr(expr)
	if err != nil {
		return false, "", false
	}

	val, detail, err := evalExpr(node, vars)
	if err != nil {
		return false, detail, false
	}
	b, ok := val.(bool)
	if !ok {
		return false, detail, false
	}
	return b, detail, true
}

func evalExpr(e ast.Expr, vars map[string]string) (interface{}, string, error) {
	switch n := e.(type) {
	case *ast.BinaryExpr:
		left, ld, err := evalExpr(n.X, vars)
		if err != nil {
			return nil, "", err
		}
		right, rd, err := evalExpr(n.Y, vars)
		if err != nil {
			return nil, "", err
		}
		detail := ld + " " + n.Op.String() + " " + rd
		switch n.Op {
		case token.LSS, token.LEQ, token.GTR, token.GEQ, token.EQL, token.NEQ:
			return compareOp(left, right, n.Op), detail, nil
		case token.LAND:
			lb, _ := left.(bool)
			rb, _ := right.(bool)
			return lb && rb, detail, nil
		case token.LOR:
			lb, _ := left.(bool)
			rb, _ := right.(bool)
			return lb || rb, detail, nil
		case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
			li, lerr := toInt64(left)
			ri, rerr := toInt64(right)
			if lerr != nil || rerr != nil {
				return nil, detail, fmt.Errorf("non-int arithmetic")
			}
			switch n.Op {
			case token.ADD:
				return li + ri, detail, nil
			case token.SUB:
				return li - ri, detail, nil
			case token.MUL:
				return li * ri, detail, nil
			case token.QUO:
				if ri == 0 {
					return nil, detail, fmt.Errorf("div zero")
				}
				return li / ri, detail, nil
			case token.REM:
				if ri == 0 {
					return nil, detail, fmt.Errorf("mod zero")
				}
				return li % ri, detail, nil
			}
		}
		return nil, detail, fmt.Errorf("unsupported binary op")
	case *ast.UnaryExpr:
		if n.Op == token.NOT {
			v, d, err := evalExpr(n.X, vars)
			if err != nil {
				return nil, "", err
			}
			b, ok := v.(bool)
			if !ok {
				return nil, d, fmt.Errorf("not bool")
			}
			return !b, "!" + d, nil
		}
		if n.Op == token.SUB {
			v, d, err := evalExpr(n.X, vars)
			if err != nil {
				return nil, "", err
			}
			i, err := toInt64(v)
			if err != nil {
				return nil, d, err
			}
			return -i, "-" + d, nil
		}
	case *ast.ParenExpr:
		return evalExpr(n.X, vars)
	case *ast.Ident:
		if n.Name == "true" {
			return true, "true", nil
		}
		if n.Name == "false" {
			return false, "false", nil
		}
		if n.Name == "nil" {
			return nil, "nil", nil
		}
		raw, ok := vars[n.Name]
		if !ok {
			return nil, n.Name, fmt.Errorf("unknown %s", n.Name)
		}
		return parseLiteralOrVar(raw, n.Name)
	case *ast.BasicLit:
		switch n.Kind {
		case token.INT:
			i, err := strconv.ParseInt(n.Value, 0, 64)
			return i, n.Value, err
		case token.FLOAT:
			f, err := strconv.ParseFloat(n.Value, 64)
			return f, n.Value, err
		case token.STRING:
			s, err := strconv.Unquote(n.Value)
			return s, n.Value, err
		}
	case *ast.CallExpr:
		if ident, ok := n.Fun.(*ast.Ident); ok && ident.Name == "len" && len(n.Args) == 1 {
			v, d, err := evalExpr(n.Args[0], vars)
			if err != nil {
				return nil, "", err
			}
			switch x := v.(type) {
			case int64:
				return x, "len(" + d + ")", nil
			case string:
				return int64(len(x)), "len(" + d + ")", nil
			default:
				if s, ok := vars[formatExprArg(n.Args[0])]; ok {
					if n := sliceLenFromDisplay(s); n >= 0 {
						return int64(n), "len(" + d + ")", nil
					}
				}
			}
		}
	case *ast.IndexExpr:
		base, bd, err := evalExpr(n.X, vars)
		if err != nil {
			// try display string for slices
			if name, ok := n.X.(*ast.Ident); ok {
				raw := vars[name.Name]
				idx, id, ierr := evalExpr(n.Index, vars)
				if ierr == nil {
					ii, ierr2 := toInt64(idx)
					if ierr2 == nil {
						if v, ok := indexInDisplay(raw, int(ii)); ok {
							return v, bd + "[" + id + "]", nil
						}
					}
				}
			}
			return nil, "", err
		}
		idx, id, err := evalExpr(n.Index, vars)
		if err != nil {
			return nil, "", err
		}
		ii, err := toInt64(idx)
		if err != nil {
			return nil, bd + "[" + id + "]", err
		}
		switch arr := base.(type) {
		case []interface{}:
			if int(ii) >= 0 && int(ii) < len(arr) {
				return arr[int(ii)], bd + "[" + id + "]", nil
			}
		case string:
			if int(ii) >= 0 && int(ii) < len(arr) {
				return int64(arr[int(ii)]), bd + "[" + id + "]", nil
			}
		}
		if name, ok := n.X.(*ast.Ident); ok {
			raw := vars[name.Name]
			if v, ok := indexInDisplay(raw, int(ii)); ok {
				return v, name.Name + "[" + id + "]", nil
			}
		}
	case *ast.SelectorExpr:
		// e.g. len on field — fallback to display
		if x, ok := n.X.(*ast.Ident); ok {
			key := x.Name + "." + n.Sel.Name
			if raw, ok := vars[key]; ok {
				return parseLiteralOrVar(raw, key)
			}
		}
	}
	return nil, "", fmt.Errorf("unsupported expr")
}

func formatExprArg(e ast.Expr) string {
	switch x := e.(type) {
	case *ast.Ident:
		return x.Name
	default:
		return ""
	}
}

func compareOp(left, right interface{}, op token.Token) bool {
	// string equality
	ls, lok := left.(string)
	rs, rok := right.(string)
	if lok && rok {
		switch op {
		case token.EQL:
			return ls == rs
		case token.NEQ:
			return ls != rs
		}
	}
	li, lerr := toInt64(left)
	ri, rerr := toInt64(right)
	if lerr == nil && rerr == nil {
		switch op {
		case token.LSS:
			return li < ri
		case token.LEQ:
			return li <= ri
		case token.GTR:
			return li > ri
		case token.GEQ:
			return li >= ri
		case token.EQL:
			return li == ri
		case token.NEQ:
			return li != ri
		}
	}
	return false
}

func toInt64(v interface{}) (int64, error) {
	switch x := v.(type) {
	case int64:
		return x, nil
	case int:
		return int64(x), nil
	case float64:
		return int64(x), nil
	case bool:
		if x {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("not int")
	}
}

func parseLiteralOrVar(raw, name string) (interface{}, string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "true" {
		return true, name, nil
	}
	if raw == "false" {
		return false, name, nil
	}
	if raw == "nil" {
		return nil, name, nil
	}
	if strings.HasPrefix(raw, "\"") {
		s, err := strconv.Unquote(raw)
		return s, name, err
	}
	if strings.HasPrefix(raw, "'") {
		return raw, name, nil
	}
	if i, err := strconv.ParseInt(raw, 10, 64); err == nil {
		return i, name, nil
	}
	if strings.HasPrefix(raw, "[") {
		return raw, name, nil
	}
	return raw, name, nil
}

func sliceLenFromDisplay(s string) int {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "[") {
		return -1
	}
	inner := strings.TrimPrefix(s, "[")
	inner = strings.TrimSuffix(inner, "]")
	if inner == "" || inner == "…" {
		return 0
	}
	parts := strings.Split(inner, ",")
	n := 0
	for _, p := range parts {
		if strings.TrimSpace(p) == "…" {
			break
		}
		n++
	}
	return n
}

func indexInDisplay(slice string, i int) (interface{}, bool) {
	slice = strings.TrimSpace(slice)
	if !strings.HasPrefix(slice, "[") {
		return nil, false
	}
	inner := strings.TrimPrefix(slice, "[")
	inner = strings.TrimSuffix(inner, "]")
	parts := strings.Split(inner, ",")
	if i < 0 || i >= len(parts) {
		return nil, false
	}
	p := strings.TrimSpace(parts[i])
	if p == "…" {
		return nil, false
	}
	if iv, err := strconv.ParseInt(p, 10, 64); err == nil {
		return iv, true
	}
	if strings.HasPrefix(p, "\"") {
		s, err := strconv.Unquote(p)
		if err == nil {
			return s, true
		}
	}
	return p, true
}

// FormatConditionResult builds display text for UI.
func FormatConditionResult(expr string, result bool, detail string, lang Lang) string {
	truth := "true"
	falsy := "false"
	if lang == LangRU {
		truth = "истина"
		falsy = "ложь"
	}
	res := truth
	if !result {
		res = falsy
	}
	if detail != "" {
		return expr + " → " + res + " (" + detail + ")"
	}
	return expr + " → " + res
}
