package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/scovl/ollanta/ollantastore/postgres"
)

// ProfilesHandler handles quality profile API endpoints.
type ProfilesHandler struct {
	profiles *postgres.ProfileRepository
	projects *postgres.ProjectRepository
}

// NewProfilesHandler creates a ProfilesHandler.
func NewProfilesHandler(profiles *postgres.ProfileRepository, projects *postgres.ProjectRepository) *ProfilesHandler {
	return &ProfilesHandler{profiles: profiles, projects: projects}
}

// List handles GET /api/v1/profiles?language=go
func (h *ProfilesHandler) List(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("language")
	list, err := h.profiles.List(r.Context(), lang)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonOK(w, http.StatusOK, list)
}

// Get handles GET /api/v1/profiles/{id}
func (h *ProfilesHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}
	p, err := h.profiles.GetByID(r.Context(), id)
	if errors.Is(err, postgres.ErrNotFound) {
		jsonError(w, http.StatusNotFound, "profile not found")
		return
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonOK(w, http.StatusOK, p)
}

// Create handles POST /api/v1/profiles
func (h *ProfilesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p postgres.QualityProfile
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := h.profiles.Create(r.Context(), &p); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	jsonOK(w, http.StatusCreated, p)
}

// Update handles PUT /api/v1/profiles/{id}
func (h *ProfilesHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var p postgres.QualityProfile
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid json")
		return
	}
	p.ID = id
	if err := h.profiles.Update(r.Context(), &p); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	jsonOK(w, http.StatusOK, p)
}

// Delete handles DELETE /api/v1/profiles/{id}
func (h *ProfilesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.profiles.Delete(r.Context(), id); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ActivateRule handles POST /api/v1/profiles/{id}/rules
func (h *ProfilesHandler) ActivateRule(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		RuleKey  string            `json:"rule_key"`
		Severity string            `json:"severity"`
		Params   map[string]string `json:"params"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := h.profiles.ActivateRule(r.Context(), id, req.RuleKey, req.Severity, req.Params); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeactivateRule handles DELETE /api/v1/profiles/{id}/rules/{rule}
func (h *ProfilesHandler) DeactivateRule(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}
	rule := routeParam(r, "rule")
	if err := h.profiles.DeactivateRule(r.Context(), id, rule); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// EffectiveRules handles GET /api/v1/profiles/{id}/effective-rules
func (h *ProfilesHandler) EffectiveRules(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}
	rules, err := h.profiles.ResolveEffectiveRules(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonOK(w, http.StatusOK, rules)
}

// AssignToProject handles POST /api/v1/projects/{key}/profiles
func (h *ProfilesHandler) AssignToProject(w http.ResponseWriter, r *http.Request) {
	key := routeParam(r, "key")
	project, err := h.projects.GetByKey(r.Context(), key)
	if errors.Is(err, postgres.ErrNotFound) {
		jsonError(w, http.StatusNotFound, "project not found")
		return
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var req struct {
		Language  string `json:"language"`
		ProfileID int64  `json:"profile_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := h.profiles.AssignToProject(r.Context(), project.ID, req.Language, req.ProfileID); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Copy handles POST /api/v1/profiles/{id}/copy
func (h *ProfilesHandler) Copy(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		jsonError(w, http.StatusBadRequest, "name is required")
		return
	}
	profile, err := h.profiles.Copy(r.Context(), id, req.Name)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonOK(w, http.StatusCreated, profile)
}

// SetDefault handles POST /api/v1/profiles/{id}/set-default
func (h *ProfilesHandler) SetDefault(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.profiles.SetDefault(r.Context(), id); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
