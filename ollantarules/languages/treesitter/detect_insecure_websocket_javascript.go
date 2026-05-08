package treesitter

import (
	"strings"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// DetectInsecureWebsocketJS detects new WebSocket('ws://...').
var DetectInsecureWebsocketJS = ollantarules.Rule{
	MetaKey: "js:detect-insecure-websocket",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(new_expression
  constructor: (identifier) @ctor
  arguments: (arguments
    (string
      (string_fragment) @url))
  (#eq? @ctor "WebSocket")
) @expr`
		matches, err := ctx.Query.Run(ctx.ParsedFile, query, ctx.Grammar)
		if err != nil {
			return nil
		}
		var issues []*domain.Issue
		seen := map[int]bool{}
		for _, m := range matches {
			expr := m.Captures["expr"]
			urlNode := m.Captures["url"]
			if expr == nil || urlNode == nil {
				continue
			}
			url := ctx.Query.Text(ctx.ParsedFile, urlNode)
			if !strings.HasPrefix(url, "ws://") {
				continue
			}
			line, _, _, _ := ctx.Query.Position(expr)
			if seen[line] {
				continue
			}
			seen[line] = true
			issue := domain.NewIssue("js:detect-insecure-websocket", ctx.Path, line)
			issue.Severity = domain.SeverityMajor
			issue.Type = domain.TypeVulnerability
			issue.Message = "Insecure WebSocket URL (ws://); use wss:// for encrypted connections"
			issues = append(issues, issue)
		}
		return issues
	},
}
