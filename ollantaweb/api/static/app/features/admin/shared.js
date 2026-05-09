import { escAttr, escHtml, fmtNum } from '../../core/utils.js';

export function operatorLabel(op) {
  const labels = { GT: 'is greater than', LT: 'is less than', GTE: 'is greater than or equal', LTE: 'is less than or equal', EQ: 'equals', NE: 'is not equal' };
  return labels[op] || op;
}

export function gateConditionLabel(metric) {
  const labels = { bugs: 'Bugs', vulnerabilities: 'Vulnerabilities', code_smells: 'Code Smells', coverage: 'Coverage', new_bugs: 'New Bugs', new_vulnerabilities: 'New Vulnerabilities', new_code_smells: 'New Code Smells', new_coverage: 'Coverage on New Code', duplicated_lines_density: 'Duplicated Lines (%)', new_duplicated_lines_density: 'Duplicated Lines on New Code (%)', new_security_hotspots_reviewed: 'Security Hotspots Reviewed', new_maintainability_rating: 'Maintainability Rating on New Code', new_reliability_rating: 'Reliability Rating on New Code', new_security_rating: 'Security Rating on New Code', security_hotspots_reviewed: 'Security Hotspots Reviewed', security_review_rating: 'Security Review Rating', reliability_remediation_effort: 'Reliability Remediation Effort', security_remediation_effort: 'Security Remediation Effort' };
  return labels[metric] || metric;
}

export function conditionMetricSuffix(value, metric) {
  const pctMetrics = ['coverage', 'new_coverage', 'duplicated_lines_density', 'new_duplicated_lines_density'];
  if (pctMetrics.includes(metric)) return value + '%';
  const ratingMetrics = ['new_maintainability_rating', 'new_reliability_rating', 'new_security_rating', 'security_review_rating'];
  if (ratingMetrics.includes(metric)) return String(value);
  return String(value);
}

export function gateMetricOptions() {
  return [
    ['bugs', 'Bugs'],
    ['vulnerabilities', 'Vulnerabilities'],
    ['code_smells', 'Code Smells'],
    ['coverage', 'Coverage'],
    ['new_bugs', 'New Bugs'],
    ['new_vulnerabilities', 'New Vulnerabilities'],
    ['new_code_smells', 'New Code Smells'],
    ['new_coverage', 'Coverage on New Code'],
    ['duplicated_lines_density', 'Duplicated Lines (%)'],
    ['new_duplicated_lines_density', 'Duplicated Lines on New Code (%)'],
  ].map(([value, label]) => `<option value="${escAttr(value)}">${escHtml(label)}</option>`).join('');
}

export function formatNewCodePeriod(ncp) {
  if (!ncp || !ncp.strategy || ncp.scope === 'inherited') return 'auto (default)';
  const value = ncp.value ? ' \u2014 ' + escHtml(ncp.value) : '';
  return escHtml(ncp.strategy) + value;
}

export function formatSeconds(value) {
  if (value == null) return '-';
  const seconds = Number(value);
  if (seconds < 60) return seconds + 's';
  if (seconds < 3600) return Math.floor(seconds / 60) + 'm';
  if (seconds < 86400) return Math.floor(seconds / 3600) + 'h';
  return Math.floor(seconds / 86400) + 'd';
}

export function formatParameterValue(value) {
  if (value === undefined || value === null || value === '') return '-';
  if (Array.isArray(value)) return value.length ? value.join(', ') : '-';
  if (typeof value === 'boolean') return value ? 'yes' : 'no';
  if (typeof value === 'object') return JSON.stringify(value);
  return String(value);
}

export function customRuleLifecycleClass(lifecycle) {
  if (lifecycle === 'published' || lifecycle === 'valid') return 'badge-ok';
  if (lifecycle === 'invalid' || lifecycle === 'disabled') return 'badge-warn';
  return '';
}

export function customRuleStatusClass(status) {
  if (status === 'passed') return 'badge-ok';
  if (status === 'failed' || status === 'requires_runtime') return 'badge-warn';
  return '';
}

export function customRuleLifecycleOptions(selected) {
  return [
    ['all', 'All'],
    ['published', 'Published'],
    ['draft', 'Draft'],
    ['disabled', 'Disabled'],
  ].map(([value, label]) => `<option value="${escAttr(value)}"${selected === value ? ' selected' : ''}>${escHtml(label)}</option>`).join('');
}

export function ruleStudioStat(label, value, tone = '') {
  return `<div class="rule-studio-stat ${tone}"><span>${escHtml(label)}</span><strong>${fmtNum(value)}</strong></div>`;
}

export function customRuleMatchesFilters(rule, filters) {
  const lifecycle = rule.lifecycle || 'draft';
  if (filters.lifecycle && filters.lifecycle !== 'all' && lifecycle !== filters.lifecycle) {
    return false;
  }
  const query = (filters.search || '').trim().toLowerCase();
  if (!query) {
    return true;
  }
  return [rule.key, rule.name, rule.language, rule.engine, rule.type, rule.severity, rule.pack_name]
    .filter(Boolean)
    .some(value => String(value).toLowerCase().includes(query));
}

