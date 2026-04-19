'use strict';

const API = '/api/v1';

let state = {
  user: null,
  view: 'login',   // 'login' | 'projects' | 'project'
  projects: [],
  currentProject: null,
  currentScan: null,
  overviewData: null,    // from /overview API
  issues: [],
  issuesTotal: 0,
  issueOffset: 0,
  issueFilter: { severity: 'all', type: 'all', status: 'all', search: '' },
  loading: false,
  loadingIssues: false,
  projectTab: 'overview',  // 'overview' | 'issues' | 'activity' | 'gate' | 'webhooks' | 'profiles'
  gateData: null,
  webhooksData: null,
  profilesData: null,
  activityData: null,
  newCodePeriod: null,
  selectedIssue: null,
};

// ── Storage ───────────────────────────────────────────────────────────────────

function getToken()  { return localStorage.getItem('olt_token'); }
function saveToken(t) { localStorage.setItem('olt_token', t); }
function clearStorage() {
  localStorage.removeItem('olt_token');
  localStorage.removeItem('olt_user');
}
function saveUser(u) { localStorage.setItem('olt_user', JSON.stringify(u)); }
function loadUser()  {
  try { return JSON.parse(localStorage.getItem('olt_user') || 'null'); }
  catch { return null; }
}

// ── API helper ────────────────────────────────────────────────────────────────

async function apiFetch(path, opts = {}) {
  const headers = { 'Content-Type': 'application/json' };
  const t = getToken();
  if (t) headers['Authorization'] = 'Bearer ' + t;
  if (opts.headers) Object.assign(headers, opts.headers);

  const res = await fetch(API + path, { ...opts, headers });

  if (res.status === 401) {
    logout();
    throw new Error('Session expired');
  }
  if (res.status === 204) return null;
  const body = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(body.error || res.statusText);
  return body;
}

// ── Formatters ────────────────────────────────────────────────────────────────

function fmtDate(d) {
  if (!d) return '\u2014';
  const date = new Date(d);
  const diff = Date.now() - date.getTime();
  if (diff < 60_000)     return 'just now';
  if (diff < 3_600_000)  return Math.floor(diff / 60_000) + 'm ago';
  if (diff < 86_400_000) return Math.floor(diff / 3_600_000) + 'h ago';
  if (diff < 604_800_000) return Math.floor(diff / 86_400_000) + 'd ago';
  return date.toLocaleDateString();
}

function fmtNum(n) {
  return (n == null ? 0 : Number(n)).toLocaleString();
}

function fmtK(n) {
  if (n == null) n = 0;
  if (n >= 1_000_000) return (n / 1_000_000).toFixed(1) + 'M';
  if (n >= 1_000) return (n / 1_000).toFixed(1) + 'k';
  return String(n);
}

function fmtPct(n) {
  if (n == null) return '\u2014';
  return Number(n).toFixed(1) + '%';
}

// ── Constants ─────────────────────────────────────────────────────────────────

const SEV_ORDER  = ['blocker','critical','major','minor','info'];
const SEV_COLOR  = { blocker:'#ef4444', critical:'#f97316', major:'#eab308', minor:'#22c55e', info:'#64748b' };
const SEV_BG     = { blocker:'rgba(239,68,68,.12)', critical:'rgba(249,115,22,.10)', major:'rgba(234,179,8,.09)', minor:'rgba(34,197,94,.09)', info:'rgba(100,116,139,.09)' };
const SEV_LABEL  = { blocker:'Blocker', critical:'Critical', major:'Major', minor:'Minor', info:'Info' };
const TYPE_ICON  = { bug:'\uD83D\uDC1B', code_smell:'\uD83C\uDF3F', vulnerability:'\uD83D\uDD12' };
const TYPE_COLOR = { bug:'#ef4444', code_smell:'#22c55e', vulnerability:'#f97316' };
const TYPE_LABEL = { bug:'Bug', code_smell:'Code Smell', vulnerability:'Vulnerability' };

// ── Render ────────────────────────────────────────────────────────────────────

function render() {
  const app = document.getElementById('app');
  if (state.view === 'login') {
    app.innerHTML = renderLogin();
    bindLogin();
    return;
  }
  app.innerHTML = renderNav() + '<main>' + renderContent() + '</main>';
  bindMain();
}

function renderNav() {
  const u = state.user || {};
  const name = u.name || u.login || 'User';
  return `<nav>
    <span class="logo">
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
      Ollanta
    </span>
    <span class="spacer"></span>
    <span class="user-info">${escHtml(name)}</span>
    <button class="logout-btn" id="logoutBtn">Sign out</button>
  </nav>`;
}

function renderContent() {
  if (state.view === 'projects') return renderDashboard();
  if (state.view === 'project')  return renderProjectDetail();
  return '';
}

// ── Login ─────────────────────────────────────────────────────────────────────

function renderLogin() {
  return `<div class="login-wrapper">
    <div class="login-card">
      <h1>
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        Ollanta
      </h1>
      <p class="subtitle">Static analysis platform</p>
      <div class="field">
        <label for="loginUser">Username</label>
        <input id="loginUser" type="text" placeholder="admin" autocomplete="username">
      </div>
      <div class="field">
        <label for="loginPass">Password</label>
        <input id="loginPass" type="password" placeholder="\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022" autocomplete="current-password">
      </div>
      <button class="btn btn-primary" id="loginBtn">Sign in</button>
      <div id="loginError" class="error-msg"></div>
    </div>
  </div>`;
}

function bindLogin() {
  const btn   = document.getElementById('loginBtn');
  const errEl = document.getElementById('loginError');
  const userEl = document.getElementById('loginUser');
  const passEl = document.getElementById('loginPass');

  async function doLogin() {
    const login    = userEl.value.trim();
    const password = passEl.value;
    if (!login || !password) { errEl.textContent = 'Enter username and password.'; return; }

    btn.disabled = true;
    btn.textContent = 'Signing in\u2026';
    errEl.textContent = '';

    try {
      const data = await apiFetch('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ login, password }),
      });
      saveToken(data.access_token);
      saveUser(data.user || {});
      state.user = data.user || {};
      await loadProjects();
    } catch (e) {
      errEl.textContent = e.message || 'Login failed.';
      btn.disabled = false;
      btn.textContent = 'Sign in';
    }
  }

  btn.addEventListener('click', doLogin);
  passEl.addEventListener('keydown', e => { if (e.key === 'Enter') doLogin(); });
  userEl.addEventListener('keydown', e => { if (e.key === 'Enter') passEl.focus(); });
}

