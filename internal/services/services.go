package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vyve/vyve-backend/internal/config"
	"github.com/vyve/vyve-backend/internal/models"
	"github.com/vyve/vyve-backend/internal/repository"
	"github.com/vyve/vyve-backend/pkg/analytics"
	"github.com/vyve/vyve-backend/pkg/notifications"
)

// ReflectionService handles reflection business logic
type ReflectionService interface {
	Create(ctx context.Context, userID uuid.UUID, req CreateReflectionRequest) (*models.Reflection, error)
	GetToday(ctx context.Context, userID uuid.UUID) (*models.Reflection, error)
}

type reflectionService struct {
	reflectionRepo repository.ReflectionRepository
}

// NewReflectionService creates a new reflection service
func NewReflectionService(reflectionRepo repository.ReflectionRepository) ReflectionService {
	return &reflectionService{
		reflectionRepo: reflectionRepo,
	}
}

type CreateReflectionRequest struct {
	Prompt     string   `json:"prompt"`
	Responses  []string `json:"responses"`
	Mood       string   `json:"mood"`
	Insights   []string `json:"insights"`
	Intentions []string `json:"intentions"`
}

func (s *reflectionService) Create(ctx context.Context, userID uuid.UUID, req CreateReflectionRequest) (*models.Reflection, error) {
	reflection := &models.Reflection{
		UserID:     userID,
		Prompt:     req.Prompt,
		Responses:  req.Responses,
		Mood:       req.Mood,
		Insights:   req.Insights,
		Intentions: req.Intentions,
	}

	if err := s.reflectionRepo.Create(ctx, reflection); err != nil {
		return nil, err
	}

	return reflection, nil
}

func (s *reflectionService) GetToday(ctx context.Context, userID uuid.UUID) (*models.Reflection, error) {
	return s.reflectionRepo.GetToday(ctx, userID)
}

// NudgeService handles nudge business logic
type NudgeService interface {
	// List and retrieve
	List(ctx context.Context, userID uuid.UUID, opts ListOptions) ([]*models.Nudge, *repository.PaginationResult, error)
	GetByID(ctx context.Context, userID, nudgeID uuid.UUID) (*models.Nudge, error)
	GetActive(ctx context.Context, userID uuid.UUID) ([]*models.Nudge, error)
	GetHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*models.Nudge, error)
	
	// Actions
	MarkSeen(ctx context.Context, userID, nudgeID uuid.UUID) error
	MarkActedOn(ctx context.Context, userID, nudgeID uuid.UUID) error
	Dismiss(ctx context.Context, userID, nudgeID uuid.UUID) error
	
	// Generation
	GenerateNudges(ctx context.Context, userID uuid.UUID) error
	GenerateForPerson(ctx context.Context, userID, personID uuid.UUID) ([]*models.Nudge, error)
	GenerateSystemNudges(ctx context.Context, userID uuid.UUID) ([]*models.Nudge, error)
}

// Note: NudgeService implementation is in nudge_service.go

// GDPRService handles GDPR compliance
type GDPRService interface {
	ExportUserData(ctx context.Context, userID uuid.UUID) (*models.DataExport, error)
	DeleteAllUserData(ctx context.Context, userID uuid.UUID) error
	AnonymizeUserData(ctx context.Context, userID uuid.UUID) error
	RecordConsent(ctx context.Context, userID uuid.UUID, consentType string, granted bool) error
	GetConsents(ctx context.Context, userID uuid.UUID) ([]*models.UserConsent, error)
	GetLatestExport(ctx context.Context, userID uuid.UUID) (*models.DataExport, error)
	GetExport(ctx context.Context, userID uuid.UUID, exportID uuid.UUID) (*models.DataExport, error)
	GetAuditLog(ctx context.Context, userID uuid.UUID) ([]*models.AuditLog, error)
}

type gdprService struct {
	repos      *repository.Repositories
	encryption bool
}

// NewGDPRService creates a new GDPR service
func NewGDPRService(repos *repository.Repositories, encryptionConfig config.EncryptionConfig) GDPRService {
	return &gdprService{
		repos:      repos,
		encryption: encryptionConfig.Enabled,
	}
}

func (s *gdprService) ExportUserData(ctx context.Context, userID uuid.UUID) (*models.DataExport, error) {
	// Check for existing pending export
	existing, err := s.repos.DataExport.GetPending(ctx, userID)
	if err == nil && existing != nil {
		return existing, nil
	}

	// Create export request
	export := &models.DataExport{
		UserID: userID,
		Status: "pending",
		Format: "json",
	}

	if err := s.repos.DataExport.Create(ctx, export); err != nil {
		return nil, err
	}

	// Process export asynchronously
	go s.processDataExport(context.Background(), export)

	return export, nil
}

