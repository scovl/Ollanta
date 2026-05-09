import { apiFetch } from '../../core/api.js';
import { state } from '../../core/state.js';
import { escAttr, escHtml, fmtNum } from '../../core/utils.js';
import { getRenderView, getShowToast } from './context.js';
import { customRuleAIModelMessage, customRuleAIModelOptions, customRuleAIModels, customRuleAIStateLabel, customRuleDefaultEngine, customRuleEngineGuide, customRuleFullKey, customRuleLifecycleClass, customRuleLifecycleOptions, customRuleMatchesFilters, customRuleStatusClass, renderCustomRuleAIProviderCard, ruleStudioStat } from './shared.js';

export async function loadCustomRulesData() {
  try {
    const [rulesData, enginesData, aiData] = await Promise.all([
      apiFetch('/custom-rules'),
      apiFetch('/rule-engines'),
      apiFetch('/custom-rules/ai/models'),
    ]);
    state.customRulesData = rulesData.items || (Array.isArray(rulesData) ? rulesData : []);
    state.customRuleEngines = Array.isArray(enginesData) ? enginesData : [];
    state.customRuleAIProviders = aiData.providers || aiData.items || [];
  } catch {
    state.customRulesData = [];
    state.customRuleEngines = [];
    state.customRuleAIProviders = [];
  }
  getRenderView()();
}

export function renderCustomRulesTab() {
  const rules = state.customRulesData;
  if (rules === null) return `<div class="loading-state"><div class="spinner"></div></div>`;
  const engines = state.customRuleEngines || [];
  const aiProviders = state.customRuleAIProviders || [];
  const filters = state.customRuleFilters || { search: '', lifecycle: 'all' };
  const allRules = rules || [];
  const visibleRules = allRules.filter(rule => customRuleMatchesFilters(rule, filters));
  const editingId = state.editingCustomRuleId;
  const editingRule = editingId ? allRules.find(rule => rule.id === editingId) : null;
  const builderOpen = editingId != null || state.builderOpen;

  return `<div class="tab-section">
    <p class="section-title" style="margin-top:24px">Rule Studio</p>
    <p style="color:var(--text-muted);font-size:13px;margin-bottom:16px">Create custom rules to enforce project-specific patterns.</p>
    ${renderCustomRuleBuilder(engines, aiProviders, editingRule, builderOpen)}
    <div class="rule-studio-toolbar" style="display:flex;gap:10px;align-items:center;margin-bottom:16px;flex-wrap:wrap">
      <span style="font-weight:600;font-size:13px">${fmtNum(visibleRules.length)} of ${fmtNum(allRules.length)} rules</span>
      <input id="customRuleSearch" class="filter-input" placeholder="Search rules" value="${escAttr(filters.search || '')}" style="width:200px">
      <select id="customRuleLifecycleFilter" class="filter-sel">
        ${customRuleLifecycleOptions(filters.lifecycle || 'all')}
      </select>
      <button class="btn-sm btn-outline" id="refreshCustomRulesBtn" type="button">Refresh</button>
    </div>
    ${visibleRules.length ? visibleRules.map(rule => renderCustomRuleRow(rule)).join('') : renderRuleStudioEmpty(allRules.length)}
  </div>`;
}

function renderRuleStudioEmpty(hasRules) {
  if (hasRules) return '<p style="color:var(--text-muted);padding:20px 0">No rules match the current filters.</p>';
  return `<div class="empty-state" style="padding:40px 0">
    <p>No custom rules yet.</p>
    <button class="btn-sm btn-outline" id="focusRuleBuilderBtn" type="button">Create your first rule</button>
  </div>`;
}

