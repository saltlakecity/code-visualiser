package testrunner

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"visualisator/internal/models"
)

var (
	pairRe  = regexp.MustCompile(`^\s*([a-zA-Z_]\w*)\s*=\s*(.+)$`)
	arrayRe = regexp.MustCompile(`^\[(.*)\]$`)
)

// ResolveArgs returns Go literal args for a test case.
func ResolveArgs(tc models.TestCase) (string, error) {
	if strings.TrimSpace(tc.Args) != "" {
		return strings.TrimSpace(tc.Args), nil
	}
	if strings.TrimSpace(tc.Input) == "" {
		return "", nil
	}
	return ParseInput(tc.Input)
}

// ParseInput converts LeetCode-style input to Go argument literals.
func ParseInput(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", nil
	}
	parts := splitInputParts(input)
	var args []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		m := pairRe.FindStringSubmatch(part)
		if len(m) != 3 {
			return "", fmt.Errorf("invalid input part: %q", part)
		}
		key := strings.ToLower(m[1])
		val := strings.TrimSpace(m[2])
		switch {
		case key == "nums" || key == "arr" || key == "array" || strings.HasPrefix(key, "nums"):
			lit, err := arrayToGoIntSlice(val)
			if err != nil {
				return "", err
			}
			args = append(args, lit)
		case key == "target" || key == "k" || key == "n" || key == "x" || key == "val":
			args = append(args, val)
		case key == "list" || key == "head":
			args = append(args, buildListNodeArg(val))
		case key == "tree" || key == "root":
			args = append(args, buildTreeNodeArg(val))
		case strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]"):
			lit, err := arrayToGoIntSlice(val)
			if err != nil {
				return "", err
			}
			args = append(args, lit)
		default:
			args = append(args, val)
		}
	}
	return strings.Join(args, ", "), nil
}

func splitInputParts(input string) []string {
	var parts []string
	var cur strings.Builder
	depth := 0
	for _, r := range input {
		switch r {
		case '[', '(', '{':
			depth++
			cur.WriteRune(r)
		case ']', ')', '}':
			depth--
			cur.WriteRune(r)
		case ',':
			if depth == 0 {
				parts = append(parts, cur.String())
				cur.Reset()
				continue
			}
			cur.WriteRune(r)
		default:
			cur.WriteRune(r)
		}
	}
	if cur.Len() > 0 {
		parts = append(parts, cur.String())
	}
	return parts
}

func arrayToGoIntSlice(val string) (string, error) {
	val = strings.TrimSpace(val)
	if !strings.HasPrefix(val, "[") {
		return val, nil
	}
	inner := strings.Trim(val, "[]")
	if inner == "" {
		return "[]int{}", nil
	}
	items := strings.Split(inner, ",")
	var elems []string
	for _, it := range items {
		it = strings.TrimSpace(it)
		if it == "" || strings.EqualFold(it, "null") {
			continue
		}
		if _, err := strconv.Atoi(it); err != nil {
			return "", fmt.Errorf("expected int in array, got %q", it)
		}
		elems = append(elems, it)
	}
	return fmt.Sprintf("[]int{%s}", strings.Join(elems, ", ")), nil
}

func buildListNodeArg(val string) string {
	val = strings.TrimSpace(val)
	if !strings.HasPrefix(val, "[") {
		return val
	}
	inner := strings.Trim(strings.TrimSpace(val), "[]")
	if inner == "" {
		return "nil"
	}
	items := strings.Split(inner, ",")
	var nodes []string
	for i, it := range items {
		it = strings.TrimSpace(it)
		if it == "" || strings.EqualFold(it, "null") {
			continue
		}
		if i < len(items)-1 {
			nodes = append(nodes, fmt.Sprintf("&ListNode{Val: %s, Next: ", it))
		} else {
			nodes = append(nodes, fmt.Sprintf("&ListNode{Val: %s}", it))
		}
	}
	if len(nodes) == 0 {
		return "nil"
	}
	chain := nodes[0]
	for i := 1; i < len(nodes); i++ {
		chain += nodes[i]
	}
	for i := 0; i < len(nodes)-1; i++ {
		chain += "}"
	}
	return chain
}

