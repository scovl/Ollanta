'use strict';

import { apiFetch, setUnauthorizedHandler } from './core/api.js';
import { parseProjectRoute } from './core/scope.js';
import { createInitialState, replaceState, resetProjectState, state } from './core/state.js';
import { clearStorage, getToken, loadUser, saveToken, saveUser } from './core/storage.js';
import { badgeClassForGateStatus, cardClassForGateStatus, escAttr, escHtml, fmtDate } from './core/utils.js';
import { closeIssueDetail } from './features/issues.js';
import { bindProjectViewControls, loadProject, renderProjectDetail } from './project-flow.js';

const BRAND_MARK_PATH = '/branding/ollanta-mark.png';

setUnauthorizedHandler(logout);

export function render() {
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
  const user = state.user || {};
  const name = user.name || user.login || 'User';
  return `<nav>
    ${renderBrandLockup()}
    <span class="spacer"></span>
    <span class="user-info">${escHtml(name)}</span>
    <button class="logout-btn" id="logoutBtn">Sign out</button>
  </nav>`;
}

function renderBrandLockup() {
  return `<span class="brand-lockup brand-lockup-inline">
    <img class="brand-mark" src="${BRAND_MARK_PATH}" alt="Ollanta" width="36" height="36">
  </span>`;
}



function renderContent() {
  if (state.view === 'projects') return renderDashboard();
  if (state.view === 'project') return renderProjectDetail();
  return '';
}

function renderLogin() {
  return `<div class="login-wrapper">
    <div class="login-card-unified">
      <div class="login-card-header">
        <img class="login-brand-mark" src="${BRAND_MARK_PATH}" alt="" aria-hidden="true" width="80" height="80">
        <div class="login-brand-text">
          <span class="login-brand-name">Ollanta</span>
          <span class="login-brand-tagline">Static Analysis Platform</span>
        </div>
        <div class="login-features">
          <div class="login-feature">
            <span class="login-feature-icon">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
            </span>
            <span>Quality gates &amp; metrics across scans</span>
          </div>
          <div class="login-feature">
            <span class="login-feature-icon">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="6" y1="3" x2="6" y2="15"/><circle cx="18" cy="6" r="3"/><circle cx="6" cy="18" r="3"/><path d="M18 9a9 9 0 0 1-9 9"/></svg>
            </span>
            <span>Multi-project, branch &amp; PR tracking</span>
          </div>
          <div class="login-feature">
            <span class="login-feature-icon">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="16 18 22 12 16 6"/><polyline points="8 6 2 12 8 18"/></svg>
            </span>
            <span>Code browsing with inline issue markers</span>
          </div>
        </div>
      </div>
      <div class="login-card-content">
        <p class="login-eyebrow">Welcome back</p>
        <h1 class="login-title">Sign in</h1>
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
    </div>
  </div>`;
}

function bindLogin() {
  const btn = document.getElementById('loginBtn');
  const errEl = document.getElementById('loginError');
  const userEl = document.getElementById('loginUser');
  const passEl = document.getElementById('loginPass');

  async function doLogin() {
    const login = userEl.value.trim();
    const password = passEl.value;
    if (!login || !password) {
      errEl.textContent = 'Enter username and password.';
      return;
    }

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
    } catch (err) {
      errEl.textContent = err.message || 'Login failed.';
      btn.disabled = false;
      btn.textContent = 'Sign in';
    }
  }

  btn.addEventListener('click', doLogin);
  passEl.addEventListener('keydown', event => { if (event.key === 'Enter') doLogin(); });
  userEl.addEventListener('keydown', event => { if (event.key === 'Enter') passEl.focus(); });
}

export async function loadProjects() {
  resetProjectState();
  state.view = 'projects';
  state.loading = true;
  history.replaceState({}, '', globalThis.location.pathname);
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

  const projects = state.projects;
  const count = projects.length;

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
      : `<div class="projects-grid">${projects.map(renderProjectCard).join('')}</div>`
    }`;
}

function renderProjectCard(project) {
  const tags = (project.tags || []).filter(Boolean);
  const tagsHtml = tags.length
    ? '<div class="tags">' + tags.map(tag => '<span class="tag">' + escHtml(tag) + '</span>').join('') + '</div>'
    : '';

  const gateStatus = project.gate_status || '';
  const gateCls = cardClassForGateStatus(gateStatus);
  const gateBadge = gateStatus ? `<span class="badge ${badgeClassForGateStatus(gateStatus)}">${escHtml(gateStatus)}</span>` : '';

  return `<div class="project-card ${gateCls}" data-key="${escAttr(project.key)}">
    <div class="card-top">
      <span class="key">${escHtml(project.key)}</span>
      ${gateBadge}
    </div>
    <div class="name" title="${escAttr(project.name || project.key)}">${escHtml(project.name || project.key)}</div>
    ${tagsHtml}
    <div class="footer">Updated ${fmtDate(project.updated_at)}</div>
  </div>`;
}

export function showToast(msg, type = 'success') {
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

function bindMain() {
  document.getElementById('logoutBtn')?.addEventListener('click', logout);
  document.getElementById('backBtn')?.addEventListener('click', () => loadProjects());
  document.querySelectorAll('.project-card').forEach(card => {
    card.addEventListener('click', () => loadProject(card.dataset.key, { project: card.dataset.key, tab: 'overview', branch: '', pullRequest: '' }));
  });
  if (state.view === 'project') {
    bindProjectViewControls();
  }
}

export function logout() {
  clearStorage();
  history.replaceState({}, '', globalThis.location.pathname);
  replaceState(createInitialState());
  render();
}

function handleGlobalKeydown(event) {
  if (event.key === 'Escape' && state.selectedIssue) {
    closeIssueDetail();
  }
}

export async function init() {
  globalThis.addEventListener('popstate', async () => {
    if (!getToken()) return;
    const route = parseProjectRoute();
    if (route.project) {
      await loadProject(route.project, route);
      return;
    }
    await loadProjects();
  });

  const token = getToken();
  if (token) {
    state.user = loadUser();
    const route = parseProjectRoute();
    if (route.project) {
      await loadProject(route.project, route);
    } else {
      await loadProjects();
    }
  } else {
    render();
  }
}

let bootstrapped = false;

export function bootBrowserApp() {
  if (bootstrapped) return;
  bootstrapped = true;
  document.addEventListener('keydown', handleGlobalKeydown);
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
    return;
  }
  void init();
}