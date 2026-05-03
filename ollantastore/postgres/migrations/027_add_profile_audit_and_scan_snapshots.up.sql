CREATE TABLE IF NOT EXISTS quality_profile_changelog (
    id            BIGSERIAL PRIMARY KEY,
    profile_id    BIGINT,
    project_id    BIGINT,
    language      TEXT NOT NULL DEFAULT '',
    action        TEXT NOT NULL,
    rule_key      TEXT NOT NULL DEFAULT '',
    old_value     TEXT NOT NULL DEFAULT '',
    new_value     TEXT NOT NULL DEFAULT '',
    actor_user_id BIGINT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_quality_profile_changelog_profile_created
    ON quality_profile_changelog(profile_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_quality_profile_changelog_project_created
    ON quality_profile_changelog(project_id, created_at DESC);

CREATE TABLE IF NOT EXISTS scan_profile_snapshots (
    id                 BIGSERIAL PRIMARY KEY,
    project_id         BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    scan_id            BIGINT NOT NULL REFERENCES scans(id) ON DELETE CASCADE,
    scope_type         TEXT NOT NULL DEFAULT 'branch',
    branch             TEXT NOT NULL DEFAULT '',
    pull_request_key   TEXT NOT NULL DEFAULT '',
    language           TEXT NOT NULL DEFAULT '',
    profile_id         BIGINT,
    profile_name       TEXT NOT NULL DEFAULT '',
    source             TEXT NOT NULL DEFAULT '',
    active_rule_count  INT NOT NULL DEFAULT 0,
    rules_hash         TEXT NOT NULL DEFAULT '',
    metadata_available BOOLEAN NOT NULL DEFAULT TRUE,
    diagnostics        JSONB NOT NULL DEFAULT '[]',
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (scan_id, language)
);

CREATE INDEX IF NOT EXISTS idx_scan_profile_snapshots_project_scope
    ON scan_profile_snapshots(project_id, scope_type, branch, pull_request_key, scan_id DESC);
CREATE INDEX IF NOT EXISTS idx_scan_profile_snapshots_scan
    ON scan_profile_snapshots(scan_id);