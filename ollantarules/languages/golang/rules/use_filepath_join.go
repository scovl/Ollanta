package rules

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// UseFilepathJoin detects string concatenation used to build file paths
// instead of filepath.Join or filepath.Join.
var UseFilepathJoin = ollantarules.Rule{
	MetaKey: "go:use-filepath-join",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		var issues []*domain.Issue
		ast.Inspect(ctx.AST, func(n ast.Node) bool {
			be, ok := n.(*ast.BinaryExpr)
			if !ok || be.Op != token.ADD {
				return true
			}
			if looksLikePath(be.X) || looksLikePath(be.Y) {
				line := ctx.FileSet.Position(be.Pos()).Line
				issue := domain.NewIssue("go:use-filepath-join", ctx.Path, line)
				issue.Severity = domain.SeverityMinor
				issue.Type = domain.TypeCodeSmell
				issue.Message = "Use filepath.Join instead of string concatenation for building paths"
				issues = append(issues, issue)
			}
			return true
		})
		return issues
	},
}

// looksLikePath reports whether an expression appears to be a path fragment.
func looksLikePath(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			s := strings.Trim(e.Value, "`\"")
			return strings.Contains(s, "/") || strings.Contains(s, "\\")
		}
	case *ast.BinaryExpr:
		return looksLikePath(e.X) || looksLikePath(e.Y)
	}
	return false
}
