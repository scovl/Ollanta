package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// DangerousSubprocessPY detects subprocess calls with shell=True.
var DangerousSubprocessPY = ollantarules.Rule{
	MetaKey: "py:dangerous-subprocess",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(call
  function: (attribute
    object: (identifier) @mod
    attribute: (identifier) @func)
  arguments: (argument_list
    (keyword_argument
      name: (identifier) @kw
      value: (true)))
  (#eq? @mod "subprocess")
  (#match? @func "^(run|call|Popen|check_output|check_call)$")
  (#eq? @kw "shell")
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
			issue := domain.NewIssue("py:dangerous-subprocess", ctx.Path, line)
			issue.Severity = domain.SeverityMajor
			issue.Type = domain.TypeVulnerability
			issue.Message = "subprocess with shell=True is dangerous; avoid shell injection by passing a list of arguments"
			issues = append(issues, issue)
		}
		return issues
	},
}
