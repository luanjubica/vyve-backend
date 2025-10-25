-- Enhance nudges table to support both AI-generated and system-generated recommendations

-- Add new columns for AI functionality
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS analysis_id UUID REFERENCES relationship_analyses(id) ON DELETE SET NULL;
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS source VARCHAR(20) DEFAULT 'system';
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS reasoning TEXT;
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS suggested_actions TEXT[];
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS conversation_starters TEXT[];
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS timing VARCHAR(50);
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS estimated_impact VARCHAR(20);

-- Enhance status tracking
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'pending';
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS accepted_at TIMESTAMP;
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS completed_at TIMESTAMP;
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS dismissed_at TIMESTAMP;

-- Add AI metadata for cost tracking
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS provider VARCHAR(50);
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS model VARCHAR(100);

-- Update message column to TEXT if not already
ALTER TABLE nudges ALTER COLUMN message TYPE TEXT;

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_nudges_analysis_id ON nudges(analysis_id);
CREATE INDEX IF NOT EXISTS idx_nudges_source ON nudges(source);
CREATE INDEX IF NOT EXISTS idx_nudges_status ON nudges(status);
CREATE INDEX IF NOT EXISTS idx_nudges_timing ON nudges(timing);
CREATE INDEX IF NOT EXISTS idx_nudges_expires_at ON nudges(expires_at) WHERE expires_at IS NOT NULL;

-- Update existing nudges to have proper status based on seen/acted_on
UPDATE nudges 
SET status = CASE 
    WHEN acted_on = true THEN 'completed'
    WHEN seen = true THEN 'seen'
    ELSE 'pending'
END
WHERE status = 'pending' OR status IS NULL;

-- Add comment for documentation
COMMENT ON COLUMN nudges.source IS 'Source of the nudge: ai (AI-generated) or system (rule-based)';
COMMENT ON COLUMN nudges.analysis_id IS 'Links to relationship_analyses if this is an AI-generated nudge';
COMMENT ON COLUMN nudges.status IS 'Current status: pending, seen, accepted, completed, dismissed';
