package ollantarules

import (
	"fmt"
	"io/fs"

	"github.com/scovl/ollanta/ollantacore/domain"
)

// Registry holds the set of registered Rule values and provides lookup methods.
type Registry struct {
	rules []Rule
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Register adds a Rule to the registry.
// It panics if Meta.Key is empty, Meta.Language is empty, or Check is nil.
func (r *Registry) Register(rule Rule) {
	if rule.Meta.Key == "" {
		panic("ollantarules: Register called with empty Meta.Key")
	}
	if rule.Meta.Language == "" {
		panic(fmt.Sprintf("ollantarules: Register called with empty Meta.Language for rule %q", rule.Meta.Key))
	}
	if rule.Check == nil {
		panic(fmt.Sprintf("ollantarules: Register called with nil Check for rule %q", rule.Meta.Key))
	}
	r.rules = append(r.rules, rule)
}

// All returns a defensive copy of all registered rules.
func (r *Registry) All() []Rule {
	out := make([]Rule, len(r.rules))
	copy(out, r.rules)
	return out
}

// FindByKey returns a pointer to the Rule with the given key, or nil if not found.
func (r *Registry) FindByKey(key string) *Rule {
	for i := range r.rules {
		if r.rules[i].Meta.Key == key {
			return &r.rules[i]
		}
	}
	return nil
}

// FindByLanguage returns all rules targeting the given language,
// including cross-language rules (Language == "*").
func (r *Registry) FindByLanguage(lang string) []Rule {
	var out []Rule
	for _, rule := range r.rules {
		if rule.Meta.Language == lang || rule.Meta.Language == "*" {
			out = append(out, rule)
		}
	}
	return out
}

// Rules converts all registered rules to domain.Rule structs using metadata only.
func (r *Registry) Rules() []*domain.Rule {
	out := make([]*domain.Rule, len(r.rules))
	for i, rule := range r.rules {
		schema := make(map[string]domain.ParamDef, len(rule.Meta.Params))
		for _, p := range rule.Meta.Params {
			schema[p.Key] = p
		}
		out[i] = &domain.Rule{
			Key:             rule.Meta.Key,
			Name:            rule.Meta.Name,
			Description:     rule.Meta.Description,
			Language:        rule.Meta.Language,
			Type:            rule.Meta.Type,
			DefaultSeverity: rule.Meta.DefaultSeverity,
			Tags:            rule.Meta.Tags,
			ParamsSchema:    schema,
		}
	}
	return out
}

// globalRegistry is the package-level registry populated by init() in rule packages.
var globalRegistry Registry

// Global returns the global registry populated at init time by rule packages.
func Global() *Registry {
	return &globalRegistry
}

// MustRegister loads JSON metadata from fsys, binds each rule's MetaKey to its
// RuleMeta, and appends the rules to the global registry. It panics if metadata
// cannot be loaded, a rule's MetaKey is missing from the JSON, or a duplicate
// key is registered.
func MustRegister(fsys fs.FS, pattern string, rules ...Rule) {
	meta, err := LoadMeta(fsys, pattern)
	if err != nil {
		panic("ollantarules.MustRegister: " + err.Error())
	}
	for _, rule := range rules {
		m, ok := meta[rule.MetaKey]
		if !ok {
			panic(fmt.Sprintf("ollantarules.MustRegister: no metadata for key %q", rule.MetaKey))
		}
		rule.Meta = m
		if globalRegistry.FindByKey(rule.Meta.Key) != nil {
			panic(fmt.Sprintf("ollantarules.MustRegister: duplicate key %q", rule.Meta.Key))
		}
		globalRegistry.Register(rule)
	}
}
