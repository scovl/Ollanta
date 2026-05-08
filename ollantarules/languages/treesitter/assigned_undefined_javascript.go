package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// AssignedUndefinedJS detects assignments of undefined which are redundant
// in JavaScript (variables are already undefined by default).
var AssignedUndefinedJS = ollantarules.Rule{
	MetaKey: "js:assigned-undefined",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(assignment_expression
  right: (identifier) @undef
  (#eq? @undef "undefined")
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
			issue := domain.NewIssue("js:assigned-undefined", ctx.Path, line)
			issue.Severity = domain.SeverityMinor
			issue.Type = domain.TypeCodeSmell
			issue.Message = "Assigning undefined is redundant; variables are undefined by default"
			issues = append(issues, issue)
		}
		return issues
	},
}
