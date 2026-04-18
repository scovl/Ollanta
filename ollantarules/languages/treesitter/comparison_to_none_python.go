package treesitter

import (
	"fmt"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// ComparisonToNonePY detects Python comparisons to None using == or != instead of
// the identity operators `is` and `is not`. PEP 8 explicitly requires identity
// comparison for None. SonarQube equivalent: python:S5727.
var ComparisonToNonePY = ollantarules.Rule{
	MetaKey: "py:comparison-to-none",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(comparison_operator
          [
            (none) @none
            (_ (none) @none)
          ]
          operators: _ @op
        ) @cmp`

		matches, err := ctx.Query.Run(ctx.ParsedFile, query, ctx.Grammar)
		if err != nil {
			return nil
		}

		var issues []*domain.Issue
		seen := map[int]bool{}
		for _, m := range matches {
			cmp := m.Captures["cmp"]
			op := m.Captures["op"]
			if cmp == nil || op == nil {
				continue
			}
			opText := ctx.Query.Text(ctx.ParsedFile, op)
			if opText != "==" && opText != "!=" {
				continue
			}
			line, _, _, _ := ctx.Query.Position(cmp)
			if seen[line] {
				continue
			}
			seen[line] = true

			better := "is"
			if opText == "!=" {
				better = "is not"
			}
			issue := domain.NewIssue("py:comparison-to-none", ctx.Path, line)
			issue.Severity = domain.SeverityMinor
			issue.Type = domain.TypeCodeSmell
			issue.Message = fmt.Sprintf(
				"Use '%s None' instead of '%s None' (PEP 8)", better, opText,
			)
			issues = append(issues, issue)
		}
		return issues
	},
}
