package rules

import (
	"go/ast"

	"github.com/scovl/ollanta/ollantacore/domain"
	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// LoopPointer detects capturing a loop variable by reference inside a goroutine.
// This is a classic Go bug where all goroutines see the final value of the loop variable.
var LoopPointer = ollantarules.Rule{
	MetaKey: "go:loop-pointer",
	Check: func(ctx *ollantarules.AnalysisContext) []*domain.Issue {
		var issues []*domain.Issue
		ast.Inspect(ctx.AST, func(n ast.Node) bool {
			forStmt, ok := n.(*ast.ForStmt)
			if !ok {
				// Also check range statements
				rangeStmt, ok := n.(*ast.RangeStmt)
				if !ok {
					return true
				}
				issues = append(issues, checkRangeStmt(ctx, rangeStmt)...)
				return true
			}
			issues = append(issues, checkForStmt(ctx, forStmt)...)
			return true
		})
		return issues
	},
}

func checkForStmt(ctx *ollantarules.AnalysisContext, stmt *ast.ForStmt) []*domain.Issue {
	var issues []*domain.Issue
	loopVars := map[string]bool{}
	if init, ok := stmt.Init.(*ast.AssignStmt); ok {
		for _, lhs := range init.Lhs {
			if id, ok := lhs.(*ast.Ident); ok {
				loopVars[id.Name] = true
			}
		}
	}
	if post, ok := stmt.Post.(*ast.IncDecStmt); ok {
		if id, ok := post.X.(*ast.Ident); ok {
			loopVars[id.Name] = true
		}
	}
	ast.Inspect(stmt.Body, func(n ast.Node) bool {
		goStmt, ok := n.(*ast.GoStmt)
		if !ok {
			return true
		}
		fnLit, ok := goStmt.Call.Fun.(*ast.FuncLit)
		if !ok {
			return true
		}
		// Check if the goroutine references any loop variable
		ast.Inspect(fnLit.Body, func(inner ast.Node) bool {
			id, ok := inner.(*ast.Ident)
			if !ok {
				return true
			}
			if loopVars[id.Name] && !isParamOfFuncLit(id, fnLit) {
				line := ctx.FileSet.Position(id.Pos()).Line
				issue := domain.NewIssue("go:loop-pointer", ctx.Path, line)
				issue.Severity = domain.SeverityMajor
				issue.Type = domain.TypeBug
				issue.Message = "Loop variable captured by goroutine closure; pass it as an argument instead"
				issues = append(issues, issue)
				return false
			}
			return true
		})
		return true
	})
	return issues
}

func checkRangeStmt(ctx *ollantarules.AnalysisContext, stmt *ast.RangeStmt) []*domain.Issue {
	var issues []*domain.Issue
	loopVars := map[string]bool{}
	if id, ok := stmt.Key.(*ast.Ident); ok && id.Name != "_" {
		loopVars[id.Name] = true
	}
	if id, ok := stmt.Value.(*ast.Ident); ok && id.Name != "_" {
		loopVars[id.Name] = true
	}
	ast.Inspect(stmt.Body, func(n ast.Node) bool {
		goStmt, ok := n.(*ast.GoStmt)
		if !ok {
			return true
		}
		fnLit, ok := goStmt.Call.Fun.(*ast.FuncLit)
		if !ok {
			return true
		}
		ast.Inspect(fnLit.Body, func(inner ast.Node) bool {
			id, ok := inner.(*ast.Ident)
			if !ok {
				return true
			}
			if loopVars[id.Name] && !isParamOfFuncLit(id, fnLit) {
				line := ctx.FileSet.Position(id.Pos()).Line
				issue := domain.NewIssue("go:loop-pointer", ctx.Path, line)
				issue.Severity = domain.SeverityMajor
				issue.Type = domain.TypeBug
				issue.Message = "Range variable captured by goroutine closure; pass it as an argument instead"
				issues = append(issues, issue)
				return false
			}
			return true
		})
		return true
	})
	return issues
}

func isParamOfFuncLit(id *ast.Ident, fnLit *ast.FuncLit) bool {
	if fnLit.Type == nil || fnLit.Type.Params == nil {
		return false
	}
	for _, p := range fnLit.Type.Params.List {
		for _, n := range p.Names {
			if n.Name == id.Name {
				return true
			}
		}
	}
	return false
}
