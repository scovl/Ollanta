package parser

import (
	"path/filepath"
	"sync"
)

// LanguageRegistry maps file extensions and language names to tree-sitter grammars.
type LanguageRegistry struct {
	mu     sync.RWMutex
	byExt  map[string]Language
	byName map[string]Language
}

// NewRegistry creates an empty LanguageRegistry.
func NewRegistry() *LanguageRegistry {
	return &LanguageRegistry{
		byExt:  make(map[string]Language),
		byName: make(map[string]Language),
	}
}

// Register adds lang for all its declared extensions and name.
func (r *LanguageRegistry) Register(lang Language) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byName[lang.Name()] = lang
	for _, ext := range lang.Extensions() {
		r.byExt[ext] = lang
	}
}

// ForExtension returns the Language registered for the given file extension.
func (r *LanguageRegistry) ForExtension(ext string) (Language, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	l, ok := r.byExt[ext]
	return l, ok
}

// ForName returns the Language registered under the given canonical name.
func (r *LanguageRegistry) ForName(name string) (Language, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	l, ok := r.byName[name]
	return l, ok
}

// ForFile looks up language by file extension.
func (r *LanguageRegistry) ForFile(path string) (Language, bool) {
	return r.ForExtension(filepath.Ext(path))
}
