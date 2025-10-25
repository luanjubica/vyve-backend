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
