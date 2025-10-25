-- Drop triggers
DROP TRIGGER IF EXISTS update_relationship_analyses_updated_at ON relationship_analyses;
DROP TRIGGER IF EXISTS update_ai_analysis_jobs_updated_at ON ai_analysis_jobs;

-- Drop function
DROP FUNCTION IF EXISTS update_ai_tables_updated_at();

-- Drop tables in reverse order (respecting foreign keys)
DROP TABLE IF EXISTS ai_analysis_jobs;
DROP TABLE IF EXISTS relationship_analyses;
