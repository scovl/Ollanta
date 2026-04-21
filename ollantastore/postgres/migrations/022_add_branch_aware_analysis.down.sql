DROP INDEX IF EXISTS idx_code_snapshot_files_scope_path;
DROP INDEX IF EXISTS idx_code_snapshot_scopes_updated;
DROP TABLE IF EXISTS code_snapshot_files;
DROP TABLE IF EXISTS code_snapshot_scopes;

DROP INDEX IF EXISTS idx_scans_project_pr_scope_date;
DROP INDEX IF EXISTS idx_scans_project_branch_scope_date;

ALTER TABLE scans
    DROP COLUMN IF EXISTS pull_request_base,
    DROP COLUMN IF EXISTS pull_request_key,
    DROP COLUMN IF EXISTS scope_type;

ALTER TABLE projects
    DROP COLUMN IF EXISTS main_branch;