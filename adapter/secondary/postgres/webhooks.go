package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/scovl/ollanta/domain/model"
	"github.com/scovl/ollanta/domain/port"
)

// WebhookRepository provides CRUD for webhooks and deliveries.
type WebhookRepository struct {
	db *DB
}

// NewWebhookRepository creates a WebhookRepository backed by db.
func NewWebhookRepository(db *DB) *WebhookRepository {
	return &WebhookRepository{db: db}
}

// compile-time interface check
var _ port.IWebhookRepo = (*WebhookRepository)(nil)

// Create inserts a new webhook.
func (r *WebhookRepository) Create(ctx context.Context, wh *model.Webhook) error {
	if wh.Events == nil {
		wh.Events = []string{}
	}
	return r.db.Pool.QueryRow(ctx, `
		INSERT INTO webhooks (project_id, name, url, secret, events, enabled)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`,
		wh.ProjectID, wh.Name, wh.URL, wh.Secret, wh.Events, wh.Enabled,
	).Scan(&wh.ID, &wh.CreatedAt, &wh.UpdatedAt)
}

// GetByID returns a single webhook.
func (r *WebhookRepository) GetByID(ctx context.Context, id int64) (*model.Webhook, error) {
	wh := &model.Webhook{}
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, project_id, name, url, secret, events, enabled, created_at, updated_at
		FROM webhooks WHERE id = $1`, id,
	).Scan(&wh.ID, &wh.ProjectID, &wh.Name, &wh.URL, &wh.Secret, &wh.Events,
		&wh.Enabled, &wh.CreatedAt, &wh.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	return wh, err
}

// ListByProject returns webhooks for a project (plus global webhooks when projectID > 0).
func (r *WebhookRepository) ListByProject(ctx context.Context, projectID int64) ([]*model.Webhook, error) {
	var rows pgx.Rows
	var err error
	if projectID == 0 {
		rows, err = r.db.Pool.Query(ctx, `
			SELECT id, project_id, name, url, secret, events, enabled, created_at, updated_at
			FROM webhooks WHERE project_id IS NULL ORDER BY name`)
	} else {
		rows, err = r.db.Pool.Query(ctx, `
			SELECT id, project_id, name, url, secret, events, enabled, created_at, updated_at
			FROM webhooks WHERE project_id = $1 OR project_id IS NULL ORDER BY name`, projectID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWebhooks(rows)
}

// Update updates a webhook's mutable fields.
func (r *WebhookRepository) Update(ctx context.Context, wh *model.Webhook) error {
	if wh.Events == nil {
		wh.Events = []string{}
	}
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE webhooks
		SET name = $1, url = $2, secret = $3, events = $4, enabled = $5, updated_at = now()
		WHERE id = $6`,
		wh.Name, wh.URL, wh.Secret, wh.Events, wh.Enabled, wh.ID)
	return err
}

// Delete removes a webhook.
func (r *WebhookRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM webhooks WHERE id = $1`, id)
	return err
}

// RecordDelivery records a delivery attempt.
func (r *WebhookRepository) RecordDelivery(ctx context.Context, d *model.WebhookDelivery) error {
	return r.db.Pool.QueryRow(ctx, `
		INSERT INTO webhook_deliveries (webhook_id, event, payload, response_code, response_body, success, attempt)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, delivered_at`,
		d.WebhookID, d.Event, d.Payload, d.ResponseCode, d.ResponseBody, d.Success, d.Attempt,
	).Scan(&d.ID, &d.DeliveredAt)
}

// ListDeliveries returns recent deliveries for a webhook (newest first).
func (r *WebhookRepository) ListDeliveries(ctx context.Context, webhookID int64, limit int) ([]*model.WebhookDelivery, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, webhook_id, event, payload, response_code, response_body, success, attempt, delivered_at
		FROM webhook_deliveries
		WHERE webhook_id = $1
		ORDER BY delivered_at DESC
		LIMIT $2`, webhookID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDeliveries(rows)
}

// ForEvent returns all enabled webhooks that subscribe to the given event.
// Used internally by the webhook dispatcher; not part of the port interface.
func (r *WebhookRepository) ForEvent(ctx context.Context, projectID int64, event string) ([]*model.Webhook, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, project_id, name, url, secret, events, enabled, created_at, updated_at
		FROM webhooks
		WHERE enabled = TRUE
		  AND (project_id IS NULL OR project_id = $1)
		  AND (events = '{}' OR $2 = ANY(events))`,
		projectID, event)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWebhooks(rows)
}

// ── helpers ──────────────────────────────────────────────────────────────────

func scanWebhooks(rows pgx.Rows) ([]*model.Webhook, error) {
	var out []*model.Webhook
	for rows.Next() {
		wh := &model.Webhook{}
		if err := rows.Scan(&wh.ID, &wh.ProjectID, &wh.Name, &wh.URL, &wh.Secret,
			&wh.Events, &wh.Enabled, &wh.CreatedAt, &wh.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, wh)
	}
	return out, rows.Err()
}

func scanDeliveries(rows pgx.Rows) ([]*model.WebhookDelivery, error) {
	var out []*model.WebhookDelivery
	for rows.Next() {
		d := &model.WebhookDelivery{}
		if err := rows.Scan(&d.ID, &d.WebhookID, &d.Event, &d.Payload,
			&d.ResponseCode, &d.ResponseBody, &d.Success, &d.Attempt, &d.DeliveredAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}
