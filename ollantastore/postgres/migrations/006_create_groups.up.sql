CREATE TABLE groups (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    is_builtin  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE group_members (
    group_id   BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    user_id    BIGINT NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (group_id, user_id)
);

INSERT INTO groups (name, description, is_builtin) VALUES
    ('ollanta-users',  'All authenticated users', TRUE),
    ('ollanta-admins', 'Instance administrators', TRUE);
