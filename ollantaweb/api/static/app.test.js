'use strict';

const test = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');

const apiPrefix = '/api/v1';

function createHarness(options = {}) {
  const search = options.search || '';
  const source = fs.readFileSync(path.join(__dirname, 'app.js'), 'utf8');
  const storage = new Map();
  const requests = [];
  const historyCalls = [];
  const elements = new Map();

  function getElement(id) {
    if (!elements.has(id)) {
      elements.set(id, {
        id,
        innerHTML: '',
        value: '',
        addEventListener() {},
        classList: { add() {}, remove() {} },
      });
    }
    return elements.get(id);
  }

  const context = {
    console,
    URLSearchParams,
    history: {
      pushState(_state, _title, url) { historyCalls.push({ type: 'push', url }); },
      replaceState(_state, _title, url) { historyCalls.push({ type: 'replace', url }); },
    },
    location: { search, pathname: '/ui' },
    window: { addEventListener() {} },
    localStorage: {
      getItem(key) { return storage.has(key) ? storage.get(key) : null; },
      setItem(key, value) { storage.set(key, String(value)); },
      removeItem(key) { storage.delete(key); },
    },
    document: {
      addEventListener() {},
        getElementById(id) { return getElement(id); },
      querySelectorAll() { return []; },
      createElement() { return { classList: { add() {}, remove() {} }, remove() {}, textContent: '' }; },
      body: { appendChild() {} },
    },
    fetch: async (url, requestOptions) => {
      requests.push(String(url));
      const response = options.fetchHandler ? await options.fetchHandler(String(url), requestOptions) : { body: {} };
      return {
        status: response.status ?? 200,
        ok: response.ok ?? true,
        json: async () => response.body ?? {},
      };
    },
    setTimeout() { return 0; },
    clearTimeout() {},
    requestAnimationFrame(callback) { callback(); },
    confirm() { return true; },
  };
  context.globalThis = context;
  vm.createContext(context);
  vm.runInContext(
    source + '\n;globalThis.__appExports = { state, normalizeScope, parseProjectRoute, buildProjectRoute, buildScopeQuery, changeScope, loadProject, loadCodeTreeData, renderCodeTab, renderProjectInformationTab, renderScopeToolbar };',
    context,
    { filename: 'app.js' },
  );
  return {
    app: context.__appExports,
    requests,
    historyCalls,
    elements,
  };
}

test('normalizeScope maps snake_case scope fields', () => {
  const { app } = createHarness();
  const scope = app.normalizeScope({
    type: 'pull_request',
    branch: 'feature/login',
    pull_request: '42',
    pull_request_base: 'main',
    default_branch: 'main',
  });

  assert.equal(scope.type, 'pull_request');
  assert.equal(scope.pullRequestKey, '42');
  assert.equal(scope.pullRequestBase, 'main');
  assert.equal(scope.defaultBranch, 'main');
});

test('parseProjectRoute reads project, tab and branch scope from query string', () => {
  const { app } = createHarness({ search: '?project=demo&tab=code&branch=release%2F1.2' });
  const route = { ...app.parseProjectRoute() };

  assert.deepEqual(route, {
    project: 'demo',
    tab: 'code',
    branch: 'release/1.2',
    pullRequest: '',
  });
});

test('buildProjectRoute preserves pull request scope in the URL', () => {
	const { app } = createHarness();
  const route = app.buildProjectRoute('demo', 'information', {
    type: 'pull_request',
    pullRequestKey: '42',
    branch: 'feature/login',
  });

  assert.equal(route, '?project=demo&tab=information&pull_request=42');
});

test('renderScopeToolbar shows the real default branch name in the branch selector', () => {
  const { app } = createHarness();

  app.state.scope = app.normalizeScope({ type: 'branch', branch: 'main', defaultBranch: 'main' });
  app.state.branchesData = [
    { name: '', is_default: false },
    { name: 'main', is_default: true },
    { name: 'release', is_default: false },
  ];
  app.state.pullRequestsData = [];

  const html = app.renderScopeToolbar();
  assert.match(html, /<option value="main">main · default<\/option>/);
  assert.doesNotMatch(html, /<option value="">Default branch<\/option>/);
});

test('changeScope refreshes scoped data and persists the branch in the URL', async () => {
  const { app, requests, historyCalls } = createHarness({
    fetchHandler(url) {
      if (url === apiPrefix + '/projects/demo/overview?branch=release') {
        return { body: { scope: { type: 'branch', branch: 'release', default_branch: 'main' } } };
      }
      if (url === apiPrefix + '/projects/demo/branches') {
        return { body: { items: [{ name: 'main', is_default: true }, { name: 'release', is_default: false }] } };
      }
      if (url === apiPrefix + '/projects/demo/pull-requests') {
        return { body: { items: [] } };
      }
      throw new Error('unexpected fetch ' + url);
    },
  });

  app.state.currentProject = { key: 'demo' };
  app.state.projectTab = 'overview';
  app.state.scope = app.normalizeScope({ type: 'branch', branch: 'main', defaultBranch: 'main' });

  await app.changeScope({ type: 'branch', branch: 'release', defaultBranch: 'main' });

  assert.deepEqual(requests, [
    apiPrefix + '/projects/demo/overview?branch=release',
    apiPrefix + '/projects/demo/branches',
    apiPrefix + '/projects/demo/pull-requests',
  ]);
  assert.equal(app.state.scope.branch, 'release');
  assert.equal(historyCalls.at(-1).url, '?project=demo&tab=overview&branch=release');
});

