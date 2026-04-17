package api

import (
	"encoding/json"
	"net/http"

	auth "github.com/scovl/ollanta/adapter/secondary/oauth"
	"github.com/scovl/ollanta/adapter/secondary/postgres"
	"github.com/scovl/ollanta/domain/model"
)

// userView is the public representation of a user (no password hash).
type userView struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Provider  string `json:"provider"`
	IsActive  bool   `json:"is_active"`
}

func toUserView(u *model.User) userView {
	return userView{
		ID:        u.ID,
		Login:     u.Login,
		Email:     u.Email,
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
		Provider:  u.Provider,
		IsActive:  u.IsActive,
	}
}

// UsersHandler handles CRUD for users.
type UsersHandler struct {
	users  *postgres.UserRepository
	tokens *postgres.TokenRepository
}

// NewUsersHandler creates a UsersHandler.
func NewUsersHandler(users *postgres.UserRepository, tokens *postgres.TokenRepository) *UsersHandler {
	return &UsersHandler{users: users, tokens: tokens}
}

// Me handles GET /api/v1/users/me.
func (h *UsersHandler) Me(w http.ResponseWriter, r *http.Request) {
	u := UserFromContext(r.Context())
	if u == nil {
		jsonError(w, http.StatusUnauthorized, "not authenticated")
		return
	}
	jsonOK(w, http.StatusOK, toUserView(u))
}

// List handles GET /api/v1/users (requires manage_users).
func (h *UsersHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.users.List(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	views := make([]userView, len(users))
	for i, u := range users {
		views[i] = toUserView(u)
	}
	jsonOK(w, http.StatusOK, map[string]interface{}{
		"users": views,
		"total": len(users),
	})
}

// Get handles GET /api/v1/users/{id} (requires manage_users).
func (h *UsersHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	u, err := h.users.GetByID(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "user not found")
		return
	}
	jsonOK(w, http.StatusOK, toUserView(u))
}

// Create handles POST /api/v1/users (requires manage_users).
func (h *UsersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.Login == "" || req.Email == "" || req.Password == "" {
		jsonError(w, http.StatusBadRequest, "login, email, and password are required")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	u := &model.User{
		Login:        req.Login,
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: hash,
		Provider:     "local",
	}
	if err := h.users.Create(r.Context(), u); err != nil {
		jsonError(w, http.StatusConflict, "login or email already exists")
		return
	}
	jsonOK(w, http.StatusCreated, toUserView(u))
}

// Update handles PUT /api/v1/users/{id} (requires manage_users).
func (h *UsersHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	var req struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	u, err := h.users.GetByID(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "user not found")
		return
	}
	if req.Name != "" {
		u.Name = req.Name
	}
	if req.Email != "" {
		u.Email = req.Email
	}
	if req.AvatarURL != "" {
		u.AvatarURL = req.AvatarURL
	}
	if err := h.users.Update(r.Context(), u); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonOK(w, http.StatusOK, toUserView(u))
}

// Deactivate handles DELETE /api/v1/users/{id} (requires manage_users).
func (h *UsersHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	if err := h.users.Deactivate(r.Context(), id); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListTokens handles GET /api/v1/users/{id}/tokens (requires manage_users).
func (h *UsersHandler) ListTokens(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	tokens, err := h.tokens.ListByUser(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonOK(w, http.StatusOK, map[string]interface{}{"tokens": tokens})
}

// DeleteToken handles DELETE /api/v1/users/{id}/tokens/{tid} (requires manage_users).
func (h *UsersHandler) DeleteToken(w http.ResponseWriter, r *http.Request) {
	_, err := parseID(r, "id")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	tokenID, err := parseID(r, "tid")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid token id")
		return
	}
	if err := h.tokens.Delete(r.Context(), tokenID); err != nil {
		jsonError(w, http.StatusNotFound, "token not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
