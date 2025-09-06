package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vyve/vyve-backend/internal/middleware"
	"github.com/vyve/vyve-backend/internal/services"
)

type userHandler struct {
	userService services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

// UserProfileResponse represents the user profile response structure
type UserProfileResponse struct {
	ID               string                 `json:"id"`
	Username         string                 `json:"username"`
	Email            string                 `json:"email"`
	EmailVerified    bool                   `json:"email_verified"`
	AvatarURL        string                 `json:"avatar_url,omitempty"`
	DisplayName      string                 `json:"display_name,omitempty"`
	Bio              string                 `json:"bio,omitempty"`
	Timezone         string                 `json:"timezone"`
	Locale           string                 `json:"locale"`
	LastLoginAt      *time.Time             `json:"last_login_at,omitempty"`
	LastActivityAt   *time.Time             `json:"last_activity_at,omitempty"`
	StreakCount      int                    `json:"streak_count"`
	LastReflectionAt *time.Time             `json:"last_reflection_at,omitempty"`
	Settings         map[string]interface{} `json:"settings"`
	OnboardingCompleted bool                 `json:"onboarding_completed"`
	CreatedAt        time.Time              `json:"created_at"`
}

func (h *userHandler) GetProfile(c *fiber.Ctx) error {
	fmt.Println("GetProfile handler called")
	
	userID, err := middleware.GetUserID(c)
	if err != nil {
		fmt.Printf("Error getting user ID: %v\n", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Unauthorized: Invalid or missing authentication token",
		})
	}

	fmt.Printf("Looking up user with ID: %s\n", userID)
	user, err := h.userService.GetByID(c.Context(), userID)
	if err != nil {
		fmt.Printf("Error getting user by ID: %v\n", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "User not found or you don't have permission to access this resource",
		})
	}

	fmt.Printf("Found user: %+v\n", user)

	// Create a clean response without sensitive data
	response := UserProfileResponse{
		ID:                user.ID.String(),
		Username:          user.Username,
		Email:             user.Email,
		EmailVerified:     user.EmailVerified,
		AvatarURL:         user.AvatarURL,
		DisplayName:       user.DisplayName,
		Bio:               user.Bio,
		Timezone:          user.Timezone,
		Locale:            user.Locale,
		LastLoginAt:       user.LastLoginAt,
		LastActivityAt:    user.LastActivityAt,
		StreakCount:       user.StreakCount,
		LastReflectionAt:  user.LastReflectionAt,
		Settings:          user.Settings,
		OnboardingCompleted: user.OnboardingCompleted,
		CreatedAt:         user.CreatedAt,
	}

	fmt.Println("Sending response with user data")
	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

func (h *userHandler) UpdateProfile(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 1"})
}

func (h *userHandler) DeleteAccount(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 2"})
}

func (h *userHandler) ChangePassword(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 3"})
}

func (h *userHandler) UploadAvatar(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 4"})
}

func (h *userHandler) GetStats(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	stats, err := h.userService.GetStats(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get stats"})
	}

	return c.JSON(stats)
}

func (h *userHandler) GetSettings(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	settings, err := h.userService.GetSettings(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get settings"})
	}

	return c.JSON(settings)
}



func (h *userHandler) UpdateSettings(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Parse settings directly from request body
	var settings map[string]interface{}
	if err := c.BodyParser(&settings); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.userService.UpdateSettings(c.Context(), userID, settings); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update settings"})
	}

	// Return the updated settings
	updatedSettings, err := h.userService.GetSettings(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get updated settings"})
	}

	return c.JSON(fiber.Map{"settings": updatedSettings})
}


// Implement the remaining UserHandler interface methods with stubs
func (h *userHandler) LinkOAuthAccount(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 5"})
}

func (h *userHandler) UnlinkOAuthAccount(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 6"})
}

func (h *userHandler) RegisterPushToken(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 7"})
}

func (h *userHandler) DeactivatePushToken(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 8"})
}

// Analytics methods
func (h *userHandler) GetDashboard(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 9"})
}

func (h *userHandler) GetMetrics(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 10"})
}

func (h *userHandler) GetTrends(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 11"})
}

func (h *userHandler) TrackEvent(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 12"})
}

func (h *userHandler) GetEvents(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 13"})
}

func (h *userHandler) GetDailyMetrics(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 14"})
}

// Search methods
func (h *userHandler) GlobalSearch(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 15"})
}

func (h *userHandler) GetSearchSuggestions(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 16"})
}

// Notification methods
func (h *userHandler) GetNotificationPreferences(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 17"})
}

func (h *userHandler) UpdateNotificationPreferences(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 18"})
}

func (h *userHandler) SendTestNotification(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 19"})
}

// Admin methods
func (h *userHandler) AdminListUsers(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 20"})
}

func (h *userHandler) AdminGetUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 21"})
}

func (h *userHandler) AdminUpdateUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 22"})
}

func (h *userHandler) AdminDeleteUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 23"})
}

func (h *userHandler) AdminSuspendUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 24"})
}

func (h *userHandler) AdminUnsuspendUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 25"})
}

// System methods
func (h *userHandler) GetSystemStats(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 26"})
}

func (h *userHandler) GetSystemHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 27"})
}

func (h *userHandler) ClearCache(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 28"})
}

func (h *userHandler) GetSystemLogs(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 29"})
}

func (h *userHandler) GetSystemConfig(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 30"})
}

func (h *userHandler) UpdateSystemConfig(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 31"})
}

// Analytics methods
func (h *userHandler) GetAnalyticsOverview(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 32"})
}

func (h *userHandler) GetUserAnalytics(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 33"})
}

func (h *userHandler) GetEngagementAnalytics(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 34"})
}

func (h *userHandler) GetRetentionAnalytics(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented 35"})
}