// ── Dashboard ─────────────────────────────────────────────────────────────────

async function loadProjects() {
  state.view    = 'projects';
  state.loading = true;
  render();

  try {
    const data = await apiFetch('/projects?limit=100');
    state.projects = data.items || [];
  } catch {
    state.projects = [];
  }

  state.loading = false;
  render();
}

function renderDashboard() {
  if (state.loading) {
    return `<div class="loading-state"><div class="spinner"></div></div>`;
  }

  const ps = state.projects;
  const count = ps.length;

  return `
    <div class="page-header">
      <h2>Projects <span style="font-size:14px;color:var(--text-muted);font-weight:400">(${count})</span></h2>
      <p>All projects registered on this platform</p>
    </div>
    ${count === 0
      ? `<div class="empty-state">
           <div class="empty-icon">\uD83D\uDCC2</div>
           <p>No projects yet. Run a scan to register the first project.</p>
         </div>`
      : `<div class="projects-grid">${ps.map(renderProjectCard).join('')}</div>`
    }`;
}

function renderProjectCard(p) {
  const tags = (p.tags || []).filter(Boolean);
  const tagsHtml = tags.length
    ? `<div class="tags">${tags.map(t => `<span class="tag">${escHtml(t)}</span>`).join('')}</div>`
    : '';

  // Gate status from latest scan info (if available in project list)
  const gs = p.gate_status || '';
  const gateCls = gs === 'OK' ? 'card-gate-ok' : gs === 'ERROR' ? 'card-gate-error' : gs === 'WARN' ? 'card-gate-warn' : '';
  const gateBadge = gs ? `<span class="badge ${gs === 'OK' ? 'badge-ok' : gs === 'WARN' ? 'badge-warn' : 'badge-error'}">${escHtml(gs)}</span>` : '';

  return `<div class="project-card ${gateCls}" data-key="${escAttr(p.key)}">
    <div class="card-top">
      <span class="key">${escHtml(p.key)}</span>
      ${gateBadge}
    </div>
    <div class="name" title="${escAttr(p.name || p.key)}">${escHtml(p.name || p.key)}</div>
    ${tagsHtml}
    <div class="footer">Updated ${fmtDate(p.updated_at)}</div>
  </div>`;
}

// ── Project detail ────────────────────────────────────────────────────────────

const ISSUE_PAGE = 50;

async function loadProject(key) {
  state.view           = 'project';
  state.currentProject = null;
  state.currentScan    = null;
  state.overviewData   = null;
  state.issues         = [];
  state.issuesTotal    = 0;
  state.issueOffset    = 0;
  state.issueFilter    = { severity: 'all', type: 'all', status: 'all', search: '' };
  state.projectTab     = 'overview';
  state.gateData       = null;
  state.webhooksData   = null;
  state.profilesData   = null;
  state.activityData   = null;
  state.newCodePeriod  = null;
  state.selectedIssue  = null;
  state.loading        = true;
  render();

  try {
    const [project, overview] = await Promise.all([
      apiFetch('/projects/' + encodeURIComponent(key)),
      apiFetch('/projects/' + encodeURIComponent(key) + '/overview').catch(() => null),
    ]);
    state.currentProject = project;
    state.overviewData   = overview;
    state.currentScan    = overview?.last_scan || null;
  } catch { /* ignore */ }

  state.loading = false;
  render();
}

function renderProjectDetail() {
  const backBtn = `<button class="back-btn" id="backBtn">\u2190 Projects</button>`;

  if (state.loading) {
    return backBtn + `<div class="loading-state"><div class="spinner"></div></div>`;
  }

  const p = state.currentProject;
  if (!p) {
    return backBtn + `<div class="empty-state"><p>Project not found.</p></div>`;
  }

  const s = state.currentScan;
  const gateCls = !s ? '' : s.gate_status === 'OK' ? 'badge-ok' : s.gate_status === 'WARN' ? 'badge-warn' : 'badge-error';
  const gateBadge = s && s.gate_status ? `<span class="badge ${gateCls}">${escHtml(s.gate_status)}</span>` : '';
  const desc = [p.description, (p.tags || []).filter(Boolean).join(', ')].filter(Boolean).join(' \u00B7 ');

  const tab = state.projectTab;
  const tabs = ['overview','issues','activity','gate','webhooks','profiles'];
  const tabLabels = { overview: 'Overview', issues: 'Issues', activity: 'Activity', gate: 'Quality Gate', webhooks: 'Webhooks', profiles: 'Profiles' };
  const issueCount = state.overviewData?.last_scan?.total_issues ?? '';
  const tabsHtml = `<div class="proj-tabs">${tabs.map(t => {
    let badge = '';
    if (t === 'issues' && issueCount !== '') badge = `<span class="tab-badge">${fmtK(issueCount)}</span>`;
    return `<button class="tab-btn${t===tab?' active':''}" data-tab="${t}">${tabLabels[t]}${badge}</button>`;
  }).join('')}</div>`;

  let tabContent = '';
  if (tab === 'overview') {
    tabContent = renderOverviewTab();
  } else if (tab === 'issues') {
    tabContent = `<div id="issues-section"></div>`;
  } else if (tab === 'activity') {
    tabContent = renderActivityTab();
  } else if (tab === 'gate') {
    tabContent = renderGateTab();
  } else if (tab === 'webhooks') {
    tabContent = renderWebhooksTab();
  } else if (tab === 'profiles') {
    tabContent = renderProfilesTab();
  }

  return `
    ${backBtn}
    <div class="detail-header">
      <h2>${escHtml(p.name || p.key)} ${gateBadge}</h2>
      <p>${escHtml(p.key)}${desc ? ' \u2014 ' + escHtml(desc) : ''}</p>
    </div>
    ${tabsHtml}
    ${tabContent}`;
}

// ── Overview Tab ──────────────────────────────────────────────────────────────

