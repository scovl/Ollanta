package search

import (
	"context"
	"log"

	"github.com/scovl/ollanta/domain/model"
	"github.com/scovl/ollanta/domain/port"
)

// indexJob is an internal task queued by Enqueue.
type indexJob struct {
	scanID     int64
	projectID  int64
	projectKey string
}

// SearchWorker implements application/ingest.ISearchEnqueuer.
// It asynchronously indexes issues after each scan completes.
type SearchWorker struct {
	queue   chan indexJob
	indexer *MeilisearchIndexer
	issues  port.IIssueRepo
}

// NewSearchWorker creates a SearchWorker backed by the indexer and issue repo.
// bufferSize controls the depth of the async queue.
func NewSearchWorker(indexer *MeilisearchIndexer, issues port.IIssueRepo, bufferSize int) *SearchWorker {
	return &SearchWorker{
		queue:   make(chan indexJob, bufferSize),
		indexer: indexer,
		issues:  issues,
	}
}

// Enqueue submits an index job without blocking. Drops the job silently if the queue is full.
func (w *SearchWorker) Enqueue(scanID, projectID int64, projectKey string) {
	select {
	case w.queue <- indexJob{scanID: scanID, projectID: projectID, projectKey: projectKey}:
	default:
		log.Printf("search worker: queue full, dropping index job for scan %d", scanID)
	}
}

// Start processes index jobs until ctx is cancelled. Must be called in a goroutine.
func (w *SearchWorker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case j, ok := <-w.queue:
			if !ok {
				return
			}
			w.process(ctx, j)
		}
	}
}

// Stop closes the queue, causing Start to drain and return.
func (w *SearchWorker) Stop() {
	close(w.queue)
}

func (w *SearchWorker) process(ctx context.Context, j indexJob) {
	pid := j.projectID
	sid := j.scanID
	issues, _, err := w.issues.Query(ctx, model.IssueFilter{
		ProjectID: &pid,
		ScanID:    &sid,
		Limit:     10000,
	})
	if err != nil {
		log.Printf("search worker: query issues for scan %d: %v", j.scanID, err)
		return
	}

	rows := make([]model.IssueRow, len(issues))
	for i, iss := range issues {
		rows[i] = *iss
	}

	if err := w.indexer.IndexIssues(ctx, j.projectKey, rows); err != nil {
		log.Printf("search worker: index issues for scan %d: %v", j.scanID, err)
	}
}
