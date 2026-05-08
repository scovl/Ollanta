package rules

import (
	"go/ast"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// TemplateHTMLDoesNotEscape detects template.HTML calls which disable escaping.
var TemplateHTMLDoesNotEscape = ollantarules.Rule{
	MetaKey: "go:template-html-does-not-escape",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		var issues []*domain.Issue
		ast.Inspect(ctx.AST, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if !isTemplateHTMLCall(call) {
				return true
			}
			// Only flag if the argument is not a string literal (heuristic for dynamic input)
			if len(call.Args) == 0 {
				return true
			}
			if _, isLit := call.Args[0].(*ast.BasicLit); isLit {
				return true
			}
			line := ctx.FileSet.Position(call.Pos()).Line
			issue := domain.NewIssue("go:template-html-does-not-escape", ctx.Path, line)
			issue.Severity = domain.SeverityMajor
			issue.Type = domain.TypeVulnerability
			issue.Message = "template.HTML disables HTML escaping; ensure the input is sanitized"
			issues = append(issues, issue)
			return true
		})
		return issues
	},
}

func isTemplateHTMLCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return pkg.Name == "template" && sel.Sel.Name == "HTML"
}
