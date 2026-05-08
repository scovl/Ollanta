package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// ReturnInInitPY detects return statements inside __init__ methods,
// which is invalid in Python (only None can be returned implicitly).
var ReturnInInitPY = ollantarules.Rule{
	MetaKey: "py:return-in-init",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(function_definition
  name: (identifier) @name
  (#eq? @name "__init__")
  body: (block
    (return_statement) @ret)
) @func`
		matches, err := ctx.Query.Run(ctx.ParsedFile, query, ctx.Grammar)
		if err != nil {
			return nil
		}
		var issues []*domain.Issue
		seen := map[int]bool{}
		for _, m := range matches {
			funcNode := m.Captures["func"]
			if funcNode == nil {
				continue
			}
			line, _, _, _ := ctx.Query.Position(funcNode)
			if seen[line] {
				continue
			}
			seen[line] = true
			issue := domain.NewIssue("py:return-in-init", ctx.Path, line)
			issue.Severity = domain.SeverityMajor
			issue.Type = domain.TypeBug
			issue.Message = "__init__ should not contain an explicit return statement"
			issues = append(issues, issue)
		}
		return issues
	},
}
