package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base model with common fields
type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// EnergyPattern dictionary (e.g., energizing, neutral, draining, mixed)
type EnergyPattern struct {
    Base
    Name  string `gorm:"not null;uniqueIndex" json:"name"`
    Color string `json:"color"`
}

// CommunicationMethod dictionary (e.g., whatsapp, call, text, email)
type CommunicationMethod struct {
    Base
    Name  string `gorm:"not null;uniqueIndex" json:"name"`
    Icon  string `json:"icon"`
}

// RelationshipStatus dictionary (e.g., new, stable, fading, tense)
type RelationshipStatus struct {
    Base
    Name  string `gorm:"not null;uniqueIndex" json:"name"`
    Color string `json:"color"`
}

// Intention dictionary (e.g., audit, improve, maintain, boundaries)
type Intention struct {
    Base
    Name  string `gorm:"not null;uniqueIndex" json:"name"`
    Color string `json:"color"`
}

// Category represents a normalized category for people (e.g., friend, family, colleague)
type Category struct {
    Base
    UserID uuid.UUID `gorm:"not null;index" json:"user_id"` // scope categories per user
    Name   string    `gorm:"not null;index:idx_user_name,unique" json:"name"`
    Color  string    `json:"color"`

    // Relations
    User    User     `gorm:"foreignKey:UserID" json:"-"`
    Persons []Person `gorm:"foreignKey:CategoryID" json:"-"`
}

// BeforeCreate hook to set UUID
func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// User represents a user in the system
type User struct {
	Base
	Username         string           `gorm:"uniqueIndex;not null" json:"username"`
	Email            string           `gorm:"uniqueIndex;not null" json:"email"`
	EmailVerified    bool             `gorm:"default:false" json:"email_verified"`
	PasswordHash     string           `json:"-"`
	AvatarURL        string           `json:"avatar_url"`
	DisplayName      string           `json:"display_name"`
	Bio              string           `json:"bio"`
	Timezone         string           `gorm:"default:'UTC'" json:"timezone"`
	Locale           string           `gorm:"default:'en'" json:"locale"`
	LastLoginAt      *time.Time       `json:"last_login_at"`
	LastActivityAt   *time.Time       `json:"last_activity_at"`
	StreakCount      int              `gorm:"default:0" json:"streak_count"`
	LastReflectionAt *time.Time       `json:"last_reflection_at"`
	Settings         JSONB            `gorm:"type:jsonb" json:"settings"`
	Metadata         JSONB            `gorm:"type:jsonb" json:"metadata"`
	DataResidency    string           `gorm:"default:'us'" json:"data_residency"` // us, eu, etc.
	
	// Add these onboarding fields:
	OnboardingCompleted bool  `gorm:"default:false" json:"onboarding_completed"`
	OnboardingSteps     JSONB `gorm:"type:jsonb" json:"onboarding_steps"`

	// Relations
	AuthProviders []AuthProvider `json:"-"`
	People        []Person       `json:"-"`
	Interactions  []Interaction  `json:"-"`
	Reflections   []Reflection   `json:"-"`
	Nudges        []Nudge        `json:"-"`
	Events        []Event        `json:"-"`
	PushTokens    []PushToken    `json:"-"`
	Consents      []UserConsent  `json:"-"`
	AuditLogs     []AuditLog     `json:"-"`
}

