-- Default admin account. Password: "admin" (bcrypt cost 10). Change on first login.
INSERT INTO users (login, email, name, password_hash, provider)
VALUES (
    'admin',
    'admin@localhost',
    'Administrator',
    '$2a$10$uFynqN5wS8pUyMfgBTvjP.JYbx5ngtEbsaZHkKkzANtlLZXVugPCq',
    'local'
)
ON CONFLICT (login) DO NOTHING;

INSERT INTO group_members (group_id, user_id)
SELECT g.id, u.id
FROM groups g
CROSS JOIN users u
WHERE g.name IN ('ollanta-admins', 'ollanta-users')
  AND u.login = 'admin'
ON CONFLICT DO NOTHING;

-- Grant admin all global permissions individually as well (via user target)
INSERT INTO global_permissions (target, target_id, permission)
SELECT 'user', u.id, unnest(ARRAY[
    'admin',
    'create_project',
    'manage_users',
    'manage_groups',
    'execute_analysis',
    'manage_quality_gates'
])
FROM users u WHERE u.login = 'admin'
ON CONFLICT DO NOTHING;
