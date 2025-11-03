package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/vyve/vyve-backend/internal/models"
	"github.com/vyve/vyve-backend/internal/repository"
	"github.com/vyve/vyve-backend/pkg/ai"
)

// AnalysisService handles AI-powered relationship analysis
type AnalysisService interface {
	// Analysis operations
	AnalyzeRelationship(ctx context.Context, userID, personID uuid.UUID) (*models.RelationshipAnalysis, error)
	GetLatestAnalysis(ctx context.Context, userID, personID uuid.UUID) (*models.RelationshipAnalysis, error)
	GetAnalysisHistory(ctx context.Context, userID, personID uuid.UUID, limit int) ([]*models.RelationshipAnalysis, error)
	RefreshAnalysis(ctx context.Context, userID, personID uuid.UUID) (*models.RelationshipAnalysis, error)
	
	// Recommendation operations (now using Nudge model)
	GenerateRecommendations(ctx context.Context, userID, personID uuid.UUID) ([]*models.Nudge, error)
	GetActiveRecommendations(ctx context.Context, userID uuid.UUID) ([]*models.Nudge, error)
	GetRecommendationsForPerson(ctx context.Context, userID, personID uuid.UUID) ([]*models.Nudge, error)
	UpdateRecommendationStatus(ctx context.Context, userID, recommendationID uuid.UUID, status string) error
	
	// Batch operations
	BatchAnalyze(ctx context.Context, userID uuid.UUID, personIDs []uuid.UUID) (*models.AIAnalysisJob, error)
	GetJobStatus(ctx context.Context, userID, jobID uuid.UUID) (*models.AIAnalysisJob, error)
}

type analysisService struct {
	aiService      *ai.Service
	analysisRepo   repository.AnalysisRepository
	personRepo     repository.PersonRepository
	interactionRepo repository.InteractionRepository
}

// NewAnalysisService creates a new analysis service
func NewAnalysisService(
	aiService *ai.Service,
	analysisRepo repository.AnalysisRepository,
	personRepo repository.PersonRepository,
	interactionRepo repository.InteractionRepository,
) AnalysisService {
	return &analysisService{
		aiService:       aiService,
		analysisRepo:    analysisRepo,
		personRepo:      personRepo,
		interactionRepo: interactionRepo,
	}
}

