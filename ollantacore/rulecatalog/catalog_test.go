package rulecatalog

import (
	"testing"
)

func TestCweReferenceFromTags(t *testing.T) {
	tests := []struct {
		tags []string
		want string
	}{
		{[]string{"security", "xss", "cwe-1004"}, "https://cwe.mitre.org/data/definitions/1004.html"},
		{[]string{"security", "cwe-79"}, "https://cwe.mitre.org/data/definitions/79.html"},
		{[]string{"security", "crypto", "cwe-327"}, "https://cwe.mitre.org/data/definitions/327.html"},
		{[]string{"complexity", "readability"}, ""},
		{[]string{}, ""},
		{[]string{"cwe-"}, ""},
	}
	for _, tt := range tests {
		got := cweReferenceFromTags(tt.tags)
		if got != tt.want {
			t.Errorf("cweReferenceFromTags(%v) = %q, want %q", tt.tags, got, tt.want)
		}
	}
}

func TestRulesHaveReferenceURLWhenCWETagPresent(t *testing.T) {
	for _, r := range Rules() {
		if r.ReferenceURL != "" {
			continue
		}
		for _, tag := range r.Tags {
			if len(tag) > 4 && tag[:4] == "cwe-" {
				t.Errorf("rule %q has cwe tag %q but no ReferenceURL", r.Key, tag)
			}
		}
	}
}

func TestCatalogDefensiveCopies(t *testing.T) {
	rule, ok := ByKey("go:bad-tmp")
	if !ok {
		t.Fatal("go:bad-tmp not found")
	}
	original := rule.Tags
	rule.Tags = append(rule.Tags, "mutated")

	again, ok := ByKey("go:bad-tmp")
	if !ok {
		t.Fatal("go:bad-tmp not found on second read")
	}
	if len(again.Tags) != len(original) {
		t.Fatalf("catalog mutation leaked, tags = %v", again.Tags)
	}
}
