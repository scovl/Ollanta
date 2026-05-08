package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// UselessTernaryTS detects ternary expressions equivalent to the condition itself
// or its negation (e.g., cond ? true : false).
var UselessTernaryTS = ollantarules.Rule{
	MetaKey: "ts:useless-ternary",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(ternary_expression
  condition: (_)
  consequence: (_) @con
  alternative: (_) @alt
) @expr`
		matches, err := ctx.Query.Run(ctx.ParsedFile, query, ctx.Grammar)
		if err != nil {
			return nil
		}
		var issues []*domain.Issue
		seen := map[int]bool{}
		for _, m := range matches {
			expr := m.Captures["expr"]
			con := m.Captures["con"]
			alt := m.Captures["alt"]
			if expr == nil || con == nil || alt == nil {
				continue
			}
			conText := ctx.Query.Text(ctx.ParsedFile, con)
			altText := ctx.Query.Text(ctx.ParsedFile, alt)
			var message string
			if conText == "true" && altText == "false" {
				message = "Ternary 'cond ? true : false' is equivalent to 'cond'"
			} else if conText == "false" && altText == "true" {
				message = "Ternary 'cond ? false : true' is equivalent to '!cond'"
			} else {
				continue
			}
			line, _, _, _ := ctx.Query.Position(expr)
			if seen[line] {
				continue
			}
			seen[line] = true
			issue := domain.NewIssue("ts:useless-ternary", ctx.Path, line)
			issue.Severity = domain.SeverityMinor
			issue.Type = domain.TypeBug
			issue.Message = message
			issues = append(issues, issue)
		}
		return issues
	},
}
