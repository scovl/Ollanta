ALTER TABLE scan_jobs
    DROP COLUMN IF EXISTS worker_heartbeat,
    DROP COLUMN IF EXISTS project_id;
