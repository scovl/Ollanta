package rules

import (
	"go/ast"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// MissingSSLMinVersion detects tls.Config literals without an explicit MinVersion field.
var MissingSSLMinVersion = ollantarules.Rule{
	MetaKey: "go:missing-ssl-minversion",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		var issues []*domain.Issue
		ast.Inspect(ctx.AST, func(n ast.Node) bool {
			comp, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}
			if !isTLSConfigType(comp.Type) {
				return true
			}
			if hasField(comp, "MinVersion") {
				return true
			}
			line := ctx.FileSet.Position(comp.Pos()).Line
			issue := domain.NewIssue("go:missing-ssl-minversion", ctx.Path, line)
			issue.Severity = domain.SeverityMajor
			issue.Type = domain.TypeVulnerability
			issue.Message = "tls.Config should set MinVersion to avoid using outdated TLS versions"
			issues = append(issues, issue)
			return true
		})
		return issues
	},
}

func isTLSConfigType(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return pkg.Name == "tls" && sel.Sel.Name == "Config"
}

func hasField(comp *ast.CompositeLit, name string) bool {
	for _, elt := range comp.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}
		if key.Name == name {
			return true
		}
	}
	return false
}