export function customRuleFullKey(namespace, key) {
  const cleanNamespace = (namespace || 'custom').trim().toLowerCase() || 'custom';
  const cleanKey = (key || '').trim().toLowerCase();
  if (!cleanKey) return cleanNamespace + ':<rule-id>';
  if (cleanKey.includes(':')) return cleanKey;
  return cleanNamespace + ':' + cleanKey;
}

export function customRuleEngineGuide(engine) {
  if (engine === 'go-ast') return 'Choose a built-in Go AST matcher and the exact call or import to flag.';
  if (engine === 'tree-sitter') return 'Write a structural Tree-sitter query for advanced language-specific matches.';
  return 'Use a regular expression against source text.';
}

export function customRuleDefaultEngine(engines) {
  if (engines.some(item => item.engine === 'text')) return 'text';
  return engines[0]?.engine || 'text';
}

export function customRuleAIStateLabel(model) {
  if (model.status === 'connected') return model.local ? 'local' : 'connected';
  if (model.status === 'setup_required') return 'setup required';
  if (model.status === 'unavailable') return 'unavailable';
  return model.status || '';
}

export function customRuleAIModelMessage(model) {
  if (model.message) return model.message;
  if (model.status === 'connected') return model.local ? 'Local model ready.' : 'Model ready.';
  if (model.status === 'setup_required') return 'Set up this provider before generating drafts.';
  if (model.status === 'unavailable') return 'This provider is currently unavailable.';
  return '';
}

export function customRuleAIModels(providers) {
  const models = [];
  providers.forEach(provider => {
    const options = provider.model_options || (provider.models || []).map(model => ({ id: model, label: model, status: provider.status || (provider.configured ? 'connected' : 'setup_required') }));
    options.forEach(option => {
      const status = option.status || provider.status || (provider.configured ? 'connected' : 'setup_required');
      models.push({
        provider: provider.id,
        providerLabel: provider.label || provider.id,
        model: option.id || option.label,
        label: option.label || option.id,
        status,
        local: Boolean(option.local || provider.local),
        setupRequired: Boolean(option.setup_required || provider.setup_required || status !== 'connected'),
        setupURL: provider.setup_url || '#custom-rule-ai-provider-setup',
        message: option.message || provider.message || '',
        defaultModel: provider.default_model,
      });
    });
  });
  const firstConnected = models.findIndex(model => model.status === 'connected' && (!model.defaultModel || model.model === model.defaultModel));
  const fallbackConnected = models.findIndex(model => model.status === 'connected');
  let selectedIndex = 0;
  if (firstConnected >= 0) {
    selectedIndex = firstConnected;
  } else if (fallbackConnected >= 0) {
    selectedIndex = fallbackConnected;
  }
  return models.map((model, index) => ({ ...model, selected: index === selectedIndex }));
}

export function customRuleAIModelOptions(models) {
  if (!models.length) {
    return '<option value="">No AI providers available</option>';
  }
  return models.map(model => {
    const state = customRuleAIStateLabel(model);
    const selected = model.selected ? ' selected' : '';
    return `<option value="${escAttr(model.provider + '|' + model.model)}" data-provider="${escAttr(model.provider)}" data-provider-label="${escAttr(model.providerLabel)}" data-model="${escAttr(model.model)}" data-status="${escAttr(model.status)}" data-local="${model.local ? 'true' : 'false'}" data-setup-url="${escAttr(model.setupURL)}" data-message="${escAttr(customRuleAIModelMessage(model))}"${selected}>${escHtml(model.providerLabel)} / ${escHtml(model.label)}${state ? ' - ' + escHtml(state) : ''}</option>`;
  }).join('');
}

export function renderCustomRuleAIProviderCard(provider, active = false) {
  const status = provider.status || (provider.configured ? 'connected' : 'setup_required');
  const models = provider.models || [];
  const diagnostics = provider.diagnostics || [];
  const modelCountText = models.length === 1 ? '1 model available.' : fmtNum(models.length) + ' models available.';
  return `<article class="custom-rule-ai-provider-card${active ? ' active' : ''}" data-ai-provider-card="${escAttr(provider.id)}">
    <div>
      <strong>${escHtml(provider.label || provider.id)}</strong>
      <span class="badge ${status === 'connected' ? 'badge-ok' : 'badge-warn'}">${escHtml(customRuleAIStateLabel({ status, local: provider.local }))}</span>
    </div>
    <p>${escHtml(provider.message || 'Configure this provider to use it from Rule Studio.')}</p>
    ${provider.base_url ? `<p class="mono">${escHtml(provider.base_url)}</p>` : ''}
    ${models.length ? `<p>${modelCountText}</p><p class="mono">${escHtml(models.join(', '))}</p>` : ''}
    ${diagnostics.length ? `<p class="custom-rule-ai-provider-diagnostic">${escHtml(diagnostics[0])}</p>` : ''}
  </article>`;
}
