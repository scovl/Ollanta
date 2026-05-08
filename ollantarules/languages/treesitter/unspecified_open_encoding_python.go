package treesitter

import (
	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// UnspecifiedOpenEncodingPY detects open() calls that omit the encoding
// parameter, which may cause platform-dependent behavior (e.g. Windows CP1252).
var UnspecifiedOpenEncodingPY = ollantarules.Rule{
	MetaKey: "py:unspecified-open-encoding",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		// Find all open() calls
		allQuery := `(call
  function: (identifier) @fn
  (#eq? @fn "open")
) @call`
		allMatches, err := ctx.Query.Run(ctx.ParsedFile, allQuery, ctx.Grammar)
		if err != nil {
			return nil
		}

		// Find open() calls that DO have encoding= keyword
		encQuery := `(call
  function: (identifier) @fn
  (#eq? @fn "open")
  arguments: (argument_list
    (keyword_argument
      name: (identifier) @enc
      (#eq? @enc "encoding")))
) @call`
		encMatches, err := ctx.Query.Run(ctx.ParsedFile, encQuery, ctx.Grammar)
		if err != nil {
			return nil
		}

		// Build set of lines that have encoding
		hasEncoding := map[int]bool{}
		for _, m := range encMatches {
			callNode := m.Captures["call"]
			if callNode == nil {
				continue
			}
			line, _, _, _ := ctx.Query.Position(callNode)
			hasEncoding[line] = true
		}

		var issues []*domain.Issue
		seen := map[int]bool{}
		for _, m := range allMatches {
			callNode := m.Captures["call"]
			if callNode == nil {
				continue
			}
			line, _, _, _ := ctx.Query.Position(callNode)
			if hasEncoding[line] || seen[line] {
				continue
			}
			seen[line] = true
			issue := domain.NewIssue("py:unspecified-open-encoding", ctx.Path, line)
			issue.Severity = domain.SeverityMinor
			issue.Type = domain.TypeCodeSmell
			issue.Message = "open() without explicit encoding may cause platform-dependent behavior"
			issues = append(issues, issue)
		}
		return issues
	},
}
