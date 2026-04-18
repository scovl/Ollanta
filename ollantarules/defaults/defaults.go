// Package defaults provides a pre-configured Registry with all built-in
// Ollanta rules registered. Importing this package triggers init() in each
// rule package, which populates the global registry automatically.
package defaults

import (
	ollantarules "github.com/scovl/ollanta/ollantarules"

	// Blank imports trigger init()-based rule registration.
	_ "github.com/scovl/ollanta/ollantarules/languages/golang/rules"
	_ "github.com/scovl/ollanta/ollantarules/languages/treesitter"
)

// NewRegistry returns the global Registry pre-loaded with all built-in rules.
// Backward-compatible: delegates to ollantarules.Global().
func NewRegistry() *ollantarules.Registry {
	return ollantarules.Global()
}
