package rules

import (
	"go/ast"
	"strings"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// BadTmp detects creation of temporary files under /tmp using hardcoded
// paths instead of os.CreateTemp.
var BadTmp = ollantarules.Rule{
	MetaKey: "go:bad-tmp",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		var issues []*domain.Issue
		ast.Inspect(ctx.AST, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			fn, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			pkg, ok := fn.X.(*ast.Ident)
			if !ok {
				return true
			}
			// Detect os.Create, os.Open, os.WriteFile, ioutil.WriteFile with /tmp/ path
			if !isFileOp(pkg.Name, fn.Sel.Name) {
				return true
			}
			if len(call.Args) == 0 {
				return true
			}
			path := stringLiteral(call.Args[0])
			if strings.HasPrefix(path, "/tmp/") {
				line := ctx.FileSet.Position(call.Pos()).Line
				issue := domain.NewIssue("go:bad-tmp", ctx.Path, line)
				issue.Severity = domain.SeverityMajor
				issue.Type = domain.TypeVulnerability
				issue.Message = "Insecure temporary file: use os.CreateTemp instead of hardcoded /tmp/ paths"
				issues = append(issues, issue)
			}
			return true
		})
		return issues
	},
}

func isFileOp(pkg, sel string) bool {
	switch pkg {
	case "os":
		return sel == "Create" || sel == "Open" || sel == "OpenFile" || sel == "WriteFile"
	case "ioutil":
		return sel == "WriteFile"
	}
	return false
}

func stringLiteral(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return strings.Trim(e.Value, "`\"")
	}
	return ""
}
