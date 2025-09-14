package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/vyve/vyve-backend/internal/middleware"
	"github.com/vyve/vyve-backend/internal/services"
)

type interactionHandler struct {
	interactionService services.InteractionService
}

// NewInteractionHandler creates a new interaction handler
func NewInteractionHandler(interactionService services.InteractionService) InteractionHandler {
	return &interactionHandler{
		interactionService: interactionService,
	}
}

// List handles GET /interactions/
func (h *interactionHandler) List(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	opts := services.ListOptions{
		Page:  page,
		Limit: limit,
	}

	interactions, pagination, err := h.interactionService.List(c.Context(), userID, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get interactions"})
	}

	return c.JSON(fiber.Map{
		"success":    true,
		"data":       interactions,
		"pagination": pagination,
	})
}

// Create handles POST /interactions/
func (h *interactionHandler) Create(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req services.CreateInteractionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if req.PersonID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "person_id is required"})
	}

	if req.EnergyImpact == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "energy_impact is required"})
	}

	// Validate energy impact values
	validEnergyImpacts := map[string]bool{
		"energizing": true,
		"neutral":    true,
		"draining":   true,
	}
	if !validEnergyImpacts[req.EnergyImpact] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "energy_impact must be 'energizing', 'neutral', or 'draining'"})
	}

	// Validate quality if provided
	if req.Quality > 0 && (req.Quality < 1 || req.Quality > 5) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "quality must be between 1 and 5"})
	}

	interaction, err := h.interactionService.Create(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create interaction"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    interaction,
	})
}

// Get handles GET /interactions/:id
func (h *interactionHandler) Get(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	interactionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid interaction ID"})
	}

	interaction, err := h.interactionService.GetByID(c.Context(), userID, interactionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Interaction not found"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    interaction,
	})
}

// Update handles PUT /interactions/:id
func (h *interactionHandler) Update(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	interactionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid interaction ID"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate energy impact if being updated
	if energyImpact, ok := updates["energy_impact"].(string); ok {
		validEnergyImpacts := map[string]bool{
			"energizing": true,
			"neutral":    true,
			"draining":   true,
		}
		if !validEnergyImpacts[energyImpact] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "energy_impact must be 'energizing', 'neutral', or 'draining'"})
		}
	}

	// Validate quality if being updated
	if quality, ok := updates["quality"]; ok {
		if qualityFloat, ok := quality.(float64); ok {
			qualityInt := int(qualityFloat)
			if qualityInt < 1 || qualityInt > 5 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "quality must be between 1 and 5"})
			}
			updates["quality"] = qualityInt
		}
	}

	interaction, err := h.interactionService.Update(c.Context(), userID, interactionID, updates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update interaction"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    interaction,
	})
}

// Delete handles DELETE /interactions/:id
func (h *interactionHandler) Delete(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	interactionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid interaction ID"})
	}

	if err := h.interactionService.Delete(c.Context(), userID, interactionID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete interaction"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Interaction deleted successfully",
	})
}

// GetRecent handles GET /interactions/recent
func (h *interactionHandler) GetRecent(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Parse limit parameter
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	interactions, err := h.interactionService.GetRecent(c.Context(), userID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get recent interactions"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    interactions,
	})
}

// GetByDate handles GET /interactions/by-date
func (h *interactionHandler) GetByDate(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Parse date parameter
	dateStr := c.Query("date")
	if dateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "date parameter is required (YYYY-MM-DD format)"})
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid date format. Use YYYY-MM-DD"})
	}

	interactions, err := h.interactionService.GetByDate(c.Context(), userID, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get interactions by date"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    interactions,
		"date":    dateStr,
	})
}

// GetEnergyDistribution handles GET /interactions/energy-distribution
func (h *interactionHandler) GetEnergyDistribution(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Parse period parameter
	period := c.Query("period", "week")
	validPeriods := map[string]bool{
		"week":    true,
		"month":   true,
		"quarter": true,
		"year":    true,
	}

	if !validPeriods[period] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid period. Must be: week, month, quarter, or year"})
	}

	distribution, err := h.interactionService.GetEnergyDistribution(c.Context(), userID, period)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get energy distribution"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    distribution,
		"period":  period,
	})
}

// BulkCreate handles POST /interactions/bulk
func (h *interactionHandler) BulkCreate(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		Interactions []services.CreateInteractionRequest `json:"interactions" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if len(req.Interactions) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "At least one interaction is required"})
	}

	// Limit bulk operations
	if len(req.Interactions) > 50 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Maximum 50 interactions allowed per bulk request"})
	}

	// Validate each interaction
	validEnergyImpacts := map[string]bool{
		"energizing": true,
		"neutral":    true,
		"draining":   true,
	}

	for i, interaction := range req.Interactions {
		if interaction.PersonID == uuid.Nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "person_id is required for all interactions",
				"index": i,
			})
		}

		if interaction.EnergyImpact == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "energy_impact is required for all interactions",
				"index": i,
			})
		}

		if !validEnergyImpacts[interaction.EnergyImpact] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "energy_impact must be 'energizing', 'neutral', or 'draining'",
				"index": i,
			})
		}

		if interaction.Quality > 0 && (interaction.Quality < 1 || interaction.Quality > 5) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "quality must be between 1 and 5",
				"index": i,
			})
		}
	}

	interactions, err := h.interactionService.BulkCreate(c.Context(), userID, req.Interactions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create interactions"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    interactions,
		"count":   len(interactions),
	})
}
