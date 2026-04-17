CREATE TABLE IF NOT EXISTS issues (
    id             BIGSERIAL,
    scan_id        BIGINT NOT NULL,
    project_id     BIGINT NOT NULL,
    rule_key       TEXT NOT NULL,
    component_path TEXT NOT NULL,
    line           INT NOT NULL DEFAULT 0,
    column_num     INT NOT NULL DEFAULT 0,
    end_line       INT NOT NULL DEFAULT 0,
    end_column     INT NOT NULL DEFAULT 0,
    message        TEXT NOT NULL DEFAULT '',
    type           TEXT NOT NULL,
    severity       TEXT NOT NULL,
    status         TEXT NOT NULL DEFAULT 'open',
    resolution     TEXT NOT NULL DEFAULT '',
    effort_minutes INT NOT NULL DEFAULT 0,
    line_hash      TEXT NOT NULL DEFAULT '',
    tags           TEXT[] NOT NULL DEFAULT '{}',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

CREATE TABLE IF NOT EXISTS issues_default PARTITION OF issues DEFAULT;

CREATE INDEX IF NOT EXISTS idx_issues_scan     ON issues (scan_id);
CREATE INDEX IF NOT EXISTS idx_issues_project  ON issues (project_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_issues_rule     ON issues (rule_key);
CREATE INDEX IF NOT EXISTS idx_issues_severity ON issues (severity);
CREATE INDEX IF NOT EXISTS idx_issues_type     ON issues (type);
CREATE INDEX IF NOT EXISTS idx_issues_hash     ON issues (project_id, rule_key, line_hash);