function renderOverviewTab() {
  const o = state.overviewData;
  if (!o) {
    return `<div class="empty-state">
      <div class="empty-icon">\uD83D\uDD2C</div>
      <p>No scans yet for this project.<br>Run <code>ollanta</code> to submit a scan.</p>
    </div>`;
  }

  const s = o.last_scan;
  const gate = o.quality_gate || {};
  const m = o.measures || {};
  const facets = o.facets || {};
  const nc = o.new_code || {};

  // Gate hero
  const gateHeroHtml = renderGateHero(gate);

  // Metric cards
  const bugs = m.bugs || 0;
  const vulns = m.vulnerabilities || 0;
  const smells = m.code_smells || 0;
  const coverage = m.coverage;
  const dupDensity = m.duplicated_lines_density;
  const ncloc = m.ncloc || 0;

  const metricsHtml = `<div class="measures-row">
    ${metricCard('Bugs', bugs, bugs > 0 ? 'danger' : 'success', bugs > 0 ? 'card-red' : 'card-green', 'bug')}
    ${metricCard('Vulnerabilities', vulns, vulns > 0 ? 'warning' : 'success', vulns > 0 ? 'card-yellow' : 'card-green', 'vulnerability')}
    ${metricCard('Code Smells', smells, 'muted', smells > 20 ? 'card-yellow' : 'card-green', 'code_smell')}
    ${metricCardPct('Coverage', coverage, coverage == null ? 'card-neutral' : coverage >= 80 ? 'card-green' : coverage >= 60 ? 'card-yellow' : 'card-red')}
    ${metricCardPct('Duplication', dupDensity, dupDensity == null ? 'card-neutral' : dupDensity <= 3 ? 'card-green' : dupDensity <= 10 ? 'card-yellow' : 'card-red')}
    ${metricCardK('Lines of Code', ncloc, 'card-neutral')}
  </div>`;

  // Severity distribution
  const sevDist = facets.by_severity || {};
  const sevDistHtml = renderDistribution('Severity', SEV_ORDER, SEV_LABEL, SEV_COLOR, sevDist);

  // Type distribution
  const typeDist = facets.by_type || {};
  const typeOrder = ['bug','code_smell','vulnerability'];
  const typeDistHtml = renderDistribution('Type', typeOrder, TYPE_LABEL, TYPE_COLOR, typeDist);

  // Hotspot files (top 10)
  const fileDist = facets.by_file || {};
  const fileEntries = Object.entries(fileDist).sort((a,b) => b[1] - a[1]).slice(0, 10);
  const hotspotHtml = fileEntries.length > 0 ? `
    <div class="hotspot-section">
      <p class="section-title">Hotspot Files</p>
      <div class="hotspot-list">
        ${fileEntries.map(([f, c]) => {
          const short = f.replace(/\\/g, '/').split('/').slice(-3).join('/');
          return `<div class="hotspot-row" data-file="${escAttr(f)}">
            <span class="hotspot-file">${escHtml(short)}</span>
            <span class="hotspot-count">${c}</span>
          </div>`;
        }).join('')}
      </div>
    </div>` : '';

  // New code section
  const newCodeHtml = (nc.new_issues != null || nc.closed_issues != null) ? `
    <div class="new-code-section">
      <span class="new-code-badge">New Code</span>
      <div class="new-code-metrics">
        <span><span class="ncm-val" style="color:${nc.new_issues > 0 ? 'var(--warning)' : 'var(--success)'}">${fmtNum(nc.new_issues || 0)}</span> new issues</span>
        <span><span class="ncm-val" style="color:var(--success)">${fmtNum(nc.closed_issues || 0)}</span> closed</span>
      </div>
    </div>` : '';

  // Scan info
  const scanInfoHtml = s ? `
    <p class="section-title">Latest Scan</p>
    <div class="scan-info">
      <div>
        <div class="info-label">Version</div>
        <div class="info-value">${escHtml(s.version || '\u2014')}</div>
      </div>
      <div>
        <div class="info-label">Branch</div>
        <div class="info-value">${escHtml(s.branch || '\u2014')}</div>
      </div>
      <div>
        <div class="info-label">Commit</div>
        <div class="info-value mono" style="font-size:12px">${s.commit_sha ? escHtml(s.commit_sha.slice(0,8)) : '\u2014'}</div>
      </div>
      <div>
        <div class="info-label">Status</div>
        <div class="info-value">${escHtml(s.status || '\u2014')}</div>
      </div>
      <div>
        <div class="info-label">Analysis date</div>
        <div class="info-value">${s.analysis_date ? new Date(s.analysis_date).toLocaleString() : '\u2014'}</div>
      </div>
      <div>
        <div class="info-label">Elapsed</div>
        <div class="info-value">${s.elapsed_ms ? (s.elapsed_ms / 1000).toFixed(1) + 's' : '\u2014'}</div>
      </div>
    </div>` : '';

  return gateHeroHtml + newCodeHtml + metricsHtml + sevDistHtml + typeDistHtml + hotspotHtml + scanInfoHtml;
}

function renderGateHero(gate) {
  if (!gate || !gate.status || gate.status === 'NONE') {
    return `<div class="gate-hero gate-loading">
      <div class="gate-badge">
        <span class="gate-icon">\u2014</span>
        <div class="gate-text">
          <span class="gate-label">Quality Gate</span>
          <span class="gate-status-text">Not configured</span>
        </div>
      </div>
    </div>`;
  }

  const s = gate.status;
  const cls = s === 'OK' ? 'gate-passed' : s === 'WARN' ? 'gate-warn' : 'gate-failed';
  const icon = s === 'OK' ? '\u2713' : s === 'WARN' ? '!' : '\u2717';
  const text = s === 'OK' ? 'Passed' : s === 'WARN' ? 'Warning' : 'Failed';

  const conds = (gate.conditions || []).map(c => {
    return `<div class="gate-cond">
      <span class="gate-cond-metric">${escHtml(c.metric)}</span>
      <span class="gate-cond-value">${escHtml(c.operator)} ${c.threshold}</span>
    </div>`;
  }).join('');

  return `<div class="gate-hero ${cls}">
    <div class="gate-badge">
      <span class="gate-icon">${icon}</span>
      <div class="gate-text">
        <span class="gate-label">Quality Gate</span>
        <span class="gate-status-text">${text}</span>
      </div>
    </div>
    ${conds ? `<div class="gate-conditions-list">${conds}</div>` : ''}
  </div>`;
}

function metricCard(label, value, colorCls, cardCls, typeFilter) {
  return `<div class="metric-card ${cardCls} clickable" data-mc-type="${typeFilter || ''}">
    <div class="metric-value ${colorCls}">${fmtNum(value)}</div>
    <div class="metric-label">${label}</div>
    <div class="metric-hint">View issues \u203A</div>
  </div>`;
}

