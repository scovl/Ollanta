import { apiFetch } from '../../core/api.js';
import { state } from '../../core/state.js';
import { escAttr, escHtml } from '../../core/utils.js';
import { getRenderView, getShowToast } from './context.js';
import { formatNewCodePeriod } from './shared.js';

export async function loadWebhooksData() {
  const project = state.currentProject;
  if (!project) return;
  try {
    const data = await apiFetch('/webhooks?project_key=' + encodeURIComponent(project.key));
    state.webhooksData = data.items || (Array.isArray(data) ? data : []);
  } catch {
    state.webhooksData = [];
  }
  try {
    state.newCodePeriod = await apiFetch('/projects/' + encodeURIComponent(project.key) + '/new-code-period');
  } catch {
    state.newCodePeriod = null;
  }
  getRenderView()();
}

export function renderWebhooksTab() {
  const webhooks = state.webhooksData;
  if (webhooks === null) return `<div class="loading-state"><div class="spinner"></div></div>`;

  const ncp = state.newCodePeriod;
  const ncpStr = formatNewCodePeriod(ncp);

  const webhookRows = webhooks.length === 0
    ? `<div class="empty-state" style="padding:20px 0"><p>No webhooks configured.</p></div>`
    : webhooks.map(webhook => `
      <div class="webhook-row">
        <div class="webhook-info">
          <span class="webhook-name">${escHtml(webhook.name)}</span>
          <span class="webhook-url" title="${escAttr(webhook.url)}">${escHtml(webhook.url)}</span>
        </div>
        <div class="webhook-btns">
          <button class="btn-sm btn-outline test-wh-btn" data-wh-id="${webhook.id}">Test</button>
          <button class="btn-sm btn-danger del-wh-btn" data-wh-id="${webhook.id}">Delete</button>
        </div>
      </div>`).join('');

  return `<div class="tab-section">
    <p class="section-title" style="margin-top:24px">Webhooks</p>
    <div class="webhook-list">${webhookRows}</div>
    <div class="create-form">
      <h4 style="font-size:14px;font-weight:600;margin-bottom:12px">Add webhook</h4>
      <div class="form-row">
        <input id="newWhName" class="filter-input" placeholder="Name" style="width:150px">
        <input id="newWhUrl" class="filter-input" placeholder="https://\u2026" style="flex:1;min-width:200px">
        <input id="newWhSecret" class="filter-input" placeholder="Secret (optional)" style="width:160px">
        <button class="btn btn-primary" id="addWhBtn" style="width:auto;padding:6px 18px;margin-top:0">Add</button>
      </div>
    </div>

    <p class="section-title" style="margin-top:32px">New Code Period</p>
    <div class="scan-info" style="grid-template-columns:1fr auto;gap:16px;align-items:center">
      <div>
        <div class="info-label">Current strategy</div>
        <div class="info-value" id="ncpDisplay">${ncpStr}</div>
      </div>
      <div class="form-row" style="justify-content:flex-end">
        <select id="ncpStrategy" class="filter-sel">
          <option value="auto"${(!ncp || ncp.strategy === 'auto') ? ' selected' : ''}>Auto</option>
          <option value="previous_version"${ncp?.strategy === 'previous_version' ? ' selected' : ''}>Previous version</option>
          <option value="number_of_days"${ncp?.strategy === 'number_of_days' ? ' selected' : ''}>Number of days</option>
          <option value="reference_branch"${ncp?.strategy === 'reference_branch' ? ' selected' : ''}>Reference branch</option>
        </select>
        <input id="ncpValue" class="filter-input" placeholder="Value (if needed)" style="width:140px" value="${escAttr(ncp?.value || '')}">
        <button class="btn btn-primary" id="saveNcpBtn" style="width:auto;padding:6px 18px;margin-top:0">Save</button>
      </div>
    </div>
  </div>`;
}

export function bindWebhooksContent() {
  const project = state.currentProject;
  if (!project) return;
  const showToast = getShowToast();

  document.getElementById('addWhBtn')?.addEventListener('click', async () => {
    const name = document.getElementById('newWhName')?.value.trim();
    const url = document.getElementById('newWhUrl')?.value.trim();
    const secret = document.getElementById('newWhSecret')?.value.trim();
    if (!name || !url) {
      showToast('Name and URL are required.', 'error');
      return;
    }
    try {
      await apiFetch('/webhooks', {
        method: 'POST',
        body: JSON.stringify({ name, url, secret: secret || '', project_key: project.key }),
      });
      state.webhooksData = null;
      await loadWebhooksData();
    } catch (err) {
      showToast(err.message, 'error');
    }
  });

  document.querySelectorAll('.test-wh-btn').forEach(btn => {
    btn.addEventListener('click', async () => {
      btn.disabled = true;
      try {
        await apiFetch('/webhooks/' + btn.dataset.whId + '/test', { method: 'POST' });
        showToast('Test delivery sent.');
      } catch (err) {
        showToast(err.message, 'error');
      }
      btn.disabled = false;
    });
  });

  document.querySelectorAll('.del-wh-btn').forEach(btn => {
    btn.addEventListener('click', async () => {
      if (!confirm('Delete this webhook?')) return;
      btn.disabled = true;
      try {
        await apiFetch('/webhooks/' + btn.dataset.whId, { method: 'DELETE' });
        state.webhooksData = null;
        await loadWebhooksData();
      } catch (err) {
        showToast(err.message, 'error');
        btn.disabled = false;
      }
    });
  });

  document.getElementById('saveNcpBtn')?.addEventListener('click', async () => {
    const strategy = document.getElementById('ncpStrategy')?.value;
    const value = document.getElementById('ncpValue')?.value.trim();
    try {
      await apiFetch('/projects/' + encodeURIComponent(project.key) + '/new-code-period', {
        method: 'PUT',
        body: JSON.stringify({ strategy, value: value || '' }),
      });
      state.newCodePeriod = { strategy, value };
      const display = document.getElementById('ncpDisplay');
      if (display) display.textContent = strategy + (value ? ' \u2014 ' + value : '');
      showToast('New code period saved.');
    } catch (err) {
      showToast(err.message, 'error');
    }
  });
}