function renderCustomRuleRow(rule) {
  const lifecycle = rule.lifecycle || 'draft';
  const status = rule.validation_status || 'none';
  const expanded = state.expandedCustomRuleId === rule.id;
  const editing = state.editingCustomRuleId === rule.id;
  return `<article class="custom-rule-row${expanded ? ' expanded' : ''}${editing ? ' editing' : ''}" data-rule-id="${rule.id}">
    <div class="custom-rule-row-header">
      <button type="button" class="custom-rule-row-toggle" data-rule-id="${rule.id}" aria-expanded="${expanded ? 'true' : 'false'}">
        <span class="custom-rule-row-chevron" aria-hidden="true">${expanded ? '\u25BE' : '\u25B8'}</span>
        <span class="custom-rule-name">${escHtml(rule.name || rule.key)}</span>
        <span class="mono" style="font-size:11px;color:var(--text-muted)">${escHtml(rule.key)}</span>
      </button>
      <div class="custom-rule-badges">
        <span>${escHtml(rule.language || '-')}</span>
        <span class="badge ${customRuleLifecycleClass(lifecycle)}">${escHtml(lifecycle)}</span>
        <span class="badge ${customRuleStatusClass(status)}">${escHtml(status)}</span>
      </div>
      <div class="custom-rule-row-actions">
        <button class="btn-sm btn-outline edit-custom-rule-btn" data-rule-id="${rule.id}"${editing ? ' disabled' : ''}>Edit</button>
        ${lifecycle === 'draft' && status === 'passed' ? `<button class="btn-sm btn-primary publish-custom-rule-btn" data-rule-id="${rule.id}">Publish</button>` : ''}
        ${expanded ? `<button class="btn-sm btn-ghost disable-custom-rule-btn" data-rule-id="${rule.id}">Disable</button>` : ''}
      </div>
    </div>
    ${expanded ? renderCustomRuleDetails(rule) : ''}
  </article>`;
}

function renderCustomRuleDetails(rule) {
  const engineConfig = rule.engine_config || {};
  const examples = Array.isArray(rule.examples) ? rule.examples : [];
  const compliantExample = examples.find(ex => ex.compliant);
  const noncompliantExample = examples.find(ex => !ex.compliant);
  const tags = Array.isArray(rule.tags) ? rule.tags : [];
  const created = rule.created_at ? new Date(rule.created_at).toLocaleString() : '-';
  const updated = rule.updated_at ? new Date(rule.updated_at).toLocaleString() : '-';
  const published = rule.published_at ? new Date(rule.published_at).toLocaleString() : '';
  const matcherRows = Object.entries(engineConfig)
    .filter(([, value]) => value !== undefined && value !== null && value !== '')
    .map(([key, value]) => `<dt>${escHtml(key)}</dt><dd><code class="mono custom-rule-detail-code">${escHtml(String(value))}</code></dd>`)
    .join('');
  return `<div class="custom-rule-detail">
    <div class="custom-rule-detail-grid">
      <dl class="custom-rule-detail-list">
        <dt>Pack</dt><dd>${escHtml(rule.pack_name || '-')}</dd>
        <dt>Version</dt><dd>${escHtml(String(rule.version || 1))}</dd>
        <dt>Message</dt><dd>${escHtml(rule.message || '-')}</dd>
        <dt>Description</dt><dd>${escHtml(rule.description || '-')}</dd>
        <dt>Tags</dt><dd>${tags.length ? tags.map(tag => `<span class="badge badge-soft">${escHtml(tag)}</span>`).join(' ') : '-'}</dd>
        <dt>Created</dt><dd>${escHtml(created)}</dd>
        <dt>Updated</dt><dd>${escHtml(updated)}</dd>
        ${published ? `<dt>Published</dt><dd>${escHtml(published)}</dd>` : ''}
        <dt>Validation hash</dt><dd><code class="mono custom-rule-detail-code">${escHtml(rule.validation_hash || '(none)')}</code></dd>
      </dl>
      <dl class="custom-rule-detail-list">
        <dt>Engine</dt><dd>${escHtml(rule.engine || '-')}</dd>
        ${matcherRows || '<dt>Matcher</dt><dd>-</dd>'}
      </dl>
    </div>
    ${noncompliantExample ? `<div class="custom-rule-detail-block"><p class="custom-rule-detail-label">Noncompliant example</p><pre class="custom-rule-detail-code-block">${escHtml(noncompliantExample.code || '')}</pre></div>` : ''}
    ${compliantExample ? `<div class="custom-rule-detail-block"><p class="custom-rule-detail-label">Compliant example</p><pre class="custom-rule-detail-code-block">${escHtml(compliantExample.code || '')}</pre></div>` : ''}
  </div>`;
}

