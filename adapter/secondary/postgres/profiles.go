package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/scovl/ollanta/domain/model"
	"github.com/scovl/ollanta/domain/port"
)

// ProfileRepository provides CRUD access to quality_profiles and related tables.
type ProfileRepository struct {
	db *DB
}

// NewProfileRepository creates a ProfileRepository backed by db.
func NewProfileRepository(db *DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// compile-time interface check
var _ port.IProfileRepo = (*ProfileRepository)(nil)

// List returns all profiles, optionally filtered by language (empty = all).
func (r *ProfileRepository) List(ctx context.Context, language string) ([]*model.QualityProfile, error) {
	q := `SELECT id, name, language, parent_id, is_default, is_builtin, created_at, updated_at
	      FROM quality_profiles`
	args := []any{}
	if language != "" {
		q += " WHERE language = $1"
		args = append(args, language)
	}
	q += " ORDER BY name"

	rows, err := r.db.Pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProfiles(rows)
}

// GetByID returns a profile by its ID.
func (r *ProfileRepository) GetByID(ctx context.Context, id int64) (*model.QualityProfile, error) {
	p := &model.QualityProfile{}
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, name, language, parent_id, is_default, is_builtin, created_at, updated_at
		FROM quality_profiles WHERE id = $1`, id,
	).Scan(&p.ID, &p.Name, &p.Language, &p.ParentID, &p.IsDefault, &p.IsBuiltin, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	return p, err
}

// Create inserts a new quality profile. Validates inheritance depth (max 3).
func (r *ProfileRepository) Create(ctx context.Context, p *model.QualityProfile) error {
	if p.ParentID != nil {
		depth, err := r.inheritanceDepth(ctx, *p.ParentID)
		if err != nil {
			return fmt.Errorf("check inheritance depth: %w", err)
		}
		if depth >= 3 {
			return fmt.Errorf("inheritance chain exceeds maximum of 3 levels")
		}
	}
	return r.db.Pool.QueryRow(ctx, `
		INSERT INTO quality_profiles (name, language, parent_id, is_default, is_builtin)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`,
		p.Name, p.Language, p.ParentID, p.IsDefault, p.IsBuiltin,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

// Update updates name, parent_id and is_default of a profile.
func (r *ProfileRepository) Update(ctx context.Context, p *model.QualityProfile) error {
	if p.IsBuiltin {
		return fmt.Errorf("cannot update builtin profile")
	}
	if p.ParentID != nil {
		depth, err := r.inheritanceDepth(ctx, *p.ParentID)
		if err != nil {
			return fmt.Errorf("check inheritance depth: %w", err)
		}
		if depth >= 3 {
			return fmt.Errorf("inheritance chain exceeds maximum of 3 levels")
		}
	}
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE quality_profiles
		SET name = $1, parent_id = $2, is_default = $3, updated_at = now()
		WHERE id = $4`,
		p.Name, p.ParentID, p.IsDefault, p.ID)
	return err
}

// Delete removes a non-builtin profile.
func (r *ProfileRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Pool.Exec(ctx,
		`DELETE FROM quality_profiles WHERE id = $1 AND is_builtin = FALSE`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("profile not found or is builtin")
	}
	return nil
}

// Copy duplicates a profile with a new name, including all its rules.
func (r *ProfileRepository) Copy(ctx context.Context, sourceID int64, newName string) (*model.QualityProfile, error) {
	src, err := r.GetByID(ctx, sourceID)
	if err != nil {
		return nil, err
	}
	newProfile := &model.QualityProfile{
		Name:     newName,
		Language: src.Language,
		ParentID: src.ParentID,
	}
	if err := r.Create(ctx, newProfile); err != nil {
		return nil, fmt.Errorf("create copy: %w", err)
	}
	rules, err := r.listRules(ctx, sourceID)
	if err != nil {
		return nil, fmt.Errorf("read rules: %w", err)
	}
	for _, rule := range rules {
		if err := r.ActivateRule(ctx, newProfile.ID, rule.RuleKey, rule.Severity, rule.Params); err != nil {
			return nil, fmt.Errorf("copy rule: %w", err)
		}
	}
	return newProfile, nil
}

// SetDefault atomically sets a profile as default for its language.
func (r *ProfileRepository) SetDefault(ctx context.Context, id int64) error {
	p, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.db.Pool.Exec(ctx,
		`UPDATE quality_profiles SET is_default = (id = $1), updated_at = now()
		 WHERE language = $2 AND (is_default = TRUE OR id = $1)`, id, p.Language)
	return err
}

