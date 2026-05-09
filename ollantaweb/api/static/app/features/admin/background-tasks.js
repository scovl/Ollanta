import { apiFetch } from '../../core/api.js';
import { state } from '../../core/state.js';
import { escAttr, escHtml, fmtDate, fmtNum } from '../../core/utils.js';
import { getRenderView, getShowToast } from './context.js';
import { formatParameterValue, formatSeconds } from './shared.js';

const TASK_STATUSES = ['', 'queued', 'running', 'retrying', 'stale', 'failed', 'completed', 'cancelled'];
const TASK_TYPES = ['', 'scan', 'index', 'webhook'];

function currentTaskFilters() {
  const filters = state.backgroundTaskFilters || {};
  const project = state.currentProject;
  return {
    type: filters.type || '',
    status: filters.status || '',
    projectKey: filters.projectKey ?? (project?.key || ''),
    workerId: filters.workerId || '',
    scanId: filters.scanId || '',
    limit: filters.limit || 25,
    offset: filters.offset || 0,
  };
}

function buildTaskQuery(filters) {
  const params = new URLSearchParams();
  if (filters.type) params.set('type', filters.type);
  if (filters.status) params.set('status', filters.status);
  if (filters.projectKey) params.set('project_key', filters.projectKey);
  if (filters.workerId) params.set('worker_id', filters.workerId);
  if (filters.scanId) params.set('scan_id', filters.scanId);
  params.set('limit', String(filters.limit || 25));
  params.set('offset', String(filters.offset || 0));
  return params.toString();
}

export async function loadBackgroundTasksData(options = {}) {
  const filters = { ...currentTaskFilters(), ...options };
  state.backgroundTaskFilters = filters;
  state.loadingBackgroundTasks = true;
  state.backgroundTaskError = '';
  getRenderView()();
  try {
    const query = buildTaskQuery(filters);
    const [list, summary] = await Promise.all([
      apiFetch('/admin/background-tasks?' + query),
      apiFetch('/admin/background-tasks/summary?' + query),
    ]);
    state.backgroundTasksData = list;
    state.backgroundTasksSummary = summary;
  } catch (err) {
    state.backgroundTasksData = { items: [], total: 0, limit: filters.limit, offset: filters.offset };
    state.backgroundTasksSummary = null;
    state.backgroundTaskError = err.message || 'Failed to load background tasks.';
  }
  state.loadingBackgroundTasks = false;
  getRenderView()();
}

export async function loadBackgroundTaskDetail(taskId) {
  if (!taskId) return;
  state.loadingBackgroundTaskDetail = true;
  state.selectedBackgroundTask = null;
  state.backgroundTaskError = '';
  getRenderView()();
  try {
    state.selectedBackgroundTask = await apiFetch('/admin/background-tasks/' + encodeURIComponent(taskId));
  } catch (err) {
    state.backgroundTaskError = err.message || 'Failed to load task details.';
  }
  state.loadingBackgroundTaskDetail = false;
  getRenderView()();
}

export function renderAdminLinksTab() {
  return renderBackgroundTasksPage();
}

export function renderBackgroundTasksPage() {
  const filters = currentTaskFilters();
  const data = state.backgroundTasksData || { items: [], total: 0, limit: filters.limit, offset: filters.offset };
  const items = data.items || [];
  const total = data.total || 0;
  const summary = state.backgroundTasksSummary;
  const selected = state.selectedBackgroundTask;
  const canPrev = (filters.offset || 0) > 0;
  const canNext = (filters.offset || 0) + (filters.limit || 25) < total;
  const globalHeader = state.view === 'background-tasks'
    ? `<div class="page-header background-task-page-title">
        <div>
          <h2>Background Tasks</h2>
          <p>Operational queue for scan intake, indexing, and webhook delivery.</p>
        </div>
        <button class="back-btn" id="backgroundTasksBackBtn">Back to projects</button>
      </div>`
    : '';

  return `${globalHeader}<div class="tab-section background-tasks-page">
    <div class="background-task-header">
      <div>
        <p class="section-title">Queue activity</p>
        <p class="background-task-subtitle">Tasks are grouped by project context first; the technical job id is shown as secondary metadata.</p>
      </div>
      <button class="btn-sm btn-outline" id="refreshBackgroundTasksBtn">Refresh</button>
    </div>
    ${state.backgroundTaskError ? `<div class="error-msg background-task-error">${escHtml(state.backgroundTaskError)}</div>` : ''}
    ${renderTaskSummary(summary)}
    ${renderTaskFilters(filters)}
    <div class="background-task-layout">
      <section class="background-task-table-panel">
        ${state.loadingBackgroundTasks ? `<div class="loading-state"><div class="spinner"></div></div>` : renderTaskTable(items)}
        <div class="background-task-pagination">
          <span>${fmtNum(total)} tasks</span>
          <button class="btn-sm btn-outline" data-task-page="prev"${canPrev ? '' : ' disabled'}>Previous</button>
          <button class="btn-sm btn-outline" data-task-page="next"${canNext ? '' : ' disabled'}>Next</button>
        </div>
      </section>
      <aside class="background-task-detail-panel">
        ${state.loadingBackgroundTaskDetail ? `<div class="loading-state"><div class="spinner"></div></div>` : renderTaskDetail(selected)}
      </aside>
    </div>
  </div>`;
}

