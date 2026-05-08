package treesitter

import (
	"fmt"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// ListModifyIteratingPY detects modification of a list while iterating over it,
// which typically produces unexpected results or skips elements.
var ListModifyIteratingPY = ollantarules.Rule{
	MetaKey: "py:list-modify-iterating",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(for_statement
  left: (identifier) @item
  right: (identifier) @list
  body: (block
    (expression_statement
      (call
        function: (attribute
          object: (identifier) @list2
          attribute: (identifier) @method
          (#match? @method "^(remove|pop|append|insert|clear|extend)$")))))
) @loop`
		matches, err := ctx.Query.Run(ctx.ParsedFile, query, ctx.Grammar)
		if err != nil {
			return nil
		}
		var issues []*domain.Issue
		seen := map[int]bool{}
		for _, m := range matches {
			listNode := m.Captures["list"]
			list2Node := m.Captures["list2"]
			loopNode := m.Captures["loop"]
			if listNode == nil || list2Node == nil || loopNode == nil {
				continue
			}
			l1 := ctx.Query.Text(ctx.ParsedFile, listNode)
			l2 := ctx.Query.Text(ctx.ParsedFile, list2Node)
			if l1 == l2 {
				line, _, _, _ := ctx.Query.Position(loopNode)
				if seen[line] {
					continue
				}
				seen[line] = true
				issue := domain.NewIssue("py:list-modify-iterating", ctx.Path, line)
				issue.Severity = domain.SeverityMajor
				issue.Type = domain.TypeBug
				issue.Message = fmt.Sprintf("List '%s' is modified while iterating over it", l1)
				issues = append(issues, issue)
			}
		}
		return issues
	},
}
