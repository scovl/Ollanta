-- Restore the original (incorrect) state of migration 011 if rolling back.
DELETE FROM quality_profiles WHERE is_builtin = TRUE AND language = 'rust';

INSERT INTO quality_profiles (name, language, is_default, is_builtin)
VALUES
    ('Ollanta Way', 'java',   TRUE, TRUE),
    ('Ollanta Way', 'csharp', TRUE, TRUE)
ON CONFLICT (name, language) DO NOTHING;
