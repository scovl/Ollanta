import { describe, expect, it } from "vitest";

import { renderAIFixContent, renderDetailTabs, type AIFixViewState } from "./detailView";
import type { AIAgent, AIFixPreview, Issue } from "./types";

function buildIssue(): Issue {
  return {
    rule_key: "go:nil-check",
    component_path: "/tmp/main.go",
    line: 12,
    column: 1,
    end_line: 12,
    end_column: 18,
    message: "Use the AI fix flow",
    type: "code_smell",
    severity: "major",
    status: "open",
    engine_id: "ollanta",
    line_hash: "hash",
    tags: [],
    secondary_locations: [],
  };
}

function buildState(preview: AIFixPreview | null = null): AIFixViewState {
  return {
    loadingAgents: false,
    loadingPreview: false,
    applying: false,
    selectedAgentId: "mock-agent",
    statusMessage: "",
    errorMessage: "",
    preview,
  };
}

describe("detailView", () => {
  it("renders tabs with Fix with AI visible", () => {
    const html = renderDetailTabs("ai-fix");

    expect(html).toContain("Fix with AI");
    expect(html).toContain('data-detail-tab="ai-fix"');
    expect(html).toContain("detail-tab active");
  });

  it("renders empty-state when no agents are configured", () => {
    const html = renderAIFixContent(buildIssue(), buildState(), []);

    expect(html).toContain("No AI agent is configured for the local scanner.");
    expect(html).toContain("Generate a fix preview");
  });

  it("renders agent selector, preview and apply action", () => {
    const preview: AIFixPreview = {
      preview_id: "preview-1",
      agent: { id: "mock-agent", label: "Mock AI", provider: "mock", model: "deterministic" },
      status: "ready",
      summary: "Generated fix preview",
      explanation: "Preview explanation",
      diff: "@@ lines 12-12 @@\n- old\n+ new",
      file_path: "/tmp/main.go",
      start_line: 12,
      end_line: 12,
      original_snippet: "old",
      replacement: "new",
    };
    const agents: AIAgent[] = [{ id: "mock-agent", label: "Mock AI", provider: "mock", model: "deterministic" }];

    const html = renderAIFixContent(buildIssue(), buildState(preview), agents);

    expect(html).toContain('id="ai-agent-select"');
    expect(html).toContain("Generate fix");
    expect(html).toContain("Preview explanation");
    expect(html).toContain("Apply to file");
    expect(html).toContain("@@ lines 12-12 @@");
  });
});