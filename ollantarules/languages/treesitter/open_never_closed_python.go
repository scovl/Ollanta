package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// OpenNeverClosedPY detects open() calls whose return value is immediately
// discarded (not assigned and not used in a with statement).
var OpenNeverClosedPY = ollantarules.Rule{
	MetaKey: "py:open-never-closed",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(expression_statement
  (call
    function: (identifier) @fn
    (#eq? @fn "open"))
) @stmt`
		matches, err := ctx.Query.Run(ctx.ParsedFile, query, ctx.Grammar)
		if err != nil {
			return nil
		}
		var issues []*domain.Issue
		seen := map[int]bool{}
		for _, m := range matches {
			stmt := m.Captures["stmt"]
			if stmt == nil {
				continue
			}
			line, _, _, _ := ctx.Query.Position(stmt)
			if seen[line] {
				continue
			}
			seen[line] = true
			issue := domain.NewIssue("py:open-never-closed", ctx.Path, line)
			issue.Severity = domain.SeverityMinor
			issue.Type = domain.TypeCodeSmell
			issue.Message = "open() return value is discarded; use 'with' or ensure the file is closed"
			issues = append(issues, issue)
		}
		return issues
	},
}