function renderCustomRuleBuilder(engines, aiProviders, editingRule, open) {
  const engineList = engines.length ? engines : [
    { engine: 'text', name: 'Text pattern' },
    { engine: 'go-ast', name: 'Go AST pattern' },
    { engine: 'tree-sitter', name: 'Tree-sitter query' },
  ];
  const selectedEngine = editingRule?.engine || customRuleDefaultEngine(engineList);
  const engineOptions = engineList
    .map(item => `<option value="${escAttr(item.engine)}"${item.engine === selectedEngine ? ' selected' : ''}>${escHtml(item.name || item.engine)}</option>`).join('');
  const headerTitle = editingRule ? 'Edit rule' : 'New rule';
  const submitLabel = editingRule ? 'Save changes' : 'Create draft';
  const cancelButton = editingRule ? `<button class="btn-sm btn-outline" id="cancelEditCustomRuleBtn" type="button">Cancel</button>` : '';
  const prefill = editingRule || {};

  return `<details class="rule-builder-details" id="customRuleBuilder"${open ? ' open' : ''}>
    <summary class="rule-builder-summary"><h4>${escHtml(headerTitle)}</h4></summary>
    <div class="rule-builder-body">
      <div class="custom-rule-form-grid">
        <label class="custom-rule-field"><span>Name</span><input id="customRuleName" class="filter-input" placeholder="Rule name" value="${escAttr(prefill.name || '')}"></label>
        <label class="custom-rule-field"><span>Rule key</span><input id="customRuleKey" class="filter-input" placeholder="no-console-log" value="${escAttr(prefill.rule_id || '')}"></label>
        <label class="custom-rule-field"><span>Namespace</span><input id="customRuleNamespace" class="filter-input" value="${escAttr(prefill.namespace || 'custom')}"></label>
        <label class="custom-rule-field"><span>Language</span><select id="customRuleLanguage" class="filter-sel">
          <option value="go"${(prefill.language || 'go') === 'go' ? ' selected' : ''}>Go</option>
          <option value="javascript"${prefill.language === 'javascript' ? ' selected' : ''}>JavaScript</option>
          <option value="typescript"${prefill.language === 'typescript' ? ' selected' : ''}>TypeScript</option>
          <option value="python"${prefill.language === 'python' ? ' selected' : ''}>Python</option>
          <option value="rust"${prefill.language === 'rust' ? ' selected' : ''}>Rust</option>
        </select></label>
        <label class="custom-rule-field"><span>Type</span><select id="customRuleType" class="filter-sel">
          <option value="code_smell"${(prefill.type || 'code_smell') === 'code_smell' ? ' selected' : ''}>Code Smell</option>
          <option value="bug"${prefill.type === 'bug' ? ' selected' : ''}>Bug</option>
          <option value="vulnerability"${prefill.type === 'vulnerability' ? ' selected' : ''}>Vulnerability</option>
        </select></label>
        <label class="custom-rule-field"><span>Severity</span><select id="customRuleSeverity" class="filter-sel">
          <option value="major"${(prefill.severity || 'major') === 'major' ? ' selected' : ''}>Major</option>
          <option value="critical"${prefill.severity === 'critical' ? ' selected' : ''}>Critical</option>
          <option value="minor"${prefill.severity === 'minor' ? ' selected' : ''}>Minor</option>
          <option value="info"${prefill.severity === 'info' ? ' selected' : ''}>Info</option>
        </select></label>
        <label class="custom-rule-field"><span>Engine</span><select id="customRuleEngine" class="filter-sel">${engineOptions}</select></label>
        <label class="custom-rule-field full"><span>Message</span><input id="customRuleMessage" class="filter-input" placeholder="Issue message" value="${escAttr(prefill.message || '')}"></label>
      </div>
      <p class="custom-rule-field-hint" id="customRuleLanguageHint" hidden>Go AST rules are Go-only.</p>
      <div class="rule-section-label">Pattern</div>
      <div class="custom-rule-engine-guide" id="customRuleEngineGuide">Regular expression against source text.</div>
      <div class="custom-rule-form-grid">
        <label class="custom-rule-field full" data-engine-field="text"${selectedEngine === 'text' ? '' : ' hidden'}><input id="customRuleTextPattern" class="filter-input" placeholder="debugger|TODO|panic" value="${escAttr((prefill.engine_config || prefill).text_pattern || (prefill.engine_config || {}).pattern || '')}"></label>
        <label class="custom-rule-field" data-engine-field="go-ast"${selectedEngine === 'go-ast' ? '' : ' hidden'}><select id="customRuleGoASTPattern" class="filter-sel">
          <option value="forbidden_call">Forbidden call</option>
          <option value="forbidden_import">Forbidden import</option>
        </select></label>
        <label class="custom-rule-field full" data-engine-field="go-ast"${selectedEngine === 'go-ast' ? '' : ' hidden'}><input id="customRuleTarget" class="filter-input" placeholder="fmt.Println or net/http"></label>
        <label class="custom-rule-field full" data-engine-field="tree-sitter"${selectedEngine === 'tree-sitter' ? '' : ' hidden'}><textarea id="customRuleQuery" class="filter-input" style="min-height:80px" placeholder="(call_expression function: (identifier) @name)"></textarea></label>
      </div>
      <div class="rule-section-label">Examples</div>
      <div class="custom-rule-form-grid">
        <label class="custom-rule-field full"><span>Noncompliant code</span><textarea id="customRuleExample" class="filter-input" style="min-height:60px" placeholder="code that should produce an issue">${escHtml((prefill.examples || []).find(e => !e?.compliant)?.code || prefill.noncompliant_example || '')}</textarea></label>
        <label class="custom-rule-field full"><span>Compliant code</span><textarea id="customRuleCompliantExample" class="filter-input" style="min-height:60px" placeholder="code that should pass">${escHtml((prefill.examples || []).find(e => e?.compliant)?.code || prefill.compliant_example || '')}</textarea></label>
      </div>
      ${renderCustomRuleAIAssist(aiProviders)}
      <div class="rule-builder-actions">
        <button class="btn btn-primary" id="createCustomRuleBtn" type="button">${escHtml(submitLabel)}</button>
        ${cancelButton}
        ${!editingRule ? `<button class="btn-sm btn-ghost" id="closeRuleBuilderBtn" type="button">Close</button>` : ''}
      </div>
    </div>
  </details>`;
}

