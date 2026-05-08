package treesitter

import (
	"fmt"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// UncheckedReturnsPY detects discarded return values from os functions that
// should typically be checked (remove, rename, mkdir, etc.).
var UncheckedReturnsPY = ollantarules.Rule{
	MetaKey: "py:unchecked-returns",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(expression_statement
  (call
    function: (attribute
      object: (identifier) @mod
      attribute: (identifier) @func)
    (#eq? @mod "os")
    (#match? @func "^(remove|rename|mkdir|rmdir|chmod|chown|makedirs|renames|replace|link|symlink)$"))
) @stmt`
		matches, err := ctx.Query.Run(ctx.ParsedFile, query, ctx.Grammar)
		if err != nil {
			return nil
		}
		var issues []*domain.Issue
		seen := map[int]bool{}
		for _, m := range matches {
			stmt := m.Captures["stmt"]
			fn := m.Captures["func"]
			if stmt == nil || fn == nil {
				continue
			}
			line, _, _, _ := ctx.Query.Position(stmt)
			if seen[line] {
				continue
			}
			seen[line] = true
			funcName := ctx.Query.Text(ctx.ParsedFile, fn)
			issue := domain.NewIssue("py:unchecked-returns", ctx.Path, line)
			issue.Severity = domain.SeverityMinor
			issue.Type = domain.TypeBug
			issue.Message = fmt.Sprintf("Return value of os.%s is discarded; check for errors", funcName)
			issues = append(issues, issue)
		}
		return issues
	},
}
