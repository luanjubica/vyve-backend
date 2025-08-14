-- Drop functions
DROP FUNCTION IF EXISTS aggregate_daily_metrics(UUID, DATE);
DROP FUNCTION IF EXISTS log_audit_event(UUID, VARCHAR, VARCHAR, VARCHAR, JSONB, VARCHAR, TEXT, VARCHAR, VARCHAR, VARCHAR, TEXT);

-- Drop triggers
DROP TRIGGER IF EXISTS update_data_exports_updated_at ON data_exports;
DROP TRIGGER IF EXISTS update_audit_logs_updated_at ON audit_logs;
DROP TRIGGER IF EXISTS update_user_consents_updated_at ON user_consents;
DROP TRIGGER IF EXISTS update_daily_metrics_updated_at ON daily_metrics;
DROP TRIGGER IF EXISTS update_events_updated_at ON events;

-- Drop tables
DROP TABLE IF EXISTS data_exports;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS user_consents;
DROP TABLE IF EXISTS daily_metrics;
DROP TABLE IF EXISTS events;