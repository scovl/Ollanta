CREATE TABLE global_permissions (
    id         BIGSERIAL PRIMARY KEY,
    target     TEXT NOT NULL,
    target_id  BIGINT NOT NULL,
    permission TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (target, target_id, permission)
);

CREATE TABLE project_permissions (
    id         BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    target     TEXT NOT NULL,
    target_id  BIGINT NOT NULL,
    permission TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (project_id, target, target_id, permission)
);

CREATE INDEX idx_project_perms_project ON project_permissions (project_id);

INSERT INTO global_permissions (target, target_id, permission)
SELECT 'group', id, unnest(ARRAY[
    'admin',
    'create_project',
    'manage_users',
    'manage_groups',
    'execute_analysis',
    'manage_quality_gates'
])
FROM groups WHERE name = 'ollanta-admins';

INSERT INTO global_permissions (target, target_id, permission)
SELECT 'group', id, unnest(ARRAY[
    'browse',
    'execute_analysis'
])
FROM groups WHERE name = 'ollanta-users';
