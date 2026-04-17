package model

import "time"

// QualityProfile is the canonical quality profile record.
type QualityProfile struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Language  string    `json:"language"`
	ParentID  *int64    `json:"parent_id,omitempty"`
	IsDefault bool      `json:"is_default"`
	IsBuiltin bool      `json:"is_builtin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProfileRule associates a rule activation to a quality profile.
type ProfileRule struct {
	ID        int64            `json:"id"`
	ProfileID int64            `json:"profile_id"`
	RuleKey   string           `json:"rule_key"`
	Severity  string           `json:"severity"`
	Params    map[string]string `json:"params,omitempty"`
}

// EffectiveRule is the resolved rule configuration (after profile inheritance).
type EffectiveRule struct {
	RuleKey  string            `json:"rule_key"`
	Severity string            `json:"severity"`
	Params   map[string]string `json:"params,omitempty"`
}

// ProfileYAMLEntry is used for bulk-loading profile rules from YAML.
type ProfileYAMLEntry struct {
	RuleKey  string            `yaml:"rule_key"`
	Severity string            `yaml:"severity"`
	Params   map[string]string `yaml:"params,omitempty"`
}
