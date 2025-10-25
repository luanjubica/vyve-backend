# SQL Quick Reference - AI Analysis Feature

## 🚀 Quick Start: Run These Two Scripts

### 1️⃣ First Migration (AI Tables)
```bash
psql -U your_user -d vyve_dev -f migrations/000004_create_ai_analysis_tables.up.sql
```

### 2️⃣ Second Migration (Enhance Nudges)
```bash
psql -U your_user -d vyve_dev -f migrations/000005_enhance_nudges_for_ai.up.sql
```

---

## 📊 What Gets Created

### New Tables (2)
- ✅ `relationship_analyses` - AI analysis results with scores
- ✅ `ai_analysis_jobs` - Background job tracking

### Modified Tables (1)
- ✅ `nudges` - Enhanced with 13 new columns for AI support

### New Indexes (15)
- 5 on `relationship_analyses`
- 5 on `ai_analysis_jobs`  
- 5 on `nudges`

---

## ✅ Verification (Run After Migration)

```sql
-- Quick check all tables exist
SELECT table_name FROM information_schema.tables 
WHERE table_name IN ('relationship_analyses', 'ai_analysis_jobs', 'nudges');

-- Check nudges has new columns
SELECT column_name FROM information_schema.columns 
WHERE table_name = 'nudges' 
AND column_name IN ('analysis_id', 'source', 'status', 'reasoning');

-- Count indexes
SELECT COUNT(*) FROM pg_indexes 
WHERE tablename IN ('relationship_analyses', 'ai_analysis_jobs', 'nudges');
```

Expected results:
- 3 tables found
- 4 new columns on nudges
- 15+ indexes total

---

## 🔄 Rollback (If Needed)

```bash
# Rollback in reverse order
psql -U your_user -d vyve_dev -f migrations/000005_enhance_nudges_for_ai.down.sql
psql -U your_user -d vyve_dev -f migrations/000004_create_ai_analysis_tables.down.sql
```

---

## 📝 Key Schema Changes

### relationship_analyses
```sql
CREATE TABLE relationship_analyses (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    person_id UUID REFERENCES persons(id),
    
    -- Scores (0-100)
    connection_strength DECIMAL(5,2),
    engagement_quality DECIMAL(5,2),
    overall_score DECIMAL(5,2),
    
    -- AI content
    summary TEXT,
    key_insights TEXT[],
    patterns TEXT[],
    
    -- Metadata
    provider VARCHAR(50),
    tokens_used INTEGER,
    analyzed_at TIMESTAMP
);
```

### nudges (new columns)
```sql
ALTER TABLE nudges ADD COLUMN analysis_id UUID;      -- Link to AI analysis
ALTER TABLE nudges ADD COLUMN source VARCHAR(20);    -- 'ai' or 'system'
ALTER TABLE nudges ADD COLUMN reasoning TEXT;        -- Why this nudge
ALTER TABLE nudges ADD COLUMN suggested_actions TEXT[];
ALTER TABLE nudges ADD COLUMN conversation_starters TEXT[];
ALTER TABLE nudges ADD COLUMN timing VARCHAR(50);    -- now, today, this_week
ALTER TABLE nudges ADD COLUMN estimated_impact VARCHAR(20);
ALTER TABLE nudges ADD COLUMN status VARCHAR(50);    -- pending, completed, etc.
ALTER TABLE nudges ADD COLUMN provider VARCHAR(50);  -- openai, anthropic
```

---

## 🎯 Common Queries After Setup

### Get all AI-generated nudges
```sql
SELECT * FROM nudges WHERE source = 'ai' ORDER BY created_at DESC LIMIT 10;
```

### Get latest analysis for a person
```sql
SELECT * FROM relationship_analyses 
WHERE user_id = 'user-uuid' AND person_id = 'person-uuid'
ORDER BY analyzed_at DESC LIMIT 1;
```

### Check job status
```sql
SELECT id, status, progress, total_items, processed_items 
FROM ai_analysis_jobs 
WHERE user_id = 'user-uuid' 
ORDER BY created_at DESC;
```

### Count nudges by source
```sql
SELECT source, COUNT(*) 
FROM nudges 
GROUP BY source;
```

---

## 🔍 Troubleshooting

### Check if migrations ran
```sql
-- If using golang-migrate
SELECT version, dirty FROM schema_migrations;
```

### Check table structure
```sql
\d relationship_analyses
\d ai_analysis_jobs
\d nudges
```

### Check foreign keys
```sql
SELECT
    tc.table_name,
    kcu.column_name,
    ccu.table_name AS foreign_table
FROM information_schema.table_constraints tc
JOIN information_schema.key_column_usage kcu ON tc.constraint_name = kcu.constraint_name
JOIN information_schema.constraint_column_usage ccu ON ccu.constraint_name = tc.constraint_name
WHERE tc.constraint_type = 'FOREIGN KEY'
AND tc.table_name IN ('relationship_analyses', 'nudges');
```

---

## 📦 Complete Setup Checklist

- [ ] Run migration 000004 (AI tables)
- [ ] Run migration 000005 (Enhance nudges)
- [ ] Verify 3 tables exist
- [ ] Verify nudges has new columns
- [ ] Check indexes created (15+)
- [ ] Test foreign key constraints
- [ ] Add API keys to .env
- [ ] Restart application
- [ ] Test `/api/v1/people/:id/analysis` endpoint

---

## 💡 Pro Tips

1. **Always backup before migrations**
   ```bash
   pg_dump vyve_dev > backup_before_ai_migration.sql
   ```

2. **Test in development first**
   - Run on dev database
   - Verify everything works
   - Then apply to production

3. **Monitor performance**
   ```sql
   -- Check table sizes
   SELECT 
       schemaname,
       tablename,
       pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
   FROM pg_tables
   WHERE tablename IN ('relationship_analyses', 'ai_analysis_jobs', 'nudges')
   ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
   ```

4. **Keep migrations in version control**
   - Commit migration files
   - Track which version is deployed where
   - Document any manual changes

---

## 🆘 Need Help?

- **Full guide**: See `SQL_MIGRATION_GUIDE.md`
- **Architecture**: See `MERGE_RECOMMENDATIONS_NUDGES.md`
- **Configuration**: See `AI_ENV_VARS.md`
- **Implementation**: See `IMPLEMENTATION_SUMMARY.md`
