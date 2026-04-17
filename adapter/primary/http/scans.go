package api

import (
"encoding/json"
"errors"
"net/http"

"github.com/scovl/ollanta/adapter/secondary/postgres"
"github.com/scovl/ollanta/application/ingest"
"github.com/scovl/ollanta/domain/model"
)

// ScansHandler handles scan-related endpoints.
type ScansHandler struct {
scans    *postgres.ScanRepository
projects *postgres.ProjectRepository
pipeline *ingest.IngestUseCase
}

// Ingest handles POST /api/v1/scans
func (h *ScansHandler) Ingest(w http.ResponseWriter, r *http.Request) {
var req ingest.IngestRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
jsonError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
return
}
result, err := h.pipeline.Ingest(r.Context(), &req)
if err != nil {
jsonError(w, http.StatusInternalServerError, err.Error())
return
}
jsonOK(w, http.StatusCreated, result)
}

// Get handles GET /api/v1/scans/{id}.
func (h *ScansHandler) Get(w http.ResponseWriter, r *http.Request) {
id, err := parseID(r, "id")
if err != nil {
jsonError(w, http.StatusBadRequest, "invalid scan id")
return
}
scan, err := h.scans.GetByID(r.Context(), id)
if errors.Is(err, model.ErrNotFound) {
jsonError(w, http.StatusNotFound, "scan not found")
return
}
if err != nil {
jsonError(w, http.StatusInternalServerError, err.Error())
return
}
jsonOK(w, http.StatusOK, scan)
}

// ListByProject handles GET /api/v1/projects/{key}/scans.
func (h *ScansHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
key := routeParam(r, "key")
project, err := h.projects.GetByKey(r.Context(), key)
if errors.Is(err, model.ErrNotFound) {
jsonError(w, http.StatusNotFound, "project not found")
return
}
if err != nil {
jsonError(w, http.StatusInternalServerError, err.Error())
return
}
scans, err := h.scans.ListByProject(r.Context(), project.ID)
if err != nil {
jsonError(w, http.StatusInternalServerError, err.Error())
return
}
jsonOK(w, http.StatusOK, map[string]interface{}{
"items": scans,
"total": len(scans),
})
}

// Latest handles GET /api/v1/projects/{key}/scans/latest.
func (h *ScansHandler) Latest(w http.ResponseWriter, r *http.Request) {
key := routeParam(r, "key")
project, err := h.projects.GetByKey(r.Context(), key)
if errors.Is(err, model.ErrNotFound) {
jsonError(w, http.StatusNotFound, "project not found")
return
}
if err != nil {
jsonError(w, http.StatusInternalServerError, err.Error())
return
}
scan, err := h.scans.GetLatest(r.Context(), project.ID)
if errors.Is(err, model.ErrNotFound) {
jsonError(w, http.StatusNotFound, "no scans for project")
return
}
if err != nil {
jsonError(w, http.StatusInternalServerError, err.Error())
return
}
jsonOK(w, http.StatusOK, scan)
}