function renderCustomRuleAIAssist(providers) {
  const models = customRuleAIModels(providers);
  const hasModels = models.length > 0;
  const modelOptions = customRuleAIModelOptions(models);
  const selected = models.find(model => model.selected) || models[0] || null;
  const canGenerate = selected?.status === 'connected';
  const status = selected ? customRuleAIModelMessage(selected) : 'No AI providers are available yet.';
  const setupLabel = selected ? customRuleAIStateLabel(selected) || 'setup required' : 'setup required';
  return `<div class="rule-builder-section custom-rule-ai-panel">
    <p>AI assist</p>
    <div class="custom-rule-ai-grid">
      <label class="custom-rule-field"><span>Model</span><select id="customRuleAIModel" class="filter-sel"${hasModels ? '' : ' disabled'}>${modelOptions}</select></label>
      <label class="custom-rule-field full"><span>Intent</span><textarea id="customRuleAIIntent" class="filter-input custom-rule-ai-prompt" placeholder="Flag debug logging in production code"></textarea></label>
    </div>
    <div class="custom-rule-ai-actions">
      <button class="btn-sm btn-outline" id="generateCustomRuleAIBtn" type="button"${canGenerate ? '' : ' disabled'}>Generate draft</button>
      <span id="customRuleAIStatus">${escHtml(status)}</span>
      <a class="custom-rule-ai-setup-link" id="connectCustomRuleAIProviderBtn" href="#custom-rule-ai-provider-setup" role="button" data-provider="${escAttr(selected?.provider || '')}"${canGenerate || !selected ? ' hidden' : ''}>${escHtml(setupLabel)}</a>
    </div>
    ${renderCustomRuleAIProviderSetup(providers)}
  </div>`;
}

function renderCustomRuleAIProviderSetup(providers) {
  if (!providers.length) return '';
  return `<div class="custom-rule-ai-provider-setup" id="customRuleAIProviderSetup" hidden>
    <p class="custom-rule-ai-provider-title">AI provider setup</p>
    <div class="custom-rule-ai-provider-list">
      ${providers.map(provider => renderCustomRuleAIProviderCard(provider)).join('')}
    </div>
  </div>`;
}

function bindCustomRuleBuilderControls() {
  const engineSelect = document.getElementById('customRuleEngine');
  if (!engineSelect) return;
  engineSelect.addEventListener('change', syncCustomRuleBuilder);
  document.getElementById('customRuleAIModel')?.addEventListener('change', syncCustomRuleAIControls);
  document.getElementById('connectCustomRuleAIProviderBtn')?.addEventListener('click', openCustomRuleAIProviderSetup);
  document.getElementById('customRuleNamespace')?.addEventListener('input', syncCustomRuleBuilder);
  document.getElementById('customRuleKey')?.addEventListener('input', syncCustomRuleBuilder);
  syncCustomRuleBuilder();
  syncCustomRuleAIControls();
}

