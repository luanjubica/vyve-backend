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
	ID                  string                 `json:"id"`
	Username            string                 `json:"username"`
	Email               string                 `json:"email"`
	EmailVerified       bool                   `json:"email_verified"`
	AvatarURL           string                 `json:"avatar_url,omitempty"`
	DisplayName         string                 `json:"display_name,omitempty"`
	Bio                 string                 `json:"bio,omitempty"`
	Timezone            string                 `json:"timezone"`
	Locale              string                 `json:"locale"`
	LastLoginAt         *time.Time             `json:"last_login_at,omitempty"`
	LastActivityAt      *time.Time             `json:"last_activity_at,omitempty"`
	StreakCount         int                    `json:"streak_count"`
	LastReflectionAt    *time.Time             `json:"last_reflection_at,omitempty"`
	Settings            map[string]interface{} `json:"settings"`
	OnboardingCompleted bool                   `json:"onboarding_completed"`
	CreatedAt           time.Time              `json:"created_at"`
}

// GetProfile handles GET /users/me
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
		ID:                  user.ID.String(),
		Username:            user.Username,
		Email:               user.Email,
		EmailVerified:       user.EmailVerified,
		AvatarURL:           user.AvatarURL,
		DisplayName:         user.DisplayName,
		Bio:                 user.Bio,
		Timezone:            user.Timezone,
		Locale:              user.Locale,
		LastLoginAt:         user.LastLoginAt,
		LastActivityAt:      user.LastActivityAt,
		StreakCount:         user.StreakCount,
		LastReflectionAt:    user.LastReflectionAt,
		Settings:            user.Settings,
		OnboardingCompleted: user.OnboardingCompleted,
		CreatedAt:           user.CreatedAt,
	}

	fmt.Println("Sending response with user data")
	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// UpdateProfile handles PUT /users/me
func (h *userHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, err := h.userService.Update(c.Context(), userID, updates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update profile"})
	}

	response := UserProfileResponse{
		ID:                  user.ID.String(),
		Username:            user.Username,
		Email:               user.Email,
		EmailVerified:       user.EmailVerified,
		AvatarURL:           user.AvatarURL,
		DisplayName:         user.DisplayName,
		Bio:                 user.Bio,
		Timezone:            user.Timezone,
		Locale:              user.Locale,
		LastLoginAt:         user.LastLoginAt,
		LastActivityAt:      user.LastActivityAt,
		StreakCount:         user.StreakCount,
		LastReflectionAt:    user.LastReflectionAt,
		Settings:            user.Settings,
		OnboardingCompleted: user.OnboardingCompleted,
		CreatedAt:           user.CreatedAt,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetSettings handles GET /users/me/settings
func (h *userHandler) GetSettings(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	settings, err := h.userService.GetSettings(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get settings"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    settings,
	})
}

// UpdateSettings handles PUT /users/me/settings
func (h *userHandler) UpdateSettings(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Parse settings from request body
	var settings map[string]interface{}
	if err := c.BodyParser(&settings); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// If client sends a wrapped payload like {"success": true, "data": {...}}, unwrap it
	if raw, ok := settings["data"]; ok {
		if inner, ok := raw.(map[string]interface{}); ok {
			settings = inner
		}
	}
	// Drop any wrapper flags that might have been passed back from previous responses
	delete(settings, "success")

	if err := h.userService.UpdateSettings(c.Context(), userID, settings); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update settings"})
	}

	// Return the updated settings
	updatedSettings, err := h.userService.GetSettings(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get updated settings"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    updatedSettings,
	})
}

// GetStats handles GET /users/me/stats
func (h *userHandler) GetStats(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	stats, err := h.userService.GetStats(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get stats"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}

// DeleteAccount handles DELETE /users/me
func (h *userHandler) DeleteAccount(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := h.userService.Delete(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete account"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Account deleted successfully",
	})
}

// ChangePassword handles POST /users/me/change-password
func (h *userHandler) ChangePassword(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Change password not implemented yet"})
}

// UploadAvatar handles POST /users/me/upload-avatar
func (h *userHandler) UploadAvatar(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Upload avatar not implemented yet"})
}

// Implement the remaining UserHandler interface methods with stubs
func (h *userHandler) LinkOAuthAccount(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "OAuth linking not implemented yet"})
}

func (h *userHandler) UnlinkOAuthAccount(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "OAuth unlinking not implemented yet"})
}

func (h *userHandler) RegisterPushToken(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Push token registration not implemented yet"})
}

func (h *userHandler) DeactivatePushToken(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Push token deactivation not implemented yet"})
}

// Analytics methods
func (h *userHandler) GetDashboard(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Get user stats
	stats, _ := h.userService.GetStats(c.Context(), userID)

	// Return dashboard summary
	dashboard := fiber.Map{
		"user_stats": stats,
		"quick_stats": fiber.Map{
			"total_people":       stats["total_people"],
			"total_interactions": stats["total_interactions"],
			"streak":             stats["streak"],
		},
		"message": "Dashboard data aggregated",
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    dashboard,
	})
}

