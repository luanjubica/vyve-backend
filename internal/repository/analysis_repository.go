package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vyve/vyve-backend/internal/models"
	"gorm.io/gorm"
)

// AnalysisRepository handles relationship analysis data access
type AnalysisRepository interface {
	// RelationshipAnalysis operations
	CreateAnalysis(ctx context.Context, analysis *models.RelationshipAnalysis) error
	GetLatestAnalysis(ctx context.Context, userID, personID uuid.UUID) (*models.RelationshipAnalysis, error)
	GetAnalysisByID(ctx context.Context, analysisID uuid.UUID) (*models.RelationshipAnalysis, error)
	ListAnalyses(ctx context.Context, userID uuid.UUID, opts FilterOptions) ([]*models.RelationshipAnalysis, *PaginationResult, error)
	GetAnalysisHistory(ctx context.Context, userID, personID uuid.UUID, limit int) ([]*models.RelationshipAnalysis, error)
	
	// Nudge operations (AI recommendations merged into Nudge)
	CreateRecommendation(ctx context.Context, recommendation *models.Nudge) error
	GetRecommendationByID(ctx context.Context, recommendationID uuid.UUID) (*models.Nudge, error)
	ListRecommendations(ctx context.Context, userID uuid.UUID, opts FilterOptions) ([]*models.Nudge, *PaginationResult, error)
	GetActiveRecommendations(ctx context.Context, userID uuid.UUID) ([]*models.Nudge, error)
	GetRecommendationsForPerson(ctx context.Context, userID, personID uuid.UUID) ([]*models.Nudge, error)
	UpdateRecommendationStatus(ctx context.Context, recommendationID uuid.UUID, status string) error
	
	// AIAnalysisJob operations
	CreateJob(ctx context.Context, job *models.AIAnalysisJob) error
	GetJobByID(ctx context.Context, jobID uuid.UUID) (*models.AIAnalysisJob, error)
	UpdateJob(ctx context.Context, job *models.AIAnalysisJob) error
	ListPendingJobs(ctx context.Context, limit int) ([]*models.AIAnalysisJob, error)
	GetUserJobs(ctx context.Context, userID uuid.UUID, limit int) ([]*models.AIAnalysisJob, error)
}

type analysisRepository struct {
	db *gorm.DB
}

// NewAnalysisRepository creates a new analysis repository
func NewAnalysisRepository(db *gorm.DB) AnalysisRepository {
	return &analysisRepository{db: db}
}

// CreateAnalysis creates a new relationship analysis
func (r *analysisRepository) CreateAnalysis(ctx context.Context, analysis *models.RelationshipAnalysis) error {
	return r.db.WithContext(ctx).Create(analysis).Error
}

// GetLatestAnalysis gets the most recent analysis for a person
func (r *analysisRepository) GetLatestAnalysis(ctx context.Context, userID, personID uuid.UUID) (*models.RelationshipAnalysis, error) {
	var analysis models.RelationshipAnalysis
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND person_id = ?", userID, personID).
		Order("analyzed_at DESC").
		First(&analysis).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	
	return &analysis, nil
}

// GetAnalysisByID gets an analysis by ID
func (r *analysisRepository) GetAnalysisByID(ctx context.Context, analysisID uuid.UUID) (*models.RelationshipAnalysis, error) {
	var analysis models.RelationshipAnalysis
	err := r.db.WithContext(ctx).
		Preload("Recommendations").
		First(&analysis, "id = ?", analysisID).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	
	return &analysis, nil
}

// ListAnalyses lists analyses with filtering and pagination
func (r *analysisRepository) ListAnalyses(ctx context.Context, userID uuid.UUID, opts FilterOptions) ([]*models.RelationshipAnalysis, *PaginationResult, error) {
	var analyses []*models.RelationshipAnalysis
	var total int64
	
	query := r.db.WithContext(ctx).Model(&models.RelationshipAnalysis{}).Where("user_id = ?", userID)
	
	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}
	
	// Apply pagination
	offset := (opts.Page - 1) * opts.Limit
	query = query.Offset(offset).Limit(opts.Limit)
	
	// Apply ordering
	orderBy := "analyzed_at DESC"
	if opts.OrderBy != "" {
		orderBy = opts.OrderBy
		if opts.Desc {
			orderBy += " DESC"
		}
	}
	query = query.Order(orderBy)
	
	if err := query.Find(&analyses).Error; err != nil {
		return nil, nil, err
	}
	
	pagination := &PaginationResult{
		Total:       total,
		Page:        opts.Page,
		Limit:       opts.Limit,
		TotalPages:  (int(total) + opts.Limit - 1) / opts.Limit,
		HasNext:     opts.Page < (int(total)+opts.Limit-1)/opts.Limit,
		HasPrevious: opts.Page > 1,
	}
	
	return analyses, pagination, nil
}

// GetAnalysisHistory gets analysis history for a person
func (r *analysisRepository) GetAnalysisHistory(ctx context.Context, userID, personID uuid.UUID, limit int) ([]*models.RelationshipAnalysis, error) {
	var analyses []*models.RelationshipAnalysis
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND person_id = ?", userID, personID).
		Order("analyzed_at DESC").
		Limit(limit).
		Find(&analyses).Error
	
	return analyses, err
}

