package ollantaparser

import (
	"context"
	"fmt"
	"os"

	sitter "github.com/smacker/go-tree-sitter"
)

// ParsedFile holds the result of parsing a single source file with tree-sitter.
// It is the multi-language equivalent of *ast.File from go/ast.
// Call Close when done to free the internal parser memory.
type ParsedFile struct {
	Path     string
	Language string
	Source   []byte
	Tree     *sitter.Tree
	parser   *sitter.Parser
}

// Parse parses source using the given Language grammar and returns a ParsedFile.
// path is used only for diagnostics (error messages, Issue.ComponentPath).
func Parse(path string, source []byte, lang Language) (*ParsedFile, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(lang.tsLanguage())

	tree, err := parser.ParseCtx(context.Background(), nil, source)
	if err != nil {
		return nil, fmt.Errorf("ollantaparser: parsing %s: %w", path, err)
	}
	return &ParsedFile{
		Path:     path,
		Language: lang.Name(),
		Source:   source,
		Tree:     tree,
		parser:   parser,
	}, nil
}

// ParseFile reads path from disk and parses it with the given Language grammar.
func ParseFile(path string, lang Language) (*ParsedFile, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ollantaparser: reading %s: %w", path, err)
	}
	return Parse(path, source, lang)
}

// RootNode returns the root node of the syntax tree. Convenience shortcut.
func (f *ParsedFile) RootNode() *sitter.Node {
	return f.Tree.RootNode()
}

// Close frees the internal tree-sitter parser. Should be called when the
// ParsedFile is no longer needed to avoid CGo memory leaks.
func (f *ParsedFile) Close() {
	if f.parser != nil {
		f.parser.Close()
		f.parser = nil
	}
}