// AnalyzeRelationship performs AI analysis on a relationship
func (s *analysisService) AnalyzeRelationship(ctx context.Context, userID, personID uuid.UUID) (*models.RelationshipAnalysis, error) {
	log.Printf("[ANALYSIS_SERVICE] Starting analysis for person=%s, user=%s", personID, userID)
	
	// Check if AI service is available
	if s.aiService == nil {
		log.Printf("[ANALYSIS_SERVICE] ❌ AI service is nil - not configured")
		return nil, fmt.Errorf("AI service is not available - please enable FEATURE_AI_INSIGHTS and configure API keys")
	}
	
	log.Printf("[ANALYSIS_SERVICE] Fetching person details...")
	// Get person details
	person, err := s.personRepo.FindByID(ctx, personID)
	if err != nil {
		log.Printf("[ANALYSIS_SERVICE] ❌ Failed to get person: %v", err)
		return nil, fmt.Errorf("failed to get person: %w", err)
	}
	
	// Check ownership
	if person.UserID != userID {
		log.Printf("[ANALYSIS_SERVICE] ❌ Forbidden: user=%s does not own person=%s", userID, personID)
		return nil, repository.ErrForbidden
	}
	
	log.Printf("[ANALYSIS_SERVICE] Fetching recent interactions...")
	// Get recent interactions (last 30)
	interactions, err := s.personRepo.GetRecentInteractions(ctx, personID, 30)
	if err != nil {
		log.Printf("[ANALYSIS_SERVICE] ❌ Failed to get interactions: %v", err)
		return nil, fmt.Errorf("failed to get interactions: %w", err)
	}
	
	log.Printf("[ANALYSIS_SERVICE] Found %d interactions, building AI request...", len(interactions))
	// Build AI request
	aiReq := s.buildAnalysisRequest(person, interactions)
	
	log.Printf("[ANALYSIS_SERVICE] Calling AI service for analysis...")
	// Call AI service
	aiResp, err := s.aiService.Analyze(ctx, aiReq)
	if err != nil {
		log.Printf("[ANALYSIS_SERVICE] ❌ AI analysis failed: %v", err)
		return nil, fmt.Errorf("AI analysis failed: %w", err)
	}
	
	log.Printf("[ANALYSIS_SERVICE] ✅ AI analysis completed successfully")
	
	// Create analysis record
	analysis := &models.RelationshipAnalysis{
		UserID:               userID,
		PersonID:             personID,
		ConnectionStrength:   aiResp.ConnectionStrength,
		EngagementQuality:    aiResp.EngagementQuality,
		CommunicationBalance: aiResp.CommunicationBalance,
		EnergyAlignment:      aiResp.EnergyAlignment,
		RelationshipHealth:   aiResp.RelationshipHealth,
		OverallScore:         aiResp.OverallScore,
		Summary:              aiResp.Summary,
		KeyInsights:          aiResp.KeyInsights,
		Patterns:             aiResp.Patterns,
		Strengths:            aiResp.Strengths,
		Concerns:             aiResp.Concerns,
		TrendDirection:       aiResp.TrendDirection,
		Provider:             s.aiService.GetProviderName(),
		Model:                s.aiService.GetModelName(),
		TokensUsed:           aiResp.TokensUsed,
		ProcessingTimeMs:     aiResp.ProcessingTimeMs,
		InteractionsCount:    len(interactions),
		AnalyzedAt:           time.Now(),
	}
	
	log.Printf("[ANALYSIS_SERVICE] Saving analysis to database...")
	if err := s.analysisRepo.CreateAnalysis(ctx, analysis); err != nil {
		log.Printf("[ANALYSIS_SERVICE] ❌ Failed to save analysis to database: %v", err)
		return nil, fmt.Errorf("failed to save analysis: %w", err)
	}
	
	log.Printf("[ANALYSIS_SERVICE] ✅ Analysis saved successfully with ID: %s", analysis.ID)
	return analysis, nil
}

// GetLatestAnalysis gets the most recent analysis for a person
func (s *analysisService) GetLatestAnalysis(ctx context.Context, userID, personID uuid.UUID) (*models.RelationshipAnalysis, error) {
	analysis, err := s.analysisRepo.GetLatestAnalysis(ctx, userID, personID)
	if err != nil {
		if err == repository.ErrNotFound {
			// No analysis exists, create one
			return s.AnalyzeRelationship(ctx, userID, personID)
		}
		return nil, err
	}
	
	// Check if analysis is stale (older than 24 hours)
	if time.Since(analysis.AnalyzedAt) > 24*time.Hour {
		// Optionally refresh in background
		go s.AnalyzeRelationship(context.Background(), userID, personID)
	}
	
	return analysis, nil
}

// GetAnalysisHistory gets analysis history for a person
func (s *analysisService) GetAnalysisHistory(ctx context.Context, userID, personID uuid.UUID, limit int) ([]*models.RelationshipAnalysis, error) {
	return s.analysisRepo.GetAnalysisHistory(ctx, userID, personID, limit)
}

// RefreshAnalysis forces a new analysis
func (s *analysisService) RefreshAnalysis(ctx context.Context, userID, personID uuid.UUID) (*models.RelationshipAnalysis, error) {
	return s.AnalyzeRelationship(ctx, userID, personID)
}