function syncCustomRuleBuilder() {
  const engine = document.getElementById('customRuleEngine')?.value || 'text';
  document.querySelectorAll('[data-engine-field]').forEach(field => {
    const engines = (field.dataset.engineField || '').split(/\s+/).filter(Boolean);
    const visible = engines.includes(engine);
    field.hidden = !visible;
    field.querySelectorAll('input, textarea, select').forEach(control => {
      control.disabled = !visible;
    });
  });

  const language = document.getElementById('customRuleLanguage');
  const languageHint = document.getElementById('customRuleLanguageHint');
  if (language) {
    if (engine === 'go-ast') {
      language.value = 'go';
      language.disabled = true;
      if (languageHint) languageHint.hidden = false;
    } else {
      language.disabled = false;
      if (languageHint) languageHint.hidden = true;
    }
  }

  const guide = document.getElementById('customRuleEngineGuide');
  if (guide) guide.textContent = customRuleEngineGuide(engine);

  const preview = document.getElementById('customRuleKeyPreview');
  if (preview) {
    const namespace = document.getElementById('customRuleNamespace')?.value;
    const key = document.getElementById('customRuleKey')?.value;
    preview.textContent = customRuleFullKey(namespace, key);
  }
}

function syncCustomRuleAIControls() {
  const select = document.getElementById('customRuleAIModel');
  const option = select?.selectedOptions?.[0];
  const generate = document.getElementById('generateCustomRuleAIBtn');
  const connect = document.getElementById('connectCustomRuleAIProviderBtn');
  const status = document.getElementById('customRuleAIStatus');
  const modelStatus = option?.dataset.status || '';
  const canGenerate = modelStatus === 'connected';
  if (generate) generate.disabled = !canGenerate;
  if (connect) {
    connect.hidden = canGenerate || !option?.dataset.provider;
    connect.dataset.provider = option?.dataset.provider || '';
    connect.textContent = customRuleAIStateLabel({ status: modelStatus, local: option?.dataset.local === 'true' }) || 'setup required';
    connect.setAttribute?.('aria-label', 'Open AI provider setup for ' + (option?.dataset.providerLabel || option?.dataset.provider || 'selected model'));
  }
  if (status) status.textContent = option?.dataset.message || '';
}

function openCustomRuleAIProviderSetup(event) {
  event?.preventDefault?.();
  const select = document.getElementById('customRuleAIModel');
  const option = select?.selectedOptions?.[0];
  const provider = option?.dataset.provider || event?.currentTarget?.dataset?.provider;
  if (provider) state.customRuleAISetupProvider = provider;
  const aiProvidersTab = document.querySelector('.tab-btn[data-tab="ai-providers"]');
  if (aiProvidersTab) {
    aiProvidersTab.click();
    return;
  }
  const panel = document.getElementById('customRuleAIProviderSetup');
  if (!panel) return;
  panel.hidden = false;
  panel.querySelectorAll('[data-ai-provider-card]').forEach(card => {
    card.classList.toggle('active', card.dataset.aiProviderCard === provider);
  });
  panel.scrollIntoView?.({ behavior: 'smooth', block: 'nearest' });
}

function readCustomRuleBuilderDraft() {
  return {
    pack_name: document.getElementById('customRulePack')?.value.trim() || 'Rule Studio',
    namespace: document.getElementById('customRuleNamespace')?.value.trim() || 'custom',
    rule_id: document.getElementById('customRuleKey')?.value.trim() || '',
    name: document.getElementById('customRuleName')?.value.trim() || '',
    language: document.getElementById('customRuleLanguage')?.value || 'go',
    type: document.getElementById('customRuleType')?.value || 'code_smell',
    severity: document.getElementById('customRuleSeverity')?.value || 'major',
    engine: document.getElementById('customRuleEngine')?.value || 'text',
    text_pattern: document.getElementById('customRuleTextPattern')?.value.trim() || '',
    go_ast_pattern: document.getElementById('customRuleGoASTPattern')?.value || 'forbidden_call',
    target: document.getElementById('customRuleTarget')?.value.trim() || '',
    tree_sitter_query: document.getElementById('customRuleQuery')?.value.trim() || '',
    noncompliant_example: document.getElementById('customRuleExample')?.value || '',
    compliant_example: document.getElementById('customRuleCompliantExample')?.value || '',
    message: document.getElementById('customRuleMessage')?.value.trim() || '',
  };
}

