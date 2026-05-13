package domain

import "github.com/scovl/ollanta/domain/model"

type Severity = model.Severity

const (
	SeverityBlocker  = model.SeverityBlocker
	SeverityCritical = model.SeverityCritical
	SeverityMajor    = model.SeverityMajor
	SeverityMinor    = model.SeverityMinor
	SeverityInfo     = model.SeverityInfo
)

type IssueType = model.IssueType

const (
	TypeBug             = model.TypeBug
	TypeVulnerability   = model.TypeVulnerability
	TypeCodeSmell       = model.TypeCodeSmell
	TypeSecurityHotspot = model.TypeSecurityHotspot
)

type IssueQualityDomain = model.IssueQualityDomain

const (
	QualitySecurity        = model.QualitySecurity
	QualityReliability     = model.QualityReliability
	QualityMaintainability = model.QualityMaintainability
	QualityTestability     = model.QualityTestability
)

type Status = model.Status

const (
	StatusOpen      = model.StatusOpen
	StatusConfirmed = model.StatusConfirmed
	StatusClosed    = model.StatusClosed
	StatusReopened  = model.StatusReopened
)

type SecondaryLocation = model.SecondaryLocation

// Issue represents a single quality finding produced by a rule during analysis.
// It extends the domain model Issue with the full set of fields for legacy compatibility.
type Issue struct {
	RuleKey            string              `json:"rule_key"`
	ComponentPath      string              `json:"component_path"`
	Line               int                 `json:"line"`
	Column             int                 `json:"column"`
	EndLine            int                 `json:"end_line"`
	EndColumn          int                 `json:"end_column"`
	Message            string              `json:"message"`
	Type               IssueType           `json:"type"`
	Severity           Severity            `json:"severity"`
	QualityDomain      IssueQualityDomain  `json:"quality_domain,omitempty"`
	Language           string              `json:"language,omitempty"`
	Status             Status              `json:"status"`
	Resolution         string              `json:"resolution,omitempty"`
	EffortMinutes      int                 `json:"effort_minutes,omitempty"`
	EngineID           string              `json:"engine_id,omitempty"`
	LineHash           string              `json:"line_hash,omitempty"`
	Tags               []string            `json:"tags,omitempty"`
	SecondaryLocations []SecondaryLocation `json:"secondary_locations"`
}

func NewIssue(ruleKey, componentPath string, line int) *Issue {
	mi := model.NewIssue(ruleKey, componentPath, line)
	return &Issue{
		RuleKey:            mi.RuleKey,
		ComponentPath:      mi.ComponentPath,
		Line:               mi.Line,
		Status:             Status(mi.Status),
		EngineID:           mi.EngineID,
		SecondaryLocations: mi.SecondaryLocations,
	}
}

func ComputeLineHash(fileContent string, line int) string {
	return model.ComputeLineHash(fileContent, line)
}
