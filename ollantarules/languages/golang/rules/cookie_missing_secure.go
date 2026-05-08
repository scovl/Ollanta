package rules

import (
	"go/ast"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// CookieMissingSecure detects http.Cookie literals without Secure = true.
var CookieMissingSecure = ollantarules.Rule{
	MetaKey: "go:cookie-missing-secure",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		var issues []*domain.Issue
		ast.Inspect(ctx.AST, func(n ast.Node) bool {
			comp, ok := n.(*ast.CompositeLit)
			if !ok || !isCookieType(comp.Type) {
				return true
			}
			if hasField(comp, "Secure") {
				return true
			}
			line := ctx.FileSet.Position(comp.Pos()).Line
			issue := domain.NewIssue("go:cookie-missing-secure", ctx.Path, line)
			issue.Severity = domain.SeverityMajor
			issue.Type = domain.TypeVulnerability
			issue.Message = "http.Cookie should set Secure to true to ensure cookies are sent only over HTTPS"
			issues = append(issues, issue)
			return true
		})
		return issues
	},
}
