# SQL Migration Guide for AI Analysis Feature

## Overview
This guide provides the exact SQL scripts you need to run to set up the AI-powered relationship analysis feature with the unified Nudge model.

## Prerequisites
- PostgreSQL database running
- Existing tables: `users`, `persons`, `interactions`, `nudges`
- Database user with CREATE TABLE, ALTER TABLE, and CREATE INDEX permissions

---

## Migration 1: Create AI Analysis Tables

**File**: `migrations/000004_create_ai_analysis_tables.up.sql`

Run this first to create the core AI analysis infrastructure:

```sql
-- Create relationship_analyses table
CREATE TABLE IF NOT EXISTS relationship_analyses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    person_id UUID NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    
    -- Scores (0-100)
    connection_strength DECIMAL(5,2) DEFAULT 0,
    engagement_quality DECIMAL(5,2) DEFAULT 0,
    communication_balance DECIMAL(5,2) DEFAULT 0,
    energy_alignment DECIMAL(5,2) DEFAULT 0,
    relationship_health DECIMAL(5,2) DEFAULT 0,
    overall_score DECIMAL(5,2) DEFAULT 0,
    
    -- Analysis content
    summary TEXT,
    key_insights TEXT[],
    patterns TEXT[],
    strengths TEXT[],
    concerns TEXT[],
    trend_direction VARCHAR(50),
    
    -- Metadata
    provider VARCHAR(50) NOT NULL,
    model VARCHAR(100),
    tokens_used INTEGER DEFAULT 0,
    processing_time_ms INTEGER DEFAULT 0,
    version INTEGER DEFAULT 1,
    analyzed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    interactions_count INTEGER DEFAULT 0,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create indexes for relationship_analyses
CREATE INDEX idx_relationship_analyses_user_person ON relationship_analyses(user_id, person_id);
CREATE INDEX idx_relationship_analyses_user_id ON relationship_analyses(user_id);
CREATE INDEX idx_relationship_analyses_person_id ON relationship_analyses(person_id);
CREATE INDEX idx_relationship_analyses_analyzed_at ON relationship_analyses(analyzed_at DESC);
CREATE INDEX idx_relationship_analyses_deleted_at ON relationship_analyses(deleted_at);

-- Create ai_analysis_jobs table
CREATE TABLE IF NOT EXISTS ai_analysis_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Job details
    job_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    priority INTEGER DEFAULT 5,
    
    -- Target
    person_ids TEXT[],
    
    -- Progress tracking
    total_items INTEGER DEFAULT 0,
    processed_items INTEGER DEFAULT 0,
    failed_items INTEGER DEFAULT 0,
    progress DECIMAL(5,2) DEFAULT 0,
    
    -- Results
    result_data JSONB,
    error TEXT,
    
    -- Cost tracking
    total_tokens_used INTEGER DEFAULT 0,
    estimated_cost DECIMAL(10,4) DEFAULT 0,
    
    -- Timing
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create indexes for ai_analysis_jobs
CREATE INDEX idx_ai_analysis_jobs_user_id ON ai_analysis_jobs(user_id);
CREATE INDEX idx_ai_analysis_jobs_status ON ai_analysis_jobs(status);
CREATE INDEX idx_ai_analysis_jobs_priority ON ai_analysis_jobs(priority DESC);
CREATE INDEX idx_ai_analysis_jobs_created_at ON ai_analysis_jobs(created_at DESC);
CREATE INDEX idx_ai_analysis_jobs_deleted_at ON ai_analysis_jobs(deleted_at);

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_ai_tables_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_relationship_analyses_updated_at
    BEFORE UPDATE ON relationship_analyses
    FOR EACH ROW
    EXECUTE FUNCTION update_ai_tables_updated_at();

CREATE TRIGGER update_ai_analysis_jobs_updated_at
    BEFORE UPDATE ON ai_analysis_jobs
    FOR EACH ROW
    EXECUTE FUNCTION update_ai_tables_updated_at();
```

---

## Migration 2: Enhance Nudges Table for AI

**File**: `migrations/000005_enhance_nudges_for_ai.up.sql`

Run this second to add AI capabilities to the existing nudges table:

```sql
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
```

---

## Quick Execution Commands

### Using psql:
```bash
# Connect to your database
psql -U your_user -d vyve_dev

# Run migration 1
\i migrations/000004_create_ai_analysis_tables.up.sql

# Run migration 2
\i migrations/000005_enhance_nudges_for_ai.up.sql

# Verify tables were created
\dt relationship_analyses
\dt ai_analysis_jobs
\d nudges
```

