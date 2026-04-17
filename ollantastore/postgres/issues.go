package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// IssueRow is the database representation of a single issue.
type IssueRow struct {
	ID            int64
	ScanID        int64
	ProjectID     int64
	RuleKey       string
	ComponentPath string
	Line          int
	Column        int
	EndLine       int
	EndColumn     int
	Message       string
	Type          string
	Severity      string
	Status        string
	Resolution    string
	EffortMinutes int
	LineHash      string
	Tags          []string
	CreatedAt     time.Time
}

// IssueFilter specifies query parameters for listing issues.
type IssueFilter struct {
	ProjectID *int64
	ScanID    *int64
	RuleKey   *string
	Severity  *string
	Type      *string
	Status    *string
	FilePath  *string // applied as LIKE pattern against component_path
	Limit     int     // default 100, max 1000
	Offset    int
}

// IssueFacets holds aggregate distributions for the issues index.
type IssueFacets struct {
	BySeverity map[string]int
	ByType     map[string]int
	ByRule     map[string]int
}

// IssueRepository provides access to the issues table.
type IssueRepository struct {
	db *DB
}

// NewIssueRepository creates an IssueRepository backed by db.
func NewIssueRepository(db *DB) *IssueRepository {
	return &IssueRepository{db: db}
}

// BulkInsert inserts issues using the PostgreSQL COPY protocol for maximum throughput.
func (r *IssueRepository) BulkInsert(ctx context.Context, issues []IssueRow) error {
	if len(issues) == 0 {
		return nil
	}

	conn, err := r.db.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire conn for bulk insert: %w", err)
	}
	defer conn.Release()

	rows := make([][]interface{}, len(issues))
	for i, iss := range issues {
		tags := iss.Tags
		if tags == nil {
			tags = []string{}
		}
		rows[i] = []interface{}{
			iss.ScanID, iss.ProjectID, iss.RuleKey, iss.ComponentPath,
			iss.Line, iss.Column, iss.EndLine, iss.EndColumn,
			iss.Message, iss.Type, iss.Severity, iss.Status,
			iss.Resolution, iss.EffortMinutes, iss.LineHash, tags,
		}
	}

	_, err = conn.CopyFrom(
		ctx,
		pgx.Identifier{"issues"},
		[]string{
			"scan_id", "project_id", "rule_key", "component_path",
			"line", "column_num", "end_line", "end_column",
			"message", "type", "severity", "status",
			"resolution", "effort_minutes", "line_hash", "tags",
		},
		pgx.CopyFromRows(rows),
	)
	return err
}

// Query returns issues matching the filter, plus the total count before LIMIT/OFFSET.
func (r *IssueRepository) Query(ctx context.Context, f IssueFilter) ([]*IssueRow, int, error) {
	if f.Limit <= 0 {
		f.Limit = 100
	}
	if f.Limit > 1000 {
		f.Limit = 1000
	}

	conds := []string{}
	args := []interface{}{}
	n := 1

	if f.ProjectID != nil {
		conds = append(conds, fmt.Sprintf("project_id = $%d", n))
		args = append(args, *f.ProjectID)
		n++
	}
	if f.ScanID != nil {
		conds = append(conds, fmt.Sprintf("scan_id = $%d", n))
		args = append(args, *f.ScanID)
		n++
	}
	if f.RuleKey != nil {
		conds = append(conds, fmt.Sprintf("rule_key = $%d", n))
		args = append(args, *f.RuleKey)
		n++
	}
	if f.Severity != nil {
		conds = append(conds, fmt.Sprintf("severity = $%d", n))
		args = append(args, *f.Severity)
		n++
	}
	if f.Type != nil {
		conds = append(conds, fmt.Sprintf("type = $%d", n))
		args = append(args, *f.Type)
		n++
	}
	if f.Status != nil {
		conds = append(conds, fmt.Sprintf("status = $%d", n))
		args = append(args, *f.Status)
		n++
	}
	if f.FilePath != nil {
		conds = append(conds, fmt.Sprintf("component_path LIKE $%d", n))
		args = append(args, "%"+*f.FilePath+"%")
		n++
	}

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM issues " + where
	var total int
	if err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count issues: %w", err)
	}

	args = append(args, f.Limit, f.Offset)
	listQuery := fmt.Sprintf(`
		SELECT id, scan_id, project_id, rule_key, component_path,
		       line, column_num, end_line, end_column, message,
		       type, severity, status, resolution, effort_minutes,
		       line_hash, tags, created_at
		FROM issues %s
		ORDER BY created_at DESC, id DESC
		LIMIT $%d OFFSET $%d`, where, n, n+1)

	rows, err := r.db.Pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var issues []*IssueRow
	for rows.Next() {
		iss := &IssueRow{}
		if err := rows.Scan(
			&iss.ID, &iss.ScanID, &iss.ProjectID, &iss.RuleKey, &iss.ComponentPath,
			&iss.Line, &iss.Column, &iss.EndLine, &iss.EndColumn, &iss.Message,
			&iss.Type, &iss.Severity, &iss.Status, &iss.Resolution, &iss.EffortMinutes,
			&iss.LineHash, &iss.Tags, &iss.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		issues = append(issues, iss)
	}
	return issues, total, rows.Err()
}

// Facets returns count distributions by severity, type, and rule for a given scan.
func (r *IssueRepository) Facets(ctx context.Context, projectID, scanID int64) (*IssueFacets, error) {
	facets := &IssueFacets{
		BySeverity: make(map[string]int),
		ByType:     make(map[string]int),
		ByRule:     make(map[string]int),
	}

	type facetRow struct {
		key   string
		count int
	}

	queryFacet := func(column string) ([]facetRow, error) {
		q := fmt.Sprintf(`
			SELECT %s, COUNT(*) FROM issues
			WHERE project_id = $1 AND scan_id = $2
			GROUP BY %s`, column, column)
		rows, err := r.db.Pool.Query(ctx, q, projectID, scanID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var out []facetRow
		for rows.Next() {
			var fr facetRow
			if err := rows.Scan(&fr.key, &fr.count); err != nil {
				return nil, err
			}
			out = append(out, fr)
		}
		return out, rows.Err()
	}

	sev, err := queryFacet("severity")
	if err != nil {
		return nil, fmt.Errorf("facet severity: %w", err)
	}
	for _, r := range sev {
		facets.BySeverity[r.key] = r.count
	}

	typ, err := queryFacet("type")
	if err != nil {
		return nil, fmt.Errorf("facet type: %w", err)
	}
	for _, r := range typ {
		facets.ByType[r.key] = r.count
	}

	rule, err := queryFacet("rule_key")
	if err != nil {
		return nil, fmt.Errorf("facet rule: %w", err)
	}
	for _, r := range rule {
		facets.ByRule[r.key] = r.count
	}

	return facets, nil
}

// CountByProject returns the total number of issues for a project.
func (r *IssueRepository) CountByProject(ctx context.Context, projectID int64) (int, error) {
	var n int
	err := r.db.Pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM issues WHERE project_id = $1", projectID,
	).Scan(&n)
	return n, err
}
