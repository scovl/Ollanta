package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
	sitter "github.com/smacker/go-tree-sitter"
)

// DictModifyIteratingPY detects deletion from a dict while iterating over it.
var DictModifyIteratingPY = ollantarules.Rule{
	MetaKey: "py:dict-modify-iterating",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		return walkIssues(ctx, func(node *sitter.Node) *domain.Issue {
			if node.Type() != "for_statement" {
				return nil
			}
			// Extract iterator variable name
			iterName := ""
			for i := 0; i < int(node.ChildCount()); i++ {
				child := node.Child(i)
				if child.Type() == "identifier" {
					iterName = ctx.Query.Text(ctx.ParsedFile, child)
					break
				}
			}
			if iterName == "" {
				return nil
			}
			// Look for 'del dict[var]' inside the body
			var found bool
			var walk func(*sitter.Node)
			walk = func(n *sitter.Node) {
				if n == nil || found {
					return
				}
				if n.Type() == "delete_statement" || n.Type() == "delete" {
					// Check if it deletes using the iterator variable
					for j := 0; j < int(n.ChildCount()); j++ {
						sub := n.Child(j)
						if sub.Type() == "subscript" {
							for k := 0; k < int(sub.ChildCount()); k++ {
								if sub.Child(k).Type() == "identifier" && ctx.Query.Text(ctx.ParsedFile, sub.Child(k)) == iterName {
									found = true
									return
								}
							}
						}
					}
				}
				for j := 0; j < int(n.ChildCount()); j++ {
					walk(n.Child(j))
				}
			}
			walk(node)
			if !found {
				return nil
			}
			line, _, _, _ := ctx.Query.Position(node)
			issue := domain.NewIssue("py:dict-modify-iterating", ctx.Path, line)
			issue.Severity = domain.SeverityMajor
			issue.Type = domain.TypeBug
			issue.Message = "Dictionary is modified (del) while iterating over it"
			return issue
		})
	},
}