function metricCardPct(label, value, cardCls) {
  return `<div class="metric-card ${cardCls}">
    <div class="metric-value">${value != null ? fmtPct(value) : '\u2014'}</div>
    <div class="metric-label">${label}</div>
  </div>`;
}

function metricCardK(label, value, cardCls) {
  return `<div class="metric-card ${cardCls}">
    <div class="metric-value">${fmtK(value)}</div>
    <div class="metric-label">${label}</div>
  </div>`;
}

function renderDistribution(title, order, labels, colors, data) {
  const total = order.reduce((s, k) => s + (data[k] || 0), 0);
  if (total === 0) return '';

  const rows = order.map(k => {
    const n = data[k] || 0;
    const pct = total > 0 ? (n / total * 100) : 0;
    return `<div class="dist-row">
      <span class="dist-label">${labels[k] || k}</span>
      <div class="dist-bar">
        <div class="dist-fill" style="width:${pct}%;background:${colors[k] || 'var(--accent)'}"></div>
      </div>
      <span class="dist-count">${fmtNum(n)}</span>
    </div>`;
  }).join('');

  return `<div class="dist-section">
    <p class="dist-title">${title} Distribution</p>
    <div class="dist-bar-wrap">${rows}</div>
  </div>`;
}

// ── Issues Tab ────────────────────────────────────────────────────────────────

async function loadIssues(append) {
  const p = state.currentProject;
  if (!p) return;
  if (!append) state.issueOffset = 0;

  state.loadingIssues = true;
  renderIssuesSection();

  const f      = state.issueFilter;
  const scanId  = state.currentScan?.id;
  let qs = scanId
    ? `scan_id=${scanId}&limit=${ISSUE_PAGE}&offset=${state.issueOffset}`
    : `project_id=${p.id}&limit=${ISSUE_PAGE}&offset=${state.issueOffset}`;
  if (f.severity !== 'all') qs += `&severity=${encodeURIComponent(f.severity)}`;
  if (f.type     !== 'all') qs += `&type=${encodeURIComponent(f.type)}`;
  if (f.status   !== 'all') qs += `&status=${encodeURIComponent(f.status)}`;
  if (f.search)             qs += `&file=${encodeURIComponent(f.search)}`;

  try {
    const data = await apiFetch('/issues?' + qs);
    if (append) {
      state.issues = state.issues.concat(data.items || []);
    } else {
      state.issues = data.items || [];
    }
    state.issuesTotal = data.total || 0;
  } catch {
    if (!append) state.issues = [];
  }

  state.loadingIssues = false;
  renderIssuesSection();
}

function renderIssuesSection() {
  const el = document.getElementById('issues-section');
  if (!el) return;
  el.innerHTML = buildIssuesHtml();
  bindIssueControls();
}

function sevCounts(issues) {
  const c = { blocker:0, critical:0, major:0, minor:0, info:0 };
  for (const i of issues) if (i.severity in c) c[i.severity]++;
  return c;
}

function buildIssuesHtml() {
  const issues = state.issues;
  const total  = state.issuesTotal;
  const f      = state.issueFilter;

  const counts = sevCounts(issues);
  const chips  = SEV_ORDER.map(sev => {
    const n      = counts[sev];
    const active = f.severity === sev;
    return `<button class="sev-chip${active?' active':''}" data-sev="${sev}" style="--chip-color:${SEV_COLOR[sev]};--chip-bg:${SEV_BG[sev]}">
      <span class="chip-dot" style="background:${SEV_COLOR[sev]}"></span>
      ${SEV_LABEL[sev]}
      <span class="chip-count">${n}</span>
    </button>`;
  }).join('');

  const summaryBar = `<div class="sev-bar">${chips}</div>`;

  const filtersHtml = `
    <div class="issues-toolbar">
      <span class="section-title" style="margin:0">Issues
        <span style="font-size:13px;font-weight:400;color:var(--text-muted)">&nbsp;${total.toLocaleString()} total</span>
      </span>
      <div class="issues-filters">
        <select id="filterType" class="filter-sel">
          <option value="all"${f.type==='all'?' selected':''}>All types</option>
          <option value="bug"${f.type==='bug'?' selected':''}>Bug</option>
          <option value="code_smell"${f.type==='code_smell'?' selected':''}>Code Smell</option>
          <option value="vulnerability"${f.type==='vulnerability'?' selected':''}>Vulnerability</option>
        </select>
        <select id="filterStatus" class="filter-sel">
          <option value="all"${f.status==='all'?' selected':''}>All statuses</option>
          <option value="open"${f.status==='open'?' selected':''}>Open</option>
          <option value="closed"${f.status==='closed'?' selected':''}>Closed</option>
        </select>
        <input id="filterSearch" class="filter-input" type="text" placeholder="Search file or message\u2026" value="${escAttr(f.search)}">
      </div>
    </div>`;

  if (state.loadingIssues && issues.length === 0) {
    return summaryBar + filtersHtml + `<div class="loading-state"><div class="spinner"></div></div>`;
  }

  if (issues.length === 0) {
    return summaryBar + filtersHtml + `<div class="empty-state" style="padding:32px 0"><p>No issues match the current filters.</p></div>`;
  }

  const rows = issues.map((i, idx) => {
    const color = SEV_COLOR[i.severity] || '#64748b';
    const bg    = SEV_BG[i.severity]   || 'transparent';
    const icon  = TYPE_ICON[i.type] || '?';
    const file  = (i.component_path || '').replace(/\\/g, '/').split('/').slice(-3).join('/');
    const loc   = i.end_line && i.end_line !== i.line ? `${i.line}\u2013${i.end_line}` : `${i.line}`;
    const status = i.status || 'open';
    const isClosed = status === 'closed';
    let actionBtns = '';
    if (!isClosed) {
      actionBtns = `
        <button class="itbtn fp-btn" data-id="${i.id}" data-res="false_positive" title="False positive">FP</button>
        <button class="itbtn wf-btn" data-id="${i.id}" data-res="wont_fix" title="Won\u2019t fix">WF</button>
        <button class="itbtn ok-btn" data-id="${i.id}" data-res="fixed" title="Mark as fixed">\u2713</button>`;
    } else {
      actionBtns = `<button class="itbtn re-btn" data-id="${i.id}" data-res="" title="Reopen">\u21A9</button>`;
    }
    return `<tr style="--row-sev-color:${color};--row-sev-bg:${bg}" class="sev-row${isClosed?' row-closed':''}" data-issue-idx="${idx}">
      <td><span class="sev-badge" style="background:${color}">${escHtml(i.severity)}</span></td>
      <td>${icon} ${escHtml((i.type||'').replace('_',' '))}</td>
      <td class="mono" style="font-size:11px">${escHtml(i.rule_key||'')}</td>
      <td class="file-cell" title="${escAttr(i.component_path||'')}"><span class="mono">${escHtml(file)}<span style="color:var(--text-muted)">:${loc}</span></span></td>
      <td>${escHtml(i.message||'')}</td>
      <td class="actions-cell" onclick="event.stopPropagation()">${actionBtns}</td>
    </tr>`;
  }).join('');

  const hasMore = issues.length < total;
  const moreBtn = hasMore
    ? `<div style="text-align:center;padding:16px">
        <button class="btn btn-primary" id="loadMoreBtn" style="width:auto;padding:8px 24px">
          ${state.loadingIssues ? 'Loading\u2026' : `Load more (${total - issues.length} remaining)`}
        </button>
       </div>`
    : '';

  return summaryBar + filtersHtml + `
    <div class="issues-table-wrap">
      <table class="issues-table">
        <thead><tr>
          <th>Severity</th><th>Type</th><th>Rule</th><th>File</th><th>Message</th><th></th>
        </tr></thead>
        <tbody>${rows}</tbody>
      </table>
    </div>
    ${moreBtn}`;
}

