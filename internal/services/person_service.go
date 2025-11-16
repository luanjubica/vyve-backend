package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/vyve/vyve-backend/internal/models"
	"github.com/vyve/vyve-backend/internal/repository"
	"github.com/vyve/vyve-backend/pkg/analytics"
	"github.com/vyve/vyve-backend/pkg/storage"
)

// PersonService handles person/relationship business logic
type PersonService interface {
	// Core CRUD operations
	Create(ctx context.Context, userID uuid.UUID, req CreatePersonRequest) (*models.Person, error)
	GetByID(ctx context.Context, userID, personID uuid.UUID) (*models.Person, error)
	List(ctx context.Context, userID uuid.UUID, opts ListOptions) ([]*models.Person, *repository.PaginationResult, error)
	Update(ctx context.Context, userID, personID uuid.UUID, updates map[string]interface{}) (*models.Person, error)
	Delete(ctx context.Context, userID, personID uuid.UUID) error
	Restore(ctx context.Context, userID, personID uuid.UUID) error

	// Extended operations
	UpdateHealthScore(ctx context.Context, userID, personID uuid.UUID) error
	GetCategories(ctx context.Context, userID uuid.UUID) ([]string, error)
	Search(ctx context.Context, userID uuid.UUID, query string) ([]*models.Person, error)
	UploadAvatar(ctx context.Context, userID, personID uuid.UUID, fileData []byte, contentType string) (*models.Person, error)

	// People operations
	CountPeople(ctx context.Context, userID uuid.UUID) (int64, error)
}

type personService struct {
	personRepo repository.PersonRepository
	analytics  analytics.Analytics
	storage    storage.Storage
}

// NewPersonService creates a new person service
func NewPersonService(personRepo repository.PersonRepository, analyticsService analytics.Analytics, storageService storage.Storage) PersonService {
	return &personService{
		personRepo: personRepo,
		analytics:  analyticsService,
		storage:    storageService,
	}
}

// CreatePersonRequest represents a request to create a person
type CreatePersonRequest struct {
	Name                  string     `json:"name" validate:"required"`
	CategoryID            *uuid.UUID `json:"category_id"`
	Relationship          string     `json:"relationship"`
	AvatarURL             string     `json:"avatar_url"`
	CommunicationMethodID *uuid.UUID `json:"communication_method_id"`
	RelationshipStatusID  *uuid.UUID `json:"relationship_status_id"`
	IntentionID           *uuid.UUID `json:"intention_id"`
	Context               []string   `json:"context"`
	Notes                 string     `json:"notes"`
}

// ListOptions represents listing options
type ListOptions struct {
	Page     int         `json:"page"`
	Limit    int         `json:"limit"`
	Category string      `json:"category"`
	Search   string      `json:"search"`
	OrderBy  string      `json:"order_by"`
	Source   string      `json:"source"`    // For nudges: 'ai' or 'system'
	Status   string      `json:"status"`    // For nudges: status filter
	PersonID *uuid.UUID  `json:"person_id"` // For nudges: filter by person
}

// Create creates a new person
func (s *personService) Create(ctx context.Context, userID uuid.UUID, req CreatePersonRequest) (*models.Person, error) {
	person := &models.Person{
		UserID:                userID,
		Name:                  req.Name,
		Relationship:          req.Relationship,
		AvatarURL:             req.AvatarURL,
		CategoryID:            req.CategoryID,
		CommunicationMethodID: req.CommunicationMethodID,
		RelationshipStatusID:  req.RelationshipStatusID,
		IntentionID:           req.IntentionID,
		Context:               req.Context,
		Notes:                 req.Notes,
		HealthScore:           50.0, // Default health score
	}

	if err := s.personRepo.Create(ctx, person); err != nil {
		return nil, err
	}

	// Track person added event
	go analytics.TrackPersonAdded(ctx, s.analytics, userID.String(), person.ID.String(), req.Relationship)

	return person, nil
}

// GetByID gets a person by ID
func (s *personService) GetByID(ctx context.Context, userID, personID uuid.UUID) (*models.Person, error) {
	person, err := s.personRepo.FindByID(ctx, personID)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if person.UserID != userID {
		return nil, repository.ErrForbidden
	}

	return person, nil
}

// List lists people for a user
func (s *personService) List(ctx context.Context, userID uuid.UUID, opts ListOptions) ([]*models.Person, *repository.PaginationResult, error) {
	filterOpts := repository.FilterOptions{
		UserID:   userID,
		Category: opts.Category,
		Search:   opts.Search,
		OrderBy:  opts.OrderBy,
		Page:     opts.Page,
		Limit:    opts.Limit,
	}

	return s.personRepo.List(ctx, filterOpts)
}

