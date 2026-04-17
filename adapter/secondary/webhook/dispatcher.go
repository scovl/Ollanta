// Package webhook implements async webhook delivery with exponential retry,
// HMAC-SHA256 signing, and dead-letter logging.
package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/scovl/ollanta/adapter/secondary/postgres"
	"github.com/scovl/ollanta/domain/model"
)

// Event names recognised by the dispatcher.
const (
	EventScanCompleted  = "scan.completed"
	EventGateChanged    = "gate.changed"
	EventProjectCreated = "project.created"
	EventProjectDeleted = "project.deleted"
)

// retryDelays defines the exponential back-off schedule (3 attempts).
var retryDelays = []time.Duration{1 * time.Minute, 5 * time.Minute, 30 * time.Minute}

type job struct {
	webhook *model.Webhook
	event   string
	payload []byte
}

// Dispatcher delivers webhooks asynchronously with retry and dead-letter handling.
type Dispatcher struct {
	repo   *postgres.WebhookRepository
	queue  chan job
	client *http.Client
}

// NewDispatcher creates a Dispatcher with a buffered job queue.
func NewDispatcher(repo *postgres.WebhookRepository, bufferSize int) *Dispatcher {
	return &Dispatcher{
		repo:  repo,
		queue: make(chan job, bufferSize),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Start processes delivery jobs until ctx is cancelled.
func (d *Dispatcher) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case j, ok := <-d.queue:
			if !ok {
				return
			}
			go d.deliver(ctx, j)
		}
	}
}

// Stop closes the job queue.
func (d *Dispatcher) Stop() {
	close(d.queue)
}

// Dispatch enqueues a webhook event for all webhooks subscribed to that event.
func (d *Dispatcher) Dispatch(ctx context.Context, projectID int64, event string, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("webhook: marshal payload for event %s: %v", event, err)
		return
	}

	hooks, err := d.repo.ForEvent(ctx, projectID, event)
	if err != nil {
		log.Printf("webhook: query hooks for event %s: %v", event, err)
		return
	}

	for _, wh := range hooks {
		j := job{webhook: wh, event: event, payload: data}
		select {
		case d.queue <- j:
		default:
			log.Printf("webhook: queue full, dropping delivery for webhook %d event %s", wh.ID, event)
		}
	}
}

// DispatchScanEvent implements application/ingest.IWebhookDispatcher.
func (d *Dispatcher) DispatchScanEvent(ctx context.Context, projectID, scanID int64, event string) error {
	payload := map[string]any{
		"project_id": projectID,
		"scan_id":    scanID,
		"event":      event,
	}
	d.Dispatch(ctx, projectID, event, payload)
	return nil
}

func (d *Dispatcher) deliver(ctx context.Context, j job) {
	for attempt, delay := range retryDelays {
		attempt++ // 1-indexed for logging
		code, body, err := d.send(j.webhook, j.event, j.payload)

		del := &model.WebhookDelivery{
			WebhookID: j.webhook.ID,
			Event:     j.event,
			Payload:   string(j.payload),
			Success:   err == nil && code >= 200 && code < 300,
			Attempt:   attempt,
		}
		if code > 0 {
			del.ResponseCode = code
		}
		if body != "" {
			del.ResponseBody = body
		}
		if err := d.repo.RecordDelivery(ctx, del); err != nil {
			log.Printf("webhook: record delivery attempt %d for webhook %d: %v", attempt, j.webhook.ID, err)
		}

		if del.Success {
			return
		}

		if attempt >= len(retryDelays) {
			log.Printf("webhook: dead-letter webhook %d event %s after %d attempts: %v",
				j.webhook.ID, j.event, attempt, err)
			return
		}

		log.Printf("webhook: delivery failed (attempt %d/%d) webhook %d event %s: %v — retry in %s",
			attempt, len(retryDelays), j.webhook.ID, j.event, err, delay)
		select {
		case <-ctx.Done():
			return
		case <-time.After(delay):
		}
	}
}

func (d *Dispatcher) send(wh *model.Webhook, event string, payload []byte) (int, string, error) {
	req, err := http.NewRequest(http.MethodPost, wh.URL, bytes.NewReader(payload))
	if err != nil {
		return 0, "", fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Ollanta-Event", event)
	req.Header.Set("User-Agent", "ollanta-webhook/1.0")

	if wh.Secret != "" {
		sig := sign(payload, wh.Secret)
		req.Header.Set("X-Ollanta-Signature-256", "sha256="+sig)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	var bodyBuf bytes.Buffer
	if resp.ContentLength < 4096 {
		_, _ = bodyBuf.ReadFrom(resp.Body)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.StatusCode, bodyBuf.String(),
			fmt.Errorf("non-2xx response: %d", resp.StatusCode)
	}
	return resp.StatusCode, bodyBuf.String(), nil
}

func sign(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
