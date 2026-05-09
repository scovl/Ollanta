import { apiFetch } from '../../core/api.js';
import { state } from '../../core/state.js';
import { escAttr, escHtml, fmtNum } from '../../core/utils.js';
import { getRenderView, getShowToast } from './context.js';

export async function loadProfilesData() {
  const project = state.currentProject;
  try {
    const requests = [apiFetch('/profiles'), apiFetch('/rules')];
    if (project) {
      const projectKey = encodeURIComponent(project.key);
      requests.push(
        apiFetch('/projects/' + projectKey + '/profiles'),
        apiFetch('/projects/' + projectKey + '/profiles/effective'),
      );
    }
    const [profilesData, rulesData, assignmentsData, effectiveData] = await Promise.all(requests);
    state.profilesData = {
      profiles: profilesData.items || (Array.isArray(profilesData) ? profilesData : []),
      rules: rulesData.items || (Array.isArray(rulesData) ? rulesData : []),
      assignments: assignmentsData?.items || (Array.isArray(assignmentsData) ? assignmentsData : []),
      effective: effectiveData?.items || (Array.isArray(effectiveData) ? effectiveData : []),
    };
  } catch {
    state.profilesData = { profiles: [], rules: [], assignments: [], effective: [] };
  }
  getRenderView()();
}

export function renderProfilesTab() {
  const data = state.profilesData;
  if (data === null) return `<div class="loading-state"><div class="spinner"></div></div>`;
  const profiles = Array.isArray(data) ? data : (data.profiles || []);
  const assignments = Array.isArray(data) ? [] : (data.assignments || []);
  const effective = Array.isArray(data) ? [] : (data.effective || []);
  const rules = Array.isArray(data) ? [] : (data.rules || []);
  if (!profiles.length) return `<div class="empty-state" style="padding:40px 0"><p>No quality profiles found.</p></div>`;

  const byLang = {};
  for (const profile of profiles) {
    if (!byLang[profile.language]) byLang[profile.language] = [];
    byLang[profile.language].push(profile);
  }

  const assignmentByLang = Object.fromEntries(assignments.map(item => [item.language, item]));
  const effectiveByLang = Object.fromEntries(effective.map(item => [item.language, item]));
  const ruleByKey = Object.fromEntries(rules.map(rule => [rule.key, rule]));

  const sections = Object.entries(byLang).map(([lang, profs]) => `
    <div class="profile-lang-section">
      <h4 class="profile-lang-title">${escHtml(lang)}</h4>
      ${renderEffectiveProfileSummary(lang, assignmentByLang[lang], effectiveByLang[lang], ruleByKey)}
      <div class="profile-list">
        ${profs.map(profile => `
          <div class="profile-row">
            <div class="profile-info">
              <span class="profile-name">${escHtml(profile.name)}</span>
              ${profile.is_builtin ? `<span class="badge badge-ok" style="font-size:10px;margin-left:6px">Built-in</span>` : ''}
              ${profile.is_default ? `<span class="badge badge-warn" style="font-size:10px;margin-left:6px">Default</span>` : ''}
              ${profile.parser_only ? `<span class="badge" style="font-size:10px;margin-left:6px">Parser only</span>` : ''}
              <span style="color:var(--text-muted);font-size:12px;margin-left:8px">${profile.rule_count || 0} active rules</span>
            </div>
            <div class="profile-actions">
              <button class="btn-sm btn-outline export-profile-btn" data-profile-id="${profile.id}" data-profile-name="${escAttr(profile.name)}">Export</button>
              <button class="btn-sm btn-outline assign-profile-btn"
                data-profile-id="${profile.id}"
                data-profile-lang="${escAttr(profile.language)}"
                data-profile-name="${escAttr(profile.name)}">Assign to project</button>
            </div>
          </div>`).join('')}
      </div>
    </div>`).join('');

  return `<div class="tab-section">
    <p class="section-title" style="margin-top:24px">Quality Profiles</p>
    <p style="color:var(--text-muted);font-size:13px;margin-bottom:16px">Profiles define which rules are active for each language.</p>
    ${sections}
  </div>`;
}

