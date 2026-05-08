package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// PicklePY detects pickle.load or pickle.loads with non-hardcoded sources.
var PicklePY = ollantarules.Rule{
	MetaKey: "py:pickle",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(call
  function: (attribute
    object: (identifier) @mod
    attribute: (identifier) @func)
  (#eq? @mod "pickle")
  (#match? @func "^(load|loads)$")
) @call`
		matches, err := ctx.Query.Run(ctx.ParsedFile, query, ctx.Grammar)
		if err != nil {
			return nil
		}
		var issues []*domain.Issue
		seen := map[int]bool{}
		for _, m := range matches {
			call := m.Captures["call"]
			if call == nil {
				continue
			}
			line, _, _, _ := ctx.Query.Position(call)
			if seen[line] {
				continue
			}
			seen[line] = true
			issue := domain.NewIssue("py:pickle", ctx.Path, line)
			issue.Severity = domain.SeverityMajor
			issue.Type = domain.TypeVulnerability
			issue.Message = "pickle.load can execute arbitrary code; do not unpickle untrusted data"
			issues = append(issues, issue)
		}
		return issues
	},
}
