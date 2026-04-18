package rules

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// NoLargeFunctions detects Go functions and methods that exceed a configurable
// line count threshold. Long functions are harder to read, test and maintain.
var NoLargeFunctions = ollantarules.Rule{
	MetaKey: "go:no-large-functions",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		maxLines := paramInt(ctx.Params, "max_lines", 40)
		var issues []*domain.Issue

		ast.Inspect(ctx.AST, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				return true
			}
			start := ctx.FileSet.Position(fn.Pos()).Line
			end := ctx.FileSet.Position(fn.Body.End()).Line
			lines := end - start + 1
			if lines > maxLines {
				name := fn.Name.Name
				issue := domain.NewIssue("go:no-large-functions", ctx.Path, start)
				issue.EndLine = end
				issue.Severity = domain.SeverityMajor
				issue.Type = domain.TypeCodeSmell
				issue.Message = fmt.Sprintf("Function '%s' has %d lines (max: %d)", name, lines, maxLines)
				issues = append(issues, issue)
			}
			return true
		})
		return issues
	},
}

// paramInt reads an int param from ctx.Params, falling back to defaultVal.
func paramInt(params map[string]string, key string, defaultVal int) int {
	if v, ok := params[key]; ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return defaultVal
}

// lineOf returns the 1-based line number of a token.Pos.
func lineOf(fset *token.FileSet, pos token.Pos) int {
	return fset.Position(pos).Line
}
