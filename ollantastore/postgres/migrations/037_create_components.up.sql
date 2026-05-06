-- Components table: persistent project/directory/file tree.
-- UUID is deterministic (SHA256 of project_key:path:qualifier) for idempotent upserts.
CREATE TABLE IF NOT EXISTS components (
    uuid            TEXT PRIMARY KEY,
    project_id      BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    path            TEXT NOT NULL,
    qualifier       TEXT NOT NULL CHECK (qualifier IN ('TRK', 'DIR', 'FIL')),
    parent_uuid     TEXT REFERENCES components(uuid) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    last_scan_id    BIGINT REFERENCES scans(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_components_project ON components(project_id);
CREATE INDEX IF NOT EXISTS idx_components_parent ON components(parent_uuid);