function applyCustomRuleAISuggestion(suggestion) {
  if (!suggestion) return;
  setInputValue('customRulePack', suggestion.pack_name);
  applyCustomRuleID(suggestion.rule_id || suggestion.key);
  setInputValue('customRuleName', suggestion.name);
  setInputValue('customRuleType', suggestion.type);
  setInputValue('customRuleSeverity', suggestion.severity);
  setInputValue('customRuleEngine', suggestion.engine);
  syncCustomRuleBuilder();
  setInputValue('customRuleLanguage', suggestion.language);
  setInputValue('customRuleTextPattern', suggestion.text_pattern || suggestion.pattern);
  setInputValue('customRuleGoASTPattern', suggestion.go_ast_pattern);
  setInputValue('customRuleTarget', suggestion.target);
  setInputValue('customRuleQuery', suggestion.tree_sitter_query || suggestion.query);
  setInputValue('customRuleExample', suggestion.noncompliant_example);
  setInputValue('customRuleCompliantExample', suggestion.compliant_example);
  setInputValue('customRuleMessage', suggestion.message);
  syncCustomRuleBuilder();
}

function applyCustomRuleID(ruleID) {
  const value = (ruleID || '').trim();
  if (!value) return;
  const separator = value.indexOf(':');
  if (separator > 0) {
    setInputValue('customRuleNamespace', value.slice(0, separator));
    setInputValue('customRuleKey', value.slice(separator + 1));
    return;
  }
  setInputValue('customRuleKey', value);
}

function setInputValue(id, value) {
  if (value === undefined || value === null || value === '') return;
  const input = document.getElementById(id);
  if (input) input.value = value;
}

function forceInputValue(id, value) {
  const input = document.getElementById(id);
  if (!input) return;
  input.value = value === undefined || value === null ? '' : value;
}

function prefillCustomRuleBuilder(rule) {
  if (!rule) return;
  const fullKey = rule.key || '';
  const separator = fullKey.indexOf(':');
  const namespace = separator >= 0 ? fullKey.slice(0, separator) : 'custom';
  const ruleID = separator >= 0 ? fullKey.slice(separator + 1) : fullKey;
  setInputValue('customRuleNamespace', namespace);
  setInputValue('customRuleKey', ruleID);
  setInputValue('customRuleName', rule.name || '');
  setInputValue('customRuleType', rule.type || 'code_smell');
  setInputValue('customRuleSeverity', rule.severity || 'major');
  setInputValue('customRuleEngine', rule.engine || 'text');
  setInputValue('customRuleLanguage', rule.language || 'go');
  setInputValue('customRuleMessage', rule.message || '');
  const cfg = rule.engine_config || {};
  setInputValue('customRuleTextPattern', cfg.pattern || '');
  setInputValue('customRuleGoASTPattern', cfg.pattern || 'forbidden_call');
  setInputValue('customRuleTarget', cfg.target || '');
  setInputValue('customRuleQuery', cfg.query || '');
  const examples = Array.isArray(rule.examples) ? rule.examples : [];
  setInputValue('customRuleExample', examples.find(ex => !ex.compliant)?.code || '');
  setInputValue('customRuleCompliantExample', examples.find(ex => ex.compliant)?.code || '');
  syncCustomRuleBuilder();
}