// GenerateRecommendations generates AI-powered recommendations (stored as Nudges)
func (s *analysisService) GenerateRecommendations(ctx context.Context, userID, personID uuid.UUID) ([]*models.Nudge, error) {
	// Check if AI service is available
	if s.aiService == nil {
		return nil, fmt.Errorf("AI service is not available - please enable FEATURE_AI_INSIGHTS and configure API keys")
	}
	
	// Get latest analysis
	analysis, err := s.GetLatestAnalysis(ctx, userID, personID)
	if err != nil {
		return nil, fmt.Errorf("failed to get analysis: %w", err)
	}
	
	// Get person details
	person, err := s.personRepo.FindByID(ctx, personID)
	if err != nil {
		return nil, fmt.Errorf("failed to get person: %w", err)
	}
	
	// Get recent interactions
	interactions, err := s.personRepo.GetRecentInteractions(ctx, personID, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get interactions: %w", err)
	}
	
	// Build AI request
	aiReq := s.buildRecommendationRequest(person, analysis, interactions)
	
	// Call AI service
	aiResp, err := s.aiService.GenerateRecommendations(ctx, aiReq)
	if err != nil {
		return nil, fmt.Errorf("AI recommendation generation failed: %w", err)
	}
	
	// Create recommendation records (as Nudges with source='ai')
	recommendations := make([]*models.Nudge, 0, len(aiResp.Recommendations))
	for _, rec := range aiResp.Recommendations {
		recommendation := &models.Nudge{
			UserID:               userID,
			PersonID:             &personID,
			AnalysisID:           &analysis.ID,
			Source:               "ai",
			Type:                 rec.Type,
			Priority:             rec.Priority,
			Title:                rec.Title,
			Message:              rec.Description,
			Reasoning:            rec.Reasoning,
			SuggestedActions:     rec.SuggestedActions,
			ConversationStarters: rec.ConversationStarters,
			Timing:               rec.Timing,
			EstimatedImpact:      rec.EstimatedImpact,
			Status:               "pending",
			Provider:             s.aiService.GetProviderName(),
			Model:                s.aiService.GetModelName(),
		}
		
		// Set expiration based on timing
		if rec.ExpiresIn != nil {
			expiresAt := time.Now().Add(*rec.ExpiresIn)
			recommendation.ExpiresAt = &expiresAt
		} else {
			// Default expiration based on timing
			switch rec.Timing {
			case "now":
				expiresAt := time.Now().Add(24 * time.Hour)
				recommendation.ExpiresAt = &expiresAt
			case "today":
				expiresAt := time.Now().Add(48 * time.Hour)
				recommendation.ExpiresAt = &expiresAt
			case "this_week":
				expiresAt := time.Now().Add(7 * 24 * time.Hour)
				recommendation.ExpiresAt = &expiresAt
			case "this_month":
				expiresAt := time.Now().Add(30 * 24 * time.Hour)
				recommendation.ExpiresAt = &expiresAt
			}
		}
		
		if err := s.analysisRepo.CreateRecommendation(ctx, recommendation); err != nil {
			return nil, fmt.Errorf("failed to save recommendation: %w", err)
		}
		
		recommendations = append(recommendations, recommendation)
	}
	
	return recommendations, nil
}

// GetActiveRecommendations gets active recommendations for a user
func (s *analysisService) GetActiveRecommendations(ctx context.Context, userID uuid.UUID) ([]*models.Nudge, error) {
	return s.analysisRepo.GetActiveRecommendations(ctx, userID)
}

// GetRecommendationsForPerson gets recommendations for a specific person
func (s *analysisService) GetRecommendationsForPerson(ctx context.Context, userID, personID uuid.UUID) ([]*models.Nudge, error) {
	return s.analysisRepo.GetRecommendationsForPerson(ctx, userID, personID)
}

// UpdateRecommendationStatus updates the status of a recommendation
func (s *analysisService) UpdateRecommendationStatus(ctx context.Context, userID, recommendationID uuid.UUID, status string) error {
	// Get recommendation to verify ownership
	recommendation, err := s.analysisRepo.GetRecommendationByID(ctx, recommendationID)
	if err != nil {
		return err
	}
	
	if recommendation.UserID != userID {
		return repository.ErrForbidden
	}
	
	return s.analysisRepo.UpdateRecommendationStatus(ctx, recommendationID, status)
}

// BatchAnalyze creates a batch analysis job
func (s *analysisService) BatchAnalyze(ctx context.Context, userID uuid.UUID, personIDs []uuid.UUID) (*models.AIAnalysisJob, error) {
	// Create job
	personIDStrings := make([]string, len(personIDs))
	for i, id := range personIDs {
		personIDStrings[i] = id.String()
	}
	
	job := &models.AIAnalysisJob{
		UserID:     userID,
		JobType:    "batch_analysis",
		Status:     "pending",
		Priority:   5,
		PersonIDs:  personIDStrings,
		TotalItems: len(personIDs),
	}
	
	if err := s.analysisRepo.CreateJob(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}
	
	// Process in background
	go s.processBatchAnalysis(context.Background(), job)
	
	return job, nil
}

