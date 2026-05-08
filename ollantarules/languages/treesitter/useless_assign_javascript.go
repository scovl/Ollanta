package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// UselessAssignJS detects self-assignments like x = x which have no effect.
var UselessAssignJS = ollantarules.Rule{
	MetaKey: "js:useless-assign",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(assignment_expression
  left: (identifier) @left
  right: (identifier) @right
) @expr`
		matches, err := ctx.Query.Run(ctx.ParsedFile, query, ctx.Grammar)
		if err != nil {
			return nil
		}
		var issues []*domain.Issue
		seen := map[int]bool{}
		for _, m := range matches {
			left := m.Captures["left"]
			right := m.Captures["right"]
			expr := m.Captures["expr"]
			if left == nil || right == nil || expr == nil {
				continue
			}
			if ctx.Query.Text(ctx.ParsedFile, left) == ctx.Query.Text(ctx.ParsedFile, right) {
				line, _, _, _ := ctx.Query.Position(expr)
				if seen[line] {
					continue
				}
				seen[line] = true
				issue := domain.NewIssue("js:useless-assign", ctx.Path, line)
				issue.Severity = domain.SeverityMinor
				issue.Type = domain.TypeBug
				issue.Message = "Self-assignment has no effect and indicates a bug or incomplete refactoring"
				issues = append(issues, issue)
			}
		}
		return issues
	},
}
