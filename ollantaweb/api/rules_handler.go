package api

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
)

//go:embed rules_data
var rulesFS embed.FS

type ruleDetail struct {
	Key              string   `json:"key"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Language         string   `json:"language"`
	Type             string   `json:"type"`
	Severity         string   `json:"severity"`
	Tags             []string `json:"tags,omitempty"`
	Rationale        string   `json:"rationale,omitempty"`
	NoncompliantCode string   `json:"noncompliant_code,omitempty"`
	CompliantCode    string   `json:"compliant_code,omitempty"`
}

// RulesHandler serves rule metadata for the issue detail panel.
type RulesHandler struct {
	byKey map[string]*ruleDetail
	all   []*ruleDetail
}

// NewRulesHandler creates a RulesHandler by loading embedded rule JSON files.
func NewRulesHandler() *RulesHandler {
	h := &RulesHandler{byKey: make(map[string]*ruleDetail)}
	_ = fs.WalkDir(rulesFS, "rules_data", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}
		data, err := fs.ReadFile(rulesFS, path)
		if err != nil {
			return nil
		}
		var r ruleDetail
		if json.Unmarshal(data, &r) == nil && r.Key != "" {
			h.byKey[r.Key] = &r
			h.all = append(h.all, &r)
		}
		return nil
	})
	return h
}

// Get handles GET /api/v1/rules/* — returns the full metadata for a single rule.
func (h *RulesHandler) Get(w http.ResponseWriter, r *http.Request) {
	raw := strings.TrimPrefix(chi.URLParam(r, "*"), "/")
	key, _ := url.PathUnescape(raw)
	if key == "" {
		jsonError(w, http.StatusBadRequest, "missing rule key")
		return
	}
	rule, ok := h.byKey[key]
	if !ok {
		jsonError(w, http.StatusNotFound, "rule not found")
		return
	}
	jsonOK(w, http.StatusOK, rule)
}

// List handles GET /api/v1/rules — returns metadata for all registered rules.
func (h *RulesHandler) List(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("language")
	if lang == "" {
		jsonOK(w, http.StatusOK, h.all)
		return
	}
	var filtered []*ruleDetail
	for _, rule := range h.all {
		if rule.Language == lang || rule.Language == "*" {
			filtered = append(filtered, rule)
		}
	}
	jsonOK(w, http.StatusOK, filtered)
}
