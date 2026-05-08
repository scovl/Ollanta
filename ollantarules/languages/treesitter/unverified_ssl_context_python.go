package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// UnverifiedSSLContextPY detects ssl._create_unverified_context().
var UnverifiedSSLContextPY = ollantarules.Rule{
	MetaKey: "py:unverified-ssl-context",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(call
  function: (attribute
    object: (identifier) @mod
    attribute: (identifier) @func)
  (#eq? @mod "ssl")
  (#eq? @func "_create_unverified_context")
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
			issue := domain.NewIssue("py:unverified-ssl-context", ctx.Path, line)
			issue.Severity = domain.SeverityMajor
			issue.Type = domain.TypeVulnerability
			issue.Message = "ssl._create_unverified_context disables certificate validation and is vulnerable to MITM attacks"
			issues = append(issues, issue)
		}
		return issues
	},
}