function bindIssueControls() {
  document.querySelectorAll('.sev-chip').forEach(btn => {
    btn.addEventListener('click', () => {
      const sev = btn.dataset.sev;
      state.issueFilter.severity = state.issueFilter.severity === sev ? 'all' : sev;
      loadIssues();
    });
  });
  document.getElementById('filterType')?.addEventListener('change', e => {
    state.issueFilter.type = e.target.value;
    loadIssues();
  });
  document.getElementById('filterStatus')?.addEventListener('change', e => {
    state.issueFilter.status = e.target.value;
    loadIssues();
  });
  let searchTimer;
  document.getElementById('filterSearch')?.addEventListener('input', e => {
    clearTimeout(searchTimer);
    searchTimer = setTimeout(() => {
      state.issueFilter.search = e.target.value.trim();
      loadIssues();
    }, 300);
  });
  document.getElementById('loadMoreBtn')?.addEventListener('click', () => {
    state.issueOffset += ISSUE_PAGE;
    loadIssues(true);
  });

  // Row click → detail panel
  document.querySelectorAll('.issues-table tbody tr[data-issue-idx]').forEach(row => {
    row.addEventListener('click', () => {
      const idx = parseInt(row.dataset.issueIdx, 10);
      if (state.issues[idx]) openIssueDetail(state.issues[idx]);
    });
  });

  // Issue transition buttons
  document.querySelectorAll('.itbtn').forEach(btn => {
    btn.addEventListener('click', async e => {
      e.stopPropagation();
      const id  = btn.dataset.id;
      const res = btn.dataset.res;
      btn.disabled = true;
      try {
        await apiFetch('/issues/' + id + '/transition', {
          method: 'POST',
          body: JSON.stringify({ resolution: res, comment: '' }),
        });
        const idx = state.issues.findIndex(i => String(i.id) === String(id));
        if (idx !== -1) {
          const iss = state.issues[idx];
          if (res === '') {
            iss.status = 'open'; iss.resolution = '';
          } else {
            iss.status = 'closed'; iss.resolution = res;
          }
        }
        renderIssuesSection();
      } catch (err) {
        showToast(err.message, 'error');
        btn.disabled = false;
      }
    });
  });
}

// ── Issue Detail Panel ────────────────────────────────────────────────────────

function openIssueDetail(issue) {
  state.selectedIssue = issue;
  const inner = document.getElementById('detail-inner');
  const panel = document.getElementById('detail-panel');
  const overlay = document.getElementById('detail-overlay');
  if (!inner || !panel || !overlay) return;

  const i = issue;
  const file = (i.component_path || '').replace(/\\/g, '/');
  const loc = i.end_line && i.end_line !== i.line ? `${i.line}\u2013${i.end_line}` : `${i.line}`;
  const sevColor = SEV_COLOR[i.severity] || '#64748b';

  let secondaryHtml = '';
  if (i.secondary_locations && i.secondary_locations.length > 0) {
    const locs = i.secondary_locations.map((sl, idx) => {
      const sf = (sl.file_path || sl.component_path || '').replace(/\\/g, '/').split('/').slice(-2).join('/');
      return `<div class="secondary-loc">
        <span class="loc-num">${idx + 1}</span>
        <span class="loc-file">${escHtml(sf)}:${sl.line || ''}</span>
        <span class="loc-msg">${escHtml(sl.message || '')}</span>
      </div>`;
    }).join('');
    secondaryHtml = `
      <div class="detail-section-title">Secondary Locations</div>
      ${locs}`;
  }

  inner.innerHTML = `
    <button class="detail-close" id="detailClose">\u2715</button>
    <div class="detail-title">${escHtml(i.message || 'Issue')}</div>
    <div class="detail-props">
      <div class="detail-prop">
        <span class="detail-prop-label">Severity</span>
        <span class="detail-prop-value"><span class="sev-badge" style="background:${sevColor}">${escHtml(i.severity)}</span></span>
      </div>
      <div class="detail-prop">
        <span class="detail-prop-label">Type</span>
        <span class="detail-prop-value">${TYPE_ICON[i.type] || ''} ${escHtml(TYPE_LABEL[i.type] || i.type || '')}</span>
      </div>
      <div class="detail-prop">
        <span class="detail-prop-label">Rule</span>
        <span class="detail-prop-value mono" style="font-size:12px">${escHtml(i.rule_key || '')}</span>
      </div>
      <div class="detail-prop">
        <span class="detail-prop-label">Engine</span>
        <span class="detail-prop-value">${escHtml(i.engine_id || '\u2014')}</span>
      </div>
      <div class="detail-prop">
        <span class="detail-prop-label">File</span>
        <span class="detail-prop-value mono" style="font-size:12px;word-break:break-all">${escHtml(file)}</span>
      </div>
      <div class="detail-prop">
        <span class="detail-prop-label">Location</span>
        <span class="detail-prop-value mono">Line ${loc}${i.column ? ', Col ' + i.column : ''}</span>
      </div>
      <div class="detail-prop">
        <span class="detail-prop-label">Status</span>
        <span class="detail-prop-value">${escHtml(i.status || 'open')}${i.resolution ? ' \u2014 ' + escHtml(i.resolution) : ''}</span>
      </div>
      ${i.tags && i.tags.length ? `<div class="detail-prop">
        <span class="detail-prop-label">Tags</span>
        <span class="detail-prop-value">${i.tags.map(t => `<span class="tag">${escHtml(t)}</span>`).join(' ')}</span>
      </div>` : ''}
    </div>
    ${secondaryHtml}`;

  panel.classList.remove('hidden');
  overlay.classList.remove('hidden');
  requestAnimationFrame(() => {
    panel.classList.add('open');
    overlay.classList.add('open');
  });

  document.getElementById('detailClose').addEventListener('click', closeIssueDetail);
  overlay.addEventListener('click', closeIssueDetail);
}

