-- Worker heartbeat and project-level lock fields for scan jobs.
ALTER TABLE scan_jobs
    ADD COLUMN IF NOT EXISTS project_id INTEGER,
    ADD COLUMN IF NOT EXISTS worker_heartbeat TIMESTAMPTZ;
