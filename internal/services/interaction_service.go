package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vyve/vyve-backend/internal/models"
	"github.com/vyve/vyve-backend/internal/repository"
	"github.com/vyve/vyve-backend/pkg/analytics"
)

// InteractionService handles interaction (vyve) business logic
type InteractionService interface {
	Create(ctx context.Context, userID uuid.UUID, req CreateInteractionRequest) (*models.Interaction, error)
	GetByID(ctx context.Context, userID, interactionID uuid.UUID) (*models.Interaction, error)
	List(ctx context.Context, userID uuid.UUID, opts ListOptions) ([]*models.Interaction, *repository.PaginationResult, error)
	Update(ctx context.Context, userID, interactionID uuid.UUID, updates map[string]interface{}) (*models.Interaction, error)
	Delete(ctx context.Context, userID, interactionID uuid.UUID) error
	GetRecent(ctx context.Context, userID uuid.UUID, limit int) ([]*models.Interaction, error)
	GetByDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]*models.Interaction, error)
	GetEnergyDistribution(ctx context.Context, userID uuid.UUID, period string) (map[string]interface{}, error)
	BulkCreate(ctx context.Context, userID uuid.UUID, interactions []CreateInteractionRequest) ([]*models.Interaction, error)
}

type interactionService struct {
	interactionRepo repository.InteractionRepository
	personRepo      repository.PersonRepository
	analytics       analytics.Analytics
}

// NewInteractionService creates a new interaction service
func NewInteractionService(
	interactionRepo repository.InteractionRepository,
	personRepo repository.PersonRepository,
	analytics analytics.Analytics,
) InteractionService {
	return &interactionService{
		interactionRepo: interactionRepo,
		personRepo:      personRepo,
		analytics:       analytics,
	}
}

// CreateInteractionRequest represents a request to create an interaction
type CreateInteractionRequest struct {
	PersonID     uuid.UUID `json:"person_id" validate:"required"`
	EnergyImpact string    `json:"energy_impact" validate:"required,oneof=energizing neutral draining"`
	Context      []string  `json:"context"`
	Duration     int       `json:"duration"`
	Quality      int       `json:"quality" validate:"min=1,max=5"`
	Notes        string    `json:"notes"`
	Location     string    `json:"location"`
	SpecialTags  []string  `json:"special_tags"`
}

// Create creates a new interaction
func (s *interactionService) Create(ctx context.Context, userID uuid.UUID, req CreateInteractionRequest) (*models.Interaction, error) {
	interaction := &models.Interaction{
		UserID:        userID,
		PersonID:      req.PersonID,
		EnergyImpact:  req.EnergyImpact,
		Context:       req.Context,
		Duration:      req.Duration,
		Quality:       req.Quality,
		Notes:         req.Notes,
		Location:      req.Location,
		SpecialTags:   req.SpecialTags,
		InteractionAt: time.Now(),
	}

	if err := s.interactionRepo.Create(ctx, interaction); err != nil {
		return nil, err
	}

	// Track analytics event
	s.analytics.Track(ctx, analytics.Event{
		UserID:    userID.String(),
		EventType: "interaction_logged",
		Properties: map[string]interface{}{
			"person_id":     req.PersonID.String(),
			"energy_impact": req.EnergyImpact,
			"context":       req.Context,
			"quality":       req.Quality,
			"duration":      req.Duration,
		},
	})

	// Update person's last interaction and health score
	go s.updatePersonMetrics(context.Background(), req.PersonID)

	return interaction, nil
}

// GetByID gets an interaction by ID
func (s *interactionService) GetByID(ctx context.Context, userID, interactionID uuid.UUID) (*models.Interaction, error) {
	interaction, err := s.interactionRepo.FindByID(ctx, interactionID)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if interaction.UserID != userID {
		return nil, repository.ErrForbidden
	}

	return interaction, nil
}

// List lists interactions for a user
func (s *interactionService) List(ctx context.Context, userID uuid.UUID, opts ListOptions) ([]*models.Interaction, *repository.PaginationResult, error) {
	filterOpts := repository.FilterOptions{
		UserID:  userID,
		OrderBy: "interaction_at",
		Desc:    true,
		Page:    opts.Page,
		Limit:   opts.Limit,
	}

	return s.interactionRepo.List(ctx, filterOpts)
}

