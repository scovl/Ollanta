package rules

import (
	"github.com/scovl/ollanta/domain/model"
	"github.com/scovl/ollanta/domain/port"
	"github.com/scovl/ollanta/ollantarules"
	"github.com/scovl/ollanta/ollantarules/defaults"
)

// Registry wraps ollantarules.Registry to expose port.IAnalyzer slices.
type Registry struct {
	inner *ollantarules.Registry
}

// NewDefaultRegistry returns a Registry pre-loaded with all built-in rules.
func NewDefaultRegistry() *Registry {
	return &Registry{inner: defaults.NewRegistry()}
}

// NewRegistry wraps an existing ollantarules.Registry.
func NewRegistry(r *ollantarules.Registry) *Registry {
	return &Registry{inner: r}
}

// All returns all registered analyzers as port.IAnalyzer.
func (r *Registry) All() []port.IAnalyzer {
	src := r.inner.All()
	out := make([]port.IAnalyzer, len(src))
	for i, a := range src {
		out[i] = Wrap(a)
	}
	return out
}

// ForLanguage returns all analyzers targeting the given language.
func (r *Registry) ForLanguage(lang string) []port.IAnalyzer {
	src := r.inner.FindByLanguage(lang)
	out := make([]port.IAnalyzer, len(src))
	for i, a := range src {
		out[i] = Wrap(a)
	}
	return out
}

// Rules returns the static metadata of all registered rules.
func (r *Registry) Rules() []*model.Rule {
	coreRules := r.inner.Rules()
	out := make([]*model.Rule, len(coreRules))
	for i, cr := range coreRules {
		schema := make(map[string]model.ParamDef, len(cr.ParamsSchema))
		for k, p := range cr.ParamsSchema {
			schema[k] = model.ParamDef{
				Key:          p.Key,
				Description:  p.Description,
				DefaultValue: p.DefaultValue,
				Type:         p.Type,
			}
		}
		out[i] = &model.Rule{
			Key:             cr.Key,
			Name:            cr.Name,
			Description:     cr.Description,
			Language:        cr.Language,
			Type:            model.IssueType(cr.Type),
			DefaultSeverity: model.Severity(cr.DefaultSeverity),
			Tags:            cr.Tags,
			ParamsSchema:    schema,
		}
	}
	return out
}
