package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// DetectEvalJS detects calls to eval() which execute arbitrary code and
// are a major security risk.
var DetectEvalJS = ollantarules.Rule{
	MetaKey: "js:detect-eval",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(call_expression
  function: (identifier) @eval
  (#eq? @eval "eval")
) @expr`
		matches, err := ctx.Query.Run(ctx.ParsedFile, query, ctx.Grammar)
		if err != nil {
			return nil
		}
		var issues []*domain.Issue
		seen := map[int]bool{}
		for _, m := range matches {
			expr := m.Captures["expr"]
			if expr == nil {
				continue
			}
			line, _, _, _ := ctx.Query.Position(expr)
			if seen[line] {
				continue
			}
			seen[line] = true
			issue := domain.NewIssue("js:detect-eval", ctx.Path, line)
			issue.Severity = domain.SeverityCritical
			issue.Type = domain.TypeVulnerability
			issue.Message = "eval() executes arbitrary code and must not be used with untrusted input"
			issues = append(issues, issue)
		}
		return issues
	},
}
