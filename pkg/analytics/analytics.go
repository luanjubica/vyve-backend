package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amplitude "github.com/amplitude/analytics-go/amplitude"
	"github.com/google/uuid"
)

// Analytics defines the analytics interface
type Analytics interface {
	Track(ctx context.Context, event Event) error
	TrackBatch(ctx context.Context, events []Event) error
	Identify(ctx context.Context, userID string, properties map[string]interface{}) error
	GroupIdentify(ctx context.Context, groupType, groupValue string, properties map[string]interface{}) error
	Flush() error
	Close() error
}

// Event represents an analytics event
type Event struct {
	UserID       string                 `json:"user_id"`
	EventType    string                 `json:"event_type"`
	Properties   map[string]interface{} `json:"properties,omitempty"`
	SessionID    string                 `json:"session_id,omitempty"`
	DeviceID     string                 `json:"device_id,omitempty"`
	Platform     string                 `json:"platform,omitempty"`
	Version      string                 `json:"version,omitempty"`
	Timestamp    time.Time              `json:"timestamp"`
	LocationLat  float64                `json:"location_lat,omitempty"`
	LocationLng  float64                `json:"location_lng,omitempty"`
	IP           string                 `json:"ip,omitempty"`
	Revenue      float64                `json:"revenue,omitempty"`
	RevenueType  string                 `json:"revenue_type,omitempty"`
}

// EventTypes
const (
	EventUserSignUp            = "user_sign_up"
	EventUserLogin             = "user_login"
	EventUserLogout            = "user_logout"
	EventInteractionLogged     = "interaction_logged"
	EventReflectionCompleted   = "reflection_completed"
	EventNudgeGenerated        = "nudge_generated"
	EventNudgeActedOn          = "nudge_acted_on"
	EventPersonAdded           = "person_added"
	EventPersonUpdated         = "person_updated"
	EventHealthScoreChanged    = "health_score_changed"
	EventStreakUpdated         = "streak_updated"
	EventNotificationSent      = "notification_sent"
	EventNotificationOpened    = "notification_opened"
	EventDataExported          = "data_exported"
	EventSessionStarted        = "session_started"
	EventSessionEnded          = "session_ended"
)

// AmplitudeAnalytics implements Analytics using Amplitude
type AmplitudeAnalytics struct {
	client amplitude.Client
}

// NewAmplitudeAnalytics creates a new Amplitude analytics service
func NewAmplitudeAnalytics(apiKey string) Analytics {
	config := amplitude.NewConfig(apiKey)
	client := amplitude.NewClient(config)
	
	return &AmplitudeAnalytics{
		client: client,
	}
}

// Track tracks a single event
func (a *AmplitudeAnalytics) Track(ctx context.Context, event Event) error {
	options := amplitude.EventOptions{
		DeviceID:    event.DeviceID,
		Platform:    event.Platform,
		AppVersion:  event.Version,
		Time:        event.Timestamp.Unix(),
		IP:          event.IP,
		SessionID:   a.parseSessionID(event.SessionID),
	}

	// Add location if provided
	if event.LocationLat != 0 || event.LocationLng != 0 {
		options.LocationLat = event.LocationLat
		options.LocationLng = event.LocationLng
	}

	// Add revenue if provided
	if event.Revenue > 0 {
		options.Price = event.Revenue
		options.RevenueType = event.RevenueType
	}

	amplitudeEvent := amplitude.Event{
		UserID:          event.UserID,
		EventType:       event.EventType,
		EventProperties: event.Properties,
		EventOptions:    options,
	}

	a.client.Track(amplitudeEvent)
	
	log.Printf("Tracked event: %s for user: %s", event.EventType, event.UserID)
	return nil
}

