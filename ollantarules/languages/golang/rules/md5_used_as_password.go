package rules

import (
	"go/ast"
	"strings"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// MD5UsedAsPassword detects use of MD5 hashing in contexts that may relate
// to password storage. MD5 is cryptographically broken and unsuitable for
// password hashing.
var MD5UsedAsPassword = ollantarules.Rule{
	MetaKey: "go:md5-used-as-password",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		var issues []*domain.Issue
		var hasCryptoMD5 bool
		for _, imp := range ctx.AST.Imports {
			path := strings.Trim(imp.Path.Value, "`\"")
			if path == "crypto/md5" {
				hasCryptoMD5 = true
				break
			}
		}
		if !hasCryptoMD5 {
			return nil
		}
		ast.Inspect(ctx.AST, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if isMD5Call(call) {
				line := ctx.FileSet.Position(call.Pos()).Line
				issue := domain.NewIssue("go:md5-used-as-password", ctx.Path, line)
				issue.Severity = domain.SeverityMajor
				issue.Type = domain.TypeVulnerability
				issue.Message = "MD5 is cryptographically broken and must not be used for password hashing; use bcrypt, scrypt, or Argon2"
				issues = append(issues, issue)
			}
			return true
		})
		return issues
	},
}

func isMD5Call(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	if pkg.Name != "md5" {
		return false
	}
	return sel.Sel.Name == "Sum" || sel.Sel.Name == "New"
}
