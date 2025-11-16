package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/vyve/vyve-backend/internal/middleware"
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

// List handles GET /people
func (h *personHandler) List(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	category := c.Query("category")
	search := c.Query("search")
	orderBy := c.Query("order_by", "name")

	opts := services.ListOptions{
		Page:     page,
		Limit:    limit,
		Category: category,
		Search:   search,
		OrderBy:  orderBy,
	}

	people, pagination, err := h.personService.List(c.Context(), userID, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get people"})
	}

	return c.JSON(fiber.Map{
		"success":    true,
		"data":       people,
		"pagination": pagination,
	})
}

// Create handles POST /people
func (h *personHandler) Create(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req services.CreatePersonRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	person, err := h.personService.Create(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create person"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    person,
	})
}

// Get handles GET /people/:id
func (h *personHandler) Get(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	personID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid person ID"})
	}

	person, err := h.personService.GetByID(c.Context(), userID, personID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Person not found"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    person,
	})
}

// Update handles PUT /people/:id
func (h *personHandler) Update(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	personID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid person ID"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	person, err := h.personService.Update(c.Context(), userID, personID, updates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update person"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    person,
	})
}

// Delete handles DELETE /people/:id
func (h *personHandler) Delete(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	personID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid person ID"})
	}

	if err := h.personService.Delete(c.Context(), userID, personID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete person"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Person deleted successfully",
	})
}

// CountPeople handles GET /people/count
func (h *personHandler) CountPeople(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	count, err := h.personService.CountPeople(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get people count"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   count,
	})
}

// Search handles POST /people/search
func (h *personHandler) Search(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		Query string `json:"query"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	people, err := h.personService.Search(c.Context(), userID, req.Query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to search people"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    people,
	})
}

// GetCategories handles GET /people/categories
func (h *personHandler) GetCategories(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	categories, err := h.personService.GetCategories(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get categories"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    categories,
	})
}

// UploadAvatar handles POST /people/:id/upload-avatar
func (h *personHandler) UploadAvatar(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	personID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid person ID"})
	}

	// Check if the request contains a JSON body with avatar_url
	var req struct {
		AvatarURL string `json:"avatar_url"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.AvatarURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "avatar_url is required"})
	}

	// Update the person's avatar URL
	updates := map[string]interface{}{
		"avatar_url": req.AvatarURL,
	}

	person, err := h.personService.Update(c.Context(), userID, personID, updates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update avatar"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    person,
	})
}

// Stub implementations for remaining methods
func (h *personHandler) Restore(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Not implemented yet"})
}

func (h *personHandler) GetInteractions(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Not implemented yet"})
}

func (h *personHandler) GetHealthScore(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Not implemented yet"})
}

func (h *personHandler) UpdateReminder(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Not implemented yet"})
}
