package treesitter

import (
	"fmt"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// InsecureHashPY detects use of weak hash algorithms md5 and sha1 in hashlib.
var InsecureHashPY = ollantarules.Rule{
	MetaKey: "py:insecure-hash",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(call
  function: (attribute
    object: (identifier) @mod
    attribute: (identifier) @func)
  (#eq? @mod "hashlib")
  (#match? @func "^(md5|sha1)$")
) @call`
		matches, err := ctx.Query.Run(ctx.ParsedFile, query, ctx.Grammar)
		if err != nil {
			return nil
		}
		var issues []*domain.Issue
		seen := map[int]bool{}
		for _, m := range matches {
			call := m.Captures["call"]
			fn := m.Captures["func"]
			if call == nil || fn == nil {
				continue
			}
			line, _, _, _ := ctx.Query.Position(call)
			if seen[line] {
				continue
			}
			seen[line] = true
			funcName := ctx.Query.Text(ctx.ParsedFile, fn)
			issue := domain.NewIssue("py:insecure-hash", ctx.Path, line)
			issue.Severity = domain.SeverityMajor
			issue.Type = domain.TypeVulnerability
			issue.Message = fmt.Sprintf("hashlib.%s is cryptographically weak; use hashlib.sha256 or better", funcName)
			issues = append(issues, issue)
		}
		return issues
	},
}
