ALTER TABLE issues ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now();

CREATE TABLE IF NOT EXISTS issue_changelog (
    id          BIGSERIAL PRIMARY KEY,
    issue_id    BIGINT NOT NULL,
    user_id     BIGINT,
    field       TEXT NOT NULL,
    old_value   TEXT NOT NULL DEFAULT '',
    new_value   TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_issue_changelog_issue_id ON issue_changelog(issue_id);
CREATE INDEX IF NOT EXISTS idx_issue_changelog_created_at ON issue_changelog(created_at);
