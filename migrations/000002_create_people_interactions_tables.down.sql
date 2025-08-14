-- Drop triggers
DROP TRIGGER IF EXISTS update_user_streak_trigger ON reflections;
DROP TRIGGER IF EXISTS update_person_last_interaction_trigger ON interactions;
DROP TRIGGER IF EXISTS update_nudges_updated_at ON nudges;
DROP TRIGGER IF EXISTS update_reflections_updated_at ON reflections;
DROP TRIGGER IF EXISTS update_interactions_updated_at ON interactions;
DROP TRIGGER IF EXISTS update_people_updated_at ON people;

-- Drop functions
DROP FUNCTION IF EXISTS update_user_streak();
DROP FUNCTION IF EXISTS update_person_last_interaction();

-- Drop tables
DROP TABLE IF EXISTS nudges;
DROP TABLE IF EXISTS reflections;
DROP TABLE IF EXISTS interactions;
DROP TABLE IF EXISTS people;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS intentions;
DROP TABLE IF EXISTS relationship_statuses;
DROP TABLE IF EXISTS communication_methods;
DROP TABLE IF EXISTS energy_patterns;