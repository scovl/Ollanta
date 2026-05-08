package treesitter

import (
	"strings"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// AvoidPyyamlLoadPY detects yaml.load() calls without Loader=SafeLoader.
var AvoidPyyamlLoadPY = ollantarules.Rule{
	MetaKey: "py:avoid-pyyaml-load",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(call
  function: (attribute
    object: (identifier) @mod
    attribute: (identifier) @func)
  (#eq? @mod "yaml")
  (#eq? @func "load")
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
			// Check if call already has keyword argument Loader=...
			callText := ctx.Query.Text(ctx.ParsedFile, call)
			if strings.Contains(callText, "Loader=") {
				continue
			}
			line, _, _, _ := ctx.Query.Position(call)
			if seen[line] {
				continue
			}
			seen[line] = true
			issue := domain.NewIssue("py:avoid-pyyaml-load", ctx.Path, line)
			issue.Severity = domain.SeverityMajor
			issue.Type = domain.TypeVulnerability
			issue.Message = "yaml.load without Loader=SafeLoader is vulnerable to arbitrary code execution"
			issues = append(issues, issue)
		}
		return issues
	},
}