// AuthProvider represents an OAuth provider linked to a user
type AuthProvider struct {
	Base
	UserID       uuid.UUID `gorm:"not null" json:"user_id"`
	Provider     string    `gorm:"not null" json:"provider"` // google, linkedin, apple
	ProviderID   string    `gorm:"not null" json:"provider_id"`
	AccessToken  string    `json:"-"`
	RefreshToken string    `json:"-"`
	ExpiresAt    *time.Time `json:"expires_at"`
	RawData      JSONB     `gorm:"type:jsonb" json:"-"`
	
	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// Person represents a person/relationship managed by a user
type Person struct {
    Base
    UserID               uuid.UUID  `gorm:"not null;index" json:"user_id"`
    Name                 string     `gorm:"not null" json:"name"`
    CategoryID           *uuid.UUID `gorm:"index" json:"category_id,omitempty"`
    CategoryRef          *Category  `gorm:"foreignKey:CategoryID" json:"-"`
    Relationship         string     `json:"relationship"`
    AvatarURL            string     `json:"avatar_url"`
    HealthScore          float64    `gorm:"default:50" json:"health_score"`
    EnergyPatternID      *uuid.UUID `gorm:"index" json:"energy_pattern_id,omitempty"`
    EnergyPatternRef     *EnergyPattern `gorm:"foreignKey:EnergyPatternID" json:"-"`
    LastInteractionAt    *time.Time `json:"last_interaction_at"`
    InteractionCount     int        `gorm:"default:0" json:"interaction_count"`
    CommunicationMethodID *uuid.UUID `gorm:"index" json:"communication_method_id,omitempty"`
    CommunicationMethodRef *CommunicationMethod `gorm:"foreignKey:CommunicationMethodID" json:"-"`
    RelationshipStatusID *uuid.UUID `gorm:"index" json:"relationship_status_id,omitempty"`
    RelationshipStatusRef *RelationshipStatus `gorm:"foreignKey:RelationshipStatusID" json:"-"`
    IntentionID          *uuid.UUID `gorm:"index" json:"intention_id,omitempty"`
    IntentionRef         *Intention `gorm:"foreignKey:IntentionID" json:"-"`
    Context              StringArray `gorm:"type:text[]" json:"context"` // work, personal, community, etc.
    Notes                string     `gorm:"encrypted" json:"notes"` // Will be encrypted if enabled
    CustomFields         JSONB      `gorm:"type:jsonb" json:"custom_fields"`
    ReminderFrequency    string     `json:"reminder_frequency"` // daily, weekly, monthly, custom
    NextReminderAt       *time.Time `json:"next_reminder_at"`
	
	// Relations
	User         User          `gorm:"foreignKey:UserID" json:"-"`
	Interactions []Interaction `json:"-"`
}

// Interaction represents a "vyve" - an interaction with a person
type Interaction struct {
	Base
	UserID        uuid.UUID   `gorm:"not null;index" json:"user_id"`
	PersonID      uuid.UUID   `gorm:"not null;index" json:"person_id"`
	EnergyImpact  string      `gorm:"not null" json:"energy_impact"` // energizing, neutral, draining
	Context       StringArray `gorm:"type:text[]" json:"context"`
	Duration      int         `json:"duration"` // in minutes
	Quality       int         `json:"quality"`  // 1-5 scale
	Notes         string      `gorm:"encrypted" json:"notes"`
	Location      string      `json:"location"`
	SpecialTags   StringArray `gorm:"type:text[]" json:"special_tags"`
	InteractionAt time.Time   `gorm:"not null;default:now()" json:"interaction_at"`
	Metadata      JSONB       `gorm:"type:jsonb" json:"metadata"`
	
	// Relations
	User   User   `gorm:"foreignKey:UserID" json:"-"`
	Person Person `gorm:"foreignKey:PersonID" json:"-"`
}

// Reflection represents a daily reflection entry
type Reflection struct {
	Base
	UserID      uuid.UUID   `gorm:"not null;index" json:"user_id"`
	Prompt      string      `json:"prompt"`
	Responses   StringArray `gorm:"type:text[];encrypted" json:"responses"`
	Mood        string      `json:"mood"`
	EnergyLevel int         `json:"energy_level"` // 1-10 scale
	Insights    StringArray `gorm:"type:text[]" json:"insights"`
	Intentions  StringArray `gorm:"type:text[]" json:"intentions"`
	Gratitude   StringArray `gorm:"type:text[]" json:"gratitude"`
	Metadata    JSONB       `gorm:"type:jsonb" json:"metadata"`
	CompletedAt time.Time   `gorm:"not null;default:now()" json:"completed_at"`
	
	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// Nudge represents an AI-generated insight or reminder
type Nudge struct {
	Base
	UserID     uuid.UUID  `gorm:"not null;index" json:"user_id"`
	PersonID   *uuid.UUID `gorm:"index" json:"person_id,omitempty"`
	Type       string     `gorm:"not null" json:"type"` // pattern, reconnect, boundary, energy, achievement
	Title      string     `gorm:"not null" json:"title"`
	Message    string     `gorm:"not null" json:"message"`
	Priority   string     `gorm:"default:'medium'" json:"priority"` // high, medium, low
	Action     string     `json:"action"`
	ActionData JSONB      `gorm:"type:jsonb" json:"action_data"`
	Seen       bool       `gorm:"default:false" json:"seen"`
	ActedOn    bool       `gorm:"default:false" json:"acted_on"`
	SeenAt     *time.Time `json:"seen_at"`
	ActedAt    *time.Time `json:"acted_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	
	// Relations
	User   User    `gorm:"foreignKey:UserID" json:"-"`
	Person *Person `gorm:"foreignKey:PersonID" json:"-"`
}

// Event represents an analytics event
type Event struct {
	Base
	UserID     uuid.UUID `gorm:"not null;index" json:"user_id"`
	EventType  string    `gorm:"not null;index" json:"event_type"`
	Properties JSONB     `gorm:"type:jsonb" json:"properties"`
	SessionID  string    `gorm:"index" json:"session_id"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	
	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// DailyMetric represents aggregated daily metrics for a user
type DailyMetric struct {
	UserID                uuid.UUID `gorm:"primaryKey" json:"user_id"`
	Date                  time.Time `gorm:"primaryKey;type:date" json:"date"`
	InteractionsCount     int       `json:"interactions_count"`
	UniquePersonsCount    int       `json:"unique_persons_count"`
	AvgEnergyScore        float64   `json:"avg_energy_score"`
	ReflectionCompleted   bool      `json:"reflection_completed"`
	NudgesGenerated       int       `json:"nudges_generated"`
	NudgesActedOn         int       `json:"nudges_acted_on"`
	PositiveInteractions  int       `json:"positive_interactions"`
	NegativeInteractions  int       `json:"negative_interactions"`
	RelationshipsActive   int       `json:"relationships_active"`
	RelationshipsImproved int       `json:"relationships_improved"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// PushToken represents a device push notification token
type PushToken struct {
	Base
	UserID     uuid.UUID `gorm:"not null;index" json:"user_id"`
	Token      string    `gorm:"not null;uniqueIndex" json:"token"`
	Platform   string    `gorm:"not null" json:"platform"` // ios, android
	DeviceID   string    `json:"device_id"`
	DeviceInfo JSONB     `gorm:"type:jsonb" json:"device_info"`
	Active     bool      `gorm:"default:true" json:"active"`
	LastUsedAt *time.Time `json:"last_used_at"`
	
	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// RefreshToken represents a refresh token for JWT authentication
type RefreshToken struct {
	Base
	UserID    uuid.UUID `gorm:"not null;index" json:"user_id"`
	Token     string    `gorm:"not null;uniqueIndex" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Revoked   bool      `gorm:"default:false" json:"revoked"`
	RevokedAt *time.Time `json:"revoked_at"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	
	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// UserConsent represents GDPR consent records
type UserConsent struct {
	Base
	UserID       uuid.UUID `gorm:"not null;index" json:"user_id"`
	ConsentType  string    `gorm:"not null" json:"consent_type"` // marketing, analytics, cookies, etc.
	Granted      bool      `gorm:"not null" json:"granted"`
	Version      string    `json:"version"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	GrantedAt    time.Time `json:"granted_at"`
	RevokedAt    *time.Time `json:"revoked_at"`
	
	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// AuditLog represents an audit trail for GDPR compliance
type AuditLog struct {
	Base
	UserID       *uuid.UUID `gorm:"index" json:"user_id,omitempty"`
	Action       string     `gorm:"not null" json:"action"`
	EntityType   string     `json:"entity_type"`
	EntityID     string     `json:"entity_id"`
	Changes      JSONB      `gorm:"type:jsonb" json:"changes"`
	IPAddress    string     `json:"ip_address"`
	UserAgent    string     `json:"user_agent"`
	RequestID    string     `json:"request_id"`
	SessionID    string     `json:"session_id"`
	Result       string     `json:"result"` // success, failure
	ErrorMessage string     `json:"error_message,omitempty"`
	
	// Relations
	User *User `gorm:"foreignKey:UserID" json:"-"`
}

// DataExport represents a GDPR data export request
type DataExport struct {
	Base
	UserID      uuid.UUID  `gorm:"not null;index" json:"user_id"`
	Status      string     `gorm:"not null;default:'pending'" json:"status"` // pending, processing, completed, failed
	Format      string     `gorm:"default:'json'" json:"format"` // json, csv
	FileURL     string     `json:"file_url,omitempty"`
	FileSize    int64      `json:"file_size,omitempty"`
	RequestedAt time.Time  `gorm:"not null;default:now()" json:"requested_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	Error       string     `json:"error,omitempty"`
	
	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// Custom types for PostgreSQL arrays and JSONB

// StringArray is a custom type for PostgreSQL text arrays
type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "{}", nil
	}
	return s, nil
}

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}
	// Handle PostgreSQL array format
	// This is a simplified version - you might need a more robust parser
	*s = value.([]string)
	return nil
}

// JSONB is a custom type for PostgreSQL JSONB fields
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	
	return json.Unmarshal(bytes, j)
}