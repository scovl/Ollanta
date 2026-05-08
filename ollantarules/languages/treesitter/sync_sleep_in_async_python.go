package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
	sitter "github.com/smacker/go-tree-sitter"
)

// SyncSleepInAsyncPY detects time.sleep inside async functions.
var SyncSleepInAsyncPY = ollantarules.Rule{
	MetaKey: "py:sync-sleep-in-async",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		return walkIssues(ctx, func(node *sitter.Node) *domain.Issue {
			if node.Type() != "function_definition" {
				return nil
			}
			// Check if function has 'async' modifier
			hasAsync := false
			for i := 0; i < int(node.ChildCount()); i++ {
				if node.Child(i).Type() == "async" {
					hasAsync = true
					break
				}
			}
			if !hasAsync {
				return nil
			}
			// Search for time.sleep call inside body
			var found bool
			var walk func(*sitter.Node)
			walk = func(n *sitter.Node) {
				if n == nil || found {
					return
				}
				if n.Type() == "call" {
					for j := 0; j < int(n.ChildCount()); j++ {
						child := n.Child(j)
						if child.Type() == "attribute" {
							obj := ""
							fn := ""
							for k := 0; k < int(child.ChildCount()); k++ {
								c := child.Child(k)
								if c.Type() == "identifier" {
									if obj == "" {
										obj = ctx.Query.Text(ctx.ParsedFile, c)
									} else {
										fn = ctx.Query.Text(ctx.ParsedFile, c)
									}
								}
							}
							if obj == "time" && fn == "sleep" {
								found = true
								return
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
			issue := domain.NewIssue("py:sync-sleep-in-async", ctx.Path, line)
			issue.Severity = domain.SeverityMajor
			issue.Type = domain.TypeBug
			issue.Message = "time.sleep blocks the event loop in async code; use asyncio.sleep instead"
			return issue
		})
	},
}
