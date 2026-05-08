package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
	sitter "github.com/smacker/go-tree-sitter"
)

// DetectPseudoRandomBytesJS detects use of crypto.pseudoRandomBytes.
var DetectPseudoRandomBytesJS = ollantarules.Rule{
	MetaKey: "js:detect-pseudoRandomBytes",
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
					if obj == "crypto" && fn == "pseudoRandomBytes" {
						line, _, _, _ := ctx.Query.Position(node)
						issue := domain.NewIssue("js:detect-pseudoRandomBytes", ctx.Path, line)
						issue.Severity = domain.SeverityMajor
						issue.Type = domain.TypeVulnerability
						issue.Message = "crypto.pseudoRandomBytes is not cryptographically secure; use crypto.randomBytes instead"
						return issue
					}
				}
			}
			return nil
		})
	},
}
