package treesitter

import (
	"fmt"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// MutableDefaultArgumentPY detects Python function definitions that use a mutable
// object (list, dict, or set) as a default parameter value. Because Python evaluates
// default values once at function definition time, mutations persist across calls,
// causing subtle and hard-to-debug bugs.
// SonarQube equivalent: python:S5717 / pylint: W0102.
var MutableDefaultArgumentPY = ollantarules.Rule{
	MetaKey: "py:mutable-default-argument",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(default_parameter
          value: [
            (list)       @mut
            (dictionary) @mut
            (set)        @mut
          ]) @param`

		matches, err := ctx.Query.Run(ctx.ParsedFile, query, ctx.Grammar)
		if err != nil {
			return nil
		}

		var issues []*domain.Issue
		for _, m := range matches {
			mut := m.Captures["mut"]
			if mut == nil {
				continue
			}
			line, _, _, _ := ctx.Query.Position(mut)
			issue := domain.NewIssue("py:mutable-default-argument", ctx.Path, line)
			issue.Severity = domain.SeverityMajor
			issue.Type = domain.TypeBug
			issue.Message = fmt.Sprintf(
				"Mutable default argument at line %d; use None and initialise inside the function instead",
				line,
			)
			issues = append(issues, issue)
		}
		return issues
	},
}
