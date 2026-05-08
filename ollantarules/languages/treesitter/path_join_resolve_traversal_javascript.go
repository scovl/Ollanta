package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
	sitter "github.com/smacker/go-tree-sitter"
)

// PathJoinResolveTraversalJS detects path.join or path.resolve with unsanitized user input.
var PathJoinResolveTraversalJS = ollantarules.Rule{
	MetaKey: "js:path-join-resolve-traversal",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		return walkIssues(ctx, func(node *sitter.Node) *domain.Issue {
			if node.Type() != "call_expression" {
				return nil
			}
			for i := 0; i < int(node.ChildCount()); i++ {
				child := node.Child(i)
				if child.Type() == "member_expression" {
					obj := ""
					fn := ""
					for j := 0; j < int(child.ChildCount()); j++ {
						c := child.Child(j)
						if c.Type() == "identifier" || c.Type() == "property_identifier" {
							txt := ctx.Query.Text(ctx.ParsedFile, c)
							if obj == "" {
								obj = txt
							} else {
								fn = txt
							}
						}
					}
					if obj == "path" && (fn == "join" || fn == "resolve") {
						line, _, _, _ := ctx.Query.Position(node)
						issue := domain.NewIssue("js:path-join-resolve-traversal", ctx.Path, line)
						issue.Severity = domain.SeverityMajor
						issue.Type = domain.TypeVulnerability
						issue.Message = "path.join/path.resolve with untrusted input can lead to path traversal; validate inputs"
						return issue
					}
				}
			}
			return nil
		})
	},
}
