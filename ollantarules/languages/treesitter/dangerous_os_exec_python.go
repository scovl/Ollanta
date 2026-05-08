package treesitter

import (
	"fmt"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// DangerousOsExecPY detects dangerous os.exec* and os.system calls.
var DangerousOsExecPY = ollantarules.Rule{
	MetaKey: "py:dangerous-os-exec",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		query := `(call
  function: (attribute
    object: (identifier) @mod
    attribute: (identifier) @func)
  (#eq? @mod "os")
  (#match? @func "^(system|execl|execle|execlp|execv|execve|execvp|execvpe|spawnl|spawnle|spawnlp|spawnv|spawnve|spawnvp|spawnvpe)$")
) @call`
		matches, err := ctx.Query.Run(ctx.ParsedFile, query, ctx.Grammar)
		if err != nil {
			return nil
		}
		var issues []*domain.Issue
		seen := map[int]bool{}
		for _, m := range matches {
			call := m.Captures["call"]
			fn := m.Captures["func"]
			if call == nil || fn == nil {
				continue
			}
			line, _, _, _ := ctx.Query.Position(call)
			if seen[line] {
				continue
			}
			seen[line] = true
			funcName := ctx.Query.Text(ctx.ParsedFile, fn)
			issue := domain.NewIssue("py:dangerous-os-exec", ctx.Path, line)
			issue.Severity = domain.SeverityMajor
			issue.Type = domain.TypeVulnerability
			issue.Message = fmt.Sprintf("os.%s executes arbitrary commands; validate all inputs", funcName)
			issues = append(issues, issue)
		}
		return issues
	},
}
