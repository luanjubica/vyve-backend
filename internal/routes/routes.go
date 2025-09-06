package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"

	"github.com/vyve/vyve-backend/internal/config"
	"github.com/vyve/vyve-backend/internal/handlers"
	"github.com/vyve/vyve-backend/internal/middleware"
	"github.com/vyve/vyve-backend/internal/realtime"
	"github.com/vyve/vyve-backend/internal/services"
)

// Handlers holds all handler instances
type Handlers struct {
	Auth        handlers.AuthHandler
	User        handlers.UserHandler
	Person      handlers.PersonHandler // Handles both person and people operations
	Interaction handlers.InteractionHandler
	Reflection  handlers.ReflectionHandler
	Nudge       handlers.NudgeHandler
	GDPR        handlers.GDPRHandler
	Realtime    *realtime.Hub
	Onboarding  handlers.OnboardingHandler
}

// Setup sets up all routes
func Setup(app *fiber.App, h *Handlers, authService services.AuthService, cfg *config.Config) {
	// Create token validator function from auth service
	validateToken := authService.ValidateToken

	// Health check
	app.Get("/health", healthCheck)

	// Metrics endpoint (only in development)
	if cfg.IsDevelopment() {
		app.Get("/metrics", monitor.New())
	}

	// API v1 routes
	api := app.Group("/api/v1")

	// Public routes (no authentication required)
	setupPublicRoutes(api, h)

	// Protected routes (authentication required)
	protected := api.Use(middleware.AuthMiddleware(validateToken))
	setupProtectedRoutes(protected, h)

	// Admin routes (admin role required)
	admin := protected.Use(middleware.RequireRole("admin"))
	setupAdminRoutes(admin, h)

	// Real-time endpoints
	setupRealtimeRoutes(app, h, validateToken)

	// Static files (if needed)
	// app.Static("/uploads", "./uploads")
}

// setupPublicRoutes sets up public API routes
func setupPublicRoutes(api fiber.Router, h *Handlers) {
	// Authentication
	auth := api.Group("/auth")
	{
		auth.Post("/signup", h.Auth.Register)   // New signup endpoint
		auth.Post("/register", h.Auth.Register) // Keep for backward compatibility
		auth.Post("/login", h.Auth.Login)
		auth.Post("/refresh", h.Auth.RefreshToken)
		auth.Post("/logout", h.Auth.Logout)
		auth.Post("/forgot-password", h.Auth.ForgotPassword)
		auth.Post("/reset-password", h.Auth.ResetPassword)
		auth.Get("/verify-email/:token", h.Auth.VerifyEmail)

		// OAuth
		auth.Get("/google", h.Auth.GoogleAuth)
		auth.Get("/google/callback", h.Auth.GoogleCallback)
		auth.Get("/linkedin", h.Auth.LinkedInAuth)
		auth.Get("/linkedin/callback", h.Auth.LinkedInCallback)
		auth.Get("/apple", h.Auth.AppleAuth)
		auth.Post("/apple/callback", h.Auth.AppleCallback) // Apple uses POST
	}
}

