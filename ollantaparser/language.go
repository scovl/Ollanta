// Package ollantaparser wraps the tree-sitter C library via CGo and exposes a
// Go-idiomatic API for multi-language parsing. It is the sole module in the
// Ollanta workspace that carries a CGo dependency; all other modules (ollantacore,
// ollantaengine, etc.) remain CGo-free by only depending on this package.
package ollantaparser

import (
	sitter "github.com/smacker/go-tree-sitter"
)

// Language represents a registered tree-sitter grammar.
// Create instances via NewLanguage — do not implement this interface directly.
type Language interface {
	// Name returns the canonical language identifier, e.g. "javascript".
	Name() string
	// Extensions returns the file-extension patterns associated with this language,
	// e.g. [".js", ".mjs"].
	Extensions() []string
	// tsLanguage returns the compiled tree-sitter language pointer for internal use.
	tsLanguage() *sitter.Language
}

// langImpl is the package-private implementation of Language.
type langImpl struct {
	name       string
	extensions []string
	lang       *sitter.Language
}

func (l *langImpl) Name() string                 { return l.name }
func (l *langImpl) Extensions() []string         { return l.extensions }
func (l *langImpl) tsLanguage() *sitter.Language { return l.lang }

// NewLanguage creates a Language from a *sitter.Language returned by a grammar package.
// e.g.:
//
//	var JavaScript = ollantaparser.NewLanguage("javascript", []string{".js", ".mjs"}, javascript.GetLanguage())
func NewLanguage(name string, extensions []string, lang *sitter.Language) Language {
	return &langImpl{
		name:       name,
		extensions: extensions,
		lang:       lang,
	}
}
