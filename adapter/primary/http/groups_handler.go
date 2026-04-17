package api

import (
	"encoding/json"
	"net/http"

	"github.com/scovl/ollanta/adapter/secondary/postgres"
)

// groupView is the public representation of a group.
type groupView struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsBuiltin   bool   `json:"is_builtin"`
}

func toGroupView(g *postgres.Group) groupView {
	return groupView{ID: g.ID, Name: g.Name, Description: g.Description, IsBuiltin: g.IsBuiltin}
}

// GroupsHandler handles CRUD for groups and group membership.
type GroupsHandler struct {
	groups *postgres.GroupRepository
}

// NewGroupsHandler creates a GroupsHandler.
func NewGroupsHandler(groups *postgres.GroupRepository) *GroupsHandler {
	return &GroupsHandler{groups: groups}
}

// List handles GET /api/v1/groups.
func (h *GroupsHandler) List(w http.ResponseWriter, r *http.Request) {
	groups, err := h.groups.List(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	views := make([]groupView, len(groups))
	for i, g := range groups {
		views[i] = toGroupView(g)
	}
	jsonOK(w, http.StatusOK, map[string]interface{}{"groups": views})
}

// Create handles POST /api/v1/groups.
func (h *GroupsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		jsonError(w, http.StatusBadRequest, "name is required")
		return
	}
	g := &postgres.Group{Name: req.Name, Description: req.Description}
	if err := h.groups.Create(r.Context(), g); err != nil {
		jsonError(w, http.StatusConflict, "group name already exists")
		return
	}
	jsonOK(w, http.StatusCreated, toGroupView(g))
}

// Update handles PUT /api/v1/groups/{id}.
func (h *GroupsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid group id")
		return
	}
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	g := &postgres.Group{ID: id, Name: req.Name, Description: req.Description}
	if err := h.groups.Update(r.Context(), g); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonOK(w, http.StatusOK, toGroupView(g))
}

// Delete handles DELETE /api/v1/groups/{id}.
func (h *GroupsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid group id")
		return
	}
	if err := h.groups.Delete(r.Context(), id); err != nil {
		jsonError(w, http.StatusNotFound, "group not found or is built-in")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// AddMember handles POST /api/v1/groups/{id}/members.
func (h *GroupsHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	groupID, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid group id")
		return
	}
	var req struct {
		UserID int64 `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == 0 {
		jsonError(w, http.StatusBadRequest, "user_id required")
		return
	}
	if err := h.groups.AddMember(r.Context(), groupID, req.UserID); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// RemoveMember handles DELETE /api/v1/groups/{id}/members/{uid}.
func (h *GroupsHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	groupID, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid group id")
		return
	}
	userID, err := parseID(r, "uid")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	if err := h.groups.RemoveMember(r.Context(), groupID, userID); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListMembers handles GET /api/v1/groups/{id}/members.
func (h *GroupsHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	groupID, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid group id")
		return
	}
	users, err := h.groups.ListMembers(r.Context(), groupID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	views := make([]userView, len(users))
	for i, u := range users {
		views[i] = toUserView(u)
	}
	jsonOK(w, http.StatusOK, map[string]interface{}{"members": views})
}
