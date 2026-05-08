package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
	sitter "github.com/smacker/go-tree-sitter"
)

// walkIssues traverses the AST rooted at ctx.ParsedFile and calls fn on every node.
// If fn returns a non-nil issue, it is collected. Deduplicates by line number.
func walkIssues(ctx *ollantarules.AnalysisContext, fn func(*sitter.Node) *domain.Issue) []*domain.Issue {
	var issues []*domain.Issue
	seen := map[int]bool{}
	root := ctx.ParsedFile.RootNode()
	var walk func(*sitter.Node)
	walk = func(n *sitter.Node) {
		if n == nil {
			return
		}
		if issue := fn(n); issue != nil {
			if !seen[issue.Line] {
				seen[issue.Line] = true
				issues = append(issues, issue)
			}
		}
		for i := 0; i < int(n.ChildCount()); i++ {
			walk(n.Child(i))
		}
	}
	walk(root)
	return issues
}
