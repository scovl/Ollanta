package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/scovl/ollanta/domain/model"
)

func TestProfileSnapshotRepositoryLegacyUnavailableWithoutDatabase(t *testing.T) {
	snapshot := model.ProfileSnapshot{Language: model.LangGo, RulesHash: "hash", Source: model.ProfileSourceBuiltin, MetadataAvailable: true}
	if !snapshot.MetadataAvailable {
		t.Fatal("profile snapshot should expose metadata availability")
	}
}

func TestProfileSnapshotRepositoryReplaceAndHashChanges(t *testing.T) {
	db, ctx, prefix := openJobRepositoryTestDB(t)
	projectID, previousScanID := createJobTestProjectAndScan(t, db, ctx, prefix+"-profiles")
	currentScanID := createProfileSnapshotTestScan(t, db, ctx, projectID)
	repo := NewProfileSnapshotRepository(db)
	scope := model.AnalysisScope{Type: model.ScopeTypeBranch, Branch: "main"}

	if err := repo.Replace(ctx, projectID, previousScanID, scope, []model.ProfileSnapshot{{Language: model.LangGo, ProfileName: "Ollanta Way", Source: model.ProfileSourceBuiltin, RulesHash: "old", MetadataAvailable: true}}); err != nil {
		t.Fatalf("replace previous: %v", err)
	}
	if err := repo.Replace(ctx, projectID, currentScanID, scope, []model.ProfileSnapshot{{Language: model.LangGo, ProfileName: "Strict Go", Source: model.ProfileSourceLocal, RulesHash: "new", MetadataAvailable: true}}); err != nil {
		t.Fatalf("replace current: %v", err)
	}

	snapshots, available, err := repo.ListByScan(ctx, currentScanID)
	if err != nil {
		t.Fatalf("ListByScan() error = %v", err)
	}
	if !available || len(snapshots) != 1 || snapshots[0].RulesHash != "new" {
		t.Fatalf("snapshots=%+v available=%v, want current profile metadata", snapshots, available)
	}
	changes, err := repo.HashChanges(ctx, currentScanID, previousScanID)
	if err != nil {
		t.Fatalf("HashChanges() error = %v", err)
	}
	if len(changes) != 1 || changes[0].Language != model.LangGo {
		t.Fatalf("changes = %+v, want Go hash change", changes)
	}
}

func TestProfileSnapshotRepositoryLegacyMarkerSuppressesHashChanges(t *testing.T) {
	db, ctx, prefix := openJobRepositoryTestDB(t)
	projectID, previousScanID := createJobTestProjectAndScan(t, db, ctx, prefix+"-legacy")
	currentScanID := createProfileSnapshotTestScan(t, db, ctx, projectID)
	repo := NewProfileSnapshotRepository(db)
	scope := model.AnalysisScope{Type: model.ScopeTypeBranch, Branch: "main"}

	if err := repo.Replace(ctx, projectID, previousScanID, scope, nil); err != nil {
		t.Fatalf("replace legacy marker: %v", err)
	}
	if err := repo.Replace(ctx, projectID, currentScanID, scope, []model.ProfileSnapshot{{Language: model.LangGo, Source: model.ProfileSourceBuiltin, RulesHash: "new", MetadataAvailable: true}}); err != nil {
		t.Fatalf("replace current: %v", err)
	}
	_, available, err := repo.ListByScan(ctx, previousScanID)
	if err != nil {
		t.Fatalf("ListByScan() error = %v", err)
	}
	if available {
		t.Fatal("available = true, want legacy metadata unavailable")
	}
	changes, err := repo.HashChanges(ctx, currentScanID, previousScanID)
	if err != nil {
		t.Fatalf("HashChanges() error = %v", err)
	}
	if len(changes) != 0 {
		t.Fatalf("changes = %+v, want none when previous metadata unavailable", changes)
	}
}

func createProfileSnapshotTestScan(t *testing.T, db *DB, ctx context.Context, projectID int64) int64 {
	t.Helper()
	var scanID int64
	if err := db.Pool.QueryRow(ctx, `
		INSERT INTO scans (project_id, version, status, analysis_date)
		VALUES ($1, 'test', 'completed', $2)
		RETURNING id`, projectID, time.Now().UTC()).Scan(&scanID); err != nil {
		t.Fatalf("insert scan: %v", err)
	}
	return scanID
}
