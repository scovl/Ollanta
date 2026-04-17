package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/scovl/ollanta/ollantastore/postgres"
	"github.com/scovl/ollanta/ollantastore/search"
	"github.com/scovl/ollanta/ollantaweb/config"
	"github.com/scovl/ollanta/ollantaweb/ingest"
)

// NewRouter builds and returns the complete chi router for the ollantaweb server.
func NewRouter(
	cfg *config.Config,
	projects *postgres.ProjectRepository,
	scans *postgres.ScanRepository,
	issues *postgres.IssueRepository,
	measures *postgres.MeasureRepository,
	users *postgres.UserRepository,
	groups *postgres.GroupRepository,
	tokens *postgres.TokenRepository,
	sessions *postgres.SessionRepository,
	perms *postgres.PermissionRepository,
	searcher *search.MeilisearchSearcher,
	indexer *search.MeilisearchIndexer,
	pipeline *ingest.Pipeline,
) http.Handler {
	r := chi.NewRouter()

	// ── Global middleware ──────────────────────────────────────────────────
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(RequestLogger)
	r.Use(CORS)
	r.Use(MaxBody(10 << 20)) // 10 MB

	// ── Health (always public) ─────────────────────────────────────────────
	r.Get("/healthz", Liveness)
	r.Get("/readyz", Readiness)

	// ── Auth middleware ────────────────────────────────────────────────────
	authMW := NewAuthMiddleware(users, tokens, sessions, []byte(cfg.JWTSecret))

	// ── Handlers ──────────────────────────────────────────────────────────
	authH := NewAuthHandler(cfg, users, groups, sessions)
	usersH := NewUsersHandler(users, tokens)
	groupsH := NewGroupsHandler(groups)
	tokensH := NewTokensHandler(tokens, projects, perms)
	permsH := NewPermsHandler(perms, projects)

	ph := &ProjectsHandler{repo: projects}
	sh := &ScansHandler{scans: scans, projects: projects, pipeline: pipeline}
	ih := &IssuesHandler{issues: issues, projects: projects}
	mh := &MeasuresHandler{measures: measures, projects: projects}
	srh := &SearchHandler{searcher: searcher}

	// ── API v1 ────────────────────────────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {
		// Public: auth endpoints
		r.Post("/auth/login", authH.Login)
		r.Post("/auth/refresh", authH.Refresh)
		r.Get("/auth/github", authH.GitHubRedirect)
		r.Get("/auth/github/callback", authH.GitHubCallback)
		r.Get("/auth/gitlab", authH.GitLabRedirect)
		r.Get("/auth/gitlab/callback", authH.GitLabCallback)
		r.Get("/auth/google", authH.GoogleRedirect)
		r.Get("/auth/google/callback", authH.GoogleCallback)

		// Protected: all other routes require a valid JWT or API token
		r.Group(func(r chi.Router) {
			r.Use(authMW.Authenticate)

			// Auth management
			r.Post("/auth/logout", authH.Logout)

			// Current user
			r.Get("/users/me", usersH.Me)

			// Own tokens (self-service)
			r.Get("/users/me/tokens", tokensH.List)
			r.Post("/users/me/tokens", tokensH.Create)
			r.Delete("/users/me/tokens/{id}", tokensH.Delete)

			// User management (requires manage_users)
			r.Route("/users", func(r chi.Router) {
				r.Use(RequirePermission(perms, "manage_users"))
				r.Get("/", usersH.List)
				r.Post("/", usersH.Create)
				r.Get("/{id}", usersH.Get)
				r.Put("/{id}", usersH.Update)
				r.Delete("/{id}", usersH.Deactivate)
				r.Get("/{id}/tokens", usersH.ListTokens)
				r.Delete("/{id}/tokens/{tid}", usersH.DeleteToken)
			})

			// Group management (requires manage_groups)
			r.Route("/groups", func(r chi.Router) {
				r.Use(RequirePermission(perms, "manage_groups"))
				r.Get("/", groupsH.List)
				r.Post("/", groupsH.Create)
				r.Put("/{id}", groupsH.Update)
				r.Delete("/{id}", groupsH.Delete)
				r.Get("/{id}/members", groupsH.ListMembers)
				r.Post("/{id}/members", groupsH.AddMember)
				r.Delete("/{id}/members/{uid}", groupsH.RemoveMember)
			})

			// Global permission management (requires admin)
			r.Route("/permissions", func(r chi.Router) {
				r.Use(RequirePermission(perms, "admin"))
				r.Get("/global", permsH.ListGlobal)
				r.Post("/global/grant", permsH.GrantGlobal)
				r.Post("/global/revoke", permsH.RevokeGlobal)
			})

			// Projects
			r.Post("/projects", ph.Create)
			r.Get("/projects", ph.List)
			r.Get("/projects/{key}", ph.Get)

			// Project-level permissions (requires admin global or project admin)
			r.Get("/projects/{key}/permissions", permsH.ListProject)
			r.Post("/projects/{key}/permissions/grant", permsH.GrantProject)
			r.Post("/projects/{key}/permissions/revoke", permsH.RevokeProject)

			// Scans
			r.Post("/scans", sh.Ingest)
			r.Get("/scans/{id}", sh.Get)
			r.Get("/projects/{key}/scans", sh.ListByProject)
			r.Get("/projects/{key}/scans/latest", sh.Latest)

			// Issues
			r.Get("/issues", ih.List)
			r.Get("/issues/facets", ih.Facets)

			// Measures
			r.Get("/projects/{key}/measures/trend", mh.Trend)

			// Search
			r.Get("/search", srh.Search)
		})
	})

	// ── Admin (requires admin permission) ─────────────────────────────────
	r.Group(func(r chi.Router) {
		r.Use(authMW.Authenticate)
		r.Use(RequirePermission(perms, "admin"))
		r.Post("/admin/reindex", func(w http.ResponseWriter, r *http.Request) {
			go func() {
				ctx := r.Context()
				if err := indexer.ReindexAll(ctx, issues, projects); err != nil {
					_ = err
				}
			}()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"status":"reindex started"}`))
		})
	})

	// ── Frontend (SPA fallback) ────────────────────────────────────────────
	r.Handle("/*", staticHandler())

	return r
}
