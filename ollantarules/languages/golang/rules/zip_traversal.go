package rules

import (
	"go/ast"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// ZipTraversal detects archive/zip Open calls where the filename may contain ".."
// and is not sanitized before use.
var ZipTraversal = ollantarules.Rule{
	MetaKey: "go:zip",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		var issues []*domain.Issue
		ast.Inspect(ctx.AST, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if !isZipOpenCall(call) {
				return true
			}
			line := ctx.FileSet.Position(call.Pos()).Line
			issue := domain.NewIssue("go:zip", ctx.Path, line)
			issue.Severity = domain.SeverityMajor
			issue.Type = domain.TypeVulnerability
			issue.Message = "Zip file entries may contain path traversal sequences (../); validate filenames before extraction"
			issues = append(issues, issue)
			return true
		})
		return issues
	},
}

func isZipOpenCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return pkg.Name == "zip" && sel.Sel.Name == "OpenReader"
}
