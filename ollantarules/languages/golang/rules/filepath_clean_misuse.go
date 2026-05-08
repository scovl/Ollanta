package rules

import (
	"go/ast"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// FilepathCleanMisuse detects using filepath.Clean as the sole sanitization
// before opening a file, which does not prevent path traversal.
var FilepathCleanMisuse = ollantarules.Rule{
	MetaKey: "go:filepath-clean-misuse",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		var issues []*domain.Issue
		ast.Inspect(ctx.AST, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if !isFileOpenCall(call) {
				return true
			}
			for _, arg := range call.Args {
				if containsFilepathClean(arg) {
					line := ctx.FileSet.Position(call.Pos()).Line
					issue := domain.NewIssue("go:filepath-clean-misuse", ctx.Path, line)
					issue.Severity = domain.SeverityMajor
					issue.Type = domain.TypeVulnerability
					issue.Message = "filepath.Clean does not prevent path traversal; validate or restrict the path"
					issues = append(issues, issue)
					break
				}
			}
			return true
		})
		return issues
	},
}

func isFileOpenCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	if pkg.Name != "os" {
		return false
	}
	switch sel.Sel.Name {
	case "Open", "Create", "OpenFile":
		return true
	}
	return false
}

func containsFilepathClean(expr ast.Expr) bool {
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		pkg, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}
		if pkg.Name == "filepath" && sel.Sel.Name == "Clean" {
			found = true
			return false
		}
		return true
	})
	return found
}