// ActivateRule adds or updates an active rule in a profile.
func (r *ProfileRepository) ActivateRule(ctx context.Context, profileID int64, ruleKey, severity string, params map[string]string) error {
	if params == nil {
		params = map[string]string{}
	}
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO quality_profile_rules (profile_id, rule_key, severity, params)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (profile_id, rule_key) DO UPDATE
		  SET severity = EXCLUDED.severity, params = EXCLUDED.params`,
		profileID, ruleKey, severity, params)
	return err
}

// DeactivateRule removes a rule from a profile.
func (r *ProfileRepository) DeactivateRule(ctx context.Context, profileID int64, ruleKey string) error {
	return r.ActivateRule(ctx, profileID, ruleKey, "OFF", nil)
}

// AssignToProject sets the active profile for a project+language combination.
func (r *ProfileRepository) AssignToProject(ctx context.Context, projectID int64, language string, profileID int64) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO project_profiles (project_id, language, profile_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (project_id, language) DO UPDATE SET profile_id = EXCLUDED.profile_id`,
		projectID, language, profileID)
	return err
}

// ByProjectAndLanguage returns the active profile for a project+language, falling back to the default.
func (r *ProfileRepository) ByProjectAndLanguage(ctx context.Context, projectID int64, language string) (*model.QualityProfile, error) {
	var profileID int64
	err := r.db.Pool.QueryRow(ctx, `
		SELECT profile_id FROM project_profiles WHERE project_id = $1 AND language = $2`,
		projectID, language,
	).Scan(&profileID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		err = r.db.Pool.QueryRow(ctx, `
			SELECT id FROM quality_profiles WHERE language = $1 AND is_default = TRUE LIMIT 1`,
			language,
		).Scan(&profileID)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		if err != nil {
			return nil, err
		}
	}
	return r.GetByID(ctx, profileID)
}

// ResolveEffectiveRules walks the inheritance chain and returns the union of active rules.
func (r *ProfileRepository) ResolveEffectiveRules(ctx context.Context, profileID int64) ([]*model.EffectiveRule, error) {
	chain, err := r.buildInheritanceChain(ctx, profileID)
	if err != nil {
		return nil, err
	}

	merged := map[string]*model.EffectiveRule{}
	for _, pid := range chain {
		rules, err := r.listRules(ctx, pid)
		if err != nil {
			return nil, fmt.Errorf("list rules for profile %d: %w", pid, err)
		}
		for _, rule := range rules {
			if rule.Severity == "OFF" {
				delete(merged, rule.RuleKey)
			} else {
				merged[rule.RuleKey] = &model.EffectiveRule{
					RuleKey:  rule.RuleKey,
					Severity: rule.Severity,
					Params:   rule.Params,
				}
			}
		}
	}

	out := make([]*model.EffectiveRule, 0, len(merged))
	for _, r := range merged {
		out = append(out, r)
	}
	return out, nil
}

// ApplyProfileYAML applies a profile-as-code YAML payload transactionally.
func (r *ProfileRepository) ApplyProfileYAML(ctx context.Context, projectID int64, language string, entries []model.ProfileYAMLEntry) error {
	profile, err := r.ByProjectAndLanguage(ctx, projectID, language)
	if err != nil {
		return fmt.Errorf("resolve profile for project %d language %s: %w", projectID, language, err)
	}
	for _, e := range entries {
		sev := e.Severity
		if sev == "" {
			sev = "major"
		}
		if err := r.ActivateRule(ctx, profile.ID, e.RuleKey, sev, e.Params); err != nil {
			return fmt.Errorf("activate rule %s: %w", e.RuleKey, err)
		}
	}
	return nil
}

// ── helpers ──────────────────────────────────────────────────────────────────

type profileRule struct {
	ID        int64
	ProfileID int64
	RuleKey   string
	Severity  string
	Params    map[string]string
}

func (r *ProfileRepository) buildInheritanceChain(ctx context.Context, profileID int64) ([]int64, error) {
	chain := []int64{}
	visited := map[int64]bool{}
	id := profileID
	for {
		if visited[id] {
			return nil, fmt.Errorf("cycle detected in profile inheritance at id=%d", id)
		}
		visited[id] = true
		chain = append([]int64{id}, chain...)
		var parentID *int64
		err := r.db.Pool.QueryRow(ctx, `SELECT parent_id FROM quality_profiles WHERE id = $1`, id).Scan(&parentID)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		if err != nil {
			return nil, err
		}
		if parentID == nil {
			break
		}
		id = *parentID
	}
	return chain, nil
}

func (r *ProfileRepository) inheritanceDepth(ctx context.Context, profileID int64) (int, error) {
	chain, err := r.buildInheritanceChain(ctx, profileID)
	if err != nil {
		return 0, err
	}
	return len(chain), nil
}

func (r *ProfileRepository) listRules(ctx context.Context, profileID int64) ([]*profileRule, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, profile_id, rule_key, severity, params
		FROM quality_profile_rules WHERE profile_id = $1`, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*profileRule
	for rows.Next() {
		rule := &profileRule{}
		if err := rows.Scan(&rule.ID, &rule.ProfileID, &rule.RuleKey, &rule.Severity, &rule.Params); err != nil {
			return nil, err
		}
		out = append(out, rule)
	}
	return out, rows.Err()
}

func scanProfiles(rows pgx.Rows) ([]*model.QualityProfile, error) {
	var out []*model.QualityProfile
	for rows.Next() {
		p := &model.QualityProfile{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Language, &p.ParentID,
			&p.IsDefault, &p.IsBuiltin, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
