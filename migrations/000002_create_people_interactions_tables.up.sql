-- Dictionary tables (normalized single-choice values)
CREATE TABLE energy_patterns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    color VARCHAR(20),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE communication_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    icon VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE relationship_statuses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    color VARCHAR(20),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE intentions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    color VARCHAR(20),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Per-user categories
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(20),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT uq_categories_user_name UNIQUE (user_id, name)
);

-- Create people table (relationships managed by users)
CREATE TABLE people (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    relationship VARCHAR(100),
    avatar_url TEXT,
    health_score FLOAT DEFAULT 50.0,
    energy_pattern_id UUID REFERENCES energy_patterns(id) ON DELETE SET NULL,
    last_interaction_at TIMESTAMP,
    interaction_count INTEGER DEFAULT 0,
    communication_method_id UUID REFERENCES communication_methods(id) ON DELETE SET NULL,
    relationship_status_id UUID REFERENCES relationship_statuses(id) ON DELETE SET NULL,
    intention_id UUID REFERENCES intentions(id) ON DELETE SET NULL,
    context TEXT[],
    notes TEXT, -- Will be encrypted if feature is enabled
    custom_fields JSONB DEFAULT '{}',
    reminder_frequency VARCHAR(20),
    next_reminder_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_people_user_id ON people(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_people_category_id ON people(category_id);
CREATE INDEX idx_people_health_score ON people(health_score);
CREATE INDEX idx_people_energy_pattern_id ON people(energy_pattern_id);
CREATE INDEX idx_people_comm_method_id ON people(communication_method_id);
CREATE INDEX idx_people_relationship_status_id ON people(relationship_status_id);
CREATE INDEX idx_people_last_interaction ON people(last_interaction_at);
CREATE INDEX idx_people_next_reminder ON people(next_reminder_at);

-- Create interactions table (vyves)
CREATE TABLE interactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    person_id UUID NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    energy_impact VARCHAR(20) NOT NULL, -- energizing, neutral, draining
    context TEXT[],
    duration INTEGER, -- in minutes
    quality INTEGER CHECK (quality >= 1 AND quality <= 5),
    notes TEXT, -- Will be encrypted if feature is enabled
    location VARCHAR(255),
    special_tags TEXT[],
    interaction_at TIMESTAMP NOT NULL DEFAULT NOW(),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_interactions_user_id ON interactions(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_interactions_person_id ON interactions(person_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_interactions_energy_impact ON interactions(energy_impact);
CREATE INDEX idx_interactions_interaction_at ON interactions(interaction_at);
CREATE INDEX idx_interactions_user_date ON interactions(user_id, interaction_at);

-- Create reflections table
CREATE TABLE reflections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    prompt TEXT,
    responses TEXT[], -- Will be encrypted if feature is enabled
    mood VARCHAR(50),
    energy_level INTEGER CHECK (energy_level >= 1 AND energy_level <= 10),
    insights TEXT[],
    intentions TEXT[],
    gratitude TEXT[],
    metadata JSONB DEFAULT '{}',
    completed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_reflections_user_id ON reflections(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_reflections_completed_at ON reflections(completed_at);
CREATE INDEX idx_reflections_mood ON reflections(mood);
CREATE INDEX idx_reflections_user_date ON reflections(user_id, completed_at);

-- Create nudges table
CREATE TABLE nudges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    person_id UUID REFERENCES people(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- pattern, reconnect, boundary, energy, achievement
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    priority VARCHAR(20) DEFAULT 'medium', -- high, medium, low
    action VARCHAR(100),
    action_data JSONB,
    seen BOOLEAN DEFAULT FALSE,
    acted_on BOOLEAN DEFAULT FALSE,
    seen_at TIMESTAMP,
    acted_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_nudges_user_id ON nudges(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_nudges_person_id ON nudges(person_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_nudges_type ON nudges(type);
CREATE INDEX idx_nudges_priority ON nudges(priority);
CREATE INDEX idx_nudges_seen ON nudges(seen);
CREATE INDEX idx_nudges_acted_on ON nudges(acted_on);
CREATE INDEX idx_nudges_expires_at ON nudges(expires_at);

-- Add triggers for updated_at
CREATE TRIGGER update_people_updated_at BEFORE UPDATE ON people
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_interactions_updated_at BEFORE UPDATE ON interactions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reflections_updated_at BEFORE UPDATE ON reflections
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_nudges_updated_at BEFORE UPDATE ON nudges
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to update person's last interaction timestamp
CREATE OR REPLACE FUNCTION update_person_last_interaction()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE people 
    SET 
        last_interaction_at = NEW.interaction_at,
        interaction_count = interaction_count + 1
    WHERE id = NEW.person_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update person's last interaction
CREATE TRIGGER update_person_last_interaction_trigger
AFTER INSERT ON interactions
FOR EACH ROW EXECUTE FUNCTION update_person_last_interaction();

-- Function to update user's streak
CREATE OR REPLACE FUNCTION update_user_streak()
RETURNS TRIGGER AS $$
DECLARE
    last_reflection_date DATE;
    current_date DATE;
BEGIN
    SELECT DATE(last_reflection_at) INTO last_reflection_date
    FROM users WHERE id = NEW.user_id;
    
    current_date := DATE(NEW.completed_at);
    
    IF last_reflection_date IS NULL OR current_date - last_reflection_date = 1 THEN
        -- Continue or start streak
        UPDATE users 
        SET 
            streak_count = CASE 
                WHEN last_reflection_date IS NULL THEN 1
                ELSE streak_count + 1
            END,
            last_reflection_at = NEW.completed_at
        WHERE id = NEW.user_id;
    ELSIF current_date > last_reflection_date THEN
        -- Break streak
        UPDATE users 
        SET 
            streak_count = 1,
            last_reflection_at = NEW.completed_at
        WHERE id = NEW.user_id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update user's streak on reflection
CREATE TRIGGER update_user_streak_trigger
AFTER INSERT ON reflections
FOR EACH ROW EXECUTE FUNCTION update_user_streak();