function renderTaskSummary(summary) {
  const cards = [
    ['Queued', summary?.queue_depth || 0, 'queued'],
    ['Running', summary?.running_count || 0, 'running'],
    ['Failed', summary?.failed_count || 0, 'failed'],
    ['Stale', summary?.stale_count || 0, 'stale'],
    ['Retrying', summary?.retry_count || 0, 'retrying'],
    ['Recent completions', summary?.recent_completion_count || 0, 'completed'],
  ];
  return `<div class="background-task-summary-grid">
    ${cards.map(([label, value, status]) => `<button class="background-task-summary-card" data-task-status-filter="${escAttr(status)}">
      <span>${escHtml(label)}</span>
      <strong>${fmtNum(value)}</strong>
    </button>`).join('')}
  </div>`;
}

function renderTaskFilters(filters) {
  return `<div class="background-task-filters">
    <label>Type<select class="filter-sel" id="taskTypeFilter">${TASK_TYPES.map(value => `<option value="${escAttr(value)}"${filters.type === value ? ' selected' : ''}>${value ? escHtml(value) : 'All'}</option>`).join('')}</select></label>
    <label>Status<select class="filter-sel" id="taskStatusFilter">${TASK_STATUSES.map(value => `<option value="${escAttr(value)}"${filters.status === value ? ' selected' : ''}>${value ? escHtml(value) : 'All'}</option>`).join('')}</select></label>
    <label>Project<input class="filter-input" id="taskProjectFilter" value="${escAttr(filters.projectKey || '')}" placeholder="project key"></label>
    <label>Scan<input class="filter-input" id="taskScanFilter" value="${escAttr(filters.scanId || '')}" placeholder="scan id"></label>
    <label>Worker<input class="filter-input" id="taskWorkerFilter" value="${escAttr(filters.workerId || '')}" placeholder="worker id"></label>
    <button class="btn-sm btn-primary" id="applyTaskFiltersBtn">Apply</button>
    <button class="btn-sm btn-outline" id="clearTaskFiltersBtn">Clear</button>
  </div>`;
}

function renderTaskTable(items) {
  if (!items.length) return `<div class="empty-state"><p>No background tasks match the current filters.</p></div>`;
  return `<table class="background-task-table">
    <thead><tr><th>Task</th><th>Status</th><th>Type</th><th>Submitted</th><th>Started</th><th>Duration</th><th>Worker</th><th>Scan</th><th>Error</th></tr></thead>
    <tbody>${items.map(task => `<tr class="background-task-row" data-task-id="${escAttr(task.id)}">
      <td>${renderTaskName(task)}</td>
      <td><span class="task-status task-status-${escAttr(task.status)}">${escHtml(task.status)}</span></td>
      <td>${renderTaskType(task)}</td>
      <td>${fmtDate(task.created_at)}</td>
      <td>${fmtDate(task.started_at)}</td>
      <td>${formatSeconds(task.duration_seconds)}</td>
      <td>${escHtml(task.worker_id || '-')}</td>
      <td>${escHtml(task.scan_id || '-')}</td>
      <td class="task-error-cell" title="${escAttr(task.last_error || '')}">${escHtml(task.last_error || '-')}</td>
    </tr>`).join('')}</tbody>
  </table>`;
}

function renderTaskName(task) {
  const primary = task.project_key || (task.project_id ? 'Project #' + task.project_id : task.type + ' task');
  const subtitle = [task.id, task.scan_id ? 'scan #' + task.scan_id : '', task.attempts ? 'attempts ' + task.attempts : ''].filter(Boolean).join(' · ');
  return `<div class="task-name-cell">
    <strong>${escHtml(primary)}</strong>
    <span>${escHtml(subtitle || task.id)}</span>
  </div>`;
}

