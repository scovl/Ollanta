package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/scovl/ollanta/domain/model"
	"github.com/scovl/ollanta/domain/port"
)

// ScanRepository provides access to the scans table.
type ScanRepository struct {
	db *DB
}

// NewScanRepository creates a ScanRepository backed by db.
func NewScanRepository(db *DB) *ScanRepository {
	return &ScanRepository{db: db}
}

// compile-time interface check
var _ port.IScanRepo = (*ScanRepository)(nil)

// Create inserts a new scan and populates its ID and CreatedAt.
func (r *ScanRepository) Create(ctx context.Context, s *model.Scan) error {
	row := r.db.Pool.QueryRow(ctx, `
		INSERT INTO scans (
			project_id, version, branch, commit_sha, status, elapsed_ms,
			gate_status, analysis_date,
			total_files, total_lines, total_ncloc, total_comments,
			total_issues, total_bugs, total_code_smells, total_vulnerabilities,
			new_issues, closed_issues
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18
		) RETURNING id, created_at`,
		s.ProjectID, s.Version, s.Branch, s.CommitSHA, s.Status, s.ElapsedMs,
		s.GateStatus, s.AnalysisDate,
		s.TotalFiles, s.TotalLines, s.TotalNcloc, s.TotalComments,
		s.TotalIssues, s.TotalBugs, s.TotalCodeSmells, s.TotalVulnerabilities,
		s.NewIssues, s.ClosedIssues,
	)
	return row.Scan(&s.ID, &s.CreatedAt)
}

// Update persists gate_status, new_issues, and closed_issues for an existing scan.
func (r *ScanRepository) Update(ctx context.Context, s *model.Scan) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE scans
		SET gate_status = $1, new_issues = $2, closed_issues = $3
		WHERE id = $4`,
		s.GateStatus, s.NewIssues, s.ClosedIssues, s.ID,
	)
	return err
}

// GetByID retrieves a scan by its primary key. Returns model.ErrNotFound when absent.
func (r *ScanRepository) GetByID(ctx context.Context, id int64) (*model.Scan, error) {
	s := &model.Scan{}
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, project_id, version, branch, commit_sha, status, elapsed_ms,
		       gate_status, analysis_date, created_at,
		       total_files, total_lines, total_ncloc, total_comments,
		       total_issues, total_bugs, total_code_smells, total_vulnerabilities,
		       new_issues, closed_issues
		FROM scans WHERE id = $1`, id,
	).Scan(
		&s.ID, &s.ProjectID, &s.Version, &s.Branch, &s.CommitSHA,
		&s.Status, &s.ElapsedMs, &s.GateStatus, &s.AnalysisDate, &s.CreatedAt,
		&s.TotalFiles, &s.TotalLines, &s.TotalNcloc, &s.TotalComments,
		&s.TotalIssues, &s.TotalBugs, &s.TotalCodeSmells, &s.TotalVulnerabilities,
		&s.NewIssues, &s.ClosedIssues,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	return s, err
}

// GetLatest returns the most recent scan for a project. Returns model.ErrNotFound when none.
func (r *ScanRepository) GetLatest(ctx context.Context, projectID int64) (*model.Scan, error) {
	s := &model.Scan{}
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, project_id, version, branch, commit_sha, status, elapsed_ms,
		       gate_status, analysis_date, created_at,
		       total_files, total_lines, total_ncloc, total_comments,
		       total_issues, total_bugs, total_code_smells, total_vulnerabilities,
		       new_issues, closed_issues
		FROM scans
		WHERE project_id = $1
		ORDER BY analysis_date DESC
		LIMIT 1`, projectID,
	).Scan(
		&s.ID, &s.ProjectID, &s.Version, &s.Branch, &s.CommitSHA,
		&s.Status, &s.ElapsedMs, &s.GateStatus, &s.AnalysisDate, &s.CreatedAt,
		&s.TotalFiles, &s.TotalLines, &s.TotalNcloc, &s.TotalComments,
		&s.TotalIssues, &s.TotalBugs, &s.TotalCodeSmells, &s.TotalVulnerabilities,
		&s.NewIssues, &s.ClosedIssues,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	return s, err
}

// ListByProject returns all scans for a project ordered by analysis_date DESC.
func (r *ScanRepository) ListByProject(ctx context.Context, projectID int64) ([]*model.Scan, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, project_id, version, branch, commit_sha, status, elapsed_ms,
		       gate_status, analysis_date, created_at,
		       total_files, total_lines, total_ncloc, total_comments,
		       total_issues, total_bugs, total_code_smells, total_vulnerabilities,
		       new_issues, closed_issues
		FROM scans
		WHERE project_id = $1
		ORDER BY analysis_date DESC`, projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scans []*model.Scan
	for rows.Next() {
		s := &model.Scan{}
		if err := rows.Scan(
			&s.ID, &s.ProjectID, &s.Version, &s.Branch, &s.CommitSHA,
			&s.Status, &s.ElapsedMs, &s.GateStatus, &s.AnalysisDate, &s.CreatedAt,
			&s.TotalFiles, &s.TotalLines, &s.TotalNcloc, &s.TotalComments,
			&s.TotalIssues, &s.TotalBugs, &s.TotalCodeSmells, &s.TotalVulnerabilities,
			&s.NewIssues, &s.ClosedIssues,
		); err != nil {
			return nil, err
		}
		scans = append(scans, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list scans: %w", err)
	}
	return scans, nil
}
