package rules

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// BindAll detects net.Listen calls that bind to all interfaces (0.0.0.0 or :::).
var BindAll = ollantarules.Rule{
	MetaKey: "go:bind-all",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		var issues []*domain.Issue
		ast.Inspect(ctx.AST, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if !isNetListenCall(ctx.AST, call) {
				return true
			}
			if len(call.Args) < 2 {
				return true
			}
			addrArg := call.Args[1]
			addr := extractStringLiteral(addrArg)
			if addr == "" {
				return true
			}
			if bindsToAll(addr) {
				line := ctx.FileSet.Position(call.Pos()).Line
				issue := domain.NewIssue("go:bind-all", ctx.Path, line)
				issue.Severity = domain.SeverityMajor
				issue.Type = domain.TypeVulnerability
				issue.Message = "Binding to all network interfaces (0.0.0.0 or ::) exposes the service on every interface"
				issues = append(issues, issue)
			}
			return true
		})
		return issues
	},
}

func isNetListenCall(f *ast.File, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	// Heuristic: net.Listen or Listen after net import
	if pkg.Name == "net" && sel.Sel.Name == "Listen" {
		return true
	}
	// Also catch http.ListenAndServe
	if pkg.Name == "http" && sel.Sel.Name == "ListenAndServe" {
		return true
	}
	return false
}

func extractStringLiteral(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.BasicLit:
		if v.Kind == token.STRING {
			s, err := strconv.Unquote(v.Value)
			if err == nil {
				return s
			}
		}
	case *ast.BinaryExpr:
		if v.Op == token.ADD {
			left := extractStringLiteral(v.X)
			right := extractStringLiteral(v.Y)
			if left != "" && right != "" {
				return left + right
			}
		}
	}
	return ""
}

func bindsToAll(addr string) bool {
	return strings.HasPrefix(addr, "0.0.0.0:") ||
		strings.HasPrefix(addr, "[::]") ||
		strings.HasPrefix(addr, ":")
}
