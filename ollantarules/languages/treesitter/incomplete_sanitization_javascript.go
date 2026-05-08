package treesitter

import (
	"strings"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
	sitter "github.com/smacker/go-tree-sitter"
)

// IncompleteSanitizationJS detects replace calls that only remove a subset of dangerous characters.
var IncompleteSanitizationJS = ollantarules.Rule{
	MetaKey: "js:incomplete-sanitization",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		return walkIssues(ctx, func(node *sitter.Node) *domain.Issue {
			if node.Type() != "call_expression" {
				return nil
			}
			for i := 0; i < int(node.ChildCount()); i++ {
				child := node.Child(i)
				if child.Type() == "member_expression" {
					fn := ""
					for j := 0; j < int(child.ChildCount()); j++ {
						c := child.Child(j)
						if c.Type() == "property_identifier" {
							fn = ctx.Query.Text(ctx.ParsedFile, c)
						}
					}
					if fn != "replace" {
						return nil
					}
				}
			}
			// Check regex argument for incomplete pattern
			for i := 0; i < int(node.ChildCount()); i++ {
				child := node.Child(i)
				if child.Type() == "arguments" {
					for j := 0; j < int(child.ChildCount()); j++ {
						c := child.Child(j)
						if c.Type() == "regex" {
							pat := ""
							for k := 0; k < int(c.ChildCount()); k++ {
								d := c.Child(k)
								if d.Type() == "regex_pattern" {
									pat = ctx.Query.Text(ctx.ParsedFile, d)
								}
							}
							if isIncompleteSanitization(pat) {
								line, _, _, _ := ctx.Query.Position(node)
								issue := domain.NewIssue("js:incomplete-sanitization", ctx.Path, line)
								issue.Severity = domain.SeverityMajor
								issue.Type = domain.TypeVulnerability
								issue.Message = "Incomplete sanitization: replace only removes some dangerous characters and may be bypassed"
								return issue
							}
						}
					}
				}
			}
			return nil
		})
	},
}

func isIncompleteSanitization(pattern string) bool {
	trimmed := strings.TrimSpace(pattern)
	if len(trimmed) <= 2 {
		return true
	}
	if strings.HasPrefix(trimmed, "<") && !strings.Contains(trimmed[1:], "<") && !strings.Contains(trimmed, ">") {
		return true
	}
	if strings.HasPrefix(trimmed, ">") && !strings.Contains(trimmed[1:], ">") && !strings.Contains(trimmed, "<") {
		return true
	}
	return false
}