// GetJobStatus gets the status of an analysis job
func (s *analysisService) GetJobStatus(ctx context.Context, userID, jobID uuid.UUID) (*models.AIAnalysisJob, error) {
	job, err := s.analysisRepo.GetJobByID(ctx, jobID)
	if err != nil {
		return nil, err
	}
	
	if job.UserID != userID {
		return nil, repository.ErrForbidden
	}
	
	return job, nil
}

// Helper functions

func (s *analysisService) buildAnalysisRequest(person *models.Person, interactions []*models.Interaction) ai.AnalysisRequest {
	req := ai.AnalysisRequest{
		PersonName:       person.Name,
		Relationship:     person.Relationship,
		InteractionCount: person.InteractionCount,
		HealthScore:      person.HealthScore,
		Context:          person.Context,
		Notes:            person.Notes,
	}
	
	if person.LastInteractionAt != nil {
		req.LastInteraction = person.LastInteractionAt
	}
	
	// Convert interactions
	req.RecentInteractions = make([]ai.InteractionData, len(interactions))
	for i, interaction := range interactions {
		req.RecentInteractions[i] = ai.InteractionData{
			Date:         interaction.InteractionAt,
			EnergyImpact: interaction.EnergyImpact,
			Quality:      interaction.Quality,
			Duration:     interaction.Duration,
			Context:      interaction.Context,
			Notes:        interaction.Notes,
		}
	}
	
	return req
}

func (s *analysisService) buildRecommendationRequest(person *models.Person, analysis *models.RelationshipAnalysis, interactions []*models.Interaction) ai.RecommendationRequest {
	req := ai.RecommendationRequest{
		PersonName:   person.Name,
		Relationship: person.Relationship,
		Context:      person.Context,
		Analysis: &ai.AnalysisResponse{
			ConnectionStrength:   analysis.ConnectionStrength,
			EngagementQuality:    analysis.EngagementQuality,
			CommunicationBalance: analysis.CommunicationBalance,
			EnergyAlignment:      analysis.EnergyAlignment,
			RelationshipHealth:   analysis.RelationshipHealth,
			OverallScore:         analysis.OverallScore,
			Summary:              analysis.Summary,
			KeyInsights:          analysis.KeyInsights,
			Patterns:             analysis.Patterns,
			Strengths:            analysis.Strengths,
			Concerns:             analysis.Concerns,
			TrendDirection:       analysis.TrendDirection,
		},
	}
	
	if person.LastInteractionAt != nil {
		req.LastInteraction = person.LastInteractionAt
	}
	
	// Convert interactions
	req.RecentInteractions = make([]ai.InteractionData, len(interactions))
	for i, interaction := range interactions {
		req.RecentInteractions[i] = ai.InteractionData{
			Date:         interaction.InteractionAt,
			EnergyImpact: interaction.EnergyImpact,
			Quality:      interaction.Quality,
			Duration:     interaction.Duration,
			Context:      interaction.Context,
			Notes:        interaction.Notes,
		}
	}
	
	return req
}

func (s *analysisService) processBatchAnalysis(ctx context.Context, job *models.AIAnalysisJob) {
	// Update job status
	now := time.Now()
	job.Status = "processing"
	job.StartedAt = &now
	s.analysisRepo.UpdateJob(ctx, job)
	
	// Process each person
	for _, personIDStr := range job.PersonIDs {
		personID, err := uuid.Parse(personIDStr)
		if err != nil {
			job.FailedItems++
			continue
		}
		
		_, err = s.AnalyzeRelationship(ctx, job.UserID, personID)
		if err != nil {
			job.FailedItems++
		} else {
			job.ProcessedItems++
		}
		
		// Update progress
		job.Progress = float64(job.ProcessedItems+job.FailedItems) / float64(job.TotalItems) * 100
		s.analysisRepo.UpdateJob(ctx, job)
	}
	
	// Mark job as completed
	completed := time.Now()
	job.Status = "completed"
	job.CompletedAt = &completed
	s.analysisRepo.UpdateJob(ctx, job)
}
