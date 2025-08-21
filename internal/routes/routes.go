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
	Person      handlers.PersonHandler  // Handles both person and people operations
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
		auth.Post("/signup", h.Auth.Register) // New signup endpoint
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
	// User profile
	user := api.Group("/users/me")
	{
		// User profile endpoints
		user.Get("", h.User.GetProfile)           // GET /users/me
		user.Put("", h.User.UpdateProfile)         // PUT /users/me
		user.Delete("", h.User.DeleteAccount)      // DELETE /users/me
		
		// User settings endpoints
		settings := user.Group("/settings")
		{
			settings.Get("", h.User.GetSettings)     // GET /users/me/settings
			settings.Put("", h.User.UpdateSettings)   // PUT /users/me/settings
		}
		
		// User stats endpoint
		user.Get("/stats", h.User.GetStats)        // GET /users/me/stats
		
		// Authentication related endpoints
		user.Post("/change-password", h.User.ChangePassword)  // POST /users/me/change-password
		user.Post("/upload-avatar", h.User.UploadAvatar)      // POST /users/me/upload-avatar
		
		// Onboarding
		onboarding := user.Group("/onboarding")
		{
			onboarding.Get("", h.Onboarding.GetOnboardingStatus)     // GET /users/me/onboarding
			onboarding.Post("/complete", h.Onboarding.CompleteOnboarding) // POST /users/me/onboarding/complete
		}
		
		// OAuth account linking
		oauth := user.Group("/oauth")
		{
			oauth.Post("/link/:provider", h.User.LinkOAuthAccount)     // POST /users/me/oauth/link/:provider
			oauth.Delete("/unlink/:provider", h.User.UnlinkOAuthAccount) // DELETE /users/me/oauth/unlink/:provider
		}
		
		// Push notifications
		user.Post("/push-token", h.User.RegisterPushToken)            // POST /users/me/push-token
		user.Delete("/push-token/:token", h.User.DeactivatePushToken) // DELETE /users/me/push-token/:token
	}
	
	// People (relationships)
	people := api.Group("/people")
	{
		people.Get("", h.Person.List)           // GET /people
		people.Post("", h.Person.Create)        // POST /people
		people.Get("/count", h.Person.CountPeople) // GET /people/count
		people.Get("/:id", h.Person.Get)
		people.Put("/:id", h.Person.Update)
		people.Delete("/:id", h.Person.Delete)
		people.Post("/:id/restore", h.Person.Restore)
		people.Get("/:id/interactions", h.Person.GetInteractions)
		people.Get("/:id/health", h.Person.GetHealthScore)
		people.Put("/:id/reminder", h.Person.UpdateReminder)
		people.Post("/search", h.Person.Search)
		people.Get("/categories", h.Person.GetCategories)
	}
	
	// Interactions (vyves)
	interactions := api.Group("/interactions")
	{
		interactions.Get("/", h.Interaction.List)
		interactions.Post("/", h.Interaction.Create)
		interactions.Get("/:id", h.Interaction.Get)
		interactions.Put("/:id", h.Interaction.Update)
		interactions.Delete("/:id", h.Interaction.Delete)
		interactions.Get("/recent", h.Interaction.GetRecent)
		interactions.Get("/by-date", h.Interaction.GetByDate)
		interactions.Get("/energy-distribution", h.Interaction.GetEnergyDistribution)
		interactions.Post("/bulk", h.Interaction.BulkCreate)
	}
	
	// Reflections
	reflections := api.Group("/reflections")
	{
		reflections.Get("/", h.Reflection.List)
		reflections.Post("/", h.Reflection.Create)
		reflections.Get("/:id", h.Reflection.Get)
		reflections.Put("/:id", h.Reflection.Update)
		reflections.Delete("/:id", h.Reflection.Delete)
		reflections.Get("/today", h.Reflection.GetToday)
		reflections.Get("/streak", h.Reflection.GetStreak)
		reflections.Get("/prompts", h.Reflection.GetPrompts)
		reflections.Get("/moods", h.Reflection.GetMoodTrends)
	}
	
	// Nudges (AI insights)
	nudges := api.Group("/nudges")
	{
		nudges.Get("/", h.Nudge.List)
		nudges.Get("/:id", h.Nudge.Get)
		nudges.Post("/:id/seen", h.Nudge.MarkSeen)
		nudges.Post("/:id/act", h.Nudge.MarkActedOn)
		nudges.Post("/:id/dismiss", h.Nudge.Dismiss)
		nudges.Get("/active", h.Nudge.GetActive)
		nudges.Get("/history", h.Nudge.GetHistory)
		nudges.Post("/generate", h.Nudge.GenerateNudges)
	}
	
	// Analytics & Insights
	analytics := api.Group("/analytics")
	{
		analytics.Get("/dashboard", h.User.GetDashboard)
		analytics.Get("/metrics", h.User.GetMetrics)
		analytics.Get("/trends", h.User.GetTrends)
		analytics.Post("/event", h.User.TrackEvent)
		analytics.Get("/events", h.User.GetEvents)
		analytics.Get("/daily-metrics", h.User.GetDailyMetrics)
	}
	
	// GDPR & Privacy
	gdpr := api.Group("/gdpr")
	{
		gdpr.Get("/consent", h.GDPR.GetConsents)
		gdpr.Post("/consent", h.GDPR.UpdateConsent)
		gdpr.Post("/export", h.GDPR.RequestDataExport)
		gdpr.Get("/export/:id", h.GDPR.GetExportStatus)
		gdpr.Get("/export/:id/download", h.GDPR.DownloadExport)
		gdpr.Delete("/data", h.GDPR.DeleteAllData)
		gdpr.Post("/anonymize", h.GDPR.AnonymizeData)
		gdpr.Get("/audit-log", h.GDPR.GetAuditLog)
	}
	
	// Search
	search := api.Group("/search")
	{
		search.Get("/", h.User.GlobalSearch)
		search.Get("/suggestions", h.User.GetSearchSuggestions)
	}
	
	// Notifications
	notifications := api.Group("/notifications")
	{
		notifications.Get("/preferences", h.User.GetNotificationPreferences)
		notifications.Put("/preferences", h.User.UpdateNotificationPreferences)
		notifications.Post("/test", h.User.SendTestNotification)
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