export function bindCustomRulesContent() {
  const renderView = getRenderView();
  const showToast = getShowToast();

  document.getElementById('refreshCustomRulesBtn')?.addEventListener('click', () => loadCustomRulesData());

  document.getElementById('focusRuleBuilderBtn')?.addEventListener('click', () => {
    state.builderOpen = true;
    renderView();
    document.getElementById('customRuleBuilder')?.scrollIntoView({ behavior: 'smooth', block: 'start' });
  });

  document.getElementById('closeRuleBuilderBtn')?.addEventListener('click', () => {
    state.builderOpen = false;
    state.editingCustomRuleId = null;
    renderView();
  });

  const builderEl = document.getElementById('customRuleBuilder');
  if (builderEl) {
    builderEl.addEventListener('toggle', () => {
      state.builderOpen = builderEl.open;
    });
  }

  document.getElementById('customRuleSearch')?.addEventListener('input', event => {
    const filters = { search: '', lifecycle: 'all', ...state.customRuleFilters };
    filters.search = event.target.value;
    state.customRuleFilters = filters;
    renderView();
  });

  document.getElementById('customRuleLifecycleFilter')?.addEventListener('change', event => {
    const filters = { search: '', lifecycle: 'all', ...state.customRuleFilters };
    filters.lifecycle = event.target.value;
    state.customRuleFilters = filters;
    renderView();
  });

  bindCustomRuleBuilderControls();

  document.querySelectorAll('.custom-rule-row-toggle').forEach(toggle => {
    toggle.addEventListener('click', () => {
      const ruleId = Number.parseInt(toggle.dataset.ruleId || '0', 10);
      if (!ruleId) return;
      state.expandedCustomRuleId = state.expandedCustomRuleId === ruleId ? null : ruleId;
      renderView();
    });
  });

  document.querySelectorAll('.edit-custom-rule-btn').forEach(btn => {
    btn.addEventListener('click', event => {
      event.stopPropagation();
      const ruleId = Number.parseInt(btn.dataset.ruleId || '0', 10);
      if (!ruleId) return;
      const rule = (state.customRulesData || []).find(item => item.id === ruleId);
      if (!rule) return;
      state.editingCustomRuleId = ruleId;
      state.expandedCustomRuleId = ruleId;
      state.builderOpen = true;
      renderView();
      requestAnimationFrame(() => {
        prefillCustomRuleBuilder(rule);
        document.getElementById('customRuleBuilder')?.scrollIntoView({ behavior: 'smooth', block: 'start' });
      });
    });
  });

  document.getElementById('cancelEditCustomRuleBtn')?.addEventListener('click', () => {
    state.editingCustomRuleId = null;
    state.builderOpen = false;
    renderView();
  });

  document.getElementById('generateCustomRuleAIBtn')?.addEventListener('click', async event => {
    const button = event.currentTarget;
    const select = document.getElementById('customRuleAIModel');
    const option = select?.selectedOptions?.[0];
    const intent = document.getElementById('customRuleAIIntent')?.value.trim();
    if (!option?.dataset.provider || !option?.dataset.model) {
      showToast('Choose an AI model first.', 'error');
      return;
    }
    if (option.dataset.status !== 'connected') {
      openCustomRuleAIProviderSetup();
      showToast('Set up the selected AI provider before generating a draft.', 'error');
      return;
    }
    if (!intent) {
      showToast('Describe what the rule should detect.', 'error');
      document.getElementById('customRuleAIIntent')?.focus();
      return;
    }
    const status = document.getElementById('customRuleAIStatus');
    button.disabled = true;
    if (status) status.textContent = 'Generating...';
    try {
      const response = await apiFetch('/custom-rules/ai/suggest', {
        method: 'POST',
        body: JSON.stringify({
          provider: option.dataset.provider,
          model: option.dataset.model,
          intent,
          current: readCustomRuleBuilderDraft(),
        }),
      });
      applyCustomRuleAISuggestion(response.suggestion || response);
      if (status) status.textContent = 'Draft generated.';
      showToast('AI draft generated.');
    } catch (err) {
      if (status) status.textContent = err.message;
      showToast(err.message, 'error');
    }
    button.disabled = false;
  });

  document.getElementById('createCustomRuleBtn')?.addEventListener('click', async () => {
    const packName = document.getElementById('customRulePack')?.value.trim() || 'Rule Studio';
    const namespace = document.getElementById('customRuleNamespace')?.value.trim() || 'custom';
    const key = document.getElementById('customRuleKey')?.value.trim();
    const name = document.getElementById('customRuleName')?.value.trim();
    const language = document.getElementById('customRuleLanguage')?.value || 'go';
    const type = document.getElementById('customRuleType')?.value || 'code_smell';
    const severity = document.getElementById('customRuleSeverity')?.value || 'major';
    const engine = document.getElementById('customRuleEngine')?.value || 'text';
    const textPattern = document.getElementById('customRuleTextPattern')?.value.trim();
    const goASTPattern = document.getElementById('customRuleGoASTPattern')?.value || 'forbidden_call';
    const target = document.getElementById('customRuleTarget')?.value.trim();
    const query = document.getElementById('customRuleQuery')?.value.trim();
    const example = document.getElementById('customRuleExample')?.value || '';
    const compliantExample = document.getElementById('customRuleCompliantExample')?.value || '';
    const message = document.getElementById('customRuleMessage')?.value.trim();
    if (!key || !name || !example) {
      showToast('Rule key, name, and noncompliant example are required.', 'error');
      return;
    }
    if (engine === 'text' && !textPattern) {
      showToast('A regexp pattern is required.', 'error');
      return;
    }
    if (engine === 'go-ast' && !target) {
      showToast('A target (e.g. fmt.Println) is required.', 'error');
      return;
    }
    if (engine === 'tree-sitter' && !query) {
      showToast('A tree-sitter query is required.', 'error');
      return;
    }
    const engineConfig = {};
    if (engine === 'tree-sitter') engineConfig.query = query;
    if (engine === 'go-ast') { engineConfig.pattern = goASTPattern; engineConfig.target = target; }
    if (engine === 'text') engineConfig.pattern = textPattern;
    const doc = {
      version: 1,
      pack: { name: packName, namespace },
      rules: [{
        key, name, language, type, severity, engine,
        engine_config: engineConfig,
        message: message || name,
        examples: [
          { name: 'compliant', code: compliantExample, compliant: true },
          { name: 'noncompliant', code: example, compliant: false },
        ],
      }],
    };
    try {
      const editingId = state.editingCustomRuleId;
      if (editingId) {
        const draftRule = doc.rules[0];
        const fullKey = draftRule.key.includes(':') ? draftRule.key : `${namespace}:${draftRule.key}`;
        await apiFetch('/custom-rules/' + encodeURIComponent(editingId), {
          method: 'PUT',
          body: JSON.stringify({ ...draftRule, key: fullKey, pack_name: packName }),
        });
        showToast('Rule updated.');
        state.editingCustomRuleId = null;
        state.builderOpen = false;
      } else {
        await apiFetch('/custom-rules', { method: 'POST', body: JSON.stringify(doc) });
        showToast('Draft created.');
        state.builderOpen = false;
      }
      await loadCustomRulesData();
    } catch (err) {
      showToast(err.message, 'error');
    }
  });

  document.getElementById('importCustomRulesBtn')?.addEventListener('click', async () => {
    const text = document.getElementById('customRuleImportText')?.value.trim();
    if (!text) {
      showToast('Paste a rule pack first.', 'error');
      return;
    }
    try {
      await apiFetch('/custom-rules/import', { method: 'POST', body: text });
      showToast('Custom rule pack imported.');
      await loadCustomRulesData();
    } catch (err) {
      showToast(err.message, 'error');
    }
  });

  document.querySelectorAll('.validate-custom-rule-btn').forEach(btn => {
    btn.addEventListener('click', async () => {
      btn.disabled = true;
      try {
        await apiFetch('/custom-rules/' + encodeURIComponent(btn.dataset.ruleId) + '/validate', { method: 'POST' });
        showToast('Custom rule validated.');
        await loadCustomRulesData();
      } catch (err) {
        showToast(err.message, 'error');
      }
      btn.disabled = false;
    });
  });

  document.querySelectorAll('.publish-custom-rule-btn').forEach(btn => {
    btn.addEventListener('click', async () => {
      btn.disabled = true;
      try {
        await apiFetch('/custom-rules/' + encodeURIComponent(btn.dataset.ruleId) + '/publish', { method: 'POST' });
        showToast('Custom rule published.');
        state.profilesData = null;
        await loadCustomRulesData();
      } catch (err) {
        showToast(err.message, 'error');
      }
      btn.disabled = false;
    });
  });

  document.querySelectorAll('.disable-custom-rule-btn').forEach(btn => {
    btn.addEventListener('click', async () => {
      if (!confirm('Disable this custom rule?')) return;
      btn.disabled = true;
      try {
        await apiFetch('/custom-rules/' + encodeURIComponent(btn.dataset.ruleId) + '/disable', { method: 'POST' });
        showToast('Custom rule disabled.');
        state.profilesData = null;
        await loadCustomRulesData();
      } catch (err) {
        showToast(err.message, 'error');
      }
      btn.disabled = false;
    });
  });

  document.querySelectorAll('.export-custom-rule-btn').forEach(btn => {
    btn.addEventListener('click', async () => {
      btn.disabled = true;
      try {
        const doc = await apiFetch('/custom-rules/' + encodeURIComponent(btn.dataset.ruleId) + '/export');
        await navigator.clipboard?.writeText(JSON.stringify(doc, null, 2));
        showToast('Custom rule exported.');
      } catch (err) {
        showToast(err.message, 'error');
      }
      btn.disabled = false;
    });
  });

  document.querySelectorAll('.add-custom-rule-profile-btn').forEach(btn => {
    btn.addEventListener('click', async () => {
      const select = document.getElementById('custom-rule-profile-' + btn.dataset.ruleId);
      const profileID = Number.parseInt(select?.value || '0', 10);
      if (!profileID) return;
      btn.disabled = true;
      try {
        await apiFetch('/profiles/' + encodeURIComponent(profileID) + '/rules', {
          method: 'POST',
          body: JSON.stringify({ rule_key: btn.dataset.ruleKey, params: {} }),
        });
        showToast('Custom rule added to profile.');
        state.profilesData = null;
        await loadCustomRulesData();
      } catch (err) {
        showToast(err.message, 'error');
      }
      btn.disabled = false;
    });
  });
}
