package treesitter

import (
	"fmt"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// BroadExceptPY detects Python except clauses that catch Exception (or any base
// exception type) without narrowing to a specific error, silently swallowing bugs.
// SonarQube equivalent: python:S5754 / pylint: W0703.
var BroadExceptPY = ollantarules.Rule{
	MetaKey: "py:broad-except",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `[
          (except_clause) @bare
          (except_clause
            (as_pattern
              (attribute
                object: (identifier) @cls)))
          (except_clause
            (as_pattern
              (identifier) @cls))
          (except_clause
            (attribute
              object: (identifier) @cls))
          (except_clause
            (identifier) @cls)
        ] @clause`

		matches, err := ctx.Query.Run(ctx.ParsedFile, query, ctx.Grammar)
		if err != nil {
			return nil
		}

		broadTypes := map[string]bool{
			"Exception":     true,
			"BaseException": true,
		}

		var issues []*domain.Issue
		seen := map[int]bool{}
		for _, m := range matches {
			clause := m.Captures["clause"]
			if clause == nil {
				continue
			}
			startLine, _, _, _ := ctx.Query.Position(clause)
			if seen[startLine] {
				continue
			}

			cls := m.Captures["cls"]
			isBroad := cls == nil // bare except:
			if cls != nil {
				name := ctx.Query.Text(ctx.ParsedFile, cls)
				isBroad = broadTypes[name]
			}

			if isBroad {
				seen[startLine] = true
				issue := domain.NewIssue("py:broad-except", ctx.Path, startLine)
				issue.Severity = domain.SeverityMajor
				issue.Type = domain.TypeBug
				if cls == nil {
					issue.Message = "Bare 'except:' clause catches all exceptions; catch specific exceptions instead"
				} else {
					issue.Message = fmt.Sprintf(
						"Catching broad exception type '%s' hides bugs; catch specific exceptions instead",
						ctx.Query.Text(ctx.ParsedFile, cls),
					)
				}
				issues = append(issues, issue)
			}
		}
		return issues
	},
}
