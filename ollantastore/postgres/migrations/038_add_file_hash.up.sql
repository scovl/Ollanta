-- File hash for cross-path issue matching on rename detection.
ALTER TABLE issues
    ADD COLUMN IF NOT EXISTS file_hash TEXT NOT NULL DEFAULT '';
