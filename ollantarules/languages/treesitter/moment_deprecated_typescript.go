package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// MomentDeprecatedTS detects imports or usage of the deprecated moment.js
// library in TypeScript code.
var MomentDeprecatedTS = ollantarules.Rule{
	MetaKey: "ts:moment-deprecated",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `[
  (import_statement
    source: (string (string_fragment) @src)
    (#match? @src "moment"))
  (call_expression
    function: (identifier) @moment
    (#eq? @moment "moment"))
] @expr`
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
			issue := domain.NewIssue("ts:moment-deprecated", ctx.Path, line)
			issue.Severity = domain.SeverityMinor
			issue.Type = domain.TypeCodeSmell
			issue.Message = "moment.js is deprecated; consider date-fns, luxon, or native Temporal"
			issues = append(issues, issue)
		}
		return issues
	},
}
