import { apiFetch } from '../../core/api.js';
import { state } from '../../core/state.js';
import { escAttr, escHtml, fmtNum } from '../../core/utils.js';
import { getRenderView } from './context.js';
import { customRuleAIModels, renderCustomRuleAIProviderCard, ruleStudioStat } from './shared.js';

export async function loadAIProvidersData() {
  try {
    const aiData = await apiFetch('/custom-rules/ai/models');
    state.customRuleAIProviders = aiData.providers || aiData.items || [];
  } catch {
    state.customRuleAIProviders = [];
  }
  getRenderView()();
}

export function renderAIProvidersTab() {
  const providers = state.customRuleAIProviders;
  if (providers === null) return `<div class="loading-state"><div class="spinner"></div></div>`;
  const items = providers || [];
  const selectedProvider = state.customRuleAISetupProvider || '';
  const connected = items.filter(provider => provider.status === 'connected' || provider.configured).length;
  const local = items.filter(provider => provider.local).length;
  const cloud = items.length - local;

  return `<div class="tab-section ai-providers-page">
    <div class="rule-studio-hero ai-providers-hero">
      <div class="rule-studio-title-block">
        <span class="rule-studio-eyebrow">Integrations</span>
        <div class="rule-studio-title-row">
          <p class="section-title">AI Providers</p>
          <span>${fmtNum(items.length)} providers</span>
        </div>
      </div>
      <div class="rule-studio-summary">
        ${ruleStudioStat('Connected', connected, connected ? 'ok' : 'warn')}
        ${ruleStudioStat('Local', local)}
        ${ruleStudioStat('Cloud', cloud)}
        ${ruleStudioStat('Models', customRuleAIModels(items).length)}
      </div>
      <button class="btn-sm btn-outline" id="refreshAIProvidersBtn" type="button">Refresh</button>
    </div>
    <section class="rule-builder-card ai-providers-workbench">
      <div class="rule-builder-head">
        <span>Provider setup</span>
        <h4>Connect models for Rule Studio</h4>
      </div>
      <div class="custom-rule-ai-provider-list">
        ${items.length ? items.map(provider => renderCustomRuleAIProviderCard(provider, provider.id === selectedProvider)).join('') : '<div class="empty-state compact"><p>No AI providers are configured.</p></div>'}
      </div>
      <div class="custom-rule-ai-actions">
        <button class="btn-sm btn-outline" id="returnRuleStudioBtn" type="button">Back to Rule Studio</button>
      </div>
    </section>
  </div>`;
}

export function bindAIProvidersContent() {
  const renderView = getRenderView();

  document.getElementById('refreshAIProvidersBtn')?.addEventListener('click', () => loadAIProvidersData());
  document.getElementById('returnRuleStudioBtn')?.addEventListener('click', () => {
    document.querySelector('.tab-btn[data-tab="custom-rules"]')?.click();
  });
}