func (h *userHandler) GetMetrics(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Get time range from query params
	period := c.Query("period", "week") // day, week, month, year

	metrics := fiber.Map{
		"period": period,
		"interactions": fiber.Map{
			"total":   0,
			"average": 0,
		},
		"people": fiber.Map{
			"total":  0,
			"active": 0,
		},
		"engagement": fiber.Map{
			"quality_score": 0,
			"frequency":     0,
		},
	}

	// TODO: Aggregate metrics from database
	_ = userID

	return c.JSON(fiber.Map{
		"success": true,
		"data":    metrics,
	})
}

func (h *userHandler) GetTrends(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	trends := fiber.Map{
		"interaction_frequency": []interface{}{},
		"relationship_health":   []interface{}{},
		"energy_patterns":       []interface{}{},
		"message":               "Trends analysis - aggregation to be implemented",
	}

	// TODO: Calculate trends from historical data
	_ = userID

	return c.JSON(fiber.Map{
		"success": true,
		"data":    trends,
	})
}

func (h *userHandler) TrackEvent(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var event struct {
		Type       string                 `json:"type"`
		Properties map[string]interface{} `json:"properties"`
	}

	if err := c.BodyParser(&event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// TODO: Store event in analytics
	_ = userID

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Event tracked",
	})
}

func (h *userHandler) GetEvents(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// TODO: Retrieve events from database
	_ = userID

	return c.JSON(fiber.Map{
		"success": true,
		"data":    []interface{}{},
	})
}

func (h *userHandler) GetDailyMetrics(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Get date from query params, default to today
	dateStr := c.Query("date", time.Now().Format("2006-01-02"))

	dailyMetrics := fiber.Map{
		"date":         dateStr,
		"interactions": 0,
		"reflections":  0,
		"active_time":  0,
		"quality_score": 0,
	}

	// TODO: Query daily metrics from database
	_ = userID

	return c.JSON(fiber.Map{
		"success": true,
		"data":    dailyMetrics,
	})
}

// Search methods
func (h *userHandler) GlobalSearch(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Global search not implemented yet"})
}

func (h *userHandler) GetSearchSuggestions(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Search suggestions not implemented yet"})
}

// Notification methods
func (h *userHandler) GetNotificationPreferences(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Notification preferences not implemented yet"})
}

func (h *userHandler) UpdateNotificationPreferences(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Update notification preferences not implemented yet"})
}

func (h *userHandler) SendTestNotification(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Test notification not implemented yet"})
}

// Admin methods
func (h *userHandler) AdminListUsers(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Admin list users not implemented yet"})
}

func (h *userHandler) AdminGetUser(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Admin get user not implemented yet"})
}

func (h *userHandler) AdminUpdateUser(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Admin update user not implemented yet"})
}

func (h *userHandler) AdminDeleteUser(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Admin delete user not implemented yet"})
}

func (h *userHandler) AdminSuspendUser(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Admin suspend user not implemented yet"})
}

func (h *userHandler) AdminUnsuspendUser(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Admin unsuspend user not implemented yet"})
}

// System methods
func (h *userHandler) GetSystemStats(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "System stats not implemented yet"})
}

func (h *userHandler) GetSystemHealth(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "System health not implemented yet"})
}

func (h *userHandler) ClearCache(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Clear cache not implemented yet"})
}

func (h *userHandler) GetSystemLogs(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "System logs not implemented yet"})
}

func (h *userHandler) GetSystemConfig(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "System config not implemented yet"})
}

func (h *userHandler) UpdateSystemConfig(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "Update system config not implemented yet"})
}

// Analytics methods
func (h *userHandler) GetAnalyticsOverview(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Get comprehensive analytics overview
	stats, _ := h.userService.GetStats(c.Context(), userID)

	overview := fiber.Map{
		"summary": stats,
		"period":  "all_time",
		"charts": fiber.Map{
			"interactions_over_time": []interface{}{},
			"relationship_health":    []interface{}{},
			"engagement_levels":      []interface{}{},
		},
		"message": "Analytics overview - detailed aggregation to be implemented",
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    overview,
	})
}

func (h *userHandler) GetUserAnalytics(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// TODO: Get detailed user behavior analytics
	_ = userID

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"message": "User analytics - to be implemented",
		},
	})
}

func (h *userHandler) GetEngagementAnalytics(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// TODO: Calculate engagement metrics
	_ = userID

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"message": "Engagement analytics - to be implemented",
		},
	})
}

func (h *userHandler) GetRetentionAnalytics(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// TODO: Calculate retention metrics
	_ = userID

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"message": "Retention analytics - to be implemented",
		},
	})
}
