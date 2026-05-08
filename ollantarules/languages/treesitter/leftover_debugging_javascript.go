package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// LeftoverDebuggingJS detects debugger statements and alert() calls left in
// production code.
var LeftoverDebuggingJS = ollantarules.Rule{
	MetaKey: "js:leftover-debugging",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `[
  (debugger_statement) @dbg
  (call_expression
    function: (identifier) @alert
    (#eq? @alert "alert"))
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
			issue := domain.NewIssue("js:leftover-debugging", ctx.Path, line)
			issue.Severity = domain.SeverityMinor
			issue.Type = domain.TypeCodeSmell
			issue.Message = "Debugging statement left in production code"
			issues = append(issues, issue)
		}
		return issues
	},
}
