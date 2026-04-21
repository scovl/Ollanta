ALTER TABLE projects
    ADD COLUMN IF NOT EXISTS main_branch TEXT NOT NULL DEFAULT '';

ALTER TABLE scans
    ADD COLUMN IF NOT EXISTS scope_type TEXT NOT NULL DEFAULT 'branch',
    ADD COLUMN IF NOT EXISTS pull_request_key TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS pull_request_base TEXT NOT NULL DEFAULT '';

UPDATE scans
SET scope_type = 'branch'
WHERE scope_type = '';

CREATE INDEX IF NOT EXISTS idx_scans_project_branch_scope_date
    ON scans (project_id, scope_type, branch, analysis_date DESC);

CREATE INDEX IF NOT EXISTS idx_scans_project_pr_scope_date
    ON scans (project_id, scope_type, pull_request_key, analysis_date DESC);

CREATE TABLE IF NOT EXISTS code_snapshot_scopes (
    project_id       BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    scope_type       TEXT NOT NULL,
    scope_key        TEXT NOT NULL,
    scan_id          BIGINT NOT NULL REFERENCES scans(id) ON DELETE CASCADE,
    branch           TEXT NOT NULL DEFAULT '',
    pull_request_key TEXT NOT NULL DEFAULT '',
    pull_request_base TEXT NOT NULL DEFAULT '',
    total_files      INT NOT NULL DEFAULT 0,
    stored_files     INT NOT NULL DEFAULT 0,
    truncated_files  INT NOT NULL DEFAULT 0,
    omitted_files    INT NOT NULL DEFAULT 0,
    stored_bytes     INT NOT NULL DEFAULT 0,
    max_file_bytes   INT NOT NULL DEFAULT 0,
    max_total_bytes  INT NOT NULL DEFAULT 0,
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (project_id, scope_type, scope_key)
);

CREATE TABLE IF NOT EXISTS code_snapshot_files (
    project_id       BIGINT NOT NULL,
    scope_type       TEXT NOT NULL,
    scope_key        TEXT NOT NULL,
    path             TEXT NOT NULL,
    language         TEXT NOT NULL DEFAULT '',
    content          TEXT NOT NULL DEFAULT '',
    size_bytes       INT NOT NULL DEFAULT 0,
    line_count       INT NOT NULL DEFAULT 0,
    is_truncated     BOOLEAN NOT NULL DEFAULT FALSE,
    is_omitted       BOOLEAN NOT NULL DEFAULT FALSE,
    omitted_reason   TEXT NOT NULL DEFAULT '',
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (project_id, scope_type, scope_key, path),
    FOREIGN KEY (project_id, scope_type, scope_key)
        REFERENCES code_snapshot_scopes(project_id, scope_type, scope_key)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_code_snapshot_scopes_updated
    ON code_snapshot_scopes (project_id, scope_type, scope_key, updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_code_snapshot_files_scope_path
    ON code_snapshot_files (project_id, scope_type, scope_key, path);