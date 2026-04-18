// Package ollantarules defines the Rule type — the extension point for all
// static analysis rules in Ollanta. Rules are dispatched by the GoSensor (go/ast)
// or the TreeSitterSensor (tree-sitter) depending on the language of the file
// being analysed.
package ollantarules

import (
	"go/ast"
	"go/token"

	"github.com/scovl/ollanta/ollantacore/domain"
	"github.com/scovl/ollanta/ollantaparser"
)

// RuleMeta carries the declarative metadata for a rule.
// It is a plain data struct — accessible, serialisable, and comparable
// without instantiating or calling any methods.
type RuleMeta struct {
	// Key is the unique rule identifier, e.g. "go:no-large-functions".
	Key string `json:"key"`
	// Name is a human-readable rule name.
	Name string `json:"name"`
	// Description explains what the rule detects and why it matters.
	Description string `json:"description"`
	// Language is the target language, e.g. "go", "javascript", or "*" for cross-language.
	Language string `json:"language"`
	// Type is the issue category produced by this rule.
	Type domain.IssueType `json:"type"`
	// DefaultSeverity is the default severity for issues found by this rule.
	DefaultSeverity domain.Severity `json:"severity"`
	// Tags are categorisation labels.
	Tags []string `json:"tags,omitempty"`
	// Params is the list of configurable parameters with defaults.
	Params []domain.ParamDef `json:"params,omitempty"`
}

// CheckFunc is the execution signature for a rule.
// It receives the per-file context and returns any issues found.
// Implementations must be safe for concurrent use and must not retain
// any mutable state between calls.
type CheckFunc func(ctx *AnalysisContext) []*domain.Issue

// Rule is the unit of analysis in Ollanta: declarative metadata paired with
// imperative detection logic. Rules are declared as package-level var values
// and registered in a Registry.
type Rule struct {
	// MetaKey is the JSON metadata key this rule maps to (e.g. "go:cognitive-complexity").
	// Set at declaration time; used by MustRegister to bind metadata automatically.
	MetaKey string
	// Meta holds all static, serialisable rule information.
	Meta RuleMeta
	// Check executes the rule against a file context and returns issues.
	Check CheckFunc
}

// AnalysisContext carries the per-file context passed to each CheckFunc.
// Exactly one of AST or ParsedFile will be non-nil, depending on which sensor
// is invoking the rule.
type AnalysisContext struct {
	// Path is the relative path of the file being analysed.
	Path string
	// Source is the raw file content.
	Source []byte
	// Language is the canonical language identifier, e.g. "go", "javascript".
	Language string
	// Params holds the configured parameter values for this rule invocation.
	Params map[string]string

	// AST and FileSet are populated by GoSensor (Language == "go" only).
	AST     *ast.File
	FileSet *token.FileSet

	// ParsedFile and Query are populated by TreeSitterSensor (all non-Go languages).
	ParsedFile *ollantaparser.ParsedFile
	Query      *ollantaparser.QueryRunner
	// Grammar is the tree-sitter Language grammar for the file being analysed.
	// Required by Query.Run to compile S-expression queries.
	Grammar ollantaparser.Language
}
