// Package languages registers tree-sitter grammars for the parser adapter.
package languages

import (
	"github.com/scovl/ollanta/adapter/secondary/parser"
	javascript "github.com/smacker/go-tree-sitter/javascript"
	python "github.com/smacker/go-tree-sitter/python"
	rust "github.com/smacker/go-tree-sitter/rust"
	typescript "github.com/smacker/go-tree-sitter/typescript/typescript"
)

// JavaScript grammar
var JavaScript = parser.NewLanguage("javascript", []string{".js", ".mjs"}, javascript.GetLanguage())

// Python grammar
var Python = parser.NewLanguage("python", []string{".py"}, python.GetLanguage())

// TypeScript grammar
var TypeScript = parser.NewLanguage("typescript", []string{".ts", ".tsx"}, typescript.GetLanguage())

// Rust grammar
var Rust = parser.NewLanguage("rust", []string{".rs"}, rust.GetLanguage())

// DefaultRegistry returns a LanguageRegistry with all four built-in grammars.
func DefaultRegistry() *parser.LanguageRegistry {
	r := parser.NewRegistry()
	r.Register(JavaScript)
	r.Register(Python)
	r.Register(TypeScript)
	r.Register(Rust)
	return r
}
