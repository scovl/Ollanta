package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/scovl/ollanta/adapter/secondary/postgres"
	"github.com/scovl/ollanta/domain/model"
)

// ProjectsHandler handles project-related endpoints.
type ProjectsHandler struct {
	repo *postgres.ProjectRepository
}

// Create handles POST /api/v1/projects — upsert a project by key.
func (h *ProjectsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p model.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if p.Key == "" {
		jsonError(w, http.StatusBadRequest, "key is required")
		return
	}
	if err := h.repo.Upsert(r.Context(), &p); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonOK(w, http.StatusCreated, &p)
}

// Get handles GET /api/v1/projects/{key}.
func (h *ProjectsHandler) Get(w http.ResponseWriter, r *http.Request) {
	key := routeParam(r, "key")
	p, err := h.repo.GetByKey(r.Context(), key)
	if errors.Is(err, model.ErrNotFound) {
		jsonError(w, http.StatusNotFound, "project not found")
		return
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonOK(w, http.StatusOK, p)
}

// List handles GET /api/v1/projects?limit=20&offset=0.
func (h *ProjectsHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = 20
	}

	projects, err := h.repo.List(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonOK(w, http.StatusOK, map[string]interface{}{
		"items":  projects,
		"total":  len(projects),
		"limit":  limit,
		"offset": offset,
	})
}
