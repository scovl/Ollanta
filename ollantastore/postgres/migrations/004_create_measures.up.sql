CREATE TABLE IF NOT EXISTS measures (
    id             BIGSERIAL PRIMARY KEY,
    scan_id        BIGINT NOT NULL REFERENCES scans(id) ON DELETE CASCADE,
    project_id     BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    metric_key     TEXT NOT NULL,
    component_path TEXT NOT NULL DEFAULT '',
    value          DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_measures_scan    ON measures (scan_id);
CREATE INDEX IF NOT EXISTS idx_measures_project ON measures (project_id, metric_key, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_measures_trend   ON measures (project_id, metric_key, created_at);
