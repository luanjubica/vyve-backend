package notifications

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"

	"github.com/vyve/vyve-backend/internal/config"
)

// NotificationService defines the notification service interface
type NotificationService interface {
	SendPushNotification(ctx context.Context, token string, notification Notification) error
	SendBatchNotifications(ctx context.Context, tokens []string, notification Notification) error
	SendTopicNotification(ctx context.Context, topic string, notification Notification) error
	SubscribeToTopic(ctx context.Context, tokens []string, topic string) error
	UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) error
}

// Notification represents a push notification
type Notification struct {
	Title    string                 `json:"title"`
	Body     string                 `json:"body"`
	Data     map[string]string      `json:"data,omitempty"`
	ImageURL string                 `json:"image_url,omitempty"`
	Priority string                 `json:"priority,omitempty"` // high, normal
	Badge    int                    `json:"badge,omitempty"`
	Sound    string                 `json:"sound,omitempty"`
}

// FCMService implements NotificationService using Firebase Cloud Messaging
type FCMService struct {
	client *messaging.Client
}

// NewFCMService creates a new FCM service
func NewFCMService(cfg config.FCMConfig) (NotificationService, error) {
	ctx := context.Background()

	// Determine credentials source: file path (preferred) or raw JSON
	var opt option.ClientOption
	if path := os.Getenv("FCM_CREDENTIALS_FILE"); path != "" {
		opt = option.WithCredentialsFile(path)
	} else if path := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); path != "" {
		opt = option.WithCredentialsFile(path)
	} else if cfg.Key != "" {
		opt = option.WithCredentialsJSON([]byte(cfg.Key))
	} else {
		return nil, fmt.Errorf("no FCM credentials provided: set FCM_CREDENTIALS_FILE, GOOGLE_APPLICATION_CREDENTIALS, or FCM_KEY")
	}

	// Initialize Firebase app
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: cfg.ProjectID,
	}, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	// Get messaging client
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get messaging client: %w", err)
	}

	return &FCMService{
		client: client,
	}, nil
}

// SendPushNotification sends a push notification to a single device
func (f *FCMService) SendPushNotification(ctx context.Context, token string, notification Notification) error {
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title:    notification.Title,
			Body:     notification.Body,
			ImageURL: notification.ImageURL,
		},
		Data: notification.Data,
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Badge: &notification.Badge,
					Sound: notification.Sound,
				},
			},
		},
		Android: &messaging.AndroidConfig{
			Priority: f.getPriority(notification.Priority),
			Notification: &messaging.AndroidNotification{
				Sound: notification.Sound,
			},
		},
	}

	response, err := f.client.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	log.Printf("Successfully sent notification: %s", response)
	return nil
}

// SendBatchNotifications sends notifications to multiple devices
func (f *FCMService) SendBatchNotifications(ctx context.Context, tokens []string, notification Notification) error {
	messages := []*messaging.Message{}
	
	for _, token := range tokens {
		message := &messaging.Message{
			Token: token,
			Notification: &messaging.Notification{
				Title:    notification.Title,
				Body:     notification.Body,
				ImageURL: notification.ImageURL,
			},
			Data: notification.Data,
			APNS: &messaging.APNSConfig{
				Payload: &messaging.APNSPayload{
					Aps: &messaging.Aps{
						Badge: &notification.Badge,
						Sound: notification.Sound,
					},
				},
			},
			Android: &messaging.AndroidConfig{
				Priority: f.getPriority(notification.Priority),
				Notification: &messaging.AndroidNotification{
					Sound: notification.Sound,
				},
			},
		}
		messages = append(messages, message)
	}

	// Send batch (max 500 messages per batch)
	batchSize := 500
	for i := 0; i < len(messages); i += batchSize {
		end := i + batchSize
		if end > len(messages) {
			end = len(messages)
		}

		batch := messages[i:end]
		response, err := f.client.SendAll(ctx, batch)
		if err != nil {
			return fmt.Errorf("failed to send batch notifications: %w", err)
		}

		log.Printf("Batch response: %d success, %d failure", response.SuccessCount, response.FailureCount)
		
		// Handle failed tokens
		if response.FailureCount > 0 {
			for idx, resp := range response.Responses {
				if !resp.Success {
					log.Printf("Failed to send to token %s: %v", tokens[i+idx], resp.Error)
				}
			}
		}
	}

	return nil
}

// SendTopicNotification sends a notification to a topic
func (f *FCMService) SendTopicNotification(ctx context.Context, topic string, notification Notification) error {
	message := &messaging.Message{
		Topic: topic,
		Notification: &messaging.Notification{
			Title:    notification.Title,
			Body:     notification.Body,
			ImageURL: notification.ImageURL,
		},
		Data: notification.Data,
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Badge: &notification.Badge,
					Sound: notification.Sound,
				},
			},
		},
		Android: &messaging.AndroidConfig{
			Priority: f.getPriority(notification.Priority),
			Notification: &messaging.AndroidNotification{
				Sound: notification.Sound,
			},
		},
	}

	response, err := f.client.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send topic notification: %w", err)
	}

	log.Printf("Successfully sent topic notification: %s", response)
	return nil
}

// SubscribeToTopic subscribes tokens to a topic
func (f *FCMService) SubscribeToTopic(ctx context.Context, tokens []string, topic string) error {
	response, err := f.client.SubscribeToTopic(ctx, tokens, topic)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	log.Printf("Subscribe response: %d success, %d failure", response.SuccessCount, response.FailureCount)
	return nil
}

// UnsubscribeFromTopic unsubscribes tokens from a topic
func (f *FCMService) UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) error {
	response, err := f.client.UnsubscribeFromTopic(ctx, tokens, topic)
	if err != nil {
		return fmt.Errorf("failed to unsubscribe from topic: %w", err)
	}

	log.Printf("Unsubscribe response: %d success, %d failure", response.SuccessCount, response.FailureCount)
	return nil
}

// getPriority converts priority string to FCM priority
func (f *FCMService) getPriority(priority string) string {
	if priority == "high" {
		return "high"
	}
	return "normal"
}

// MockNotificationService is a mock implementation for testing
type MockNotificationService struct{}

// NewMockNotificationService creates a new mock notification service
func NewMockNotificationService() NotificationService {
	return &MockNotificationService{}
}

// SendPushNotification mock implementation
func (m *MockNotificationService) SendPushNotification(ctx context.Context, token string, notification Notification) error {
	log.Printf("Mock: Sending notification to %s: %s - %s", token, notification.Title, notification.Body)
	return nil
}

// SendBatchNotifications mock implementation
func (m *MockNotificationService) SendBatchNotifications(ctx context.Context, tokens []string, notification Notification) error {
	log.Printf("Mock: Sending batch notification to %d devices: %s - %s", len(tokens), notification.Title, notification.Body)
	return nil
}

// SendTopicNotification mock implementation
func (m *MockNotificationService) SendTopicNotification(ctx context.Context, topic string, notification Notification) error {
	log.Printf("Mock: Sending topic notification to %s: %s - %s", topic, notification.Title, notification.Body)
	return nil
}

// SubscribeToTopic mock implementation
func (m *MockNotificationService) SubscribeToTopic(ctx context.Context, tokens []string, topic string) error {
	log.Printf("Mock: Subscribing %d tokens to topic %s", len(tokens), topic)
	return nil
}

// UnsubscribeFromTopic mock implementation
func (m *MockNotificationService) UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) error {
	log.Printf("Mock: Unsubscribing %d tokens from topic %s", len(tokens), topic)
	return nil
}