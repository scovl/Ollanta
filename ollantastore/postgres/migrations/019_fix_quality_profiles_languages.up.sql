-- Remove unsupported built-in profiles (java and csharp have no scanner support
-- or rules in ollantacore/constants.ExtensionToLanguage or ollantarules/).
DELETE FROM quality_profiles
WHERE is_builtin = TRUE AND language IN ('java', 'csharp');

-- Add the missing rust profile (supported via .rs extension in ollantacore).
INSERT INTO quality_profiles (name, language, is_default, is_builtin)
VALUES ('Ollanta Way', 'rust', TRUE, TRUE)
ON CONFLICT (name, language) DO NOTHING;
