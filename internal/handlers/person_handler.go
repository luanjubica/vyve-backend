package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/vyve/vyve-backend/internal/services"
)

type personHandler struct {
	personService services.PersonService
}

// NewPersonHandler creates a new person handler
func NewPersonHandler(personService services.PersonService) PersonHandler {
	return &personHandler{
		personService: personService,
	}
}

// CountPeople handles GET /people/count
func (h *personHandler) CountPeople(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get people count
	count, err := h.personService.CountPeople(c.Context(), userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get people count",
		})
	}

	return c.JSON(fiber.Map{
		"count": count,
	})
}

// Implement other required methods from PersonHandler interface with stubs
func (h *personHandler) List(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Not implemented 135"})
}

func (h *personHandler) Create(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Not implemented 136"})
}

func (h *personHandler) Get(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Not implemented 137"})
}

func (h *personHandler) Update(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Not implemented 138"})
}

func (h *personHandler) Delete(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Not implemented 139"})
}

func (h *personHandler) Restore(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Not implemented 140"})
}

func (h *personHandler) GetInteractions(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Not implemented 141"})
}

func (h *personHandler) GetHealthScore(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Not implemented 142"})
}

func (h *personHandler) UpdateReminder(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Not implemented 143"})
}

func (h *personHandler) Search(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Not implemented 145"})
}

func (h *personHandler) GetCategories(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Not implemented 146"})
}