function renderEffectiveProfileSummary(lang, assignment, effective, ruleByKey) {
  const profile = assignment?.profile;
  const source = effective?.source || assignment?.source || 'default';
  const activeRules = effective?.rules || [];
  const parserOnly = effective?.parser_only || profile?.parser_only;
  const hash = effective?.rules_hash ? effective.rules_hash.slice(0, 10) : '-';
  const profileName = effective?.profile_name || profile?.name || 'No profile';
  const sourceBadge = source === 'assigned' ? 'badge-ok' : 'badge-warn';
  const topRules = activeRules.slice(0, 8);
  const remainingRules = activeRules.length - topRules.length;
  const remainingRow = remainingRules > 0 ? `<div class="profile-rule-row muted">${fmtNum(remainingRules)} more rules</div>` : '';
  const emptyText = parserOnly ? 'Parser available; no bundled rules yet.' : 'No active rules.';
  const rulesBlock = topRules.length
    ? `<div class="profile-rule-table">
      ${topRules.map(rule => renderEffectiveRule(rule, ruleByKey)).join('')}
      ${remainingRow}
    </div>`
    : `<div class="profile-rule-empty">${emptyText}</div>`;
  return `<div class="profile-effective-box">
    <div class="profile-effective-head">
      <div>
        <span class="profile-active-label">Active profile</span>
        <strong>${escHtml(profileName)}</strong>
        <span class="badge ${sourceBadge}" style="font-size:10px;margin-left:6px">${escHtml(source)}</span>
        ${parserOnly ? `<span class="badge" style="font-size:10px;margin-left:6px">Parser only</span>` : ''}
      </div>
      <div class="profile-effective-metrics">
        <span>${fmtNum(effective?.active_rule_count || activeRules.length)} rules</span>
        <span class="mono">${escHtml(hash)}</span>
      </div>
    </div>
    ${rulesBlock}
  </div>`;
}

function renderEffectiveRule(rule, ruleByKey) {
  const meta = ruleByKey[rule.rule_key] || {};
  const tags = Array.isArray(meta.tags) ? meta.tags.slice(0, 3).join(', ') : '';
  return `<div class="profile-rule-row">
    <span class="profile-rule-name">${escHtml(meta.name || rule.rule_key)}</span>
    <span class="profile-rule-meta">${escHtml(rule.severity || meta.severity || '-')}</span>
    <span class="profile-rule-meta">${escHtml(meta.type || '-')}</span>
    <span class="profile-rule-meta">${escHtml(rule.origin || '')}</span>
    <span class="profile-rule-tags">${escHtml(tags)}</span>
  </div>`;
}

export function bindProfilesContent() {
  const project = state.currentProject;
  if (!project) return;
  const showToast = getShowToast();

  document.querySelectorAll('.assign-profile-btn').forEach(btn => {
    btn.addEventListener('click', async () => {
      const id = btn.dataset.profileId;
      const lang = btn.dataset.profileLang;
      const name = btn.dataset.profileName;
      btn.disabled = true;
      try {
        await apiFetch('/projects/' + encodeURIComponent(project.key) + '/profiles', {
          method: 'POST',
          body: JSON.stringify({ profile_id: Number.parseInt(id, 10), language: lang }),
        });
        showToast('Profile "' + name + '" assigned.');
        state.profilesData = null;
        await loadProfilesData();
      } catch (err) {
        showToast(err.message, 'error');
      }
      btn.disabled = false;
    });
  });

  document.querySelectorAll('.export-profile-btn').forEach(btn => {
    btn.addEventListener('click', async () => {
      const id = btn.dataset.profileId;
      const name = btn.dataset.profileName;
      btn.disabled = true;
      try {
        const doc = await apiFetch('/profiles/' + encodeURIComponent(id) + '/export');
        const text = JSON.stringify(doc, null, 2);
        await navigator.clipboard?.writeText(text);
        showToast('Profile "' + name + '" exported.');
      } catch (err) {
        showToast(err.message, 'error');
      }
      btn.disabled = false;
    });
  });
}
