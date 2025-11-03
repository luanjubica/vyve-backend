package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vyve/vyve-backend/internal/models"
	"github.com/vyve/vyve-backend/internal/repository"
	"github.com/vyve/vyve-backend/pkg/analytics"
	"github.com/vyve/vyve-backend/pkg/notifications"
)

type nudgeServiceImpl struct {
	nudgeRepo     repository.NudgeRepository
	personRepo    repository.PersonRepository
	notifications notifications.NotificationService
	analytics     analytics.Analytics
}

// NewNudgeService creates a new nudge service
func NewNudgeService(nudgeRepo repository.NudgeRepository, notifications notifications.NotificationService, analyticsService analytics.Analytics) NudgeService {
	return &nudgeServiceImpl{
		nudgeRepo:     nudgeRepo,
		notifications: notifications,
		analytics:     analyticsService,
	}
}

// List gets nudges with filtering and pagination
func (s *nudgeServiceImpl) List(ctx context.Context, userID uuid.UUID, opts ListOptions) ([]*models.Nudge, *repository.PaginationResult, error) {
	filters := repository.FilterOptions{
		UserID:  userID,
		Page:    opts.Page,
		Limit:   opts.Limit,
		OrderBy: "created_at",
		Desc:    true,
	}

	// Apply filters from opts
	if opts.Status != "" {
		filters.Status = opts.Status
	}
	if opts.Source != "" {
		filters.Type = opts.Source // Using Type field for source filter
	}

	return s.nudgeRepo.List(ctx, filters)
}

// GetByID gets a specific nudge by ID
func (s *nudgeServiceImpl) GetByID(ctx context.Context, userID, nudgeID uuid.UUID) (*models.Nudge, error) {
	nudge, err := s.nudgeRepo.FindByID(ctx, nudgeID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if nudge.UserID != userID {
		return nil, repository.ErrForbidden
	}

	return nudge, nil
}

// GetActive gets active nudges for a user
func (s *nudgeServiceImpl) GetActive(ctx context.Context, userID uuid.UUID) ([]*models.Nudge, error) {
	return s.nudgeRepo.GetActive(ctx, userID)
}

// GetHistory gets nudge history for a user
func (s *nudgeServiceImpl) GetHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*models.Nudge, error) {
	filters := repository.FilterOptions{
		UserID:  userID,
		Page:    1,
		Limit:   limit,
		OrderBy: "created_at",
		Desc:    true,
		Status:  "completed", // Can only filter by one status at a time
	}

	nudges, _, err := s.nudgeRepo.List(ctx, filters)
	return nudges, err
}

// MarkSeen marks a nudge as seen
func (s *nudgeServiceImpl) MarkSeen(ctx context.Context, userID, nudgeID uuid.UUID) error {
	// Verify ownership
	nudge, err := s.GetByID(ctx, userID, nudgeID)
	if err != nil {
		return err
	}

	now := time.Now()
	nudge.Seen = true
	nudge.SeenAt = &now
	if nudge.Status == "pending" {
		nudge.Status = "seen"
	}

	err = s.nudgeRepo.Update(ctx, nudge)
	if err != nil {
		return err
	}

	// Track nudge seen event
	go analytics.TrackNudgeSeen(ctx, s.analytics, userID.String(), nudgeID.String(), nudge.Type, nudge.Source)

	return nil
}

// MarkActedOn marks a nudge as acted upon
func (s *nudgeServiceImpl) MarkActedOn(ctx context.Context, userID, nudgeID uuid.UUID) error {
	// Verify ownership
	nudge, err := s.GetByID(ctx, userID, nudgeID)
	if err != nil {
		return err
	}

	now := time.Now()
	nudge.ActedOn = true
	nudge.ActedAt = &now
	nudge.Status = "completed"
	nudge.CompletedAt = &now

	err = s.nudgeRepo.Update(ctx, nudge)
	if err != nil {
		return err
	}

	// Track nudge acted on event
	go analytics.TrackNudgeAction(ctx, s.analytics, userID.String(), nudgeID.String(), nudge.Type)

	return nil
}

// Dismiss dismisses a nudge
func (s *nudgeServiceImpl) Dismiss(ctx context.Context, userID, nudgeID uuid.UUID) error {
	// Verify ownership
	nudge, err := s.GetByID(ctx, userID, nudgeID)
	if err != nil {
		return err
	}

	now := time.Now()
	nudge.Status = "dismissed"
	nudge.DismissedAt = &now

	err = s.nudgeRepo.Update(ctx, nudge)
	if err != nil {
		return err
	}

	// Track nudge dismissed event
	go analytics.TrackNudgeDismissed(ctx, s.analytics, userID.String(), nudgeID.String(), nudge.Type)

	return nil
}

// GenerateForPerson generates AI nudges for a specific person (calls analysis service)
func (s *nudgeServiceImpl) GenerateForPerson(ctx context.Context, userID, personID uuid.UUID) ([]*models.Nudge, error) {
	// This would typically call the AnalysisService.GenerateRecommendations
	// For now, return empty to avoid circular dependency
	// The proper way is to call this from the analysis handler
	return nil, nil
}

// GenerateSystemNudges generates system-based nudges for all relationships
func (s *nudgeServiceImpl) GenerateSystemNudges(ctx context.Context, userID uuid.UUID) ([]*models.Nudge, error) {
	// This would analyze all relationships and create system nudges
	// based on rules like:
	// - No interaction in X days -> "reconnect" nudge
	// - Low quality interactions -> "improve quality" nudge
	// - Unbalanced communication -> "balance" nudge
	
	// For now, return empty as this is a complex feature
	// that requires the person and interaction repositories
	return nil, nil
}

// GenerateNudges is the original interface method (kept for compatibility)
func (s *nudgeServiceImpl) GenerateNudges(ctx context.Context, userID uuid.UUID) error {
	_, err := s.GenerateSystemNudges(ctx, userID)
	return err
}
