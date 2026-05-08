package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// MarshalPY detects marshal.load or marshal.loads usage.
var MarshalPY = ollantarules.Rule{
	MetaKey: "py:marshal",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(call
  function: (attribute
    object: (identifier) @mod
    attribute: (identifier) @func)
  (#eq? @mod "marshal")
  (#match? @func "^(load|loads|dump|dumps)$")
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
			issue := domain.NewIssue("py:marshal", ctx.Path, line)
			issue.Severity = domain.SeverityMajor
			issue.Type = domain.TypeVulnerability
			issue.Message = "marshal.load can execute arbitrary code; do not unmarshal untrusted data"
			issues = append(issues, issue)
		}
		return issues
	},
}
