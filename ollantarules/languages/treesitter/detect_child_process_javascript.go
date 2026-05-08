package treesitter

import (
	"fmt"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
	sitter "github.com/smacker/go-tree-sitter"
)

// DetectChildProcessJS detects calls to exec/execSync/spawn/spawnSync on child_process.
var DetectChildProcessJS = ollantarules.Rule{
	MetaKey: "js:detect-child-process",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		return walkIssues(ctx, func(node *sitter.Node) *domain.Issue {
			if node.Type() != "call_expression" {
				return nil
			}
			for i := 0; i < int(node.ChildCount()); i++ {
				child := node.Child(i)
				if child.Type() == "member_expression" {
					obj := ""
					fn := ""
					for j := 0; j < int(child.ChildCount()); j++ {
						c := child.Child(j)
						if c.Type() == "identifier" || c.Type() == "property_identifier" {
							txt := ctx.Query.Text(ctx.ParsedFile, c)
							if obj == "" {
								obj = txt
							} else {
								fn = txt
							}
						}
					}
					if (obj == "child_process" || obj == "cp") &&
						(fn == "exec" || fn == "execSync" || fn == "spawn" || fn == "spawnSync") {
						line, _, _, _ := ctx.Query.Position(node)
						issue := domain.NewIssue("js:detect-child-process", ctx.Path, line)
						issue.Severity = domain.SeverityMajor
						issue.Type = domain.TypeVulnerability
						issue.Message = fmt.Sprintf("%s can execute arbitrary commands; validate all inputs", fn)
						return issue
					}
				}
			}
			return nil
		})
	},
}
