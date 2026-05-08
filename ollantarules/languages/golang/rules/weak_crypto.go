package rules

import (
	"go/ast"
	"strings"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// WeakCrypto detects use of weak cryptographic algorithms (DES, RC4, 3DES).
var WeakCrypto = ollantarules.Rule{
	MetaKey: "go:weak-crypto",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		var issues []*domain.Issue
		weakImports := map[string]bool{
			"crypto/des": true,
			"crypto/rc4": true,
		}
		weakFuncs := map[string]bool{
			"NewCipher": true, // des.NewCipher
			"New":       true, // rc4.New
		}

		ast.Inspect(ctx.AST, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.ImportSpec:
				path := importPath(node)
				if weakImports[path] {
					line := ctx.FileSet.Position(node.Pos()).Line
					issue := domain.NewIssue("go:weak-crypto", ctx.Path, line)
					issue.Severity = domain.SeverityMajor
					issue.Type = domain.TypeVulnerability
					issue.Message = "Weak cryptographic algorithm imported: " + path
					issues = append(issues, issue)
				}
			case *ast.CallExpr:
				if isWeakCryptoCall(node, weakFuncs) {
					line := ctx.FileSet.Position(node.Pos()).Line
					issue := domain.NewIssue("go:weak-crypto", ctx.Path, line)
					issue.Severity = domain.SeverityMajor
					issue.Type = domain.TypeVulnerability
					issue.Message = "Use of weak cryptographic algorithm detected"
					issues = append(issues, issue)
				}
			}
			return true
		})
		return issues
	},
}

func importPath(spec *ast.ImportSpec) string {
	if spec.Path == nil {
		return ""
	}
	path := strings.Trim(spec.Path.Value, `"`)
	return path
}

func isWeakCryptoCall(call *ast.CallExpr, weakFuncs map[string]bool) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if !weakFuncs[sel.Sel.Name] {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return pkg.Name == "des" || pkg.Name == "rc4"
}
