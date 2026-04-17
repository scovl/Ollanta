// Package parser wraps tree-sitter via CGo to parse source files into
// concrete ParsedFile/QueryRunner instances, exposing them through the
// domain/port.IParser interface as opaque any values.
package parser

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/scovl/ollanta/domain/port"
)

// ParsedFile holds the result of parsing a single source file with tree-sitter.
type ParsedFile struct {
	Path     string
	Language string
	Source   []byte
	Tree     *sitter.Tree
	parser   *sitter.Parser
}

// Parse parses source using the given Language grammar and returns a ParsedFile.
func Parse(path string, source []byte, lang Language) (*ParsedFile, error) {
	p := sitter.NewParser()
	p.SetLanguage(lang.tsLanguage())

	tree, err := p.ParseCtx(context.Background(), nil, source)
	if err != nil {
		return nil, fmt.Errorf("parser: parsing %s: %w", path, err)
	}
	return &ParsedFile{
		Path:     path,
		Language: lang.Name(),
		Source:   source,
		Tree:     tree,
		parser:   p,
	}, nil
}

// ParseFile reads path from disk and parses it with the given Language grammar.
func ParseFile(path string, lang Language) (*ParsedFile, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("parser: reading %s: %w", path, err)
	}
	return Parse(path, source, lang)
}

// RootNode returns the root node of the syntax tree.
func (f *ParsedFile) RootNode() *sitter.Node {
	return f.Tree.RootNode()
}

// Close frees the internal tree-sitter parser.
func (f *ParsedFile) Close() {
	if f.parser != nil {
		f.parser.Close()
		f.parser = nil
	}
}

// ── Parser implements port.IParser ───────────────────────────────────────────

// Parser implements domain/port.IParser using the tree-sitter registry.
type Parser struct {
	registry *LanguageRegistry
}

// NewParser creates a Parser backed by the given language registry.
func NewParser(registry *LanguageRegistry) *Parser {
	return &Parser{registry: registry}
}

// compile-time interface check
var _ port.IParser = (*Parser)(nil)

// ParseFile reads and parses path, returning a *ParsedFile wrapped as any.
// The language param overrides registry lookup when non-empty.
func (p *Parser) ParseFile(path, language string) (any, error) {
	lang, ok := p.registry.ForFile(path)
	if !ok && language != "" {
		lang, ok = p.registry.ForName(language)
	}
	if !ok {
		return nil, fmt.Errorf("parser: no grammar registered for %s", filepath.Ext(path))
	}
	return ParseFile(path, lang)
}

// ParseSource parses src bytes for the given language, returning a *ParsedFile wrapped as any.
func (p *Parser) ParseSource(path string, src []byte, language string) (any, error) {
	lang, ok := p.registry.ForName(language)
	if !ok {
		lang, ok = p.registry.ForFile(path)
	}
	if !ok {
		return nil, fmt.Errorf("parser: no grammar registered for language %q", language)
	}
	return Parse(path, src, lang)
}
