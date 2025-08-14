-- Create events table for analytics
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_type VARCHAR(100) NOT NULL,
    properties JSONB DEFAULT '{}',
    session_id VARCHAR(255),
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_events_user_id ON events(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_events_event_type ON events(event_type);
CREATE INDEX idx_events_session_id ON events(session_id);
CREATE INDEX idx_events_created_at ON events(created_at);
CREATE INDEX idx_events_user_type_date ON events(user_id, event_type, created_at);

-- Create daily_metrics table for aggregated analytics
CREATE TABLE daily_metrics (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    interactions_count INTEGER DEFAULT 0,
    unique_persons_count INTEGER DEFAULT 0,
    avg_energy_score FLOAT DEFAULT 0,
    reflection_completed BOOLEAN DEFAULT FALSE,
    nudges_generated INTEGER DEFAULT 0,
    nudges_acted_on INTEGER DEFAULT 0,
    positive_interactions INTEGER DEFAULT 0,
    negative_interactions INTEGER DEFAULT 0,
    relationships_active INTEGER DEFAULT 0,
    relationships_improved INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, date)
);

-- Create indexes
CREATE INDEX idx_daily_metrics_user_id ON daily_metrics(user_id);
CREATE INDEX idx_daily_metrics_date ON daily_metrics(date);
CREATE INDEX idx_daily_metrics_user_date ON daily_metrics(user_id, date DESC);

-- Create user_consents table for GDPR compliance
CREATE TABLE user_consents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    consent_type VARCHAR(50) NOT NULL, -- marketing, analytics, cookies, data_processing
    granted BOOLEAN NOT NULL,
    version VARCHAR(20),
    ip_address VARCHAR(45),
    user_agent TEXT,
    granted_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_user_consents_user_id ON user_consents(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_consents_type ON user_consents(consent_type);
CREATE INDEX idx_user_consents_granted ON user_consents(granted);

-- Create audit_logs table for GDPR compliance
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50),
    entity_id VARCHAR(255),
    changes JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    request_id VARCHAR(255),
    session_id VARCHAR(255),
    result VARCHAR(20), -- success, failure
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_request_id ON audit_logs(request_id);

-- Create data_exports table for GDPR data export requests
CREATE TABLE data_exports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, processing, completed, failed
    format VARCHAR(10) DEFAULT 'json', -- json, csv
    file_url TEXT,
    file_size BIGINT,
    requested_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    expires_at TIMESTAMP,
    error TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_data_exports_user_id ON data_exports(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_data_exports_status ON data_exports(status);
CREATE INDEX idx_data_exports_expires_at ON data_exports(expires_at);

-- Add triggers for updated_at
CREATE TRIGGER update_events_updated_at BEFORE UPDATE ON events
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_daily_metrics_updated_at BEFORE UPDATE ON daily_metrics
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_consents_updated_at BEFORE UPDATE ON user_consents
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_audit_logs_updated_at BEFORE UPDATE ON audit_logs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_data_exports_updated_at BEFORE UPDATE ON data_exports
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to log audit events
CREATE OR REPLACE FUNCTION log_audit_event(
    p_user_id UUID,
    p_action VARCHAR,
    p_entity_type VARCHAR,
    p_entity_id VARCHAR,
    p_changes JSONB,
    p_ip_address VARCHAR,
    p_user_agent TEXT,
    p_request_id VARCHAR,
    p_session_id VARCHAR,
    p_result VARCHAR,
    p_error_message TEXT DEFAULT NULL
) RETURNS UUID AS $$
DECLARE
    audit_id UUID;
BEGIN
    INSERT INTO audit_logs (
        user_id, action, entity_type, entity_id, changes,
        ip_address, user_agent, request_id, session_id,
        result, error_message
    ) VALUES (
        p_user_id, p_action, p_entity_type, p_entity_id, p_changes,
        p_ip_address, p_user_agent, p_request_id, p_session_id,
        p_result, p_error_message
    ) RETURNING id INTO audit_id;
    
    RETURN audit_id;
END;
$$ LANGUAGE plpgsql;

-- Function to aggregate daily metrics
CREATE OR REPLACE FUNCTION aggregate_daily_metrics(p_user_id UUID, p_date DATE)
RETURNS VOID AS $$
BEGIN
    INSERT INTO daily_metrics (
        user_id, date,
        interactions_count,
        unique_persons_count,
        avg_energy_score,
        reflection_completed,
        nudges_generated,
        nudges_acted_on,
        positive_interactions,
        negative_interactions,
        relationships_active
    )
    SELECT 
        p_user_id,
        p_date,
        COUNT(DISTINCT i.id),
        COUNT(DISTINCT i.person_id),
        AVG(CASE 
            WHEN i.energy_impact = 'energizing' THEN 100
            WHEN i.energy_impact = 'neutral' THEN 50
            WHEN i.energy_impact = 'draining' THEN 0
            ELSE 50
        END),
        EXISTS(SELECT 1 FROM reflections WHERE user_id = p_user_id AND DATE(completed_at) = p_date),
        (SELECT COUNT(*) FROM nudges WHERE user_id = p_user_id AND DATE(created_at) = p_date),
        (SELECT COUNT(*) FROM nudges WHERE user_id = p_user_id AND DATE(acted_at) = p_date),
        COUNT(CASE WHEN i.energy_impact = 'energizing' THEN 1 END),
        COUNT(CASE WHEN i.energy_impact = 'draining' THEN 1 END),
        (SELECT COUNT(DISTINCT person_id) FROM interactions 
         WHERE user_id = p_user_id 
         AND interaction_at >= p_date - INTERVAL '30 days'
         AND interaction_at <= p_date)
    FROM interactions i
    WHERE i.user_id = p_user_id
    AND DATE(i.interaction_at) = p_date
    GROUP BY p_user_id, p_date
    ON CONFLICT (user_id, date) DO UPDATE SET
        interactions_count = EXCLUDED.interactions_count,
        unique_persons_count = EXCLUDED.unique_persons_count,
        avg_energy_score = EXCLUDED.avg_energy_score,
        reflection_completed = EXCLUDED.reflection_completed,
        nudges_generated = EXCLUDED.nudges_generated,
        nudges_acted_on = EXCLUDED.nudges_acted_on,
        positive_interactions = EXCLUDED.positive_interactions,
        negative_interactions = EXCLUDED.negative_interactions,
        relationships_active = EXCLUDED.relationships_active,
        updated_at = NOW();
END;
$$ LANGUAGE plpgsql;