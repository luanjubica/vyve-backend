-- Drop triggers
DROP TRIGGER IF EXISTS update_push_tokens_updated_at ON push_tokens;
DROP TRIGGER IF EXISTS update_refresh_tokens_updated_at ON refresh_tokens;
DROP TRIGGER IF EXISTS update_auth_providers_updated_at ON auth_providers;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop tables
DROP TABLE IF EXISTS push_tokens;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS auth_providers;
DROP TABLE IF EXISTS users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();