// Package port defines the inbound and outbound interfaces (Ports) of the domain layer.
package port

import (
	"context"

	"github.com/scovl/ollanta/domain/model"
)

// IProfileRepo is the outbound port for quality profile persistence.
type IProfileRepo interface {
	List(ctx context.Context, language string) ([]*model.QualityProfile, error)
	GetByID(ctx context.Context, id int64) (*model.QualityProfile, error)
	Create(ctx context.Context, p *model.QualityProfile) error
	Update(ctx context.Context, p *model.QualityProfile) error
	Delete(ctx context.Context, id int64) error
	ActivateRule(ctx context.Context, profileID int64, ruleKey, severity string, params map[string]string) error
	DeactivateRule(ctx context.Context, profileID int64, ruleKey string) error
	AssignToProject(ctx context.Context, projectID int64, language string, profileID int64) error
	ByProjectAndLanguage(ctx context.Context, projectID int64, language string) (*model.QualityProfile, error)
	ResolveEffectiveRules(ctx context.Context, profileID int64) ([]*model.EffectiveRule, error)
	ApplyProfileYAML(ctx context.Context, projectID int64, language string, entries []model.ProfileYAMLEntry) error
}