func buildTreeNodeArg(val string) string {
	val = strings.TrimSpace(val)
	if !strings.HasPrefix(val, "[") {
		return val
	}
	items := strings.Split(strings.Trim(val, "[]"), ",")
	for i, it := range items {
		items[i] = strings.TrimSpace(it)
	}
	if len(items) == 0 || items[0] == "" || strings.EqualFold(items[0], "null") {
		return "nil"
	}
	return fmt.Sprintf("buildTree([]interface{}{%s})", quoteTreeItems(items))
}

func quoteTreeItems(items []string) string {
	var parts []string
	for _, it := range items {
		it = strings.TrimSpace(it)
		if it == "" || strings.EqualFold(it, "null") {
			parts = append(parts, "nil")
		} else if _, err := strconv.Atoi(it); err == nil {
			parts = append(parts, it)
		} else {
			parts = append(parts, fmt.Sprintf("%q", it))
		}
	}
	return strings.Join(parts, ", ")
}

// BuildMain generates main.go for a test case.
func BuildMain(entryFn string, tc models.TestCase, fullMain bool) (string, error) {
	if fullMain {
		return "", nil
	}
	entryFn = strings.TrimSpace(entryFn)
	if entryFn == "" {
		entryFn = "Solution"
	}
	args, err := ResolveArgs(tc)
	if err != nil {
		return "", err
	}
	setup := strings.TrimSpace(tc.Setup)
	var buf strings.Builder
	buf.WriteString("package main\n\nimport \"fmt\"\n\n")
	if setup != "" {
		buf.WriteString(setup)
		buf.WriteString("\n\n")
	}
	if strings.Contains(strings.ToLower(setup+tc.Input), "listnode") ||
		strings.Contains(strings.ToLower(setup), "listnode") {
		buf.WriteString(listNodeHelpers())
		buf.WriteString("\n")
	}
	if strings.Contains(strings.ToLower(setup+tc.Input), "treenode") ||
		strings.Contains(strings.ToLower(setup), "treenode") {
		buf.WriteString(treeNodeHelpers())
		buf.WriteString("\n")
	}
	buf.WriteString("func main() {\n")
	if args != "" {
		buf.WriteString(fmt.Sprintf("\tfmt.Println(%s(%s))\n", entryFn, args))
	} else {
		buf.WriteString(fmt.Sprintf("\tfmt.Println(%s())\n", entryFn))
	}
	buf.WriteString("}\n")
	return buf.String(), nil
}

func listNodeHelpers() string {
	return `// helpers for list tests
func intSliceToList(vals []int) *ListNode {
	if len(vals) == 0 {
		return nil
	}
	head := &ListNode{Val: vals[0]}
	cur := head
	for _, v := range vals[1:] {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return head
}`
}

func treeNodeHelpers() string {
	return `// helpers for tree tests
func buildTree(items []interface{}) *TreeNode {
	if len(items) == 0 || items[0] == nil {
		return nil
	}
	rootVal, ok := items[0].(int)
	if !ok {
		if f, ok := items[0].(float64); ok {
			rootVal = int(f)
		} else {
			return nil
		}
	}
	root := &TreeNode{Val: rootVal}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(items) {
		node := queue[0]
		queue = queue[1:]
		if i < len(items) && items[i] != nil {
			if v, ok := items[i].(int); ok {
				node.Left = &TreeNode{Val: v}
				queue = append(queue, node.Left)
			} else if f, ok := items[i].(float64); ok {
				node.Left = &TreeNode{Val: int(f)}
				queue = append(queue, node.Left)
			}
		}
		i++
		if i < len(items) && items[i] != nil {
			if v, ok := items[i].(int); ok {
				node.Right = &TreeNode{Val: v}
				queue = append(queue, node.Right)
			} else if f, ok := items[i].(float64); ok {
				node.Right = &TreeNode{Val: int(f)}
				queue = append(queue, node.Right)
			}
		}
		i++
	}
	return root
}`
}

// CompareOutput checks if actual stdout matches expected.
func CompareOutput(actual, expected string) bool {
	return normalizeOutput(actual) == normalizeOutput(expected)
}

func normalizeOutput(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Join(strings.Fields(s), " ")
	// normalize slice formatting: [0 1] vs [0,1]
	s = strings.ReplaceAll(s, ", ", " ")
	s = strings.ReplaceAll(s, ",", " ")
	s = strings.Join(strings.Fields(s), " ")
	return s
}

// RunBinary executes compiled binary and returns stdout.
func RunBinary(binaryPath string) (string, error) {
	return runBinary(binaryPath)
}