function closeIssueDetail() {
  const panel = document.getElementById('detail-panel');
  const overlay = document.getElementById('detail-overlay');
  if (!panel || !overlay) return;
  panel.classList.remove('open');
  overlay.classList.remove('open');
  setTimeout(() => {
    panel.classList.add('hidden');
    overlay.classList.add('hidden');
  }, 250);
  state.selectedIssue = null;
}

// ── Activity Tab ──────────────────────────────────────────────────────────────

function renderActivityTab() {
  const data = state.activityData;
  if (data === null) return `<div class="loading-state"><div class="spinner"></div></div>`;
  if (!data || !data.length) return `<div class="empty-state" style="padding:40px 0"><p>No scan activity yet.</p></div>`;

  const entries = data.map((entry, idx) => {
    const isLast = idx === data.length - 1;
    const dotCls = entry.gate_status === 'OK' ? 'dot-ok' : entry.gate_status === 'ERROR' ? 'dot-error' : entry.gate_status === 'WARN' ? 'dot-warn' : '';

    const eventsHtml = (entry.events || []).map(ev => {
      let cls = 'ev-version';
      if (ev.category === 'QUALITY_GATE') cls = ev.value === 'OK' ? 'ev-gate' : 'ev-gate-fail';
      if (ev.category === 'ISSUE_SPIKE') cls = 'ev-spike';
      if (ev.category === 'FIRST_ANALYSIS') cls = 'ev-first';
      return `<span class="activity-event ${cls}">${escHtml(ev.name)}</span>`;
    }).join('');

    return `<div class="activity-entry">
      <div class="activity-dot-col">
        <div class="activity-dot ${dotCls}"></div>
        ${!isLast ? '<div class="activity-line"></div>' : ''}
      </div>
      <div class="activity-body">
        <div class="activity-date">${entry.analysis_date ? new Date(entry.analysis_date).toLocaleString() : '\u2014'}</div>
        <div class="activity-main">
          <span class="activity-version">${escHtml(entry.version || 'No version')}</span>
          ${entry.branch ? `<span class="activity-branch">${escHtml(entry.branch)}</span>` : ''}
          ${entry.gate_status ? `<span class="badge ${entry.gate_status === 'OK' ? 'badge-ok' : entry.gate_status === 'WARN' ? 'badge-warn' : 'badge-error'}">${escHtml(entry.gate_status)}</span>` : ''}
        </div>
        <div class="activity-metrics">
          <span><span class="am-val">${fmtNum(entry.total_issues)}</span> issues</span>
          <span><span class="am-val" style="color:${entry.new_issues > 0 ? 'var(--warning)' : 'var(--text)'}">${fmtNum(entry.new_issues)}</span> new</span>
          <span><span class="am-val" style="color:var(--success)">${fmtNum(entry.closed_issues)}</span> closed</span>
        </div>
        ${eventsHtml ? `<div class="activity-events">${eventsHtml}</div>` : ''}
      </div>
    </div>`;
  }).join('');

  return `<div class="activity-timeline">${entries}</div>`;
}

async function loadActivityData() {
  const p = state.currentProject;
  if (!p) return;
  try {
    const data = await apiFetch('/projects/' + encodeURIComponent(p.key) + '/activity?limit=30');
    state.activityData = data.items || [];
  } catch { state.activityData = []; }
  render();
}

// ── Tab switching ─────────────────────────────────────────────────────────────

async function switchTab(tab) {
  state.projectTab = tab;
  render();
  if (tab === 'issues') {
    if (state.issues.length === 0 && !state.loadingIssues) {
      await loadIssues();
    } else {
      renderIssuesSection();
    }
    return;
  }
  if (tab === 'activity' && state.activityData === null) {
    await loadActivityData(); return;
  }
  if (tab === 'gate' && state.gateData === null) {
    await loadGateData(); return;
  }
  if (tab === 'webhooks' && state.webhooksData === null) {
    await loadWebhooksData(); return;
  }
  if (tab === 'profiles' && state.profilesData === null) {
    await loadProfilesData(); return;
  }
  bindTabContent();
}

// ── Gate tab ──────────────────────────────────────────────────────────────────

async function loadGateData() {
  try {
    const data = await apiFetch('/quality-gates');
    state.gateData = data.items || (Array.isArray(data) ? data : []);
  } catch { state.gateData = []; }
  render();
  bindTabContent();
}

function renderGateTab() {
  const gates = state.gateData;
  if (gates === null) return `<div class="loading-state"><div class="spinner"></div></div>`;
  if (!gates.length) return `<div class="empty-state" style="padding:40px 0"><p>No quality gates configured.</p></div>`;

  const rows = gates.map(g => `
    <div class="gate-card">
      <div class="gate-header">
        <div>
          <span class="gate-name">${escHtml(g.name)}</span>
          ${g.is_default ? `<span class="badge badge-ok" style="font-size:11px;margin-left:8px">Default</span>` : ''}
        </div>
        <div class="gate-actions">
          <button class="btn-sm btn-outline assign-gate-btn" data-gate-id="${g.id}" data-gate-name="${escAttr(g.name)}">Assign to project</button>
          <button class="btn-sm btn-ghost expand-gate-btn" data-gate-id="${g.id}">Conditions \u25BE</button>
        </div>
      </div>
      <div class="gate-conditions hidden" id="gate-cond-${g.id}">
        <div class="loading-inline">Loading\u2026</div>
      </div>
    </div>`).join('');

  return `<div class="tab-section">
    <p class="section-title" style="margin-top:24px">Quality Gates</p>
    <p style="color:var(--text-muted);font-size:13px;margin-bottom:16px">Conditions that must pass for a project analysis to be considered successful.</p>
    <div class="gate-list">${rows}</div>
  </div>`;
}

