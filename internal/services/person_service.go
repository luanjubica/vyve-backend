package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vyve/vyve-backend/internal/models"
	"github.com/vyve/vyve-backend/internal/repository"
)

// PersonService handles person/relationship business logic
type PersonService interface {
	Create(ctx context.Context, userID uuid.UUID, req CreatePersonRequest) (*models.Person, error)
	GetByID(ctx context.Context, userID, personID uuid.UUID) (*models.Person, error)
	List(ctx context.Context, userID uuid.UUID, opts ListOptions) ([]*models.Person, *repository.PaginationResult, error)
	Update(ctx context.Context, userID, personID uuid.UUID, updates map[string]interface{}) (*models.Person, error)
	Delete(ctx context.Context, userID, personID uuid.UUID) error
	Restore(ctx context.Context, userID, personID uuid.UUID) error
	UpdateHealthScore(ctx context.Context, userID, personID uuid.UUID) error
	GetCategories(ctx context.Context, userID uuid.UUID) ([]string, error)
	Search(ctx context.Context, userID uuid.UUID, query string) ([]*models.Person, error)
}

type personService struct {
	personRepo repository.PersonRepository
}

// NewPersonService creates a new person service
func NewPersonService(personRepo repository.PersonRepository) PersonService {
	return &personService{
		personRepo: personRepo,
	}
}

// CreatePersonRequest represents a request to create a person
type CreatePersonRequest struct {
	Name                string   `json:"name" validate:"required"`
	Category            string   `json:"category"`
	Relationship        string   `json:"relationship"`
	CommunicationMethod string   `json:"communication_method"`
	RelationshipStatus  string   `json:"relationship_status"`
	Intention           string   `json:"intention"`
	Context             []string `json:"context"`
	Notes               string   `json:"notes"`
}

// ListOptions represents listing options
type ListOptions struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Category string `json:"category"`
	Search   string `json:"search"`
	OrderBy  string `json:"order_by"`
}

// Create creates a new person
func (s *personService) Create(ctx context.Context, userID uuid.UUID, req CreatePersonRequest) (*models.Person, error) {
	person := &models.Person{
		UserID:              userID,
		Name:                req.Name,
		Category:            req.Category,
		Relationship:        req.Relationship,
		CommunicationMethod: req.CommunicationMethod,
		RelationshipStatus:  req.RelationshipStatus,
		Intention:           req.Intention,
		Context:             req.Context,
		Notes:               req.Notes,
		HealthScore:         50.0, // Default health score
	}

	if err := s.personRepo.Create(ctx, person); err != nil {
		return nil, err
	}

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
	if category, ok := updates["category"].(string); ok {
		person.Category = category
	}
	if relationship, ok := updates["relationship"].(string); ok {
		person.Relationship = relationship
	}
	if notes, ok := updates["notes"].(string); ok {
		person.Notes = notes
	}

	if err := s.personRepo.Update(ctx, person); err != nil {
		return nil, err
	}

	return person, nil
}

// Delete deletes a person
func (s *personService) Delete(ctx context.Context, userID, personID uuid.UUID) error {
	person, err := s.GetByID(ctx, userID, personID)
	if err != nil {
		return err
	}

	return s.personRepo.Delete(ctx, person.ID)
}

// Restore restores a deleted person
func (s *personService) Restore(ctx context.Context, userID, personID uuid.UUID) error {
	// TODO: Implement restore logic
	return errors.New("not implemented")
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

// Helper function removed; calculateHealthScore is defined in interaction_service.go