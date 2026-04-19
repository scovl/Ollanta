package api

import (
	"errors"
	"net/http"

	"github.com/scovl/ollanta/ollantastore/postgres"
)

// OverviewHandler serves the project overview (dashboard) endpoint.
// Inspired by SonarQube's api/navigation/component — a single call
// that returns everything the frontend needs to render the project dashboard.
type OverviewHandler struct {
	projects *postgres.ProjectRepository
	scans    *postgres.ScanRepository
	issues   *postgres.IssueRepository
	measures *postgres.MeasureRepository
	gates    *postgres.GateRepository
}

// overviewResponse is the single-call dashboard payload.
type overviewResponse struct {
	Project     *postgres.Project     `json:"project"`
	LastScan    *postgres.Scan        `json:"last_scan,omitempty"`
	QualityGate *overviewGate         `json:"quality_gate,omitempty"`
	Measures    map[string]float64    `json:"measures"`
	Facets      *postgres.IssueFacets `json:"facets,omitempty"`
	NewCode     *overviewNewCode      `json:"new_code,omitempty"`
}

type overviewGate struct {
	Status     string                    `json:"status"`
	Conditions []*postgres.GateCondition `json:"conditions,omitempty"`
}

type overviewNewCode struct {
	NewIssues    int `json:"new_issues"`
	ClosedIssues int `json:"closed_issues"`
}

// Overview handles GET /api/v1/projects/{key}/overview.
//
// Returns the project dashboard in a single response: project metadata,
// latest scan, quality gate status, key measures, issue facets, and
// new code summary. Modelled after SonarQube's unified dashboard call.
func (h *OverviewHandler) Overview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := routeParam(r, "key")

	// ── Project ────────────────────────────────────────────────────────
	project, err := h.projects.GetByKey(ctx, key)
	if errors.Is(err, postgres.ErrNotFound) {
		jsonError(w, http.StatusNotFound, "project not found")
		return
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := overviewResponse{
		Project:  project,
		Measures: make(map[string]float64),
	}

	// ── Latest scan ────────────────────────────────────────────────────
	scan, err := h.scans.GetLatest(ctx, project.ID)
	if err != nil && !errors.Is(err, postgres.ErrNotFound) {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if scan != nil {
		resp.LastScan = scan
		resp.NewCode = &overviewNewCode{
			NewIssues:    scan.NewIssues,
			ClosedIssues: scan.ClosedIssues,
		}

		// ── Facets for latest scan ─────────────────────────────────
		facets, err := h.issues.Facets(ctx, project.ID, scan.ID)
		if err == nil {
			resp.Facets = facets
		}
	}

	// ── Quality gate ───────────────────────────────────────────────────
	gate, conds, err := h.gates.ForProject(ctx, project.ID)
	if err == nil && gate != nil {
		status := "NONE"
		if scan != nil && scan.GateStatus != "" {
			status = scan.GateStatus
		}
		resp.QualityGate = &overviewGate{
			Status:     status,
			Conditions: conds,
		}
	}

	// ── Key measures (latest values) ───────────────────────────────────
	metricKeys := []string{
		"files", "lines", "ncloc", "comments",
		"bugs", "code_smells", "vulnerabilities",
		"coverage", "duplicated_lines_density",
	}
	for _, mk := range metricKeys {
		m, err := h.measures.GetLatest(ctx, project.ID, mk)
		if err == nil && m != nil {
			resp.Measures[mk] = m.Value
		}
	}
	// Also fill from scan totals if measures table is empty
	if scan != nil {
		if _, ok := resp.Measures["files"]; !ok {
			resp.Measures["files"] = float64(scan.TotalFiles)
		}
		if _, ok := resp.Measures["lines"]; !ok {
			resp.Measures["lines"] = float64(scan.TotalLines)
		}
		if _, ok := resp.Measures["ncloc"]; !ok {
			resp.Measures["ncloc"] = float64(scan.TotalNcloc)
		}
		if _, ok := resp.Measures["bugs"]; !ok {
			resp.Measures["bugs"] = float64(scan.TotalBugs)
		}
		if _, ok := resp.Measures["code_smells"]; !ok {
			resp.Measures["code_smells"] = float64(scan.TotalCodeSmells)
		}
		if _, ok := resp.Measures["vulnerabilities"]; !ok {
			resp.Measures["vulnerabilities"] = float64(scan.TotalVulnerabilities)
		}
	}

	jsonOK(w, http.StatusOK, resp)
}
