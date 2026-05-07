// Package port defines the inbound and outbound interfaces (Ports) of the domain layer.
package port

import (
	"context"

	"github.com/scovl/ollanta/domain/model"
)

// IProfileRepo is the outbound port for quality profile persistence.
// IProfileRuleRepo, IProjectProfileRepo and IProfileChangelogRepo provide
// segregated views for consumers that only need a subset of operations.
type IProfileRepo interface {
	List(ctx context.Context, language string) ([]*model.QualityProfile, error)
	GetByID(ctx context.Context, id int64) (*model.QualityProfile, error)
	Create(ctx context.Context, p *model.QualityProfile) error
	Update(ctx context.Context, p *model.QualityProfile) error
	Delete(ctx context.Context, id int64) error
	Copy(ctx context.Context, sourceID int64, newName string) (*model.QualityProfile, error)
	SetDefault(ctx context.Context, id int64) error
	ActivateRule(ctx context.Context, profileID int64, ruleKey, severity string, params map[string]string) error
	DeactivateRule(ctx context.Context, profileID int64, ruleKey string) error
	AssignToProject(ctx context.Context, projectID int64, language string, profileID int64) error
	ByProjectAndLanguage(ctx context.Context, projectID int64, language string) (*model.QualityProfile, error)
	ResolveEffectiveRules(ctx context.Context, profileID int64) ([]*model.EffectiveRule, error)
	ProjectProfiles(ctx context.Context, projectID int64) ([]*model.ProjectQualityProfile, error)
	ProjectEffectiveProfiles(ctx context.Context, projectID int64) ([]*model.EffectiveQualityProfile, error)
	ProfileChangelog(ctx context.Context, profileID int64, limit, offset int) ([]model.ProfileChangelogEntry, int, error)
	ApplyProfileRules(ctx context.Context, profileID int64, entries []model.ProfileYAMLEntry) error
	ApplyProfileYAML(ctx context.Context, projectID int64, language string, entries []model.ProfileYAMLEntry) error
}

// IProfileRuleRepo manages rule activation/deactivation on a quality profile.
type IProfileRuleRepo interface {
	ActivateRule(ctx context.Context, profileID int64, ruleKey, severity string, params map[string]string) error
	DeactivateRule(ctx context.Context, profileID int64, ruleKey string) error
}

// IProjectProfileRepo manages quality profile assignment and resolution per project.
type IProjectProfileRepo interface {
	AssignToProject(ctx context.Context, projectID int64, language string, profileID int64) error
	ByProjectAndLanguage(ctx context.Context, projectID int64, language string) (*model.QualityProfile, error)
	ResolveEffectiveRules(ctx context.Context, profileID int64) ([]*model.EffectiveRule, error)
	ProjectProfiles(ctx context.Context, projectID int64) ([]*model.ProjectQualityProfile, error)
	ProjectEffectiveProfiles(ctx context.Context, projectID int64) ([]*model.EffectiveQualityProfile, error)
}

// IProfileChangelogRepo provides the audit changelog for quality profile modifications.
type IProfileChangelogRepo interface {
	ProfileChangelog(ctx context.Context, profileID int64, limit, offset int) ([]model.ProfileChangelogEntry, int, error)
}

// IProfileSnapshotRepo stores the quality profile policy snapshot attached to each scan.
type IProfileSnapshotRepo interface {
	Replace(ctx context.Context, projectID, scanID int64, scope model.AnalysisScope, snapshots []model.ProfileSnapshot) error
	ListByScan(ctx context.Context, scanID int64) ([]model.ProfileSnapshot, bool, error)
	HashChanges(ctx context.Context, currentScanID, previousScanID int64) ([]model.ProfileHashChange, error)
}