// ── Webhooks tab ──────────────────────────────────────────────────────────────

async function loadWebhooksData() {
  const p = state.currentProject;
  try {
    const data = await apiFetch('/webhooks' + (p ? '?project_key=' + encodeURIComponent(p.key) : ''));
    state.webhooksData = data.items || (Array.isArray(data) ? data : []);
  } catch { state.webhooksData = []; }
  try {
    state.newCodePeriod = await apiFetch('/projects/' + encodeURIComponent(p.key) + '/new-code-period');
  } catch { state.newCodePeriod = null; }
  render();
  bindTabContent();
}

function renderWebhooksTab() {
  const whs = state.webhooksData;
  if (whs === null) return `<div class="loading-state"><div class="spinner"></div></div>`;

  const ncp = state.newCodePeriod;
  const ncpStr = ncp
    ? escHtml(ncp.strategy) + (ncp.value ? ' \u2014 ' + escHtml(ncp.value) : '')
    : 'auto (default)';

  const whRows = whs.length === 0
    ? `<div class="empty-state" style="padding:20px 0"><p>No webhooks configured.</p></div>`
    : whs.map(w => `
      <div class="webhook-row">
        <div class="webhook-info">
          <span class="webhook-name">${escHtml(w.name)}</span>
          <span class="webhook-url" title="${escAttr(w.url)}">${escHtml(w.url)}</span>
        </div>
        <div class="webhook-btns">
          <button class="btn-sm btn-outline test-wh-btn" data-wh-id="${w.id}">Test</button>
          <button class="btn-sm btn-danger del-wh-btn" data-wh-id="${w.id}">Delete</button>
        </div>
      </div>`).join('');

  return `<div class="tab-section">
    <p class="section-title" style="margin-top:24px">Webhooks</p>
    <div class="webhook-list">${whRows}</div>
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
          <option value="auto"${(!ncp||ncp.strategy==='auto')?' selected':''}>Auto</option>
          <option value="previous_version"${ncp?.strategy==='previous_version'?' selected':''}>Previous version</option>
          <option value="number_of_days"${ncp?.strategy==='number_of_days'?' selected':''}>Number of days</option>
          <option value="reference_branch"${ncp?.strategy==='reference_branch'?' selected':''}>Reference branch</option>
        </select>
        <input id="ncpValue" class="filter-input" placeholder="Value (if needed)" style="width:140px" value="${escAttr(ncp?.value||'')}">
        <button class="btn btn-primary" id="saveNcpBtn" style="width:auto;padding:6px 18px;margin-top:0">Save</button>
      </div>
    </div>
  </div>`;
}

// ── Profiles tab ──────────────────────────────────────────────────────────────

async function loadProfilesData() {
  try {
    const data = await apiFetch('/profiles');
    state.profilesData = data.items || (Array.isArray(data) ? data : []);
  } catch { state.profilesData = []; }
  render();
  bindTabContent();
}

function renderProfilesTab() {
  const profiles = state.profilesData;
  if (profiles === null) return `<div class="loading-state"><div class="spinner"></div></div>`;
  if (!profiles.length) return `<div class="empty-state" style="padding:40px 0"><p>No quality profiles found.</p></div>`;

  const byLang = {};
  for (const pr of profiles) {
    if (!byLang[pr.language]) byLang[pr.language] = [];
    byLang[pr.language].push(pr);
  }

  const sections = Object.entries(byLang).map(([lang, profs]) => `
    <div class="profile-lang-section">
      <h4 class="profile-lang-title">${escHtml(lang)}</h4>
      <div class="profile-list">
        ${profs.map(pr => `
          <div class="profile-row">
            <div class="profile-info">
              <span class="profile-name">${escHtml(pr.name)}</span>
              ${pr.is_builtin ? `<span class="badge badge-ok" style="font-size:10px;margin-left:6px">Built-in</span>` : ''}
              ${pr.is_default ? `<span class="badge badge-warn" style="font-size:10px;margin-left:6px">Default</span>` : ''}
              <span style="color:var(--text-muted);font-size:12px;margin-left:8px">${pr.rule_count||0} rules</span>
            </div>
            <button class="btn-sm btn-outline assign-profile-btn"
              data-profile-id="${pr.id}"
              data-profile-lang="${escAttr(pr.language)}"
              data-profile-name="${escAttr(pr.name)}">Assign to project</button>
          </div>`).join('')}
      </div>
    </div>`).join('');

  return `<div class="tab-section">
    <p class="section-title" style="margin-top:24px">Quality Profiles</p>
    <p style="color:var(--text-muted);font-size:13px;margin-bottom:16px">Profiles define which rules are active for each language.</p>
    ${sections}
  </div>`;
}

// ── Tab content event binding ─────────────────────────────────────────────────

