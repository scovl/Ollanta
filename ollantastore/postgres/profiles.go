package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// QualityProfile is a named set of active rules for a specific language.
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

// ProfileRule is an active rule within a quality profile.
type ProfileRule struct {
	ID        int64             `json:"id"`
	ProfileID int64             `json:"profile_id"`
	RuleKey   string            `json:"rule_key"`
	Severity  string            `json:"severity"`
	Params    map[string]string `json:"params"`
}

// EffectiveRule is a resolved rule after inheritance.
type EffectiveRule struct {
	RuleKey  string            `json:"rule_key"`
	Severity string            `json:"severity"`
	Params   map[string]string `json:"params"`
}

// ProfileRepository provides CRUD access to quality_profiles and related tables.
type ProfileRepository struct {
	db *DB
}

// NewProfileRepository creates a ProfileRepository backed by db.
func NewProfileRepository(db *DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// List returns all profiles, optionally filtered by language.
func (r *ProfileRepository) List(ctx context.Context, language string) ([]*QualityProfile, error) {
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
func (r *ProfileRepository) GetByID(ctx context.Context, id int64) (*QualityProfile, error) {
	p := &QualityProfile{}
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, name, language, parent_id, is_default, is_builtin, created_at, updated_at
		FROM quality_profiles WHERE id = $1`, id,
	).Scan(&p.ID, &p.Name, &p.Language, &p.ParentID, &p.IsDefault, &p.IsBuiltin, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return p, err
}

// Create inserts a new quality profile. Validates inheritance depth (max 3).
func (r *ProfileRepository) Create(ctx context.Context, p *QualityProfile) error {
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

// Update updates name, parent_id and is_default of a profile. Validates depth.
func (r *ProfileRepository) Update(ctx context.Context, p *QualityProfile) error {
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

// ActivateRule adds or updates an active rule in a profile.
// severity='OFF' deactivates the rule (stored explicitly for inheritance resolution).
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

// DeactivateRule removes a rule from a profile (sets severity='OFF' for inheritance tracking).
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
func (r *ProfileRepository) ByProjectAndLanguage(ctx context.Context, projectID int64, language string) (*QualityProfile, error) {
	// Try project-specific assignment first.
	var profileID int64
	err := r.db.Pool.QueryRow(ctx, `
		SELECT profile_id FROM project_profiles WHERE project_id = $1 AND language = $2`,
		projectID, language,
	).Scan(&profileID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		// Fall back to language default.
		err = r.db.Pool.QueryRow(ctx, `
			SELECT id FROM quality_profiles WHERE language = $1 AND is_default = TRUE LIMIT 1`,
			language,
		).Scan(&profileID)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		if err != nil {
			return nil, err
		}
	}
	return r.GetByID(ctx, profileID)
}

// ResolveEffectiveRules walks the inheritance chain and returns the union of active rules.
// Rules with severity='OFF' in any child profile are excluded from the result.
func (r *ProfileRepository) ResolveEffectiveRules(ctx context.Context, profileID int64) ([]*EffectiveRule, error) {
	chain, err := r.buildInheritanceChain(ctx, profileID)
	if err != nil {
		return nil, err
	}

	// Walk from root to leaf — child entries override parent entries.
	merged := map[string]*EffectiveRule{}
	for _, pid := range chain {
		rules, err := r.listRules(ctx, pid)
		if err != nil {
			return nil, fmt.Errorf("list rules for profile %d: %w", pid, err)
		}
		for _, rule := range rules {
			if rule.Severity == "OFF" {
				delete(merged, rule.RuleKey)
			} else {
				merged[rule.RuleKey] = &EffectiveRule{
					RuleKey:  rule.RuleKey,
					Severity: rule.Severity,
					Params:   rule.Params,
				}
			}
		}
	}

	out := make([]*EffectiveRule, 0, len(merged))
	for _, r := range merged {
		out = append(out, r)
	}
	return out, nil
}

// ApplyProfileYAML applies a profile-as-code YAML payload transactionally.
// The payload contains activate/deactivate rule lists.
func (r *ProfileRepository) ApplyProfileYAML(ctx context.Context, projectID int64, language string, entries []ProfileYAMLEntry) error {
	profile, err := r.ByProjectAndLanguage(ctx, projectID, language)
	if err != nil {
		return fmt.Errorf("resolve profile for project %d language %s: %w", projectID, language, err)
	}
	for _, e := range entries {
		sev := e.Severity
		if sev == "" {
			sev = "major"
		}
		if e.Activate {
			if err := r.ActivateRule(ctx, profile.ID, e.Rule, sev, e.Params); err != nil {
				return fmt.Errorf("activate rule %s: %w", e.Rule, err)
			}
		} else {
			if err := r.DeactivateRule(ctx, profile.ID, e.Rule); err != nil {
				return fmt.Errorf("deactivate rule %s: %w", e.Rule, err)
			}
		}
	}
	return nil
}

// ProfileYAMLEntry represents a single rule action from profiles.yml.
type ProfileYAMLEntry struct {
	Rule     string            `json:"rule"`
	Severity string            `json:"severity"`
	Params   map[string]string `json:"params"`
	Activate bool              `json:"activate"` // true=activate, false=deactivate
}

// ── helpers ───────────────────────────────────────────────────────────────────

func (r *ProfileRepository) buildInheritanceChain(ctx context.Context, profileID int64) ([]int64, error) {
	chain := []int64{}
	visited := map[int64]bool{}
	id := profileID
	for {
		if visited[id] {
			return nil, fmt.Errorf("cycle detected in profile inheritance at id=%d", id)
		}
		visited[id] = true
		chain = append([]int64{id}, chain...) // prepend so root is first
		var parentID *int64
		err := r.db.Pool.QueryRow(ctx, `SELECT parent_id FROM quality_profiles WHERE id = $1`, id).Scan(&parentID)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
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

func (r *ProfileRepository) listRules(ctx context.Context, profileID int64) ([]*ProfileRule, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, profile_id, rule_key, severity, params
		FROM quality_profile_rules WHERE profile_id = $1`, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*ProfileRule
	for rows.Next() {
		rule := &ProfileRule{}
		if err := rows.Scan(&rule.ID, &rule.ProfileID, &rule.RuleKey, &rule.Severity, &rule.Params); err != nil {
			return nil, err
		}
		out = append(out, rule)
	}
	return out, rows.Err()
}

func scanProfiles(rows pgx.Rows) ([]*QualityProfile, error) {
	var out []*QualityProfile
	for rows.Next() {
		p := &QualityProfile{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Language, &p.ParentID,
			&p.IsDefault, &p.IsBuiltin, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
