package rules

import (
	"go/ast"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// CookieMissingHttponly detects http.Cookie literals without HttpOnly = true.
var CookieMissingHttponly = ollantarules.Rule{
	MetaKey: "go:cookie-missing-httponly",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		var issues []*domain.Issue
		ast.Inspect(ctx.AST, func(n ast.Node) bool {
			comp, ok := n.(*ast.CompositeLit)
			if !ok || !isCookieType(comp.Type) {
				return true
			}
			if hasField(comp, "HttpOnly") {
				return true
			}
			line := ctx.FileSet.Position(comp.Pos()).Line
			issue := domain.NewIssue("go:cookie-missing-httponly", ctx.Path, line)
			issue.Severity = domain.SeverityMajor
			issue.Type = domain.TypeVulnerability
			issue.Message = "http.Cookie should set HttpOnly to true to mitigate XSS attacks"
			issues = append(issues, issue)
			return true
		})
		return issues
	},
}

func isCookieType(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return pkg.Name == "http" && sel.Sel.Name == "Cookie"
}
