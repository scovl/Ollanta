package parser

import (
	sitter "github.com/smacker/go-tree-sitter"
)

// Language represents a registered tree-sitter grammar.
type Language interface {
	Name() string
	Extensions() []string
	tsLanguage() *sitter.Language
}

type langImpl struct {
	name       string
	extensions []string
	lang       *sitter.Language
}

func (l *langImpl) Name() string                 { return l.name }
func (l *langImpl) Extensions() []string         { return l.extensions }
func (l *langImpl) tsLanguage() *sitter.Language { return l.lang }

// NewLanguage creates a Language from a *sitter.Language.
func NewLanguage(name string, extensions []string, lang *sitter.Language) Language {
	return &langImpl{name: name, extensions: extensions, lang: lang}
}
