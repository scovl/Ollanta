DELETE FROM group_members WHERE user_id IN (SELECT id FROM users WHERE login = 'admin');
DELETE FROM global_permissions WHERE target = 'user' AND target_id IN (SELECT id FROM users WHERE login = 'admin');
DELETE FROM users WHERE login = 'admin';
