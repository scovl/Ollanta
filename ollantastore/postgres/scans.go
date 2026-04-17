package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// Scan is the canonical scan record stored in PostgreSQL.
type Scan struct {
	ID                   int64
	ProjectID            int64
	Version              string
	Branch               string
	CommitSHA            string
	Status               string
	ElapsedMs            int64
	GateStatus           string
	AnalysisDate         time.Time
	CreatedAt            time.Time
	TotalFiles           int
	TotalLines           int
	TotalNcloc           int
	TotalComments        int
	TotalIssues          int
	TotalBugs            int
	TotalCodeSmells      int
	TotalVulnerabilities int
	NewIssues            int
	ClosedIssues         int
}

// ScanRepository provides access to the scans table.
type ScanRepository struct {
	db *DB
}

// NewScanRepository creates a ScanRepository backed by db.
func NewScanRepository(db *DB) *ScanRepository {
	return &ScanRepository{db: db}
}

// Create inserts a new scan and populates its ID and CreatedAt.
func (r *ScanRepository) Create(ctx context.Context, s *Scan) error {
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
func (r *ScanRepository) Update(ctx context.Context, s *Scan) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE scans
		SET gate_status = $1, new_issues = $2, closed_issues = $3
		WHERE id = $4`,
		s.GateStatus, s.NewIssues, s.ClosedIssues, s.ID,
	)
	return err
}

// GetByID retrieves a scan by its primary key. Returns ErrNotFound when absent.
func (r *ScanRepository) GetByID(ctx context.Context, id int64) (*Scan, error) {
	s := &Scan{}
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
		return nil, ErrNotFound
	}
	return s, err
}

// GetLatest returns the most recent scan for a project. Returns ErrNotFound when none.
func (r *ScanRepository) GetLatest(ctx context.Context, projectID int64) (*Scan, error) {
	s := &Scan{}
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
		return nil, ErrNotFound
	}
	return s, err
}

// ListByProject returns a page of scans for a project ordered by analysis_date DESC,
// plus the total count.
func (r *ScanRepository) ListByProject(ctx context.Context, projectID int64, limit, offset int) ([]*Scan, int, error) {
	if limit <= 0 {
		limit = 20
	}

	var total int
	if err := r.db.Pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM scans WHERE project_id = $1", projectID,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count scans: %w", err)
	}

	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, project_id, version, branch, commit_sha, status, elapsed_ms,
		       gate_status, analysis_date, created_at,
		       total_files, total_lines, total_ncloc, total_comments,
		       total_issues, total_bugs, total_code_smells, total_vulnerabilities,
		       new_issues, closed_issues
		FROM scans
		WHERE project_id = $1
		ORDER BY analysis_date DESC
		LIMIT $2 OFFSET $3`, projectID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var scans []*Scan
	for rows.Next() {
		s := &Scan{}
		if err := rows.Scan(
			&s.ID, &s.ProjectID, &s.Version, &s.Branch, &s.CommitSHA,
			&s.Status, &s.ElapsedMs, &s.GateStatus, &s.AnalysisDate, &s.CreatedAt,
			&s.TotalFiles, &s.TotalLines, &s.TotalNcloc, &s.TotalComments,
			&s.TotalIssues, &s.TotalBugs, &s.TotalCodeSmells, &s.TotalVulnerabilities,
			&s.NewIssues, &s.ClosedIssues,
		); err != nil {
			return nil, 0, err
		}
		scans = append(scans, s)
	}
	return scans, total, rows.Err()
}
