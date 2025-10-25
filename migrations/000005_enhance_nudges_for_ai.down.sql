-- Rollback: Remove AI enhancements from nudges table

-- Drop indexes
DROP INDEX IF EXISTS idx_nudges_analysis_id;
DROP INDEX IF EXISTS idx_nudges_source;
DROP INDEX IF EXISTS idx_nudges_status;
DROP INDEX IF EXISTS idx_nudges_timing;
DROP INDEX IF EXISTS idx_nudges_expires_at;

-- Remove AI-specific columns
ALTER TABLE nudges DROP COLUMN IF EXISTS analysis_id;
ALTER TABLE nudges DROP COLUMN IF EXISTS source;
ALTER TABLE nudges DROP COLUMN IF EXISTS reasoning;
ALTER TABLE nudges DROP COLUMN IF EXISTS suggested_actions;
ALTER TABLE nudges DROP COLUMN IF EXISTS conversation_starters;
ALTER TABLE nudges DROP COLUMN IF EXISTS timing;
ALTER TABLE nudges DROP COLUMN IF EXISTS estimated_impact;

-- Remove enhanced status tracking
ALTER TABLE nudges DROP COLUMN IF EXISTS status;
ALTER TABLE nudges DROP COLUMN IF EXISTS accepted_at;
ALTER TABLE nudges DROP COLUMN IF EXISTS completed_at;
ALTER TABLE nudges DROP COLUMN IF EXISTS dismissed_at;

-- Remove AI metadata
ALTER TABLE nudges DROP COLUMN IF EXISTS provider;
ALTER TABLE nudges DROP COLUMN IF EXISTS model;