function renderTaskType(task) {
  const labels = {
    scan: 'Scan processing',
    index: 'Search indexing',
    webhook: 'Webhook delivery',
  };
  return `<span class="task-type-chip task-type-${escAttr(task.type)}">${escHtml(labels[task.type] || task.type)}</span>`;
}

function renderTaskDetail(task) {
  if (!task) return `<div class="empty-state"><p>Select a task to inspect diagnostics and supported actions.</p></div>`;
  const details = task.details || {};
  const detailRows = Object.entries(details).map(([key, value]) => `<div class="info-row"><span>${escHtml(key)}</span><strong>${escHtml(String(value))}</strong></div>`).join('');
  const actions = task.supported_actions || [];
  const scannerParameters = renderScannerParameters(task.scanner_parameters);
  return `<div class="background-task-detail">
    <div class="task-detail-head">
      <div class="task-name-cell">
        <strong>${escHtml(task.project_key || task.id)}</strong>
        <span>${escHtml(task.id)}</span>
      </div>
      <span class="task-status task-status-${escAttr(task.status)}">${escHtml(task.status)}</span>
    </div>
    <div class="scan-info task-detail-grid">
      <div><div class="info-label">Type</div><div class="info-value">${escHtml(task.type)}</div></div>
      <div><div class="info-label">Internal status</div><div class="info-value">${escHtml(task.internal_status)}</div></div>
      <div><div class="info-label">Created</div><div class="info-value">${fmtDate(task.created_at)}</div></div>
      <div><div class="info-label">Started</div><div class="info-value">${fmtDate(task.started_at)}</div></div>
      <div><div class="info-label">Completed</div><div class="info-value">${fmtDate(task.completed_at)}</div></div>
      <div><div class="info-label">Next retry</div><div class="info-value">${fmtDate(task.next_attempt_at)}</div></div>
    </div>
    ${scannerParameters}
    ${task.last_error ? `<div class="task-diagnostic"><span>Last error</span><pre>${escHtml(task.last_error)}</pre></div>` : ''}
    <div class="task-detail-section"><h4>Details</h4>${detailRows || '<p class="muted">No type-specific details.</p>'}</div>
    <div class="task-detail-actions">
      ${actions.map(action => `<button class="btn-sm ${action === 'cancel' ? 'btn-danger' : 'btn-outline'}" data-task-action="${escAttr(action)}" data-task-id="${escAttr(task.id)}">${escHtml(action)}</button>`).join('') || '<span class="muted">No actions available for this state.</span>'}
    </div>
  </div>`;
}

function renderScannerParameters(parameters) {
  if (!parameters || !Object.keys(parameters).length) {
    return `<div class="task-detail-section scanner-params-section">
      <h4>Scanner parameters</h4>
      <p class="muted">Scanner parameters were not captured for this task. Run a new scan with the updated scanner to populate this section.</p>
    </div>`;
  }

  const options = parameters.scanner_options || {};
  const scope = parameters.analysis_scope || {};
  const tests = options.tests || parameters.test_signals || {};
  const rows = [
    ['Config file', options.config_path],
    ['Project directory', options.project_dir],
    ['Project key', options.project_key || scope.project_key],
    ['Sources', options.sources],
    ['Exclusions', options.exclusions],
    ['Format', options.format],
    ['Branch', options.branch || scope.branch],
    ['Commit', options.commit_sha || scope.commit_sha],
    ['Pull request', options.pull_request_key || scope.pull_request_key],
    ['Pull request branch', options.pull_request_branch],
    ['Pull request base', options.pull_request_base || scope.pull_request_base],
    ['Server URL', options.server],
    ['Wait for server job', options.server_wait],
    ['Wait timeout', options.server_wait_timeout],
    ['Wait poll', options.server_wait_poll],
    ['Local UI', options.local_ui],
    ['Local UI bind', options.bind],
    ['Local UI port', options.port],
    ['Debug', options.debug],
  ];

  return `<div class="task-detail-section scanner-params-section">
    <h4>Scanner parameters</h4>
    <div class="scanner-params-grid">${rows.map(([label, value]) => renderParameterRow(label, value)).join('')}</div>
    ${renderTestParameters(tests)}
  </div>`;
}

