package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/scovl/ollanta/adapter/secondary/postgres"
	"github.com/scovl/ollanta/domain/model"
)

// NewCodePeriodHandler handles new code period API endpoints.
type NewCodePeriodHandler struct {
	periods  *postgres.NewCodePeriodRepository
	projects *postgres.ProjectRepository
}

// NewNewCodePeriodHandler creates a NewCodePeriodHandler.
func NewNewCodePeriodHandler(periods *postgres.NewCodePeriodRepository, projects *postgres.ProjectRepository) *NewCodePeriodHandler {
	return &NewCodePeriodHandler{periods: periods, projects: projects}
}

// GetGlobal handles GET /api/v1/new-code-periods/global
func (h *NewCodePeriodHandler) GetGlobal(w http.ResponseWriter, r *http.Request) {
	ncp, err := h.periods.GetGlobal(r.Context())
	if errors.Is(err, model.ErrNotFound) {
		jsonError(w, http.StatusNotFound, "not found")
		return
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonOK(w, http.StatusOK, ncp)
}

// SetGlobal handles PUT /api/v1/new-code-periods/global
func (h *NewCodePeriodHandler) SetGlobal(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Strategy string `json:"strategy"`
		Value    string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := h.periods.SetGlobal(r.Context(), req.Strategy, req.Value); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetForProject handles GET /api/v1/projects/{key}/new-code-period
func (h *NewCodePeriodHandler) GetForProject(w http.ResponseWriter, r *http.Request) {
	project, err := h.resolveProject(r)
	if err != nil {
		jsonError(w, http.StatusNotFound, "project not found")
		return
	}
	ncp, err := h.periods.GetForProject(r.Context(), project.ID)
	if errors.Is(err, model.ErrNotFound) {
		jsonError(w, http.StatusNotFound, "not found")
		return
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonOK(w, http.StatusOK, ncp)
}

// SetForProject handles PUT /api/v1/projects/{key}/new-code-period
func (h *NewCodePeriodHandler) SetForProject(w http.ResponseWriter, r *http.Request) {
	project, err := h.resolveProject(r)
	if err != nil {
		jsonError(w, http.StatusNotFound, "project not found")
		return
	}
	var req struct {
		Strategy string `json:"strategy"`
		Value    string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := h.periods.SetForProject(r.Context(), project.ID, req.Strategy, req.Value); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteForProject handles DELETE /api/v1/projects/{key}/new-code-period
func (h *NewCodePeriodHandler) DeleteForProject(w http.ResponseWriter, r *http.Request) {
	project, err := h.resolveProject(r)
	if err != nil {
		jsonError(w, http.StatusNotFound, "project not found")
		return
	}
	if err := h.periods.DeleteForProject(r.Context(), project.ID); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *NewCodePeriodHandler) resolveProject(r *http.Request) (*model.Project, error) {
	key := routeParam(r, "key")
	return h.projects.GetByKey(r.Context(), key)
}
