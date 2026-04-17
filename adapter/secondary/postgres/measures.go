package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/scovl/ollanta/domain/model"
	"github.com/scovl/ollanta/domain/port"
)

// MeasureRepository provides access to the measures table.
type MeasureRepository struct {
	db *DB
}

// NewMeasureRepository creates a MeasureRepository backed by db.
func NewMeasureRepository(db *DB) *MeasureRepository {
	return &MeasureRepository{db: db}
}

// compile-time interface check
var _ port.IMeasureRepo = (*MeasureRepository)(nil)

// BulkInsert inserts all rows using individual INSERTs within a single transaction.
func (r *MeasureRepository) BulkInsert(ctx context.Context, measures []model.MeasureRow) error {
	if len(measures) == 0 {
		return nil
	}
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	for _, m := range measures {
		if _, err := tx.Exec(ctx, `
			INSERT INTO measures (scan_id, project_id, metric_key, component_path, value)
			VALUES ($1, $2, $3, $4, $5)`,
			m.ScanID, m.ProjectID, m.MetricKey, m.ComponentPath, m.Value,
		); err != nil {
			return fmt.Errorf("insert measure %s: %w", m.MetricKey, err)
		}
	}
	return tx.Commit(ctx)
}

// GetLatest returns the most recent project-level value for a metric key.
func (r *MeasureRepository) GetLatest(ctx context.Context, projectID int64, metricKey string) (*model.MeasureRow, error) {
	m := &model.MeasureRow{}
	err := r.db.Pool.QueryRow(ctx, `
		SELECT m.id, m.scan_id, m.project_id, m.metric_key, m.component_path, m.value, m.created_at
		FROM measures m
		JOIN scans s ON s.id = m.scan_id
		WHERE m.project_id = $1
		  AND m.metric_key = $2
		  AND m.component_path = ''
		ORDER BY s.analysis_date DESC
		LIMIT 1`, projectID, metricKey,
	).Scan(&m.ID, &m.ScanID, &m.ProjectID, &m.MetricKey, &m.ComponentPath, &m.Value, &m.CreatedAt)
	return m, err
}

// Trend returns a time-ordered series of project-level metric values between from and to.
func (r *MeasureRepository) Trend(ctx context.Context, projectID int64, metricKey string, from, to time.Time) ([]model.TrendPoint, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT s.analysis_date, m.value
		FROM measures m
		JOIN scans s ON s.id = m.scan_id
		WHERE m.project_id = $1
		  AND m.metric_key = $2
		  AND m.component_path = ''
		  AND s.analysis_date BETWEEN $3 AND $4
		ORDER BY s.analysis_date ASC`, projectID, metricKey, from, to,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []model.TrendPoint
	for rows.Next() {
		var pt model.TrendPoint
		if err := rows.Scan(&pt.AnalysisDate, &pt.Value); err != nil {
			return nil, err
		}
		points = append(points, pt)
	}
	return points, rows.Err()
}