test('loadProject preserves branch scope in reloads and refreshes project information', async () => {
  const { app, requests, historyCalls } = createHarness({
    fetchHandler(url) {
      if (url === apiPrefix + '/projects/demo') {
        return { body: { key: 'demo', name: 'Demo', main_branch: 'main' } };
      }
      if (url === apiPrefix + '/projects/demo/overview?branch=release') {
        return {
          body: {
            scope: { type: 'branch', branch: 'release', default_branch: 'main' },
            last_scan: { analysis_date: '2026-04-21T12:00:00Z' },
          },
        };
      }
      if (url === apiPrefix + '/projects/demo/branches') {
        return { body: { items: [{ name: 'main', is_default: true }, { name: 'release', is_default: false }] } };
      }
      if (url === apiPrefix + '/projects/demo/pull-requests') {
        return { body: { items: [] } };
      }
      if (url === apiPrefix + '/projects/demo/information?branch=release') {
        return {
          body: {
            project: { key: 'demo', name: 'Demo', main_branch: 'main' },
            scope: { type: 'branch', branch: 'release', default_branch: 'main' },
            measures: { files: 4, lines: 24, ncloc: 18, issues: 2 },
            code_snapshot: { stored_files: 1, total_files: 1, max_file_bytes: 128, max_total_bytes: 512 },
          },
        };
      }
      throw new Error('unexpected fetch ' + url);
    },
  });

  await app.loadProject('demo', { project: 'demo', tab: 'information', branch: 'release', pullRequest: '' });

  assert.equal(app.state.projectTab, 'information');
  assert.equal(app.state.scope.branch, 'release');
  assert.ok(requests.includes(apiPrefix + '/projects/demo/overview?branch=release'));
  assert.ok(requests.includes(apiPrefix + '/projects/demo/information?branch=release'));
  assert.equal(historyCalls.at(-1).url, '?project=demo&tab=information&branch=release');
  assert.match(app.renderProjectInformationTab(), /release/);
});

test('changeScope refreshes project information for a selected pull request', async () => {
  const { app, requests, historyCalls } = createHarness({
    fetchHandler(url) {
      if (url === apiPrefix + '/projects/demo/overview?pull_request=128') {
        return { body: { scope: { type: 'pull_request', branch: 'feature/login', pull_request_key: '128', pull_request_base: 'main', default_branch: 'main' } } };
      }
      if (url === apiPrefix + '/projects/demo/branches') {
        return { body: { items: [{ name: 'main', is_default: true }] } };
      }
      if (url === apiPrefix + '/projects/demo/pull-requests') {
        return { body: { items: [{ key: '128', branch: 'feature/login', base_branch: 'main' }] } };
      }
      if (url === apiPrefix + '/projects/demo/information?pull_request=128') {
        return {
          body: {
            project: { key: 'demo', name: 'Demo', main_branch: 'main' },
            scope: { type: 'pull_request', branch: 'feature/login', pull_request_key: '128', pull_request_base: 'main', default_branch: 'main' },
            measures: { files: 2, lines: 12, ncloc: 10, issues: 1 },
            code_snapshot: { stored_files: 1, total_files: 1, max_file_bytes: 128, max_total_bytes: 512 },
          },
        };
      }
      throw new Error('unexpected fetch ' + url);
    },
  });

  app.state.currentProject = { key: 'demo' };
  app.state.projectTab = 'information';
  app.state.scope = app.normalizeScope({ type: 'branch', branch: 'main', defaultBranch: 'main' });

  await app.changeScope({ type: 'pull_request', pullRequestKey: '128', branch: 'feature/login', pullRequestBase: 'main', defaultBranch: 'main' });

  assert.ok(requests.includes(apiPrefix + '/projects/demo/information?pull_request=128'));
  assert.equal(app.state.scope.pullRequestKey, '128');
  assert.equal(historyCalls.at(-1).url, '?project=demo&tab=information&pull_request=128');
  assert.match(app.renderProjectInformationTab(), /Pull request/i);
});

test('loadCodeTreeData fetches the scoped file and renderCodeTab shows issue markers', async () => {
  const { app, requests } = createHarness({
    fetchHandler(url) {
      if (url === apiPrefix + '/projects/demo/code/tree?branch=release') {
        return {
          body: {
            code_snapshot: { stored_files: 1, total_files: 1 },
            items: [{ path: 'src/app.go', line_count: 2, size_bytes: 18, is_omitted: false }],
          },
        };
      }
      if (url === apiPrefix + '/projects/demo/code/file?path=src%2Fapp.go&branch=release') {
        return {
          body: {
            file: {
              path: 'src/app.go',
              language: 'go',
              content: 'line one\nline two',
              line_count: 2,
              size_bytes: 18,
            },
            issues: [{ rule_key: 'go:no-large-functions', severity: 'major', line: 2, message: 'Too large' }],
          },
        };
      }
      throw new Error('unexpected fetch ' + url);
    },
  });

  app.state.currentProject = { key: 'demo' };
  app.state.scope = app.normalizeScope({ type: 'branch', branch: 'release', defaultBranch: 'main' });

  await app.loadCodeTreeData();

  assert.deepEqual(requests, [
    apiPrefix + '/projects/demo/code/tree?branch=release',
    apiPrefix + '/projects/demo/code/file?path=src%2Fapp.go&branch=release',
  ]);
  const html = app.renderCodeTab();
  assert.match(html, /src\/app\.go/);
  assert.match(html, /go:no-large-functions/);
  assert.match(html, /has-issue/);
  assert.match(html, /sev-major/);
});