package analyzer

import (
	"fmt"
	"strings"

	"visualisator/internal/models"
)

// NarrateContext carries per-session narration state across frames.
type NarrateContext struct {
	LoopCounts map[string]map[int]int // stmt line -> visit count
	FrameIndex int
}

// NewNarrateContext creates narration state.
func NewNarrateContext() *NarrateContext {
	return &NarrateContext{LoopCounts: make(map[string]map[int]int)}
}

// NarrateStep builds explanation for the current debugger frame.
func NarrateStep(
	lang Lang,
	sm *SourceMap,
	nc *NarrateContext,
	editorLine int,
	function string,
	prevVars, currVars []models.Variable,
	frameIndex int,
) *FrameExplanation {
	if nc != nil {
		nc.FrameIndex = frameIndex
	}

	changes := ComputeVarDeltas(prevVars, currVars)
	lookup := VarLookup(currVars)

	lineText := ""
	stmt := StatementInfo{Kind: "other"}
	if sm != nil {
		if t := lineTextAt(sm.Source, editorLine); t != "" {
			lineText = t
		}
		if info, ok := sm.StmtAt(editorLine); ok {
			stmt = info
			if stmt.Text != "" {
				lineText = stmt.Text
			}
		}
	}

	exp := &FrameExplanation{
		LineText: lineText,
		StmtKind: stmt.Kind,
		Changes:  changes,
	}

	loopIter := 0
	if nc != nil && (stmt.Kind == "for" || stmt.Kind == "range") {
		key := stmt.Kind
		if nc.LoopCounts[key] == nil {
			nc.LoopCounts[key] = make(map[int]int)
		}
		nc.LoopCounts[key][stmt.Line]++
		loopIter = nc.LoopCounts[key][stmt.Line]
		exp.LoopIter = loopIter
	}

	if stmt.ConditionExpr != "" {
		if res, detail, ok := EvalCondition(stmt.ConditionExpr, lookup); ok {
			exp.Condition = FormatConditionResult(stmt.ConditionExpr, res, detail, lang)
		}
	}

	exp.Summary = buildSummary(lang, sm, stmt, function, changes, exp.Condition, loopIter, frameIndex, lookup)

	return exp
}

// NarrateStepShort returns a one-line summary for step history lists.
func NarrateStepShort(exp *FrameExplanation) string {
	if exp == nil {
		return ""
	}
	if exp.Summary != "" {
		s := exp.Summary
		if len(s) > 72 {
			return s[:69] + "…"
		}
		return s
	}
	return exp.LineText
}

func lineTextAt(source string, line int) string {
	lines := sourceLines(source)
	return lineText(lines, line)
}