// setupProtectedRoutes sets up protected API routes
func setupProtectedRoutes(api fiber.Router, h *Handlers) {
	// User profile - FIXED: Removed duplicate routes
	user := api.Group("/users/me")
	{
		// Core user profile endpoints
		user.Get("", h.User.GetProfile)       // GET /users/me
		user.Put("", h.User.UpdateProfile)    // PUT /users/me
		user.Delete("", h.User.DeleteAccount) // DELETE /users/me

		// User settings endpoints
		user.Get("/settings", h.User.GetSettings)    // GET /users/me/settings
		user.Put("/settings", h.User.UpdateSettings) // PUT /users/me/settings

		// User stats and other endpoints
		user.Get("/stats", h.User.GetStats) // GET /users/me/stats

		// Authentication related endpoints
		user.Post("/change-password", h.User.ChangePassword) // POST /users/me/change-password
		user.Post("/upload-avatar", h.User.UploadAvatar)     // POST /users/me/upload-avatar

		// Onboarding
		user.Get("/onboarding", h.Onboarding.GetOnboardingStatus)          // GET /users/me/onboarding
		user.Post("/onboarding/complete", h.Onboarding.CompleteOnboarding) // POST /users/me/onboarding/complete

		// OAuth account linking
		user.Post("/oauth/link/:provider", h.User.LinkOAuthAccount)       // POST /users/me/oauth/link/:provider
		user.Delete("/oauth/unlink/:provider", h.User.UnlinkOAuthAccount) // DELETE /users/me/oauth/unlink/:provider

		// Push notifications
		user.Post("/push-token", h.User.RegisterPushToken)            // POST /users/me/push-token
		user.Delete("/push-token/:token", h.User.DeactivatePushToken) // DELETE /users/me/push-token/:token
	}

	// People (relationships) - FIXED: Proper ordering to avoid conflicts
	people := api.Group("/people")
	{
		// Special endpoints MUST come before /:id routes
		people.Get("/count", h.Person.CountPeople)        // GET /people/count
		people.Get("/categories", h.Person.GetCategories) // GET /people/categories
		people.Post("/search", h.Person.Search)           // POST /people/search

		// General CRUD operations
		people.Get("", h.Person.List)    // GET /people
		people.Post("", h.Person.Create) // POST /people

		// Individual person operations (MUST come after special endpoints)
		people.Get("/:id", h.Person.Get)                          // GET /people/:id
		people.Put("/:id", h.Person.Update)                       // PUT /people/:id
		people.Delete("/:id", h.Person.Delete)                    // DELETE /people/:id
		people.Post("/:id/restore", h.Person.Restore)             // POST /people/:id/restore
		people.Get("/:id/interactions", h.Person.GetInteractions) // GET /people/:id/interactions
		people.Get("/:id/health", h.Person.GetHealthScore)        // GET /people/:id/health
		people.Put("/:id/reminder", h.Person.UpdateReminder)      // PUT /people/:id/reminder
	}

	// Interactions (vyves)
	interactions := api.Group("/interactions")
	{
		interactions.Get("", h.Interaction.List)                                      // GET /interactions
		interactions.Post("", h.Interaction.Create)                                   // POST /interactions
		interactions.Get("/recent", h.Interaction.GetRecent)                          // GET /interactions/recent
		interactions.Get("/by-date", h.Interaction.GetByDate)                         // GET /interactions/by-date
		interactions.Get("/energy-distribution", h.Interaction.GetEnergyDistribution) // GET /interactions/energy-distribution
		interactions.Post("/bulk", h.Interaction.BulkCreate)                          // POST /interactions/bulk

		// Individual interaction operations
		interactions.Get("/:id", h.Interaction.Get)       // GET /interactions/:id
		interactions.Put("/:id", h.Interaction.Update)    // PUT /interactions/:id
		interactions.Delete("/:id", h.Interaction.Delete) // DELETE /interactions/:id
	}

	// Reflections
	reflections := api.Group("/reflections")
	{
		reflections.Get("", h.Reflection.List)                // GET /reflections
		reflections.Post("", h.Reflection.Create)             // POST /reflections
		reflections.Get("/today", h.Reflection.GetToday)      // GET /reflections/today
		reflections.Get("/streak", h.Reflection.GetStreak)    // GET /reflections/streak
		reflections.Get("/prompts", h.Reflection.GetPrompts)  // GET /reflections/prompts
		reflections.Get("/moods", h.Reflection.GetMoodTrends) // GET /reflections/moods

		// Individual reflection operations
		reflections.Get("/:id", h.Reflection.Get)       // GET /reflections/:id
		reflections.Put("/:id", h.Reflection.Update)    // PUT /reflections/:id
		reflections.Delete("/:id", h.Reflection.Delete) // DELETE /reflections/:id
	}

	// Nudges (AI insights)
	nudges := api.Group("/nudges")
	{
		nudges.Get("", h.Nudge.List)                     // GET /nudges
		nudges.Get("/active", h.Nudge.GetActive)         // GET /nudges/active
		nudges.Get("/history", h.Nudge.GetHistory)       // GET /nudges/history
		nudges.Post("/generate", h.Nudge.GenerateNudges) // POST /nudges/generate

		// Individual nudge operations
		nudges.Get("/:id", h.Nudge.Get)              // GET /nudges/:id
		nudges.Post("/:id/seen", h.Nudge.MarkSeen)   // POST /nudges/:id/seen
		nudges.Post("/:id/act", h.Nudge.MarkActedOn) // POST /nudges/:id/act
		nudges.Post("/:id/dismiss", h.Nudge.Dismiss) // POST /nudges/:id/dismiss
	}

	// Analytics & Insights
	analytics := api.Group("/analytics")
	{
		analytics.Get("/dashboard", h.User.GetDashboard)        // GET /analytics/dashboard
		analytics.Get("/metrics", h.User.GetMetrics)            // GET /analytics/metrics
		analytics.Get("/trends", h.User.GetTrends)              // GET /analytics/trends
		analytics.Post("/event", h.User.TrackEvent)             // POST /analytics/event
		analytics.Get("/events", h.User.GetEvents)              // GET /analytics/events
		analytics.Get("/daily-metrics", h.User.GetDailyMetrics) // GET /analytics/daily-metrics
	}

	// GDPR & Privacy
	gdpr := api.Group("/gdpr")
	{
		gdpr.Get("/consent", h.GDPR.GetConsents)                // GET /gdpr/consent
		gdpr.Post("/consent", h.GDPR.UpdateConsent)             // POST /gdpr/consent
		gdpr.Post("/export", h.GDPR.RequestDataExport)          // POST /gdpr/export
		gdpr.Get("/export/:id", h.GDPR.GetExportStatus)         // GET /gdpr/export/:id
		gdpr.Get("/export/:id/download", h.GDPR.DownloadExport) // GET /gdpr/export/:id/download
		gdpr.Delete("/data", h.GDPR.DeleteAllData)              // DELETE /gdpr/data
		gdpr.Post("/anonymize", h.GDPR.AnonymizeData)           // POST /gdpr/anonymize
		gdpr.Get("/audit-log", h.GDPR.GetAuditLog)              // GET /gdpr/audit-log
	}

	// Search
	search := api.Group("/search")
	{
		search.Get("", h.User.GlobalSearch)                     // GET /search
		search.Get("/suggestions", h.User.GetSearchSuggestions) // GET /search/suggestions
	}

	// Notifications
	notifications := api.Group("/notifications")
	{
		notifications.Get("/preferences", h.User.GetNotificationPreferences)    // GET /notifications/preferences
		notifications.Put("/preferences", h.User.UpdateNotificationPreferences) // PUT /notifications/preferences
		notifications.Post("/test", h.User.SendTestNotification)                // POST /notifications/test
	}
}

