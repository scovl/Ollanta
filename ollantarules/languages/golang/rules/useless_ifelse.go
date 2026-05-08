package rules

import (
	"go/ast"
	"strconv"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// UselessIfElse detects if statements with constant true/false conditions
// which are dead code or indicate a logic error.
var UselessIfElse = ollantarules.Rule{
	MetaKey: "go:useless-ifelse",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		var issues []*domain.Issue
		ast.Inspect(ctx.AST, func(n ast.Node) bool {
			ifStmt, ok := n.(*ast.IfStmt)
			if !ok {
				return true
			}
			if isConstBool(ifStmt.Cond, true) || isConstBool(ifStmt.Cond, false) {
				line := ctx.FileSet.Position(ifStmt.Pos()).Line
				issue := domain.NewIssue("go:useless-ifelse", ctx.Path, line)
				issue.Severity = domain.SeverityMinor
				issue.Type = domain.TypeCodeSmell
				issue.Message = "If statement with a constant boolean condition is dead code"
				issues = append(issues, issue)
			}
			return true
		})
		return issues
	},
}

// isConstBool reports whether expr is the boolean literal val.
func isConstBool(expr ast.Expr, val bool) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		if val {
			return e.Name == "true"
		}
		return e.Name == "false"
	case *ast.BasicLit:
		b, err := strconv.ParseBool(e.Value)
		return err == nil && b == val
	}
	return false
}
