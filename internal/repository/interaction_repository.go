package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/vyve/vyve-backend/internal/models"
)

// InteractionRepository defines interaction data access interface
type InteractionRepository interface {
	Create(ctx context.Context, interaction *models.Interaction) error
	Update(ctx context.Context, interaction *models.Interaction) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Interaction, error)
	List(ctx context.Context, opts FilterOptions) ([]*models.Interaction, *PaginationResult, error)
	GetRecent(ctx context.Context, userID uuid.UUID, limit int) ([]*models.Interaction, error)
	GetByPerson(ctx context.Context, personID uuid.UUID) ([]*models.Interaction, error)
	GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*models.Interaction, error)
	GetEnergyDistribution(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error)
	GetDailyCount(ctx context.Context, userID uuid.UUID, date time.Time) (int64, error)
	GetAverageQuality(ctx context.Context, userID uuid.UUID) (float64, error)
	BulkCreate(ctx context.Context, interactions []*models.Interaction) error
}

type interactionRepository struct {
	BaseRepository
}

// NewInteractionRepository creates a new interaction repository
func NewInteractionRepository(db *gorm.DB) InteractionRepository {
	return &interactionRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create creates a new interaction
func (r *interactionRepository) Create(ctx context.Context, interaction *models.Interaction) error {
	return r.db.WithContext(ctx).Create(interaction).Error
}

// Update updates an interaction
func (r *interactionRepository) Update(ctx context.Context, interaction *models.Interaction) error {
	return r.db.WithContext(ctx).Save(interaction).Error
}

// Delete soft deletes an interaction
func (r *interactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Interaction{}, "id = ?", id).Error
}

// FindByID finds an interaction by ID
func (r *interactionRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Interaction, error) {
	var interaction models.Interaction
	err := r.db.WithContext(ctx).
		Preload("Person").
		First(&interaction, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInteractionNotFound
		}
		return nil, err
	}
	return &interaction, nil
}

// List lists interactions with pagination
func (r *interactionRepository) List(ctx context.Context, opts FilterOptions) ([]*models.Interaction, *PaginationResult, error) {
	query := r.db.WithContext(ctx).Model(&models.Interaction{}).Preload("Person")
	
	// Apply filters
	if opts.UserID != uuid.Nil {
		query = query.Where("user_id = ?", opts.UserID)
	}
	
	if opts.StartDate != nil {
		query = query.Where("interaction_at >= ?", opts.StartDate)
	}
	
	if opts.EndDate != nil {
		query = query.Where("interaction_at <= ?", opts.EndDate)
	}
	
	// Apply ordering
	if opts.OrderBy != "" {
		if opts.Desc {
			query = query.Order(opts.OrderBy + " DESC")
		} else {
			query = query.Order(opts.OrderBy + " ASC")
		}
	} else {
		query = query.Order("interaction_at DESC")
	}
	
	var interactions []*models.Interaction
	result, err := Paginate(ctx, query, opts.Page, opts.Limit, &interactions)
	if err != nil {
		return nil, nil, err
	}
	
	return interactions, result, nil
}

// GetRecent gets recent interactions for a user
func (r *interactionRepository) GetRecent(ctx context.Context, userID uuid.UUID, limit int) ([]*models.Interaction, error) {
	var interactions []*models.Interaction
	err := r.db.WithContext(ctx).
		Preload("Person").
		Where("user_id = ?", userID).
		Order("interaction_at DESC").
		Limit(limit).
		Find(&interactions).Error
	return interactions, err
}

// GetByPerson gets interactions for a specific person
func (r *interactionRepository) GetByPerson(ctx context.Context, personID uuid.UUID) ([]*models.Interaction, error) {
	var interactions []*models.Interaction
	err := r.db.WithContext(ctx).
		Where("person_id = ?", personID).
		Order("interaction_at DESC").
		Find(&interactions).Error
	return interactions, err
}

// GetByDateRange gets interactions within a date range
func (r *interactionRepository) GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*models.Interaction, error) {
	var interactions []*models.Interaction
	err := r.db.WithContext(ctx).
		Preload("Person").
		Where("user_id = ? AND interaction_at >= ? AND interaction_at <= ?", 
			userID, startDate, endDate).
		Order("interaction_at DESC").
		Find(&interactions).Error
	return interactions, err
}

// GetEnergyDistribution gets energy distribution statistics
func (r *interactionRepository) GetEnergyDistribution(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error) {
	var results []struct {
		EnergyImpact string
		Count        int
		Percentage   float64
	}
	
	// Get counts by energy impact
	err := r.db.WithContext(ctx).
		Model(&models.Interaction{}).
		Select("energy_impact, COUNT(*) as count").
		Where("user_id = ? AND interaction_at >= ? AND interaction_at <= ?", 
			userID, startDate, endDate).
		Group("energy_impact").
		Scan(&results).Error
	
	if err != nil {
		return nil, err
	}
	
	// Calculate total and percentages
	total := 0
	for _, r := range results {
		total += r.Count
	}
	
	distribution := make(map[string]interface{})
	for _, r := range results {
		if total > 0 {
			r.Percentage = float64(r.Count) / float64(total) * 100
		}
		distribution[r.EnergyImpact] = map[string]interface{}{
			"count":      r.Count,
			"percentage": r.Percentage,
		}
	}
	
	distribution["total"] = total
	return distribution, nil
}

// GetDailyCount gets the count of interactions for a specific day
func (r *interactionRepository) GetDailyCount(ctx context.Context, userID uuid.UUID, date time.Time) (int64, error) {
	var count int64
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)
	
	err := r.db.WithContext(ctx).
		Model(&models.Interaction{}).
		Where("user_id = ? AND interaction_at >= ? AND interaction_at < ?", 
			userID, startOfDay, endOfDay).
		Count(&count).Error
	
	return count, err
}

// GetAverageQuality gets the average quality score
func (r *interactionRepository) GetAverageQuality(ctx context.Context, userID uuid.UUID) (float64, error) {
	var avg float64
	err := r.db.WithContext(ctx).
		Model(&models.Interaction{}).
		Where("user_id = ? AND quality > 0", userID).
		Select("AVG(quality)").
		Scan(&avg).Error
	
	return avg, err
}

// BulkCreate creates multiple interactions
func (r *interactionRepository) BulkCreate(ctx context.Context, interactions []*models.Interaction) error {
	return r.db.WithContext(ctx).CreateInBatches(interactions, 100).Error
}