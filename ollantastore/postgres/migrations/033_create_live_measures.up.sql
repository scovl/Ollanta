-- Current-value measures updated atomically every scan.
-- Avoids expensive subqueries with MAX(scan_id) for overview and dashboards.
CREATE TABLE IF NOT EXISTS live_measures (
    project_id      BIGINT    NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    component_path  TEXT      NOT NULL DEFAULT '',
    metric_key      TEXT      NOT NULL,
    value           NUMERIC   NOT NULL,
    text_value      TEXT,
    scan_id         BIGINT    REFERENCES scans(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (project_id, component_path, metric_key)
);
