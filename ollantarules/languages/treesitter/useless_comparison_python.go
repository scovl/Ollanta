package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// UselessComparisonPY detects deterministic comparisons in Python such as
// comparing a string literal with an integer or other incompatible types.
// This rule focuses on obvious always-true/always-false patterns.
var UselessComparisonPY = ollantarules.Rule{
	MetaKey: "py:useless-comparison",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(comparison
  left: (integer) @left
  right: (string) @right) @comp`
		matches, err := ctx.Query.Run(ctx.ParsedFile, query, ctx.Grammar)
		if err != nil {
			return nil
		}
		var issues []*domain.Issue
		seen := map[int]bool{}
		for _, m := range matches {
			comp := m.Captures["comp"]
			if comp == nil {
				continue
			}
			line, _, _, _ := ctx.Query.Position(comp)
			if seen[line] {
				continue
			}
			seen[line] = true
			issue := domain.NewIssue("py:useless-comparison", ctx.Path, line)
			issue.Severity = domain.SeverityMinor
			issue.Type = domain.TypeBug
			issue.Message = "Comparison between incompatible literal types is always deterministic"
			issues = append(issues, issue)
		}
		return issues
	},
}