// TrackBatch tracks multiple events
func (a *AmplitudeAnalytics) TrackBatch(ctx context.Context, events []Event) error {
	for _, event := range events {
		if err := a.Track(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// Identify identifies a user with properties
func (a *AmplitudeAnalytics) Identify(ctx context.Context, userID string, properties map[string]interface{}) error {
	identify := amplitude.Identify{}
	for key, value := range properties {
		identify.Set(key, value)
	}

	a.client.Identify(identify, amplitude.EventOptions{UserID: userID})
	
	log.Printf("Identified user: %s", userID)
	return nil
}

// GroupIdentify identifies a group with properties
func (a *AmplitudeAnalytics) GroupIdentify(ctx context.Context, groupType, groupValue string, properties map[string]interface{}) error {
	identify := amplitude.Identify{}
	for key, value := range properties {
		identify.Set(key, value)
	}

	a.client.GroupIdentify(groupType, groupValue, identify, amplitude.EventOptions{})
	
	log.Printf("Group identified: %s = %s", groupType, groupValue)
	return nil
}

// Flush flushes pending events
func (a *AmplitudeAnalytics) Flush() error {
	// Amplitude client handles this automatically
	return nil
}

// Close closes the analytics client
func (a *AmplitudeAnalytics) Close() error {
	// Amplitude client handles cleanup automatically
	return nil
}

// parseSessionID converts string session ID to int
func (a *AmplitudeAnalytics) parseSessionID(sessionID string) int {
	if sessionID == "" {
		return int(time.Now().Unix())
	}
	// Try to parse as UUID and use timestamp portion
	if _, err := uuid.Parse(sessionID); err == nil {
		return int(time.Now().Unix())
	}
	return int(time.Now().Unix())
}

// DatabaseAnalytics implements Analytics by storing events in database
type DatabaseAnalytics struct {
	// TODO: Add database repository
}

// NewDatabaseAnalytics creates a new database analytics service
func NewDatabaseAnalytics() Analytics {
	return &DatabaseAnalytics{}
}

// Track tracks a single event
func (d *DatabaseAnalytics) Track(ctx context.Context, event Event) error {
	// Store event in database
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	
	log.Printf("Database Analytics: %s", string(eventJSON))
	// TODO: Store in database
	return nil
}

// TrackBatch tracks multiple events
func (d *DatabaseAnalytics) TrackBatch(ctx context.Context, events []Event) error {
	for _, event := range events {
		if err := d.Track(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// Identify identifies a user with properties
func (d *DatabaseAnalytics) Identify(ctx context.Context, userID string, properties map[string]interface{}) error {
	log.Printf("Database Analytics - User identified: %s", userID)
	// TODO: Store in database
	return nil
}

// GroupIdentify identifies a group with properties
func (d *DatabaseAnalytics) GroupIdentify(ctx context.Context, groupType, groupValue string, properties map[string]interface{}) error {
	log.Printf("Database Analytics - Group identified: %s = %s", groupType, groupValue)
	// TODO: Store in database
	return nil
}

// Flush flushes pending events
func (d *DatabaseAnalytics) Flush() error {
	// No-op for database analytics
	return nil
}

// Close closes the analytics client
func (d *DatabaseAnalytics) Close() error {
	// No-op for database analytics
	return nil
}

// Helper functions for common analytics operations

// TrackInteraction tracks an interaction event
func TrackInteraction(ctx context.Context, analytics Analytics, userID, personID, energyImpact string, quality int) error {
	return analytics.Track(ctx, Event{
		UserID:    userID,
		EventType: EventInteractionLogged,
		Properties: map[string]interface{}{
			"person_id":     personID,
			"energy_impact": energyImpact,
			"quality":       quality,
		},
		Timestamp: time.Now(),
	})
}

// TrackReflection tracks a reflection completion event
func TrackReflection(ctx context.Context, analytics Analytics, userID, mood string, streakCount int) error {
	return analytics.Track(ctx, Event{
		UserID:    userID,
		EventType: EventReflectionCompleted,
		Properties: map[string]interface{}{
			"mood":         mood,
			"streak_count": streakCount,
		},
		Timestamp: time.Now(),
	})
}

// TrackNudgeAction tracks when a user acts on a nudge
func TrackNudgeAction(ctx context.Context, analytics Analytics, userID, nudgeID, nudgeType string) error {
	return analytics.Track(ctx, Event{
		UserID:    userID,
		EventType: EventNudgeActedOn,
		Properties: map[string]interface{}{
			"nudge_id":   nudgeID,
			"nudge_type": nudgeType,
		},
		Timestamp: time.Now(),
	})
}

// TrackHealthScoreChange tracks health score changes
func TrackHealthScoreChange(ctx context.Context, analytics Analytics, userID, personID string, oldScore, newScore float64) error {
	return analytics.Track(ctx, Event{
		UserID:    userID,
		EventType: EventHealthScoreChanged,
		Properties: map[string]interface{}{
			"person_id": personID,
			"old_score": oldScore,
			"new_score": newScore,
			"change":    newScore - oldScore,
		},
		Timestamp: time.Now(),
	})
}