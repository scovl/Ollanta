package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/scovl/ollanta/ollantastore/postgres"
)

// ActivityHandler serves the project activity timeline.
// Inspired by SonarQube's api/project_analyses/search — each scan becomes
// an activity entry decorated with notable events (quality gate changes,
// version bumps, issue spikes).
type ActivityHandler struct {
	scans    *postgres.ScanRepository
	projects *postgres.ProjectRepository
}

type activityEntry struct {
	ScanID       int64           `json:"scan_id"`
	AnalysisDate time.Time       `json:"analysis_date"`
	Version      string          `json:"version,omitempty"`
	Branch       string          `json:"branch,omitempty"`
	GateStatus   string          `json:"gate_status"`
	TotalIssues  int             `json:"total_issues"`
	NewIssues    int             `json:"new_issues"`
	ClosedIssues int             `json:"closed_issues"`
	Events       []activityEvent `json:"events"`
}

type activityEvent struct {
	Category string `json:"category"` // "QUALITY_GATE", "VERSION", "ISSUE_SPIKE", "FIRST_ANALYSIS"
	Name     string `json:"name"`
	Value    string `json:"value,omitempty"`
}

// Activity handles GET /api/v1/projects/{key}/activity?limit=20&offset=0
//
// Returns a chronological timeline of scans with notable events highlighted.
// Events are derived by comparing consecutive scans (quality gate changes,
// version bumps, issue spikes).
func (h *ActivityHandler) Activity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := routeParam(r, "key")

	project, err := h.projects.GetByKey(ctx, key)
	if errors.Is(err, postgres.ErrNotFound) {
		jsonError(w, http.StatusNotFound, "project not found")
		return
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = 20
	}

	scans, total, err := h.scans.ListByProject(ctx, project.ID, limit+1, offset)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Build activity entries; detect events by comparing consecutive scans.
	entries := make([]activityEntry, 0, len(scans))
	for i, s := range scans {
		if i >= limit {
			break // extra scan was only for comparison
		}

		entry := activityEntry{
			ScanID:       s.ID,
			AnalysisDate: s.AnalysisDate,
			Version:      s.Version,
			Branch:       s.Branch,
			GateStatus:   s.GateStatus,
			TotalIssues:  s.TotalIssues,
			NewIssues:    s.NewIssues,
			ClosedIssues: s.ClosedIssues,
		}

		// First analysis event
		if i == len(scans)-1 && offset == 0 {
			entry.Events = append(entry.Events, activityEvent{
				Category: "FIRST_ANALYSIS",
				Name:     "First analysis",
			})
		}

		// Compare with the next (older) scan in the list
		if i+1 < len(scans) {
			prev := scans[i+1]

			// Quality gate change
			if s.GateStatus != prev.GateStatus && s.GateStatus != "" {
				entry.Events = append(entry.Events, activityEvent{
					Category: "QUALITY_GATE",
					Name:     "Quality Gate " + s.GateStatus,
					Value:    prev.GateStatus + " → " + s.GateStatus,
				})
			}

			// Version bump
			if s.Version != prev.Version && s.Version != "" {
				entry.Events = append(entry.Events, activityEvent{
					Category: "VERSION",
					Name:     "Version " + s.Version,
					Value:    s.Version,
				})
			}

			// Issue spike (>50% increase in new issues)
			if prev.TotalIssues > 0 && s.NewIssues > 0 {
				increase := float64(s.NewIssues) / float64(prev.TotalIssues)
				if increase > 0.5 {
					entry.Events = append(entry.Events, activityEvent{
						Category: "ISSUE_SPIKE",
						Name:     "Issue spike detected",
						Value:    strconv.Itoa(s.NewIssues) + " new issues",
					})
				}
			}
		}

		entries = append(entries, entry)
	}

	jsonOK(w, http.StatusOK, map[string]interface{}{
		"items":  entries,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}
