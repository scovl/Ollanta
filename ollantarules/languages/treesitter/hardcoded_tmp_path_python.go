package treesitter

import (
	"strings"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// HardcodedTmpPathPY detects hardcoded /tmp/ paths in string literals,
// which may be insecure across platforms.
var HardcodedTmpPathPY = ollantarules.Rule{
	MetaKey: "py:hardcoded-tmp-path",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(string
  (string_content) @content
) @str`
		matches, err := ctx.Query.Run(ctx.ParsedFile, query, ctx.Grammar)
		if err != nil {
			return nil
		}
		var issues []*domain.Issue
		seen := map[int]bool{}
		for _, m := range matches {
			contentNode := m.Captures["content"]
			strNode := m.Captures["str"]
			if contentNode == nil || strNode == nil {
				continue
			}
			text := ctx.Query.Text(ctx.ParsedFile, contentNode)
			if strings.HasPrefix(text, "/tmp/") {
				line, _, _, _ := ctx.Query.Position(strNode)
				if seen[line] {
					continue
				}
				seen[line] = true
				issue := domain.NewIssue("py:hardcoded-tmp-path", ctx.Path, line)
				issue.Severity = domain.SeverityMinor
				issue.Type = domain.TypeCodeSmell
				issue.Message = "Hardcoded /tmp/ path is not portable and may be insecure"
				issues = append(issues, issue)
			}
		}
		return issues
	},
}
