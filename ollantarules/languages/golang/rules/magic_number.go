package rules

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// MagicNumber flags numeric literals used directly in code outside of constant
// declarations, variable initialisations, and a small set of common neutral values.
// Using named constants instead improves readability and maintainability.
// SonarQube equivalent: squid:S109.

// defaultAuthorized are the numeric literals considered neutral and not flagged by default.
var defaultAuthorized = map[string]bool{"0": true, "1": true, "2": true, "-1": true}

// MagicNumber flags numeric literals that should be extracted into named constants.
var MagicNumber = ollantarules.Rule{
	MetaKey: "go:magic-number",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		var issues []*domain.Issue

		ast.Inspect(ctx.AST, func(n ast.Node) bool {
			switch decl := n.(type) {
			case *ast.GenDecl:
				if decl.Tok == token.CONST || decl.Tok == token.VAR {
					return false // skip — literal in const/var is fine
				}
			case *ast.BasicLit:
				if decl.Kind != token.INT && decl.Kind != token.FLOAT {
					return true
				}
				val := decl.Value
				if defaultAuthorized[val] {
					return true
				}
				line := ctx.FileSet.Position(decl.Pos()).Line
				issue := domain.NewIssue("go:magic-number", ctx.Path, line)
				issue.Severity = domain.SeverityMinor
				issue.Type = domain.TypeCodeSmell
				issue.Message = fmt.Sprintf("Magic number %s; extract to a named constant", val)
				issues = append(issues, issue)
			}
			return true
		})
		return issues
	},
}
