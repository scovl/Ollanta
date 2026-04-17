package api

import (
	"encoding/json"
	"net/http"

	"github.com/scovl/ollanta/ollantastore/postgres"
)

// PermsHandler handles grant/revoke of global and project-level permissions.
type PermsHandler struct {
	perms    *postgres.PermissionRepository
	projects *postgres.ProjectRepository
}

// NewPermsHandler creates a PermsHandler.
func NewPermsHandler(perms *postgres.PermissionRepository, projects *postgres.ProjectRepository) *PermsHandler {
	return &PermsHandler{perms: perms, projects: projects}
}

// permRequest is the JSON body for grant/revoke operations.
type permRequest struct {
	Target     string `json:"target"`     // "user" or "group"
	TargetID   int64  `json:"target_id"`
	Permission string `json:"permission"`
}

func decodePermRequest(r *http.Request) (*permRequest, bool) {
	var req permRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, false
	}
	if (req.Target != "user" && req.Target != "group") || req.TargetID == 0 || req.Permission == "" {
		return nil, false
	}
	return &req, true
}

// ListGlobal handles GET /api/v1/permissions/global.
func (h *PermsHandler) ListGlobal(w http.ResponseWriter, r *http.Request) {
	perms, err := h.perms.ListGlobal(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonOK(w, http.StatusOK, map[string]interface{}{"permissions": perms})
}

// GrantGlobal handles POST /api/v1/permissions/global/grant.
func (h *PermsHandler) GrantGlobal(w http.ResponseWriter, r *http.Request) {
	req, ok := decodePermRequest(r)
	if !ok {
		jsonError(w, http.StatusBadRequest, "target, target_id, and permission required")
		return
	}
	if err := h.perms.GrantGlobal(r.Context(), req.Target, req.TargetID, req.Permission); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// RevokeGlobal handles POST /api/v1/permissions/global/revoke.
func (h *PermsHandler) RevokeGlobal(w http.ResponseWriter, r *http.Request) {
	req, ok := decodePermRequest(r)
	if !ok {
		jsonError(w, http.StatusBadRequest, "target, target_id, and permission required")
		return
	}
	if err := h.perms.RevokeGlobal(r.Context(), req.Target, req.TargetID, req.Permission); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListProject handles GET /api/v1/projects/{key}/permissions.
func (h *PermsHandler) ListProject(w http.ResponseWriter, r *http.Request) {
	key := routeParam(r, "key")
	proj, err := h.projects.GetByKey(r.Context(), key)
	if err != nil {
		jsonError(w, http.StatusNotFound, "project not found")
		return
	}
	perms, err := h.perms.ListProject(r.Context(), proj.ID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonOK(w, http.StatusOK, map[string]interface{}{"permissions": perms})
}

// GrantProject handles POST /api/v1/projects/{key}/permissions/grant.
func (h *PermsHandler) GrantProject(w http.ResponseWriter, r *http.Request) {
	key := routeParam(r, "key")
	proj, err := h.projects.GetByKey(r.Context(), key)
	if err != nil {
		jsonError(w, http.StatusNotFound, "project not found")
		return
	}
	req, ok := decodePermRequest(r)
	if !ok {
		jsonError(w, http.StatusBadRequest, "target, target_id, and permission required")
		return
	}
	if err := h.perms.GrantProject(r.Context(), proj.ID, req.Target, req.TargetID, req.Permission); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// RevokeProject handles POST /api/v1/projects/{key}/permissions/revoke.
func (h *PermsHandler) RevokeProject(w http.ResponseWriter, r *http.Request) {
	key := routeParam(r, "key")
	proj, err := h.projects.GetByKey(r.Context(), key)
	if err != nil {
		jsonError(w, http.StatusNotFound, "project not found")
		return
	}
	req, ok := decodePermRequest(r)
	if !ok {
		jsonError(w, http.StatusBadRequest, "target, target_id, and permission required")
		return
	}
	if err := h.perms.RevokeProject(r.Context(), proj.ID, req.Target, req.TargetID, req.Permission); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