// CreateRecommendation creates a new AI recommendation (stored as Nudge)
func (r *analysisRepository) CreateRecommendation(ctx context.Context, recommendation *models.Nudge) error {
	return r.db.WithContext(ctx).Create(recommendation).Error
}

// GetRecommendationByID gets a recommendation by ID
func (r *analysisRepository) GetRecommendationByID(ctx context.Context, recommendationID uuid.UUID) (*models.Nudge, error) {
	var recommendation models.Nudge
	err := r.db.WithContext(ctx).Where("source = ?", "ai").First(&recommendation, "id = ?", recommendationID).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	
	return &recommendation, nil
}

// ListRecommendations lists recommendations with filtering and pagination
func (r *analysisRepository) ListRecommendations(ctx context.Context, userID uuid.UUID, opts FilterOptions) ([]*models.Nudge, *PaginationResult, error) {
	var recommendations []*models.Nudge
	var total int64
	
	query := r.db.WithContext(ctx).Model(&models.Nudge{}).Where("user_id = ? AND source = ?", userID, "ai")
	
	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}
	
	// Apply pagination
	offset := (opts.Page - 1) * opts.Limit
	query = query.Offset(offset).Limit(opts.Limit)
	
	// Apply ordering
	orderBy := "created_at DESC"
	if opts.OrderBy != "" {
		orderBy = opts.OrderBy
		if opts.Desc {
			orderBy += " DESC"
		}
	}
	query = query.Order(orderBy)
	
	if err := query.Find(&recommendations).Error; err != nil {
		return nil, nil, err
	}
	
	pagination := &PaginationResult{
		Total:       total,
		Page:        opts.Page,
		Limit:       opts.Limit,
		TotalPages:  (int(total) + opts.Limit - 1) / opts.Limit,
		HasNext:     opts.Page < (int(total)+opts.Limit-1)/opts.Limit,
		HasPrevious: opts.Page > 1,
	}
	
	return recommendations, pagination, nil
}

// GetActiveRecommendations gets active (pending/accepted) recommendations for a user
func (r *analysisRepository) GetActiveRecommendations(ctx context.Context, userID uuid.UUID) ([]*models.Nudge, error) {
	var recommendations []*models.Nudge
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND source = ? AND status IN (?, ?) AND (expires_at IS NULL OR expires_at > ?)", 
			userID, "ai", "pending", "accepted", time.Now()).
		Order("priority DESC, created_at DESC").
		Find(&recommendations).Error
	
	return recommendations, err
}

// GetRecommendationsForPerson gets recommendations for a specific person
func (r *analysisRepository) GetRecommendationsForPerson(ctx context.Context, userID, personID uuid.UUID) ([]*models.Nudge, error) {
	var recommendations []*models.Nudge
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND person_id = ? AND source = ?", userID, personID, "ai").
		Order("created_at DESC").
		Find(&recommendations).Error
	
	return recommendations, err
}

// UpdateRecommendationStatus updates the status of a recommendation
func (r *analysisRepository) UpdateRecommendationStatus(ctx context.Context, recommendationID uuid.UUID, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	
	now := time.Now()
	switch status {
	case "accepted":
		updates["accepted_at"] = now
	case "completed":
		updates["completed_at"] = now
	case "dismissed":
		updates["dismissed_at"] = now
	}
	
	return r.db.WithContext(ctx).
		Model(&models.Nudge{}).
		Where("id = ? AND source = ?", recommendationID, "ai").
		Updates(updates).Error
}

// CreateJob creates a new AI analysis job
func (r *analysisRepository) CreateJob(ctx context.Context, job *models.AIAnalysisJob) error {
	return r.db.WithContext(ctx).Create(job).Error
}

// GetJobByID gets a job by ID
func (r *analysisRepository) GetJobByID(ctx context.Context, jobID uuid.UUID) (*models.AIAnalysisJob, error) {
	var job models.AIAnalysisJob
	err := r.db.WithContext(ctx).First(&job, "id = ?", jobID).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	
	return &job, nil
}

// UpdateJob updates a job
func (r *analysisRepository) UpdateJob(ctx context.Context, job *models.AIAnalysisJob) error {
	return r.db.WithContext(ctx).Save(job).Error
}

// ListPendingJobs lists pending jobs ordered by priority
func (r *analysisRepository) ListPendingJobs(ctx context.Context, limit int) ([]*models.AIAnalysisJob, error) {
	var jobs []*models.AIAnalysisJob
	err := r.db.WithContext(ctx).
		Where("status = ?", "pending").
		Order("priority DESC, created_at ASC").
		Limit(limit).
		Find(&jobs).Error
	
	return jobs, err
}

// GetUserJobs gets recent jobs for a user
func (r *analysisRepository) GetUserJobs(ctx context.Context, userID uuid.UUID, limit int) ([]*models.AIAnalysisJob, error) {
	var jobs []*models.AIAnalysisJob
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&jobs).Error
	
	return jobs, err
}
