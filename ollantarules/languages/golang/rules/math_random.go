package rules

import (
	"go/ast"
	"strings"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// MathRandom detects use of math/rand for security-sensitive operations.
// math/rand is not cryptographically secure; use crypto/rand instead.
var MathRandom = ollantarules.Rule{
	MetaKey: "go:math-random",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		var issues []*domain.Issue
		var hasMathRand bool
		// Check imports
		for _, imp := range ctx.AST.Imports {
			path := strings.Trim(imp.Path.Value, "`\"")
			if path == "math/rand" {
				hasMathRand = true
				break
			}
		}
		if !hasMathRand {
			return nil
		}
		ast.Inspect(ctx.AST, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if isMathRandCall(call) {
				line := ctx.FileSet.Position(call.Pos()).Line
				issue := domain.NewIssue("go:math-random", ctx.Path, line)
				issue.Severity = domain.SeverityMajor
				issue.Type = domain.TypeVulnerability
				issue.Message = "math/rand is not cryptographically secure; use crypto/rand for tokens, passwords, or IDs"
				issues = append(issues, issue)
			}
			return true
		})
		return issues
	},
}

func isMathRandCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	if pkg.Name != "rand" {
		return false
	}
	switch sel.Sel.Name {
	case "Intn", "Float64", "Int", "Int63", "Uint32", "Uint64", "ExpFloat64", "NormFloat64":
		return true
	}
	return false
}
