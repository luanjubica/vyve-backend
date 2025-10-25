package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vyve/vyve-backend/internal/middleware"
	"github.com/vyve/vyve-backend/internal/services"
)

// AnalysisHandler handles AI analysis endpoints
type AnalysisHandler interface {
	// Analysis endpoints
	GetPersonAnalysis(c *fiber.Ctx) error
	RefreshPersonAnalysis(c *fiber.Ctx) error
	GetAnalysisHistory(c *fiber.Ctx) error

	// Recommendation endpoints
	GetPersonRecommendations(c *fiber.Ctx) error
	GetActiveRecommendations(c *fiber.Ctx) error
	UpdateRecommendationStatus(c *fiber.Ctx) error

	// Batch operations
	BatchAnalyze(c *fiber.Ctx) error
	GetJobStatus(c *fiber.Ctx) error

	// Overall insights
	GetOverallInsights(c *fiber.Ctx) error
}

type analysisHandler struct {
	analysisService services.AnalysisService
}

// NewAnalysisHandler creates a new analysis handler
func NewAnalysisHandler(analysisService services.AnalysisService) AnalysisHandler {
	return &analysisHandler{
		analysisService: analysisService,
	}
}

// GetPersonAnalysis gets the latest AI analysis for a person
// GET /api/v1/people/:id/analysis
func (h *analysisHandler) GetPersonAnalysis(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	personID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid person ID",
		})
	}

	analysis, err := h.analysisService.GetLatestAnalysis(c.Context(), userID, personID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get analysis",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"analysis": analysis,
	})
}

// RefreshPersonAnalysis triggers a new analysis for a person
// POST /api/v1/people/:id/analysis/refresh
func (h *analysisHandler) RefreshPersonAnalysis(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	personID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid person ID",
		})
	}

	analysis, err := h.analysisService.RefreshAnalysis(c.Context(), userID, personID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to refresh analysis",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"analysis": analysis,
		"message":  "Analysis refreshed successfully",
	})
}

// GetAnalysisHistory gets analysis history for a person
// GET /api/v1/people/:id/analysis/history
func (h *analysisHandler) GetAnalysisHistory(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	personID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid person ID",
		})
	}

	limit := c.QueryInt("limit", 10)
	if limit > 50 {
		limit = 50
	}

	history, err := h.analysisService.GetAnalysisHistory(c.Context(), userID, personID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get analysis history",
		})
	}

	return c.JSON(fiber.Map{
		"history": history,
		"count":   len(history),
	})
}

// GetPersonRecommendations gets AI recommendations for a person
// GET /api/v1/people/:id/recommendations
func (h *analysisHandler) GetPersonRecommendations(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	personID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid person ID",
		})
	}

	recommendations, err := h.analysisService.GetRecommendationsForPerson(c.Context(), userID, personID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get recommendations",
		})
	}

	// If no recommendations exist, generate them
	if len(recommendations) == 0 {
		recommendations, err = h.analysisService.GenerateRecommendations(c.Context(), userID, personID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate recommendations",
			})
		}
	}

	return c.JSON(fiber.Map{
		"recommendations": recommendations,
		"count":           len(recommendations),
	})
}

// GetActiveRecommendations gets all active recommendations for the user
// GET /api/v1/recommendations
func (h *analysisHandler) GetActiveRecommendations(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	recommendations, err := h.analysisService.GetActiveRecommendations(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get recommendations",
		})
	}

	return c.JSON(fiber.Map{
		"recommendations": recommendations,
		"count":           len(recommendations),
	})
}

// UpdateRecommendationStatus updates the status of a recommendation
// POST /api/v1/recommendations/:id/status
func (h *analysisHandler) UpdateRecommendationStatus(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	recommendationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid recommendation ID",
		})
	}

	var req struct {
		Status string `json:"status" validate:"required,oneof=accepted completed dismissed"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.analysisService.UpdateRecommendationStatus(c.Context(), userID, recommendationID, req.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update recommendation status",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Recommendation status updated successfully",
	})
}

// BatchAnalyze triggers batch analysis for multiple people
// POST /api/v1/analytics/batch-analyze
func (h *analysisHandler) BatchAnalyze(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		PersonIDs []string `json:"person_ids"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Parse person IDs
	personIDs := make([]uuid.UUID, 0, len(req.PersonIDs))
	for _, idStr := range req.PersonIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}
		personIDs = append(personIDs, id)
	}

	if len(personIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No valid person IDs provided",
		})
	}

	job, err := h.analysisService.BatchAnalyze(c.Context(), userID, personIDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create batch analysis job",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"job":     job,
		"message": "Batch analysis job created",
	})
}

// GetJobStatus gets the status of a batch analysis job
// GET /api/v1/analytics/jobs/:id
func (h *analysisHandler) GetJobStatus(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	jobID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid job ID",
		})
	}

	job, err := h.analysisService.GetJobStatus(c.Context(), userID, jobID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get job status",
		})
	}

	return c.JSON(fiber.Map{
		"job": job,
	})
}

// GetOverallInsights gets overall relationship insights for the user
// GET /api/v1/analytics/insights
func (h *analysisHandler) GetOverallInsights(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Get all active recommendations
	recommendations, err := h.analysisService.GetActiveRecommendations(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get insights",
		})
	}

	// Aggregate insights
	insights := fiber.Map{
		"total_recommendations": len(recommendations),
		"high_priority":         0,
		"medium_priority":       0,
		"low_priority":          0,
		"by_type":               make(map[string]int),
	}

	for _, rec := range recommendations {
		switch rec.Priority {
		case "high":
			insights["high_priority"] = insights["high_priority"].(int) + 1
		case "medium":
			insights["medium_priority"] = insights["medium_priority"].(int) + 1
		case "low":
			insights["low_priority"] = insights["low_priority"].(int) + 1
		}

		byType := insights["by_type"].(map[string]int)
		byType[rec.Type]++
	}

	return c.JSON(fiber.Map{
		"insights":        insights,
		"recommendations": recommendations,
	})
}
