package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/scovl/ollanta/domain/model"
	"github.com/scovl/ollanta/domain/port"
)

// GateRepository provides CRUD access to quality_gates and related tables.
type GateRepository struct {
	db *DB
}

// NewGateRepository creates a GateRepository backed by db.
func NewGateRepository(db *DB) *GateRepository {
	return &GateRepository{db: db}
}

// compile-time interface check
var _ port.IGateRepo = (*GateRepository)(nil)

// List returns all quality gates.
func (r *GateRepository) List(ctx context.Context) ([]*model.QualityGate, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, name, is_default, is_builtin, small_changeset_lines, created_at, updated_at
		FROM quality_gates ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanGates(rows)
}

// GetByID returns a gate by its ID.
func (r *GateRepository) GetByID(ctx context.Context, id int64) (*model.QualityGate, error) {
	g := &model.QualityGate{}
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, name, is_default, is_builtin, small_changeset_lines, created_at, updated_at
		FROM quality_gates WHERE id = $1`, id,
	).Scan(&g.ID, &g.Name, &g.IsDefault, &g.IsBuiltin, &g.SmallChangesetLines, &g.CreatedAt, &g.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	return g, err
}

// Create inserts a new quality gate.
func (r *GateRepository) Create(ctx context.Context, g *model.QualityGate) error {
	return r.db.Pool.QueryRow(ctx, `
		INSERT INTO quality_gates (name, is_default, small_changeset_lines)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`,
		g.Name, g.IsDefault, g.SmallChangesetLines,
	).Scan(&g.ID, &g.CreatedAt, &g.UpdatedAt)
}

// Update updates a non-builtin gate.
func (r *GateRepository) Update(ctx context.Context, g *model.QualityGate) error {
	if g.IsBuiltin {
		return fmt.Errorf("cannot update builtin gate")
	}
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE quality_gates
		SET name = $1, is_default = $2, small_changeset_lines = $3, updated_at = now()
		WHERE id = $4`,
		g.Name, g.IsDefault, g.SmallChangesetLines, g.ID)
	return err
}

// Delete removes a non-builtin gate.
func (r *GateRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Pool.Exec(ctx,
		`DELETE FROM quality_gates WHERE id = $1 AND is_builtin = FALSE`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("gate not found or is builtin")
	}
	return nil
}

// Conditions returns all conditions for a gate.
func (r *GateRepository) Conditions(ctx context.Context, gateID int64) ([]model.GateCondition, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, gate_id, metric, operator, threshold, on_new_code
		FROM gate_conditions WHERE gate_id = $1 ORDER BY metric`, gateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.GateCondition
	for rows.Next() {
		c := model.GateCondition{}
		if err := rows.Scan(&c.ID, &c.GateID, &c.Metric, &c.Operator, &c.Threshold, &c.OnNewCode); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// AddCondition inserts a new condition to a gate.
func (r *GateRepository) AddCondition(ctx context.Context, c *model.GateCondition) error {
	return r.db.Pool.QueryRow(ctx, `
		INSERT INTO gate_conditions (gate_id, metric, operator, threshold, on_new_code)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		c.GateID, c.Metric, c.Operator, c.Threshold, c.OnNewCode,
	).Scan(&c.ID)
}

// RemoveCondition deletes a condition by ID.
func (r *GateRepository) RemoveCondition(ctx context.Context, conditionID int64) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM gate_conditions WHERE id = $1`, conditionID)
	return err
}

// AssignToProject sets the active gate for a project.
func (r *GateRepository) AssignToProject(ctx context.Context, projectID, gateID int64) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO project_gates (project_id, gate_id)
		VALUES ($1, $2)
		ON CONFLICT (project_id) DO UPDATE SET gate_id = EXCLUDED.gate_id`,
		projectID, gateID)
	return err
}

// ForProject returns the gate assigned to a project, falling back to the default gate.
func (r *GateRepository) ForProject(ctx context.Context, projectID int64) (*model.QualityGate, error) {
	var gateID int64
	err := r.db.Pool.QueryRow(ctx,
		`SELECT gate_id FROM project_gates WHERE project_id = $1`, projectID,
	).Scan(&gateID)
	if errors.Is(err, pgx.ErrNoRows) {
		err = r.db.Pool.QueryRow(ctx,
			`SELECT id FROM quality_gates WHERE is_default = TRUE LIMIT 1`,
		).Scan(&gateID)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, gateID)
}

// ── helpers ──────────────────────────────────────────────────────────────────

func scanGates(rows pgx.Rows) ([]*model.QualityGate, error) {
	var out []*model.QualityGate
	for rows.Next() {
		g := &model.QualityGate{}
		if err := rows.Scan(&g.ID, &g.Name, &g.IsDefault, &g.IsBuiltin,
			&g.SmallChangesetLines, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}