function bindTabContent() {
  const p = state.currentProject;
  if (!p) return;

  // Gate tab
  document.querySelectorAll('.expand-gate-btn').forEach(btn => {
    btn.addEventListener('click', async () => {
      const id = btn.dataset.gateId;
      const box = document.getElementById('gate-cond-' + id);
      if (!box) return;
      const hidden = box.classList.toggle('hidden');
      btn.textContent = hidden ? 'Conditions \u25BE' : 'Conditions \u25B4';
      if (!hidden && box.innerHTML.includes('Loading')) {
        try {
          const gate = await apiFetch('/quality-gates/' + id);
          const conds = gate.conditions || [];
          if (!conds.length) {
            box.innerHTML = '<p style="color:var(--text-muted);padding:8px 0;font-size:13px">No conditions defined.</p>';
          } else {
            box.innerHTML = `<table class="conditions-table">
              <thead><tr><th>Metric</th><th>Operator</th><th>Threshold</th><th>New Code Only</th></tr></thead>
              <tbody>${conds.map(c => `<tr>
                <td>${escHtml(c.metric)}</td>
                <td>${escHtml(c.operator)}</td>
                <td class="mono">${escHtml(String(c.value))}</td>
                <td>${c.on_new_code ? '\u2713' : ''}</td>
              </tr>`).join('')}</tbody>
            </table>`;
          }
        } catch { box.innerHTML = '<p style="color:var(--danger);font-size:13px">Failed to load conditions.</p>'; }
      }
    });
  });

  document.querySelectorAll('.assign-gate-btn').forEach(btn => {
    btn.addEventListener('click', async () => {
      const id = btn.dataset.gateId;
      const name = btn.dataset.gateName;
      btn.disabled = true;
      try {
        await apiFetch('/projects/' + encodeURIComponent(p.key) + '/quality-gate', {
          method: 'POST',
          body: JSON.stringify({ gate_id: parseInt(id, 10) }),
        });
        showToast('Gate "' + name + '" assigned.');
      } catch (err) { showToast(err.message, 'error'); }
      btn.disabled = false;
    });
  });

  // Webhooks tab
  document.getElementById('addWhBtn')?.addEventListener('click', async () => {
    const name   = document.getElementById('newWhName')?.value.trim();
    const url    = document.getElementById('newWhUrl')?.value.trim();
    const secret = document.getElementById('newWhSecret')?.value.trim();
    if (!name || !url) { showToast('Name and URL are required.', 'error'); return; }
    try {
      await apiFetch('/webhooks', {
        method: 'POST',
        body: JSON.stringify({ name, url, secret: secret || '', project_key: p.key }),
      });
      state.webhooksData = null;
      await loadWebhooksData();
    } catch (err) { showToast(err.message, 'error'); }
  });

  document.querySelectorAll('.test-wh-btn').forEach(btn => {
    btn.addEventListener('click', async () => {
      btn.disabled = true;
      try {
        await apiFetch('/webhooks/' + btn.dataset.whId + '/test', { method: 'POST' });
        showToast('Test delivery sent.');
      } catch (err) { showToast(err.message, 'error'); }
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
      } catch (err) { showToast(err.message, 'error'); btn.disabled = false; }
    });
  });

  document.getElementById('saveNcpBtn')?.addEventListener('click', async () => {
    const strategy = document.getElementById('ncpStrategy')?.value;
    const value    = document.getElementById('ncpValue')?.value.trim();
    try {
      await apiFetch('/projects/' + encodeURIComponent(p.key) + '/new-code-period', {
        method: 'PUT',
        body: JSON.stringify({ strategy, value: value || '' }),
      });
      state.newCodePeriod = { strategy, value };
      const display = document.getElementById('ncpDisplay');
      if (display) display.textContent = strategy + (value ? ' \u2014 ' + value : '');
      showToast('New code period saved.');
    } catch (err) { showToast(err.message, 'error'); }
  });

  // Profiles tab
  document.querySelectorAll('.assign-profile-btn').forEach(btn => {
    btn.addEventListener('click', async () => {
      const id   = btn.dataset.profileId;
      const lang = btn.dataset.profileLang;
      const name = btn.dataset.profileName;
      btn.disabled = true;
      try {
        await apiFetch('/projects/' + encodeURIComponent(p.key) + '/profiles', {
          method: 'POST',
          body: JSON.stringify({ profile_id: parseInt(id, 10), language: lang }),
        });
        showToast('Profile "' + name + '" assigned.');
      } catch (err) { showToast(err.message, 'error'); }
      btn.disabled = false;
    });
  });
}

// ── Toast notifications ───────────────────────────────────────────────────────

function showToast(msg, type) {
  type = type || 'success';
  const el = document.createElement('div');
  el.className = 'toast toast-' + type;
  el.textContent = msg;
  document.body.appendChild(el);
  setTimeout(() => el.classList.add('toast-show'), 10);
  setTimeout(() => {
    el.classList.remove('toast-show');
    setTimeout(() => el.remove(), 300);
  }, 3500);
}

// ── Event binding ─────────────────────────────────────────────────────────────

function bindMain() {
  document.getElementById('logoutBtn')?.addEventListener('click', logout);
  document.getElementById('backBtn')?.addEventListener('click', () => loadProjects());
  document.querySelectorAll('.project-card').forEach(card => {
    card.addEventListener('click', () => loadProject(card.dataset.key));
  });
  if (state.view === 'project') {
    document.querySelectorAll('.tab-btn').forEach(btn => {
      btn.addEventListener('click', () => switchTab(btn.dataset.tab));
    });

    // Overview metric cards → switch to issues tab with filter
    document.querySelectorAll('.metric-card.clickable').forEach(btn => {
      btn.addEventListener('click', () => {
        const type = btn.dataset.mcType;
        if (type) {
          state.issueFilter.type = type === 'all' ? 'all' : type;
          state.issueFilter.severity = 'all';
          state.issueFilter.status = 'all';
          state.issueFilter.search = '';
          state.issues = [];
          switchTab('issues');
        }
      });
    });

    // Hotspot file click → switch to issues tab filtered by file
    document.querySelectorAll('.hotspot-row').forEach(row => {
      row.addEventListener('click', () => {
        const file = row.dataset.file;
        if (file) {
          state.issueFilter.search = file;
          state.issueFilter.type = 'all';
          state.issueFilter.severity = 'all';
          state.issueFilter.status = 'all';
          state.issues = [];
          switchTab('issues');
        }
      });
    });

    if (state.projectTab !== 'overview' && state.projectTab !== 'issues') bindTabContent();
  }
}

function logout() {
  clearStorage();
  state = {
    user: null, view: 'login', projects: [], currentProject: null, currentScan: null,
    overviewData: null, issues: [], issuesTotal: 0, issueOffset: 0,
    issueFilter: { severity: 'all', type: 'all', status: 'all', search: '' },
    loading: false, loadingIssues: false,
    projectTab: 'overview', gateData: null, webhooksData: null, profilesData: null,
    activityData: null, newCodePeriod: null, selectedIssue: null,
  };
  render();
}

// ── Keyboard shortcuts ────────────────────────────────────────────────────────

document.addEventListener('keydown', e => {
  if (e.key === 'Escape') {
    if (state.selectedIssue) closeIssueDetail();
  }
});

// ── Security helpers ──────────────────────────────────────────────────────────

function escHtml(s) {
  if (s == null) return '';
  return String(s)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;');
}

function escAttr(s) { return escHtml(s); }

// ── Boot ──────────────────────────────────────────────────────────────────────

async function init() {
  const t = getToken();
  if (t) {
    state.user = loadUser();
    await loadProjects();
  } else {
    render();
  }
}

document.addEventListener('DOMContentLoaded', init);
