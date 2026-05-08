package treesitter

import (
	"strings"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
	sitter "github.com/smacker/go-tree-sitter"
)

// SpawnGitCloneJS detects spawn/exec of 'git clone' with untrusted arguments.
var SpawnGitCloneJS = ollantarules.Rule{
	MetaKey: "js:spawn-git-clone",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		return walkIssues(ctx, func(node *sitter.Node) *domain.Issue {
			if node.Type() != "call_expression" {
				return nil
			}
			fnName := ""
			for i := 0; i < int(node.ChildCount()); i++ {
				child := node.Child(i)
				if child.Type() == "identifier" {
					fnName = ctx.Query.Text(ctx.ParsedFile, child)
					break
				}
			}
			if fnName != "spawn" && fnName != "exec" && fnName != "execSync" {
				return nil
			}
			// Look for 'git' and 'clone' in arguments
			hasGit := false
			hasClone := false
			var walk func(*sitter.Node)
			walk = func(n *sitter.Node) {
				if n == nil {
					return
				}
				if n.Type() == "string_fragment" || n.Type() == "string" {
					txt := strings.Trim(ctx.Query.Text(ctx.ParsedFile, n), "'\"")
					if txt == "git" {
						hasGit = true
					}
					if txt == "clone" {
						hasClone = true
					}
				}
				for j := 0; j < int(n.ChildCount()); j++ {
					walk(n.Child(j))
				}
			}
			walk(node)
			if !hasGit || !hasClone {
				return nil
			}
			line, _, _, _ := ctx.Query.Position(node)
			issue := domain.NewIssue("js:spawn-git-clone", ctx.Path, line)
			issue.Severity = domain.SeverityMajor
			issue.Type = domain.TypeVulnerability
			issue.Message = "Spawning git clone with untrusted arguments can lead to command injection"
			return issue
		})
	},
}
