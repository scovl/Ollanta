package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
	sitter "github.com/smacker/go-tree-sitter"
)

// UselessEqEqJS detects self-comparisons in JavaScript/TypeScript like x == x or x != x.
var UselessEqEqJS = ollantarules.Rule{
	MetaKey: "js:useless-eqeq",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		return walkIssues(ctx, func(node *sitter.Node) *domain.Issue {
			if node.Type() != "binary_expression" {
				return nil
			}
			if node.ChildCount() < 3 {
				return nil
			}
			left := node.Child(0)
			op := node.Child(1)
			right := node.Child(2)
			if left == nil || op == nil || right == nil {
				return nil
			}
			opText := ctx.Query.Text(ctx.ParsedFile, op)
			if opText != "==" && opText != "!=" && opText != "===" && opText != "!==" {
				return nil
			}
			if ctx.Query.Text(ctx.ParsedFile, left) != ctx.Query.Text(ctx.ParsedFile, right) {
				return nil
			}
			line, _, _, _ := ctx.Query.Position(node)
			issue := domain.NewIssue("js:useless-eqeq", ctx.Path, line)
			issue.Severity = domain.SeverityMinor
			issue.Type = domain.TypeBug
			issue.Message = "Useless self-comparison: this expression is always true or always false"
			return issue
		})
	},
}
