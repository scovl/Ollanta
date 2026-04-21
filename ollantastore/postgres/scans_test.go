package postgres

import (
	"strings"
	"testing"

	domainmodel "github.com/scovl/ollanta/domain/model"
)

const legacyBlankBranchCondition = "OR branch = ''"

func TestChooseDefaultBranchPrefersMain(t *testing.T) {
	t.Parallel()

	got := chooseDefaultBranch([]string{"release", "main", "develop"})
	if got != "main" {
		t.Fatalf("chooseDefaultBranch() = %q, want main", got)
	}
}

func TestChooseDefaultBranchFallsBackToMaster(t *testing.T) {
	t.Parallel()

	got := chooseDefaultBranch([]string{"release", "master"})
	if got != "master" {
		t.Fatalf("chooseDefaultBranch() = %q, want master", got)
	}
}

func TestChooseDefaultBranchFallsBackToMostRecentObservedBranch(t *testing.T) {
	t.Parallel()

	got := chooseDefaultBranch([]string{"release", "hotfix"})
	if got != "release" {
		t.Fatalf("chooseDefaultBranch() = %q, want release", got)
	}
}

func TestBuildScopeConditionIncludesLegacyBlankBranchForDefaultBranch(t *testing.T) {
	t.Parallel()

	condition, args := buildScopeCondition(domainmodel.AnalysisScope{Type: domainmodel.ScopeTypeBranch, Branch: "main"}, "main", 2)
	if !strings.Contains(condition, legacyBlankBranchCondition) {
		t.Fatalf("condition = %q, want legacy blank-branch fallback", condition)
	}
	if len(args) != 2 || args[0] != domainmodel.ScopeTypeBranch || args[1] != "main" {
		t.Fatalf("args = %#v, want [branch main]", args)
	}
}

func TestBuildScopeConditionUsesResolvedDefaultBranchWhenBranchIsOmitted(t *testing.T) {
	t.Parallel()

	condition, args := buildScopeCondition(domainmodel.AnalysisScope{Type: domainmodel.ScopeTypeBranch}, "main", 2)
	if !strings.Contains(condition, legacyBlankBranchCondition) {
		t.Fatalf("condition = %q, want legacy blank-branch fallback", condition)
	}
	if len(args) != 2 || args[1] != "main" {
		t.Fatalf("args = %#v, want resolved default branch main", args)
	}
}

func TestBuildScopeConditionKeepsNonDefaultBranchesIsolated(t *testing.T) {
	t.Parallel()

	condition, args := buildScopeCondition(domainmodel.AnalysisScope{Type: domainmodel.ScopeTypeBranch, Branch: "release"}, "main", 2)
	if strings.Contains(condition, legacyBlankBranchCondition) {
		t.Fatalf("condition = %q, did not expect legacy blank-branch fallback", condition)
	}
	if len(args) != 2 || args[1] != "release" {
		t.Fatalf("args = %#v, want branch release", args)
	}
}

func TestBuildScopeConditionUsesPullRequestKey(t *testing.T) {
	t.Parallel()

	condition, args := buildScopeCondition(domainmodel.AnalysisScope{Type: domainmodel.ScopeTypePullRequest, PullRequestKey: "128"}, "main", 3)
	if condition != "scope_type = $3 AND pull_request_key = $4" {
		t.Fatalf("condition = %q, want pull request condition", condition)
	}
	if len(args) != 2 || args[0] != domainmodel.ScopeTypePullRequest || args[1] != "128" {
		t.Fatalf("args = %#v, want [pull_request 128]", args)
	}
}
