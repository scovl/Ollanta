package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/scovl/ollanta/adapter/secondary/postgres"
	"github.com/scovl/ollanta/domain/model"
)

// GatesHandler handles quality gate API endpoints.
type GatesHandler struct {
	gates    *postgres.GateRepository
	projects *postgres.ProjectRepository
}

// NewGatesHandler creates a GatesHandler.
func NewGatesHandler(gates *postgres.GateRepository, projects *postgres.ProjectRepository) *GatesHandler {
	return &GatesHandler{gates: gates, projects: projects}
}

// List handles GET /api/v1/quality-gates
func (h *GatesHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.gates.List(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonOK(w, http.StatusOK, list)
}

// Get handles GET /api/v1/quality-gates/{id}
func (h *GatesHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}
	gate, err := h.gates.GetByID(r.Context(), id)
	if errors.Is(err, model.ErrNotFound) {
		jsonError(w, http.StatusNotFound, "gate not found")
		return
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	conditions, err := h.gates.Conditions(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonOK(w, http.StatusOK, map[string]any{"gate": gate, "conditions": conditions})
}

// Create handles POST /api/v1/quality-gates
func (h *GatesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var g model.QualityGate
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if g.SmallChangesetLines == 0 {
		g.SmallChangesetLines = 20
	}
	if err := h.gates.Create(r.Context(), &g); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	jsonOK(w, http.StatusCreated, g)
}

// Update handles PUT /api/v1/quality-gates/{id}
func (h *GatesHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var g model.QualityGate
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid json")
		return
	}
	g.ID = id
	if err := h.gates.Update(r.Context(), &g); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	jsonOK(w, http.StatusOK, g)
}

// Delete handles DELETE /api/v1/quality-gates/{id}
func (h *GatesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.gates.Delete(r.Context(), id); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// AddCondition handles POST /api/v1/quality-gates/{id}/conditions
func (h *GatesHandler) AddCondition(w http.ResponseWriter, r *http.Request) {
	gateID, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var c model.GateCondition
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid json")
		return
	}
	c.GateID = gateID
	if err := h.gates.AddCondition(r.Context(), &c); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonOK(w, http.StatusCreated, c)
}

// RemoveCondition handles DELETE /api/v1/quality-gates/{id}/conditions/{cid}
func (h *GatesHandler) RemoveCondition(w http.ResponseWriter, r *http.Request) {
	cid, err := parseID(r, "cid")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid condition id")
		return
	}
	if err := h.gates.RemoveCondition(r.Context(), cid); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// AssignToProject handles POST /api/v1/projects/{key}/quality-gate
func (h *GatesHandler) AssignToProject(w http.ResponseWriter, r *http.Request) {
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
	var req struct {
		GateID int64 `json:"gate_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := h.gates.AssignToProject(r.Context(), project.ID, req.GateID); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