func (s *gdprService) DeleteAllUserData(ctx context.Context, userID uuid.UUID) error {
	// Note: This is a simplified implementation
	// In a full implementation, you would need to add DeleteByUser methods
	// to each repository or iterate and delete records individually

	// For now, just delete the user (which should cascade due to foreign keys)
	// Or implement soft delete and mark all user data as deleted

	// TODO: Implement proper cascade deletion:
	// - Delete AI analysis jobs
	// - Delete relationship analyses and recommendations
	// - Delete nudges
	// - Delete interactions
	// - Delete people
	// - Delete reflections
	// - Delete events
	// - Delete consents
	// - Delete audit logs
	// - Delete data exports
	// - Delete refresh tokens

	// Finally, delete user
	return s.repos.User.Delete(ctx, userID)
}

func (s *gdprService) AnonymizeUserData(ctx context.Context, userID uuid.UUID) error {
	// Get user
	user, err := s.repos.User.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Anonymize user data
	user.Username = "user_" + userID.String()[:8] + "_anonymized"
	user.Email = "anonymized_" + userID.String() + "@deleted.local"
	user.DisplayName = "Anonymized User"
	user.Bio = ""
	user.AvatarURL = ""
	user.PasswordHash = "" // Clear password

	if err := s.repos.User.Update(ctx, user); err != nil {
		return err
	}

	// Record anonymization in audit log
	// Note: Repository method may not exist yet

	return nil
}

func (s *gdprService) RecordConsent(ctx context.Context, userID uuid.UUID, consentType string, granted bool) error {
	consent := &models.UserConsent{
		UserID:      userID,
		ConsentType: consentType,
		Granted:     granted,
		GrantedAt:   time.Now(),
	}

	// Check if consent already exists
	existing, err := s.repos.Consent.GetByType(ctx, userID, consentType)
	if err == nil && existing != nil {
		// Update existing consent
		existing.Granted = granted
		if !granted {
			now := time.Now()
			existing.RevokedAt = &now
		} else {
			existing.RevokedAt = nil
			existing.GrantedAt = time.Now()
		}
		return s.repos.Consent.Update(ctx, existing)
	}

	return s.repos.Consent.Create(ctx, consent)
}

// Helper method to process data export asynchronously
func (s *gdprService) processDataExport(ctx context.Context, export *models.DataExport) {
	// Update status to processing
	export.Status = "processing"
	s.repos.DataExport.Update(ctx, export)

	// Gather all user data
	data := make(map[string]interface{})

	// Get user profile
	user, err := s.repos.User.FindByID(ctx, export.UserID)
	if err == nil {
		data["user"] = user
	}

	// Get people
	people, err := s.repos.Person.FindByUserID(ctx, export.UserID)
	if err == nil {
		data["people"] = people
	}

	// Get interactions
	interactions, err := s.repos.Interaction.GetRecent(ctx, export.UserID, 10000)
	if err == nil {
		data["interactions"] = interactions
	}

	// Get consents
	consents, _ := s.repos.Consent.GetByUser(ctx, export.UserID)
	data["consents"] = consents

	// In a real implementation, you would:
	// 1. Serialize data to JSON/CSV
	// 2. Upload to storage (S3, etc.)
	// 3. Update export with file URL and size
	// 4. Set expiration time

	// For now, mark as completed
	now := time.Now()
	expiresAt := now.Add(7 * 24 * time.Hour) // Expires in 7 days
	export.Status = "completed"
	export.CompletedAt = &now
	export.ExpiresAt = &expiresAt
	export.FileURL = "/api/v1/gdpr/export/" + export.ID.String() + "/download"

	s.repos.DataExport.Update(ctx, export)
}

func (s *gdprService) GetConsents(ctx context.Context, userID uuid.UUID) ([]*models.UserConsent, error) {
	return s.repos.Consent.GetByUser(ctx, userID)
}

func (s *gdprService) GetLatestExport(ctx context.Context, userID uuid.UUID) (*models.DataExport, error) {
	return s.repos.DataExport.GetPending(ctx, userID)
}

func (s *gdprService) GetExport(ctx context.Context, userID uuid.UUID, exportID uuid.UUID) (*models.DataExport, error) {
	export, err := s.repos.DataExport.FindByID(ctx, exportID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if export.UserID != userID {
		return nil, repository.ErrForbidden
	}

	return export, nil
}

func (s *gdprService) GetAuditLog(ctx context.Context, userID uuid.UUID) ([]*models.AuditLog, error) {
	return s.repos.AuditLog.GetByUser(ctx, userID)
}

// Background worker functions

// SendDailyReminders sends daily reminder notifications
func SendDailyReminders(userRepo repository.UserRepository, notificationService notifications.NotificationService) {
	// Stub implementation
}

// GenerateNudges generates AI-powered nudges for users
func GenerateNudges(repos *repository.Repositories, analytics analytics.Analytics, notificationService notifications.NotificationService) {
	// Stub implementation
}

// AggregateMetrics aggregates daily metrics
func AggregateMetrics(repos *repository.Repositories, analytics analytics.Analytics) {
	// Stub implementation
}