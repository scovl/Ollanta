package treesitter

import (
	"fmt"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
	sitter "github.com/smacker/go-tree-sitter"
)

// MissingHashWithEqPY detects classes that define __eq__ but not __hash__.
var MissingHashWithEqPY = ollantarules.Rule{
	MetaKey: "py:missing-hash-with-eq",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		return walkIssues(ctx, func(node *sitter.Node) *domain.Issue {
			if node.Type() != "class_definition" {
				return nil
			}
			// Find class name
			var className string
			for i := 0; i < int(node.ChildCount()); i++ {
				if node.Child(i).Type() == "identifier" {
					className = ctx.Query.Text(ctx.ParsedFile, node.Child(i))
					break
				}
			}
			if className == "" {
				return nil
			}
			// Search body for __eq__ and __hash__
			hasEq := false
			hasHash := false
			var walk func(*sitter.Node)
			walk = func(n *sitter.Node) {
				if n == nil {
					return
				}
				if n.Type() == "function_definition" {
					for j := 0; j < int(n.ChildCount()); j++ {
						c := n.Child(j)
						if c.Type() == "identifier" {
							name := ctx.Query.Text(ctx.ParsedFile, c)
							if name == "__eq__" {
								hasEq = true
							}
							if name == "__hash__" {
								hasHash = true
							}
						}
					}
				}
				for j := 0; j < int(n.ChildCount()); j++ {
					walk(n.Child(j))
				}
			}
			walk(node)
			if !hasEq || hasHash {
				return nil
			}
			line, _, _, _ := ctx.Query.Position(node)
			issue := domain.NewIssue("py:missing-hash-with-eq", ctx.Path, line)
			issue.Severity = domain.SeverityMinor
			issue.Type = domain.TypeCodeSmell
			issue.Message = fmt.Sprintf("Class '%s' defines __eq__ but not __hash__; instances will be unhashable", className)
			return issue
		})
	},
}
