package treesitter

import (
	"strings"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
	sitter "github.com/smacker/go-tree-sitter"
)

// UseDefusedXmlPY detects use of xml.etree.ElementTree.parse instead of defusedxml.
var UseDefusedXmlPY = ollantarules.Rule{
	MetaKey: "py:use-defused-xml",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		// Check if xml.etree.ElementTree is imported anywhere in the file
		hasETImport := false
		var walkAll func(*sitter.Node)
		walkAll = func(n *sitter.Node) {
			if n == nil {
				return
			}
			txt := strings.TrimSpace(ctx.Query.Text(ctx.ParsedFile, n))
			if strings.Contains(txt, "xml.etree.ElementTree") {
				hasETImport = true
			}
			for i := 0; i < int(n.ChildCount()); i++ {
				walkAll(n.Child(i))
			}
		}
		walkAll(ctx.ParsedFile.RootNode())
		if !hasETImport {
			return nil
		}
		return walkIssues(ctx, func(node *sitter.Node) *domain.Issue {
			if node.Type() != "call" {
				return nil
			}
			for i := 0; i < int(node.ChildCount()); i++ {
				child := node.Child(i)
				if child.Type() == "attribute" {
					obj := ""
					fn := ""
					for j := 0; j < int(child.ChildCount()); j++ {
						c := child.Child(j)
						if c.Type() == "identifier" {
							if obj == "" {
								obj = ctx.Query.Text(ctx.ParsedFile, c)
							} else {
								fn = ctx.Query.Text(ctx.ParsedFile, c)
							}
						}
					}
					if fn == "parse" && (obj == "ElementTree" || obj == "ET") {
						line, _, _, _ := ctx.Query.Position(node)
						issue := domain.NewIssue("py:use-defused-xml", ctx.Path, line)
						issue.Severity = domain.SeverityMajor
						issue.Type = domain.TypeVulnerability
						issue.Message = "xml.etree.ElementTree is vulnerable to XML bombs; use defusedxml.ElementTree instead"
						return issue
					}
				}
			}
			return nil
		})
	},
}
