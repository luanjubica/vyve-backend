package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vyve/vyve-backend/internal/middleware"
	"github.com/vyve/vyve-backend/internal/services"
)

type nudgeHandler struct {
	nudgeService services.NudgeService
}

// NewNudgeHandler creates a new nudge handler
func NewNudgeHandler(nudgeService services.NudgeService) NudgeHandler {
	return &nudgeHandler{
		nudgeService: nudgeService,
	}
}

// List handles GET /nudges
func (h *nudgeHandler) List(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	source := c.Query("source")       // Filter by source: 'ai' or 'system'
	status := c.Query("status")       // Filter by status: 'pending', 'seen', 'completed', etc.
	personID := c.Query("person_id")  // Filter by person

	opts := services.ListOptions{
		Page:   page,
		Limit:  limit,
		Source: source,
		Status: status,
	}

	// Add person filter if provided
	if personID != "" {
		pid, err := uuid.Parse(personID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid person_id"})
		}
		opts.PersonID = &pid
	}

	nudges, pagination, err := h.nudgeService.List(c.Context(), userID, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get nudges",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"nudges":     nudges,
		"pagination": pagination,
	})
}

// Get handles GET /nudges/:id
func (h *nudgeHandler) Get(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	nudgeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid nudge ID"})
	}

	nudge, err := h.nudgeService.GetByID(c.Context(), userID, nudgeID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Nudge not found"})
	}

	return c.JSON(fiber.Map{
		"nudge": nudge,
	})
}

// GetActive handles GET /nudges/active
func (h *nudgeHandler) GetActive(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	source := c.Query("source") // Optional filter by source

	// If source filter is provided, use List with filters
	if source != "" {
		opts := services.ListOptions{
			Page:   1,
			Limit:  100,
			Source: source,
			Status: "pending", // Active means pending
		}
		nudges, _, err := h.nudgeService.List(c.Context(), userID, opts)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get active nudges"})
		}
		return c.JSON(fiber.Map{
			"nudges": nudges,
			"count":  len(nudges),
		})
	}

	// Otherwise use the dedicated GetActive method
	nudges, err := h.nudgeService.GetActive(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get active nudges"})
	}

	return c.JSON(fiber.Map{
		"nudges": nudges,
		"count":  len(nudges),
	})
}

// GetHistory handles GET /nudges/history
func (h *nudgeHandler) GetHistory(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	if limit > 100 {
		limit = 100
	}

	nudges, err := h.nudgeService.GetHistory(c.Context(), userID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get nudge history"})
	}

	return c.JSON(fiber.Map{
		"nudges": nudges,
		"count":  len(nudges),
	})
}

// MarkSeen handles POST /nudges/:id/seen
func (h *nudgeHandler) MarkSeen(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	nudgeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid nudge ID"})
	}

	if err := h.nudgeService.MarkSeen(c.Context(), userID, nudgeID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to mark nudge as seen"})
	}

	return c.JSON(fiber.Map{
		"message": "Nudge marked as seen",
	})
}

// MarkActedOn handles POST /nudges/:id/act
func (h *nudgeHandler) MarkActedOn(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	nudgeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid nudge ID"})
	}

	if err := h.nudgeService.MarkActedOn(c.Context(), userID, nudgeID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to mark nudge as acted on"})
	}

	return c.JSON(fiber.Map{
		"message": "Nudge marked as acted on",
	})
}

// Dismiss handles DELETE /nudges/:id
func (h *nudgeHandler) Dismiss(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	nudgeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid nudge ID"})
	}

	if err := h.nudgeService.Dismiss(c.Context(), userID, nudgeID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to dismiss nudge"})
	}

	return c.JSON(fiber.Map{
		"message": "Nudge dismissed",
	})
}

// GenerateNudges handles POST /nudges/generate
func (h *nudgeHandler) GenerateNudges(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		PersonID *string `json:"person_id"` // Optional - generate for specific person
	}

	// Try to parse body, but ignore errors if body is empty
	// This allows both empty POST and POST with JSON body
	if len(c.Body()) > 0 {
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
		}
	}

	// If person_id provided, generate for that person
	if req.PersonID != nil {
		personID, err := uuid.Parse(*req.PersonID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid person_id"})
		}

		nudges, err := h.nudgeService.GenerateForPerson(c.Context(), userID, personID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to generate nudges",
				"details": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"nudges":  nudges,
			"count":   len(nudges),
			"message": "Nudges generated successfully",
		})
	}

	// Otherwise, generate system nudges for all relationships
	nudges, err := h.nudgeService.GenerateSystemNudges(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to generate nudges",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"nudges":  nudges,
		"count":   len(nudges),
		"message": "System nudges generated successfully",
	})
}
