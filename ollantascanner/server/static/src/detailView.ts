import { esc } from "./html";
import type { AIAgent, AIFixPreview, Issue } from "./types";

export type DetailTabKey = "details" | "rule" | "ai-fix";

export interface AIFixViewState {
  loadingAgents: boolean;
  loadingPreview: boolean;
  applying: boolean;
  selectedAgentId: string;
  statusMessage: string;
  errorMessage: string;
  preview: AIFixPreview | null;
}

export function renderDetailTabs(activeTab: DetailTabKey): string {
  const tabs = [
    { key: "details", label: "Details" },
    { key: "rule", label: "Rule" },
    { key: "ai-fix", label: "Fix with AI" },
  ];

  return tabs
    .map(tab => `<button class="detail-tab${activeTab === tab.key ? " active" : ""}" data-detail-tab="${tab.key}">${tab.label}</button>`)
    .join("");
}

export function renderAIFixContent(issue: Issue, state: AIFixViewState, agents: AIAgent[]): string {
  const locationSuffix = issue.end_line && issue.end_line !== issue.line ? `-${issue.end_line}` : "";
  const agentSection = renderAIFixAgentSection(state, agents);
  const previewSection = renderAIFixPreviewSection(state);

  return `
    <div class="detail-section">
      <div class="detail-section-title">Fix with AI</div>
      <div class="detail-msg ai-fix-callout">Ollanta prepara o contexto da issue, envia apenas o trecho relevante para o agente escolhido e mostra um preview antes de qualquer escrita no seu código.</div>
    </div>

    <div class="detail-section">
      <div class="detail-field detail-field-stack">
        <span class="detail-field-label">Target</span>
        <span class="detail-field-value detail-mono-block">${esc(issue.component_path)}:${issue.line}${locationSuffix}</span>
      </div>
      <div class="detail-field detail-field-stack">
        <span class="detail-field-label">Issue</span>
        <span class="detail-field-value">${esc(issue.message)}</span>
      </div>
    </div>

    <div class="detail-section">
      <div class="detail-section-title">Agent</div>
      ${agentSection}
      ${state.statusMessage ? `<div class="ai-fix-status ai-fix-status-ok">${esc(state.statusMessage)}</div>` : ""}
      ${state.errorMessage ? `<div class="ai-fix-status ai-fix-status-error">${esc(state.errorMessage)}</div>` : ""}
    </div>

    <div class="detail-section">
      <div class="detail-section-title">Preview</div>
      ${previewSection}
    </div>
  `;
}

function renderAIFixAgentSection(state: AIFixViewState, agents: AIAgent[]): string {
  if (state.loadingAgents) {
    return `<div class="detail-loading">Loading AI agents…</div>`;
  }
  if (agents.length === 0) {
    return `<div class="detail-empty">No AI agent is configured for the local scanner.</div>`;
  }

  const selectOptions = agents
    .map(agent => `<option value="${esc(agent.id)}"${state.selectedAgentId === agent.id ? " selected" : ""}>${esc(agent.label)} · ${esc(agent.model)}</option>`)
    .join("");
  const generateLabel = state.loadingPreview ? "Generating…" : "Generate fix";
  const generateDisabled = state.loadingPreview ? " disabled" : "";

  return `<div class="ai-fix-controls">
      <select id="ai-agent-select" class="ai-fix-select">${selectOptions}</select>
      <button id="ai-generate-fix" class="ai-fix-button"${generateDisabled}>${generateLabel}</button>
    </div>`;
}

function renderAIFixPreviewSection(state: AIFixViewState): string {
  if (!state.preview) {
    return `<div class="detail-empty">Generate a fix preview to inspect the patch before Ollanta edits your local file.</div>`;
  }

  const previewSummary = state.preview.summary || "Generated fix preview";
  const explanation = state.preview.explanation
    ? `<div class="rule-rationale">${esc(state.preview.explanation)}</div>`
    : "";
  const applyLabel = state.applying ? "Applying…" : "Apply to file";
  const applyDisabled = state.applying ? " disabled" : "";

  return `
    <div class="ai-fix-preview-meta">
      <div><strong>Agent:</strong> ${esc(state.preview.agent.label)}</div>
      <div><strong>Summary:</strong> ${esc(previewSummary)}</div>
    </div>
    ${explanation}
    <pre class="rule-code ai-fix-diff"><code>${esc(state.preview.diff)}</code></pre>
    <div class="ai-fix-actions">
      <button id="ai-apply-fix" class="ai-fix-button ai-fix-button-primary"${applyDisabled}>${applyLabel}</button>
    </div>
  `;
}

