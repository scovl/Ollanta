-- Daily rollup of measure values for efficient trend queries.
-- Each scan upserts into this table, merging values for the same day.
CREATE TABLE IF NOT EXISTS measure_daily_aggregates (
    project_id      BIGINT    NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    metric_key      TEXT      NOT NULL,
    date            DATE      NOT NULL,
    value_avg       NUMERIC   NOT NULL,
    value_max       NUMERIC   NOT NULL,
    value_min       NUMERIC   NOT NULL,
    sample_count    INTEGER   NOT NULL DEFAULT 1,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (project_id, metric_key, date)
);