### Using migrate CLI:
```bash
# Run all pending migrations
migrate -path migrations -database "postgresql://user:pass@localhost:5432/vyve_dev?sslmode=disable" up

# Or run specific number of migrations
migrate -path migrations -database "postgresql://user:pass@localhost:5432/vyve_dev?sslmode=disable" up 2
```

### Using DBeaver/pgAdmin:
1. Open SQL Editor
2. Copy and paste the SQL from migration 1
3. Execute
4. Copy and paste the SQL from migration 2
5. Execute

---

## Verification Queries

After running the migrations, verify everything is set up correctly:

```sql
-- Check relationship_analyses table
SELECT 
    table_name, 
    column_name, 
    data_type 
FROM information_schema.columns 
WHERE table_name = 'relationship_analyses'
ORDER BY ordinal_position;

-- Check ai_analysis_jobs table
SELECT 
    table_name, 
    column_name, 
    data_type 
FROM information_schema.columns 
WHERE table_name = 'ai_analysis_jobs'
ORDER BY ordinal_position;

-- Check nudges table has new columns
SELECT 
    column_name, 
    data_type,
    column_default
FROM information_schema.columns 
WHERE table_name = 'nudges'
AND column_name IN ('analysis_id', 'source', 'reasoning', 'suggested_actions', 
                     'conversation_starters', 'timing', 'estimated_impact', 
                     'status', 'provider', 'model')
ORDER BY column_name;

-- Check indexes
SELECT 
    indexname, 
    indexdef 
FROM pg_indexes 
WHERE tablename IN ('relationship_analyses', 'ai_analysis_jobs', 'nudges')
AND indexname LIKE '%analysis%' OR indexname LIKE '%source%' OR indexname LIKE '%status%'
ORDER BY tablename, indexname;

-- Check foreign key constraints
SELECT
    tc.table_name, 
    kcu.column_name,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name
FROM information_schema.table_constraints AS tc
JOIN information_schema.key_column_usage AS kcu
    ON tc.constraint_name = kcu.constraint_name
JOIN information_schema.constraint_column_usage AS ccu
    ON ccu.constraint_name = tc.constraint_name
WHERE tc.constraint_type = 'FOREIGN KEY'
AND tc.table_name IN ('relationship_analyses', 'ai_analysis_jobs', 'nudges');
```

---

## Rollback Instructions

If you need to rollback the changes:

### Rollback Migration 2 (Nudges Enhancement):
```sql
-- File: migrations/000005_enhance_nudges_for_ai.down.sql

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
```

### Rollback Migration 1 (AI Tables):
```sql
-- File: migrations/000004_create_ai_analysis_tables.down.sql

-- Drop triggers
DROP TRIGGER IF EXISTS update_relationship_analyses_updated_at ON relationship_analyses;
DROP TRIGGER IF EXISTS update_ai_analysis_jobs_updated_at ON ai_analysis_jobs;

-- Drop function
DROP FUNCTION IF EXISTS update_ai_tables_updated_at();

-- Drop tables in reverse order (respecting foreign keys)
DROP TABLE IF EXISTS ai_analysis_jobs;
DROP TABLE IF EXISTS relationship_analyses;
```

---

## Summary

**Tables Created:**
1. `relationship_analyses` - Stores AI analysis results
2. `ai_analysis_jobs` - Tracks background analysis jobs

**Tables Modified:**
1. `nudges` - Enhanced to support both AI and system recommendations

**New Indexes:** 15 total
- 5 on `relationship_analyses`
- 5 on `ai_analysis_jobs`
- 5 on `nudges`

**New Columns on nudges:**
- `analysis_id` - Links to AI analysis
- `source` - 'ai' or 'system'
- `reasoning` - Why this recommendation
- `suggested_actions` - Action steps array
- `conversation_starters` - Opening lines array
- `timing` - When to act
- `estimated_impact` - Expected impact
- `status` - Unified status tracking
- `accepted_at`, `completed_at`, `dismissed_at` - Status timestamps
- `provider`, `model` - AI metadata

---

## Next Steps After Migration

1. ✅ Run both migrations
2. ✅ Verify tables and columns exist
3. ✅ Check indexes are created
4. ✅ Test foreign key constraints
5. Configure environment variables (see `AI_ENV_VARS.md`)
6. Restart the application
7. Test AI analysis endpoints

---

## Troubleshooting

**Error: relation "users" does not exist**
- Make sure your base migrations have run first
- Check that users, persons, and interactions tables exist

**Error: column "analysis_id" already exists**
- Migration 2 has already been run
- Check current state with: `\d nudges`

**Error: permission denied**
- Ensure database user has CREATE TABLE and ALTER TABLE permissions
- Run as superuser or database owner

**Performance Issues**
- All necessary indexes are created by the migrations
- If queries are slow, run `ANALYZE relationship_analyses; ANALYZE nudges;`
