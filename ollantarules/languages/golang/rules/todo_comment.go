package rules

import (
	"fmt"
	"strings"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// TodoComment flags TODO, FIXME, HACK, and XXX comment markers left in production code.
// These markers indicate incomplete or problematic code that should be tracked in an
// issue tracker rather than silently left in source. SonarQube equivalent: squid:S1135.

// todoMarkers are the comment prefixes that indicate incomplete or deferred work.
var todoMarkers = []string{"TODO", "FIXME", "HACK", "XXX"}

// TodoComment flags TODO/FIXME/HACK/XXX comments that should be tracked in an issue tracker.
var TodoComment = ollantarules.Rule{
	MetaKey: "go:todo-comment",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		var issues []*domain.Issue
		for _, cg := range ctx.AST.Comments {
			for _, c := range cg.List {
				text := strings.ToUpper(c.Text)
				for _, marker := range todoMarkers {
					if strings.Contains(text, marker) {
						line := ctx.FileSet.Position(c.Slash).Line
						issue := domain.NewIssue("go:todo-comment", ctx.Path, line)
						issue.Severity = domain.SeverityInfo
						issue.Type = domain.TypeCodeSmell
						issue.Message = fmt.Sprintf("Complete the task associated with this %q comment", marker)
						issues = append(issues, issue)
						break
					}
				}
			}
		}
		return issues
	},
}
