package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/scovl/ollanta/domain/model"
	"github.com/scovl/ollanta/domain/port"
)

// ProjectRepository provides CRUD access to the projects table.
type ProjectRepository struct {
	db *DB
}

// NewProjectRepository creates a ProjectRepository backed by db.
func NewProjectRepository(db *DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// compile-time interface check
var _ port.IProjectRepo = (*ProjectRepository)(nil)

// Upsert inserts a new project or updates name/description/tags/updated_at on key conflict.
// The ID field is populated on return.
func (r *ProjectRepository) Upsert(ctx context.Context, p *model.Project) error {
	if p.Tags == nil {
		p.Tags = []string{}
	}
	row := r.db.Pool.QueryRow(ctx, `
		INSERT INTO projects (key, name, description, tags, updated_at)
		VALUES ($1, $2, $3, $4, now())
		ON CONFLICT (key) DO UPDATE
		  SET name        = EXCLUDED.name,
		      description = EXCLUDED.description,
		      tags        = EXCLUDED.tags,
		      updated_at  = now()
		RETURNING id, created_at, updated_at`,
		p.Key, p.Name, p.Description, p.Tags,
	)
	return row.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

// Create inserts a new project and populates its ID and timestamps.
func (r *ProjectRepository) Create(ctx context.Context, p *model.Project) error {
	if p.Tags == nil {
		p.Tags = []string{}
	}
	row := r.db.Pool.QueryRow(ctx, `
		INSERT INTO projects (key, name, description, tags)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`,
		p.Key, p.Name, p.Description, p.Tags,
	)
	return row.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

// GetByKey retrieves a project by its unique key. Returns model.ErrNotFound when absent.
func (r *ProjectRepository) GetByKey(ctx context.Context, key string) (*model.Project, error) {
	p := &model.Project{}
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, key, name, description, tags, created_at, updated_at
		FROM projects WHERE key = $1`, key,
	).Scan(&p.ID, &p.Key, &p.Name, &p.Description, &p.Tags, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	return p, err
}

// GetByID retrieves a project by its primary key. Returns model.ErrNotFound when absent.
func (r *ProjectRepository) GetByID(ctx context.Context, id int64) (*model.Project, error) {
	p := &model.Project{}
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, key, name, description, tags, created_at, updated_at
		FROM projects WHERE id = $1`, id,
	).Scan(&p.ID, &p.Key, &p.Name, &p.Description, &p.Tags, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	return p, err
}

// List returns all projects ordered by created_at DESC.
func (r *ProjectRepository) List(ctx context.Context) ([]*model.Project, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, key, name, description, tags, created_at, updated_at
		FROM projects
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*model.Project
	for rows.Next() {
		p := &model.Project{}
		if err := rows.Scan(&p.ID, &p.Key, &p.Name, &p.Description,
			&p.Tags, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

// Delete removes a project and cascades to scans and issues.
func (r *ProjectRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Pool.Exec(ctx, "DELETE FROM projects WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

// scanProject scans a row into a model.Project.
func scanProject(row interface{ Scan(...any) error }) (*model.Project, error) {
	p := &model.Project{}
	err := row.Scan(&p.ID, &p.Key, &p.Name, &p.Description, &p.Tags, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	// Ensure non-nil timestamps for zero-value safety.
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Time{}
	}
	return p, err
}
