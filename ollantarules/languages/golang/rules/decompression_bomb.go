package rules

import (
	"go/ast"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// DecompressionBomb detects gzip.NewReader and zlib.NewReader calls
// without size limits, which can lead to decompression bombs.
var DecompressionBomb = ollantarules.Rule{
	MetaKey: "go:decompression-bomb",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		var issues []*domain.Issue
		ast.Inspect(ctx.AST, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if isDecompressionCall(call) {
				line := ctx.FileSet.Position(call.Pos()).Line
				issue := domain.NewIssue("go:decompression-bomb", ctx.Path, line)
				issue.Severity = domain.SeverityMajor
				issue.Type = domain.TypeVulnerability
				issue.Message = "Decompression without size limits can lead to denial of service"
				issues = append(issues, issue)
			}
			return true
		})
		return issues
	},
}

func isDecompressionCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	switch pkg.Name {
	case "gzip":
		return sel.Sel.Name == "NewReader"
	case "zlib":
		return sel.Sel.Name == "NewReader"
	case "flate":
		return sel.Sel.Name == "NewReader"
	case "lzw":
		return sel.Sel.Name == "NewReader"
	}
	return false
}
