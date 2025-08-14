package handlers

import "github.com/gofiber/fiber/v2"

// UserHandler defines user handler interface
type UserHandler interface {
	// Profile
	GetProfile(c *fiber.Ctx) error
	UpdateProfile(c *fiber.Ctx) error
	DeleteAccount(c *fiber.Ctx) error
	ChangePassword(c *fiber.Ctx) error
	UploadAvatar(c *fiber.Ctx) error
	GetStats(c *fiber.Ctx) error
	GetSettings(c *fiber.Ctx) error
	UpdateSettings(c *fiber.Ctx) error
	
	// OAuth
	LinkOAuthAccount(c *fiber.Ctx) error
	UnlinkOAuthAccount(c *fiber.Ctx) error
	
	// Push notifications
	RegisterPushToken(c *fiber.Ctx) error
	DeactivatePushToken(c *fiber.Ctx) error
	
	// Analytics
	GetDashboard(c *fiber.Ctx) error
	GetMetrics(c *fiber.Ctx) error
	GetTrends(c *fiber.Ctx) error
	TrackEvent(c *fiber.Ctx) error
	GetEvents(c *fiber.Ctx) error
	GetDailyMetrics(c *fiber.Ctx) error
	
	// Search
	GlobalSearch(c *fiber.Ctx) error
	GetSearchSuggestions(c *fiber.Ctx) error
	
	// Notifications
	GetNotificationPreferences(c *fiber.Ctx) error
	UpdateNotificationPreferences(c *fiber.Ctx) error
	SendTestNotification(c *fiber.Ctx) error
	
	// Admin
	AdminListUsers(c *fiber.Ctx) error
	AdminGetUser(c *fiber.Ctx) error
	AdminUpdateUser(c *fiber.Ctx) error
	AdminDeleteUser(c *fiber.Ctx) error
	AdminSuspendUser(c *fiber.Ctx) error
	AdminUnsuspendUser(c *fiber.Ctx) error
	
	// System
	GetSystemStats(c *fiber.Ctx) error
	GetSystemHealth(c *fiber.Ctx) error
	ClearCache(c *fiber.Ctx) error
	GetSystemLogs(c *fiber.Ctx) error
	GetSystemConfig(c *fiber.Ctx) error
	UpdateSystemConfig(c *fiber.Ctx) error
	
	// Analytics Admin
	GetAnalyticsOverview(c *fiber.Ctx) error
	GetUserAnalytics(c *fiber.Ctx) error
	GetEngagementAnalytics(c *fiber.Ctx) error
	GetRetentionAnalytics(c *fiber.Ctx) error
}

// PersonHandler defines person handler interface
type PersonHandler interface {
	List(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Restore(c *fiber.Ctx) error
	GetInteractions(c *fiber.Ctx) error
	GetHealthScore(c *fiber.Ctx) error
	UpdateReminder(c *fiber.Ctx) error
	Search(c *fiber.Ctx) error
	GetCategories(c *fiber.Ctx) error
}

// InteractionHandler defines interaction handler interface
type InteractionHandler interface {
	List(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	GetRecent(c *fiber.Ctx) error
	GetByDate(c *fiber.Ctx) error
	GetEnergyDistribution(c *fiber.Ctx) error
	BulkCreate(c *fiber.Ctx) error
}

// ReflectionHandler defines reflection handler interface
type ReflectionHandler interface {
	List(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	GetToday(c *fiber.Ctx) error
	GetStreak(c *fiber.Ctx) error
	GetPrompts(c *fiber.Ctx) error
	GetMoodTrends(c *fiber.Ctx) error
}

// NudgeHandler defines nudge handler interface
type NudgeHandler interface {
	List(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	MarkSeen(c *fiber.Ctx) error
	MarkActedOn(c *fiber.Ctx) error
	Dismiss(c *fiber.Ctx) error
	GetActive(c *fiber.Ctx) error
	GetHistory(c *fiber.Ctx) error
	GenerateNudges(c *fiber.Ctx) error
}

// GDPRHandler defines GDPR handler interface
type GDPRHandler interface {
	GetConsents(c *fiber.Ctx) error
	UpdateConsent(c *fiber.Ctx) error
	RequestDataExport(c *fiber.Ctx) error
	GetExportStatus(c *fiber.Ctx) error
	DownloadExport(c *fiber.Ctx) error
	DeleteAllData(c *fiber.Ctx) error
	AnonymizeData(c *fiber.Ctx) error
	GetAuditLog(c *fiber.Ctx) error
}
