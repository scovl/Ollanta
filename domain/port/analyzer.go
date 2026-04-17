// Package port defines the inbound and outbound interfaces (Ports) of the domain layer.
// All interfaces in this package are purely behavioral: they carry no external library
// imports beyond the Go standard library and the domain/model package.
package port

import (
	"context"
	"go/ast"
	"go/token"

	"github.com/scovl/ollanta/domain/model"
)

// AnalysisContext carries everything an analyzer needs to inspect a single file.
//
// CGo-specific fields (ParsedFile, Query, Grammar) are typed as `any` so that the
// domain/ module remains free of CGo or tree-sitter imports. Adapters in
// adapter/secondary/rules/ perform a type assertion to the concrete types.
type AnalysisContext struct {
	// Path is the file path relative to the project root.
	Path string
	// Language is the canonical language identifier (e.g. "go", "javascript").
	Language string
	// Source is the raw source content of the file.
	Source []byte
	// ParsedFile is the tree-sitter ParsedFile; type is `any` to keep domain CGo-free.
	ParsedFile any
	// Query is the compiled tree-sitter query runner; type is `any`.
	Query any
	// Grammar is the tree-sitter Language grammar; type is `any`.
	Grammar any
	// GoFile is the Go AST for .go files; nil for other languages.
	GoFile *ast.File
	// GoFileSet is the token.FileSet for GoFile; nil for other languages.
	GoFileSet *token.FileSet
}

// IAnalyzer is the port that every language-specific rule must satisfy.
// Rules live in adapter/secondary/rules/ and are referenced only via this interface.
type IAnalyzer interface {
	// Key returns the unique rule key (e.g. "go:S1000").
	Key() string
	// Name returns the human-readable rule name.
	Name() string
	// Description returns the full rule description in Markdown.
	Description() string
	// Language returns the canonical language identifier this rule targets,
	// or "*" for cross-language rules.
	Language() string
	// Type returns the issue type (bug, vulnerability, code_smell, security_hotspot).
	Type() model.IssueType
	// DefaultSeverity returns the default severity assigned to issues raised by this rule.
	DefaultSeverity() model.Severity
	// Tags returns the taxonomy tags for this rule.
	Tags() []string
	// Params returns the configurable parameter schemas for this rule.
	Params() map[string]model.ParamDef
	// Check analyzes the file described by ctx and appends any issues to the provided slice.
	Check(ctx context.Context, ac AnalysisContext, issues *[]*model.Issue) error
}