// setupAdminRoutes sets up admin API routes
func setupAdminRoutes(api fiber.Router, h *Handlers) {
	// User management
	users := api.Group("/users")
	{
		users.Get("/", h.User.AdminListUsers)
		users.Get("/:id", h.User.AdminGetUser)
		users.Put("/:id", h.User.AdminUpdateUser)
		users.Delete("/:id", h.User.AdminDeleteUser)
		users.Post("/:id/suspend", h.User.AdminSuspendUser)
		users.Post("/:id/unsuspend", h.User.AdminUnsuspendUser)
	}

	// System
	system := api.Group("/system")
	{
		system.Get("/stats", h.User.GetSystemStats)
		system.Get("/health", h.User.GetSystemHealth)
		system.Post("/cache/clear", h.User.ClearCache)
		system.Get("/logs", h.User.GetSystemLogs)
		system.Get("/config", h.User.GetSystemConfig)
		system.Put("/config", h.User.UpdateSystemConfig)
	}

	// Analytics
	analytics := api.Group("/analytics")
	{
		analytics.Get("/overview", h.User.GetAnalyticsOverview)
		analytics.Get("/users", h.User.GetUserAnalytics)
		analytics.Get("/engagement", h.User.GetEngagementAnalytics)
		analytics.Get("/retention", h.User.GetRetentionAnalytics)
	}
}

// setupRealtimeRoutes sets up real-time communication routes
func setupRealtimeRoutes(app *fiber.App, h *Handlers, validateToken middleware.TokenValidator) {
	// WebSocket endpoint
	app.Get("/ws", middleware.WebSocketUpgrade(), h.Realtime.HandleWebSocket)

	// Server-Sent Events endpoint (requires authentication)
	app.Get("/sse", middleware.AuthMiddleware(validateToken), h.Realtime.HandleSSE)
}

// healthCheck handles health check requests
func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "healthy",
		"service": "vyve-api",
		"version": "1.0.0",
		"time":    c.Context().Time().Unix(),
	})
}
