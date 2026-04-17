CREATE TABLE tokens (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         TEXT NOT NULL,
    token_hash   TEXT NOT NULL UNIQUE,
    token_type   TEXT NOT NULL,
    project_id   BIGINT REFERENCES projects(id) ON DELETE CASCADE,
    last_used_at TIMESTAMPTZ,
    expires_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_tokens_user ON tokens (user_id);
CREATE INDEX idx_tokens_hash ON tokens (token_hash);