// Update updates an interaction
func (s *interactionService) Update(ctx context.Context, userID, interactionID uuid.UUID, updates map[string]interface{}) (*models.Interaction, error) {
	interaction, err := s.GetByID(ctx, userID, interactionID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if energyImpact, ok := updates["energy_impact"].(string); ok {
		interaction.EnergyImpact = energyImpact
	}
	if notes, ok := updates["notes"].(string); ok {
		interaction.Notes = notes
	}
	if quality, ok := updates["quality"].(int); ok {
		interaction.Quality = quality
	}
	if duration, ok := updates["duration"].(int); ok {
		interaction.Duration = duration
	}

	if err := s.interactionRepo.Update(ctx, interaction); err != nil {
		return nil, err
	}

	// Update person's health score if energy impact changed
	if _, ok := updates["energy_impact"]; ok {
		go s.updatePersonMetrics(context.Background(), interaction.PersonID)
	}

	return interaction, nil
}

// Delete deletes an interaction
func (s *interactionService) Delete(ctx context.Context, userID, interactionID uuid.UUID) error {
	interaction, err := s.GetByID(ctx, userID, interactionID)
	if err != nil {
		return err
	}

	if err := s.interactionRepo.Delete(ctx, interaction.ID); err != nil {
		return err
	}

	// Update person's health score
	go s.updatePersonMetrics(context.Background(), interaction.PersonID)

	return nil
}

// GetRecent gets recent interactions for a user
func (s *interactionService) GetRecent(ctx context.Context, userID uuid.UUID, limit int) ([]*models.Interaction, error) {
	return s.interactionRepo.GetRecent(ctx, userID, limit)
}

// GetByDate gets interactions for a specific date
func (s *interactionService) GetByDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]*models.Interaction, error) {
	return s.interactionRepo.GetByDateRange(ctx, userID, date, date.AddDate(0, 0, 1))
}

// GetEnergyDistribution gets energy distribution statistics
func (s *interactionService) GetEnergyDistribution(ctx context.Context, userID uuid.UUID, period string) (map[string]interface{}, error) {
	// Calculate date range based on period
	endDate := time.Now()
	var startDate time.Time
	
	switch period {
	case "week":
		startDate = endDate.AddDate(0, 0, -7)
	case "month":
		startDate = endDate.AddDate(0, -1, 0)
	case "quarter":
		startDate = endDate.AddDate(0, -3, 0)
	case "year":
		startDate = endDate.AddDate(-1, 0, 0)
	default:
		startDate = endDate.AddDate(0, 0, -7)
	}

	return s.interactionRepo.GetEnergyDistribution(ctx, userID, startDate, endDate)
}

// BulkCreate creates multiple interactions
func (s *interactionService) BulkCreate(ctx context.Context, userID uuid.UUID, requests []CreateInteractionRequest) ([]*models.Interaction, error) {
	interactions := make([]*models.Interaction, len(requests))
	
	for i, req := range requests {
		interaction, err := s.Create(ctx, userID, req)
		if err != nil {
			return nil, err
		}
		interactions[i] = interaction
	}

	return interactions, nil
}

// updatePersonMetrics updates person's health score and last interaction
func (s *interactionService) updatePersonMetrics(ctx context.Context, personID uuid.UUID) {
	// Update last interaction timestamp
	if err := s.personRepo.UpdateLastInteraction(ctx, personID); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	// Increment interaction count
	if err := s.personRepo.IncrementInteractionCount(ctx, personID); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	// Get recent interactions for health score calculation
	interactions, err := s.personRepo.GetRecentInteractions(ctx, personID, 10)
	if err != nil {
		// Log error but don't fail the operation
		return
	}

	// Calculate and update health score
	score := calculateHealthScore(interactions)
	if err := s.personRepo.UpdateHealthScore(ctx, personID, score); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}
}

// Helper function to calculate health score
func calculateHealthScore(interactions []*models.Interaction) float64 {
	if len(interactions) == 0 {
		return 50.0
	}

	weights := map[string]float64{
		"energizing": 100.0,
		"neutral":    50.0,
		"draining":   0.0,
	}

	totalScore := 0.0
	totalWeight := 0.0

	for i, interaction := range interactions {
		weight := 1.0 / float64(i+1) // Recent interactions have more weight
		score := weights[interaction.EnergyImpact]
		totalScore += score * weight
		totalWeight += weight
	}

	if totalWeight > 0 {
		return totalScore / totalWeight
	}

	return 50.0
}