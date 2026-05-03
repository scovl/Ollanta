package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/scovl/ollanta/domain/model"
)

func TestProfilesHandlerImportAcceptsYAML(t *testing.T) {
	repo := &fakeProfileRepo{profile: &model.QualityProfile{ID: 7, Language: model.LangGo, Name: "Strict Go"}}
	handler := NewProfilesHandler(repo, nil)
	body := `version: 1
language: go
rules:
  - key: go:no-large-functions
    severity: critical
    params:
      max_lines: "30"
  - key: go:todo-comment
    active: false
`
	req := profileRequestWithID(http.MethodPost, "/api/v1/profiles/7/import", "7", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/yaml")
	rr := httptest.NewRecorder()

	handler.Import(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s, want 200", rr.Code, rr.Body.String())
	}
	if len(repo.appliedEntries) != 2 {
		t.Fatalf("applied entries = %+v, want 2", repo.appliedEntries)
	}
	if repo.appliedEntries[0].RuleKey != "go:no-large-functions" || repo.appliedEntries[0].Params["max_lines"] != "30" {
		t.Fatalf("first entry = %+v, want imported active rule", repo.appliedEntries[0])
	}
	if repo.appliedEntries[1].Activate {
		t.Fatalf("second entry = %+v, want disabled rule", repo.appliedEntries[1])
	}
}

func TestProfilesHandlerExportIncludesDisabledRuleState(t *testing.T) {
	repo := &fakeProfileRepo{
		profile: &model.QualityProfile{ID: 7, Language: model.LangGo, Name: "Strict Go"},
		effectiveRules: []*model.EffectiveRule{
			{RuleKey: "go:no-large-functions", Severity: string(model.SeverityCritical)},
			{RuleKey: "go:todo-comment", Severity: "OFF", Disabled: true},
		},
	}
	handler := NewProfilesHandler(repo, nil)
	rr := httptest.NewRecorder()

	handler.Export(rr, profileRequestWithID(http.MethodGet, "/api/v1/profiles/7/export", "7", nil))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s, want 200", rr.Code, rr.Body.String())
	}
	var doc profileCodeDocument
	if err := json.NewDecoder(rr.Body).Decode(&doc); err != nil {
		t.Fatalf("decode export: %v", err)
	}
	if len(doc.Rules) != 2 || doc.Rules[1].Active == nil || *doc.Rules[1].Active {
		t.Fatalf("exported rules = %+v, want disabled active=false marker", doc.Rules)
	}
}

func TestProfilesHandlerChangelogReturnsPagination(t *testing.T) {
	repo := &fakeProfileRepo{changelog: []model.ProfileChangelogEntry{{ID: 3, ProfileID: 7, Action: "import"}}}
	handler := NewProfilesHandler(repo, nil)
	rr := httptest.NewRecorder()

	handler.Changelog(rr, profileRequestWithID(http.MethodGet, "/api/v1/profiles/7/changelog?limit=5&offset=1", "7", nil))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s, want 200", rr.Code, rr.Body.String())
	}
	var response struct {
		Items  []model.ProfileChangelogEntry `json:"items"`
		Total  int                           `json:"total"`
		Limit  int                           `json:"limit"`
		Offset int                           `json:"offset"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("decode changelog: %v", err)
	}
	if response.Total != 1 || response.Limit != 5 || response.Offset != 1 || response.Items[0].Action != "import" {
		t.Fatalf("response = %+v, want paginated changelog", response)
	}
}

func profileRequestWithID(method, target, id string, body *strings.Reader) *http.Request {
	var reader interface{ Read([]byte) (int, error) }
	if body != nil {
		reader = body
	}
	req := httptest.NewRequest(method, target, reader)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

type fakeProfileRepo struct {
	profile        *model.QualityProfile
	effectiveRules []*model.EffectiveRule
	appliedEntries []model.ProfileYAMLEntry
	changelog      []model.ProfileChangelogEntry
}

func (r *fakeProfileRepo) List(context.Context, string) ([]*model.QualityProfile, error) {
	return nil, nil
}
func (r *fakeProfileRepo) GetByID(context.Context, int64) (*model.QualityProfile, error) {
	return r.profile, nil
}
func (r *fakeProfileRepo) Create(context.Context, *model.QualityProfile) error { return nil }
func (r *fakeProfileRepo) Update(context.Context, *model.QualityProfile) error { return nil }
func (r *fakeProfileRepo) Delete(context.Context, int64) error                 { return nil }
func (r *fakeProfileRepo) Copy(context.Context, int64, string) (*model.QualityProfile, error) {
	return nil, nil
}
func (r *fakeProfileRepo) SetDefault(context.Context, int64) error { return nil }
func (r *fakeProfileRepo) ActivateRule(context.Context, int64, string, string, map[string]string) error {
	return nil
}
func (r *fakeProfileRepo) DeactivateRule(context.Context, int64, string) error         { return nil }
func (r *fakeProfileRepo) AssignToProject(context.Context, int64, string, int64) error { return nil }
func (r *fakeProfileRepo) ByProjectAndLanguage(context.Context, int64, string) (*model.QualityProfile, error) {
	return nil, nil
}
func (r *fakeProfileRepo) ResolveEffectiveRules(context.Context, int64) ([]*model.EffectiveRule, error) {
	return r.effectiveRules, nil
}
func (r *fakeProfileRepo) ProjectProfiles(context.Context, int64) ([]*model.ProjectQualityProfile, error) {
	return nil, nil
}
func (r *fakeProfileRepo) ProjectEffectiveProfiles(context.Context, int64) ([]*model.EffectiveQualityProfile, error) {
	return nil, nil
}
func (r *fakeProfileRepo) ProfileChangelog(_ context.Context, _ int64, limit, offset int) ([]model.ProfileChangelogEntry, int, error) {
	return r.changelog, len(r.changelog), nil
}
func (r *fakeProfileRepo) ApplyProfileRules(_ context.Context, _ int64, entries []model.ProfileYAMLEntry) error {
	r.appliedEntries = append([]model.ProfileYAMLEntry(nil), entries...)
	return nil
}
func (r *fakeProfileRepo) ApplyProfileYAML(context.Context, int64, string, []model.ProfileYAMLEntry) error {
	return nil
}