func buildSummary(
	lang Lang,
	sm *SourceMap,
	stmt StatementInfo,
	function string,
	changes []VarDelta,
	condition string,
	loopIter int,
	frameIndex int,
	lookup map[string]string,
) string {
	fn := shortFuncName(function)

	if frameIndex == 0 && fn != "" {
		args := formatArgSnapshot(lookup)
		if lang == LangRU {
			if args != "" {
				return fmt.Sprintf("Вход в %s(%s)", fn, args)
			}
			return fmt.Sprintf("Вход в %s", fn)
		}
		if args != "" {
			return fmt.Sprintf("Enter %s(%s)", fn, args)
		}
		return fmt.Sprintf("Enter %s", fn)
	}

	switch stmt.Kind {
	case "assign":
		if len(changes) > 0 {
			c := changes[0]
			if lang == LangRU {
				return fmt.Sprintf("Присваиваем %s: %s → %s", c.Name, c.OldValue, c.NewValue)
			}
			return fmt.Sprintf("Assign %s: %s → %s", c.Name, c.OldValue, c.NewValue)
		}
		if lang == LangRU {
			return "Присваивание: " + stmt.Text
		}
		return "Assignment: " + stmt.Text

	case "inc":
		if len(changes) > 0 {
			c := changes[0]
			if lang == LangRU {
				return fmt.Sprintf("%s: %s → %s", c.Name, c.OldValue, c.NewValue)
			}
			return fmt.Sprintf("%s: %s → %s", c.Name, c.OldValue, c.NewValue)
		}
		if lang == LangRU {
			return "Изменение счётчика: " + stmt.Text
		}
		return "Increment/decrement: " + stmt.Text

	case "if":
		if condition != "" {
			if lang == LangRU {
				return "Проверяем условие: " + condition
			}
			return "Check condition: " + condition
		}
		if lang == LangRU {
			return "Условие if: " + stmt.Text
		}
		return "If statement: " + stmt.Text

	case "for":
		if loopIter > 0 {
			if lang == LangRU {
				return fmt.Sprintf("Итерация цикла #%d — %s", loopIter, stmt.Text)
			}
			return fmt.Sprintf("Loop iteration #%d — %s", loopIter, stmt.Text)
		}
		if lang == LangRU {
			return "Цикл: " + stmt.Text
		}
		return "Loop: " + stmt.Text

	case "range":
		if loopIter > 0 {
			if lang == LangRU {
				return fmt.Sprintf("Итерация range #%d — %s", loopIter, stmt.Text)
			}
			return fmt.Sprintf("Range iteration #%d — %s", loopIter, stmt.Text)
		}
		if lang == LangRU {
			return "Range: " + stmt.Text
		}
		return "Range: " + stmt.Text

	case "return":
		retVal := ""
		if len(changes) > 0 {
			retVal = changes[len(changes)-1].NewValue
		}
		for _, c := range changes {
			if strings.Contains(c.Name, "~r") || c.Name == "(return)" {
				retVal = c.NewValue
			}
		}
		if retVal == "" {
			for _, v := range lookup {
				if v != "" && v != "—" {
					retVal = v
					break
				}
			}
		}
		if lang == LangRU {
			if retVal != "" {
				return "Возврат из функции: " + retVal
			}
			return "Возврат из функции"
		}
		if retVal != "" {
			return "Return from function: " + retVal
		}
		return "Return from function"

	case "call":
		if lang == LangRU {
			if len(changes) > 0 {
				return "Вызов: " + stmt.Text + " (данные изменились)"
			}
			return "Вызов: " + stmt.Text
		}
		if len(changes) > 0 {
			return "Call: " + stmt.Text + " (state changed)"
		}
		return "Call: " + stmt.Text

	default:
		if len(changes) > 0 {
			if lang == LangRU {
				return fmt.Sprintf("Шаг на строке %d — изменено %d переменных", stmt.Line, len(changes))
			}
			return fmt.Sprintf("Step at line %d — %d variable(s) changed", stmt.Line, len(changes))
		}
		if stmt.Text != "" {
			if lang == LangRU {
				return "Выполняется: " + stmt.Text
			}
			return "Executing: " + stmt.Text
		}
		if lang == LangRU {
			return fmt.Sprintf("Шаг %d", frameIndex+1)
		}
		return fmt.Sprintf("Step %d", frameIndex+1)
	}
}

func shortFuncName(full string) string {
	if full == "" {
		return ""
	}
	if i := strings.LastIndex(full, "."); i >= 0 {
		return full[i+1:]
	}
	return full
}

func formatArgSnapshot(lookup map[string]string) string {
	// Common LeetCode param names first
	order := []string{"nums", "target", "s", "head", "root", "n", "k", "matrix", "grid"}
	var parts []string
	seen := make(map[string]bool)
	for _, name := range order {
		if v, ok := lookup[name]; ok && v != "—" {
			parts = append(parts, name+"="+v)
			seen[name] = true
		}
	}
	for name, v := range lookup {
		if seen[name] || strings.Contains(name, ".") {
			continue
		}
		if v == "—" {
			continue
		}
		parts = append(parts, name+"="+v)
		if len(parts) >= 4 {
			break
		}
	}
	return strings.Join(parts, ", ")
}

// OutlineSummary returns a static description before run.
func OutlineSummary(lang Lang, o FunctionOutline) string {
	if o.Name == "" {
		if lang == LangRU {
			return "Вставьте функцию и нажмите «Запуск» для пошагового разбора."
		}
		return "Paste a function and press Run for a step-by-step breakdown."
	}
	if lang == LangRU {
		s := fmt.Sprintf("Функция %s(%s)", o.Name, o.Params)
		if o.Returns != "" {
			s += " → " + o.Returns
		}
		if o.DemoCall != "" {
			s += ". Демо-запуск: " + o.DemoCall
		}
		return s
	}
	s := fmt.Sprintf("Function %s(%s)", o.Name, o.Params)
	if o.Returns != "" {
		s += " → " + o.Returns
	}
	if o.DemoCall != "" {
		s += ". Demo run: " + o.DemoCall
	}
	return s
}
