package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyve/vyve-backend/internal/middleware"
	"github.com/vyve/vyve-backend/internal/services"
)

type dictionaryHandler struct {
	service services.DictionaryService
}

func NewDictionaryHandler(service services.DictionaryService) DictionaryHandler {
	return &dictionaryHandler{service: service}
}

// Categories returns user-scoped categories
func (h *dictionaryHandler) Categories(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	items, err := h.service.ListCategories(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get categories"})
	}
	return c.JSON(fiber.Map{"success": true, "data": items})
}

func (h *dictionaryHandler) CommunicationMethods(c *fiber.Ctx) error {
	items, err := h.service.ListCommunicationMethods(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get communication methods"})
	}
	return c.JSON(fiber.Map{"success": true, "data": items})
}

func (h *dictionaryHandler) RelationshipStatuses(c *fiber.Ctx) error {
	items, err := h.service.ListRelationshipStatuses(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get relationship statuses"})
	}
	return c.JSON(fiber.Map{"success": true, "data": items})
}

func (h *dictionaryHandler) Intentions(c *fiber.Ctx) error {
	items, err := h.service.ListIntentions(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get intentions"})
	}
	return c.JSON(fiber.Map{"success": true, "data": items})
}

func (h *dictionaryHandler) EnergyPatterns(c *fiber.Ctx) error {
	items, err := h.service.ListEnergyPatterns(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get energy patterns"})
	}
	return c.JSON(fiber.Map{"success": true, "data": items})
}