// Update updates a person
func (s *personService) Update(ctx context.Context, userID, personID uuid.UUID, updates map[string]interface{}) (*models.Person, error) {
	person, err := s.GetByID(ctx, userID, personID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	// TODO: Add validation and proper field mapping
	if name, ok := updates["name"].(string); ok {
		person.Name = name
	}
	if v, ok := updates["category_id"].(string); ok {
		if id, err := uuid.Parse(v); err == nil {
			person.CategoryID = &id
		}
	}
	if relationship, ok := updates["relationship"].(string); ok {
		person.Relationship = relationship
	}
	if v, ok := updates["communication_method_id"].(string); ok {
		if id, err := uuid.Parse(v); err == nil {
			person.CommunicationMethodID = &id
		}
	}
	if v, ok := updates["relationship_status_id"].(string); ok {
		if id, err := uuid.Parse(v); err == nil {
			person.RelationshipStatusID = &id
		}
	}
	if v, ok := updates["intention_id"].(string); ok {
		if id, err := uuid.Parse(v); err == nil {
			person.IntentionID = &id
		}
	}
	if notes, ok := updates["notes"].(string); ok {
		person.Notes = notes
	}

	if err := s.personRepo.Update(ctx, person); err != nil {
		return nil, err
	}

	// Track person updated event
	go analytics.TrackPersonUpdated(ctx, s.analytics, userID.String(), personID.String(), updates)

	return person, nil
}

// Delete deletes a person
func (s *personService) Delete(ctx context.Context, userID, personID uuid.UUID) error {
	person, err := s.GetByID(ctx, userID, personID)
	if err != nil {
		return err
	}

	err = s.personRepo.Delete(ctx, person.ID)
	if err != nil {
		return err
	}

	// Track person deleted event
	go analytics.TrackPersonDeleted(ctx, s.analytics, userID.String(), personID.String())

	return nil
}

// Restore restores a deleted person
func (s *personService) Restore(ctx context.Context, userID, personID uuid.UUID) error {
	// TODO: Implement proper soft-delete restoration
	// This requires adding a Restore method to the PersonRepository interface
	// For now, return a helpful error
	return errors.New("restore functionality requires database-level soft-delete restoration - not yet implemented")
}

// CountPeople returns the total number of people for a user
func (s *personService) CountPeople(ctx context.Context, userID uuid.UUID) (int64, error) {
	log.Printf("🎯 PersonService.CountPeople called with userID: %s", userID)

	people, err := s.personRepo.FindByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return int64(len(people)), nil
}

// UpdateHealthScore updates a person's health score
func (s *personService) UpdateHealthScore(ctx context.Context, userID, personID uuid.UUID) error {
	person, err := s.GetByID(ctx, userID, personID)
	if err != nil {
		return err
	}

	// Get recent interactions
	interactions, err := s.personRepo.GetRecentInteractions(ctx, personID, 10)
	if err != nil {
		return err
	}

	// Calculate health score
	score := calculateHealthScore(interactions)
	person.HealthScore = score

	return s.personRepo.Update(ctx, person)
}

// GetCategories gets unique categories for a user
func (s *personService) GetCategories(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return s.personRepo.GetCategories(ctx, userID)
}

// Search searches people for a user
func (s *personService) Search(ctx context.Context, userID uuid.UUID, query string) ([]*models.Person, error) {
	return s.personRepo.Search(ctx, userID, query, 20)
}

// UploadAvatar uploads an avatar for a person and returns the updated person
func (s *personService) UploadAvatar(ctx context.Context, userID, personID uuid.UUID, fileData []byte, contentType string) (*models.Person, error) {
	// Verify person belongs to user
	person, err := s.personRepo.FindByID(ctx, personID)
	if err != nil {
		return nil, err
	}

	if person.UserID != userID {
		return nil, repository.ErrForbidden
	}

	// Check if storage is available
	if s.storage == nil {
		return nil, errors.New("storage service not available")
	}

	// Determine file extension from content type
	extension := ".jpg"
	contentTypeLower := strings.ToLower(contentType)
	if strings.Contains(contentTypeLower, "png") {
		extension = ".png"
	} else if strings.Contains(contentTypeLower, "gif") {
		extension = ".gif"
	} else if strings.Contains(contentTypeLower, "webp") {
		extension = ".webp"
	}

	// Generate storage key: avatars/people/{userID}/{personID}/{uuid}{extension}
	fileID := uuid.New().String()
	key := fmt.Sprintf("avatars/people/%s/%s/%s%s", userID.String(), personID.String(), fileID, extension)

	// Upload to storage
	avatarURL, err := s.storage.Upload(ctx, key, fileData, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload avatar: %w", err)
	}

	// Update person with new avatar URL
	updates := map[string]interface{}{
		"avatar_url": avatarURL,
	}

	updatedPerson, err := s.Update(ctx, userID, personID, updates)
	if err != nil {
		// Try to delete uploaded file if database update fails
		_ = s.storage.Delete(ctx, key)
		return nil, fmt.Errorf("failed to update person: %w", err)
	}

	return updatedPerson, nil
}

// Helper function removed; calculateHealthScore is defined in interaction_service.go
