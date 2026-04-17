CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    login         TEXT NOT NULL UNIQUE,
    email         TEXT NOT NULL UNIQUE,
    name          TEXT NOT NULL DEFAULT '',
    password_hash TEXT NOT NULL DEFAULT '',
    avatar_url    TEXT NOT NULL DEFAULT '',
    provider      TEXT NOT NULL DEFAULT 'local',
    provider_id   TEXT NOT NULL DEFAULT '',
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    last_login_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_users_email    ON users (email);
CREATE INDEX idx_users_provider ON users (provider, provider_id);
