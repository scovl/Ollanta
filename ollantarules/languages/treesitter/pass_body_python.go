package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// PassBodyPY detects functions or classes whose body is a single pass statement
// without a docstring, indicating an incomplete implementation.
var PassBodyPY = ollantarules.Rule{
	MetaKey: "py:pass-body",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `[
  (function_definition
    body: (block . (pass_statement) .) @pass)
  (class_definition
    body: (block . (pass_statement) .) @pass)
] @def`
		matches, err := ctx.Query.Run(ctx.ParsedFile, query, ctx.Grammar)
		if err != nil {
			return nil
		}
		var issues []*domain.Issue
		seen := map[int]bool{}
		for _, m := range matches {
			defNode := m.Captures["def"]
			if defNode == nil {
				continue
			}
			line, _, _, _ := ctx.Query.Position(defNode)
			if seen[line] {
				continue
			}
			seen[line] = true
			issue := domain.NewIssue("py:pass-body", ctx.Path, line)
			issue.Severity = domain.SeverityMinor
			issue.Type = domain.TypeCodeSmell
			issue.Message = "Empty pass body without docstring — incomplete implementation?"
			issues = append(issues, issue)
		}
		return issues
	},
}
