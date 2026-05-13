package domain

import "github.com/scovl/ollanta/domain/model"

type ParamDef = model.ParamDef

type Threshold = model.Threshold

// Rule represents the static metadata of an analysis rule.
// It extends the domain model Rule with display-oriented fields (Rationale, code examples)
// that are consumed by the catalog and API presentation layer.
type Rule struct {
	Key              string              `json:"key"`
	Name             string              `json:"name"`
	Description      string              `json:"description"`
	Language         string              `json:"language"`
	Type             model.IssueType     `json:"type"`
	DefaultSeverity  model.Severity      `json:"default_severity"`
	Tags             []string            `json:"tags,omitempty"`
	ParamsSchema     map[string]ParamDef `json:"params_schema,omitempty"`
	Threshold        *Threshold          `json:"threshold,omitempty"`
	Rationale        string              `json:"rationale,omitempty"`
	NoncompliantCode string              `json:"noncompliant_code,omitempty"`
	CompliantCode    string              `json:"compliant_code,omitempty"`
	ReferenceURL     string              `json:"reference_url,omitempty"`
}
