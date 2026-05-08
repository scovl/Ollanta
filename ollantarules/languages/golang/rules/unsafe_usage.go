package rules

import (
	"go/ast"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// UnsafeUsage detects use of unsafe.Pointer or unsafe.Slice.
var UnsafeUsage = ollantarules.Rule{
	MetaKey: "go:unsafe",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		var issues []*domain.Issue
		ast.Inspect(ctx.AST, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if ok && isUnsafeCall(call) {
				line := ctx.FileSet.Position(call.Pos()).Line
				issue := domain.NewIssue("go:unsafe", ctx.Path, line)
				issue.Severity = domain.SeverityMajor
				issue.Type = domain.TypeVulnerability
				issue.Message = "Use of unsafe package bypasses Go's type safety; verify necessity and correctness"
				issues = append(issues, issue)
			}
			return true
		})
		return issues
	},
}

func isUnsafeCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	if pkg.Name != "unsafe" {
		return false
	}
	switch sel.Sel.Name {
	case "Pointer", "Slice", "SliceData", "String", "StringData", "Alignof", "Offsetof", "Sizeof":
		return true
	}
	return false
}
