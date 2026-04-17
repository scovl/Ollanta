CREATE TABLE IF NOT EXISTS scans (
    id                    BIGSERIAL PRIMARY KEY,
    project_id            BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    version               TEXT NOT NULL DEFAULT '',
    branch                TEXT NOT NULL DEFAULT '',
    commit_sha            TEXT NOT NULL DEFAULT '',
    status                TEXT NOT NULL DEFAULT 'completed',
    elapsed_ms            BIGINT NOT NULL DEFAULT 0,
    gate_status           TEXT NOT NULL DEFAULT '',
    analysis_date         TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
    total_files           INT NOT NULL DEFAULT 0,
    total_lines           INT NOT NULL DEFAULT 0,
    total_ncloc           INT NOT NULL DEFAULT 0,
    total_comments        INT NOT NULL DEFAULT 0,
    total_issues          INT NOT NULL DEFAULT 0,
    total_bugs            INT NOT NULL DEFAULT 0,
    total_code_smells     INT NOT NULL DEFAULT 0,
    total_vulnerabilities INT NOT NULL DEFAULT 0,
    new_issues            INT NOT NULL DEFAULT 0,
    closed_issues         INT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_scans_project_date ON scans (project_id, analysis_date DESC);