function renderParameterRow(label, value) {
  return `<div class="info-row"><span>${escHtml(label)}</span><strong>${escHtml(formatParameterValue(value))}</strong></div>`;
}

function renderTestParameters(tests) {
  if (!tests || !Object.keys(tests).length) return '';
  const summary = tests.summary || {};
  const modules = Array.isArray(tests.modules) ? tests.modules : [];
  const rows = [
    ['Tests enabled', tests.enabled ?? summary.enabled],
    ['Mode', tests.mode],
    ['Discover modules', tests.discover],
    ['Run commands', tests.run],
    ['Command policy', tests.command_policy],
    ['Max report age', tests.max_report_age],
    ['Max depth', tests.max_depth],
    ['Max candidates', tests.max_candidates],
    ['Max report bytes', tests.max_report_bytes],
  ];
  const moduleCards = modules.length
    ? `<div class="scanner-test-modules">${modules.map(module => renderTestModuleParameters(module)).join('')}</div>`
    : '';
  return `<div class="scanner-test-params">
    <h5>Test signal parameters</h5>
    <div class="scanner-params-grid">${rows.map(([label, value]) => renderParameterRow(label, value)).join('')}</div>
    ${moduleCards}
  </div>`;
}

function renderTestModuleParameters(module) {
  const rows = [
    ['Root', module.root],
    ['Language', module.language],
    ['Role', module.architecture_role],
    ['Policy', module.test_policy],
    ['Command', module.command],
    ['Artifact root', module.artifact_root],
    ['Report root', module.report_root],
    ['Coverage reports', module.coverage_reports],
    ['Test reports', module.test_reports],
    ['Mutation reports', module.mutation_reports],
    ['Native reports', module.native_reports],
    ['Owner', module.owner],
    ['Team', module.team],
    ['Integration required', module.integration_required],
  ];
  const title = module.name || module.root || 'test module';
  return `<div class="scanner-test-module">
    <strong>${escHtml(title)}</strong>
    <div class="scanner-params-grid">${rows.map(([label, value]) => renderParameterRow(label, value)).join('')}</div>
  </div>`;
}

export function bindBackgroundTasksContent() {
  const showToast = getShowToast();

  document.getElementById('refreshBackgroundTasksBtn')?.addEventListener('click', () => loadBackgroundTasksData());
  document.getElementById('applyTaskFiltersBtn')?.addEventListener('click', () => {
    loadBackgroundTasksData({
      type: document.getElementById('taskTypeFilter')?.value || '',
      status: document.getElementById('taskStatusFilter')?.value || '',
      projectKey: document.getElementById('taskProjectFilter')?.value.trim() || '',
      scanId: document.getElementById('taskScanFilter')?.value.trim() || '',
      workerId: document.getElementById('taskWorkerFilter')?.value.trim() || '',
      offset: 0,
    });
  });
  document.getElementById('clearTaskFiltersBtn')?.addEventListener('click', () => loadBackgroundTasksData({ type: '', status: '', projectKey: '', scanId: '', workerId: '', offset: 0 }));
  document.querySelectorAll('[data-task-status-filter]').forEach(btn => {
    btn.addEventListener('click', () => loadBackgroundTasksData({ status: btn.dataset.taskStatusFilter || '', offset: 0 }));
  });
  document.querySelectorAll('[data-task-page]').forEach(btn => {
    btn.addEventListener('click', () => {
      const filters = currentTaskFilters();
      const delta = btn.dataset.taskPage === 'next' ? filters.limit : -filters.limit;
      loadBackgroundTasksData({ offset: Math.max(0, filters.offset + delta) });
    });
  });
  document.querySelectorAll('[data-task-id].background-task-row').forEach(row => {
    row.addEventListener('click', () => loadBackgroundTaskDetail(row.dataset.taskId));
  });
  document.querySelectorAll('[data-task-action]').forEach(btn => {
    btn.addEventListener('click', async () => {
      const action = btn.dataset.taskAction;
      const id = btn.dataset.taskId;
      if (action === 'cancel' && !confirm('Cancel this queued background task?')) return;
      btn.disabled = true;
      try {
        state.selectedBackgroundTask = await apiFetch('/admin/background-tasks/' + encodeURIComponent(id) + '/' + action, { method: 'POST' });
        await loadBackgroundTasksData();
        showToast('Task ' + action + ' accepted.');
      } catch (err) {
        showToast(err.message || 'Task action failed.', 'error');
        btn.disabled = false;
      }
    });
  });
}
