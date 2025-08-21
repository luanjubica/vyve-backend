package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/vyve/vyve-backend/internal/services"
)

// OnboardingHandler defines the HTTP handlers for onboarding-related endpoints
type OnboardingHandler interface {
	// GetOnboardingStatus handles GET /users/me/onboarding
	GetOnboardingStatus(c *fiber.Ctx) error
	// CompleteOnboarding handles POST /users/me/onboarding/complete
	CompleteOnboarding(c *fiber.Ctx) error
}

type onboardingHandler struct {
	userService services.UserService
}

// NewOnboardingHandler creates a new onboarding handler
func NewOnboardingHandler(userService services.UserService) OnboardingHandler {
	return &onboardingHandler{
		userService: userService,
	}
}

// GetOnboardingStatus handles GET /users/me/onboarding
func (h *onboardingHandler) GetOnboardingStatus(c *fiber.Ctx) error {
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

	// Get onboarding status
	status, err := h.userService.GetOnboardingStatus(c.Context(), userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get onboarding status",
		})
	}

	return c.JSON(status)
}

// CompleteOnboardingRequest represents the request body for completing onboarding
type CompleteOnboardingRequest struct {
	Completed bool   `json:"completed"`
	Step     string `json:"step,omitempty"`
}

// CompleteOnboarding handles POST /users/me/onboarding/complete
func (h *onboardingHandler) CompleteOnboarding(c *fiber.Ctx) error {
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

	// Parse request body
	var req CompleteOnboardingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update onboarding status
	status, err := h.userService.UpdateOnboardingStatus(c.Context(), userUUID, req.Completed, req.Step)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update onboarding status",
		})
	}

	return c.JSON(status)
}
