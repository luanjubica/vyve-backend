package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"

	"github.com/vyve/vyve-backend/internal/config"
	"github.com/vyve/vyve-backend/internal/handlers"
	"github.com/vyve/vyve-backend/internal/models"
	"github.com/vyve/vyve-backend/internal/realtime"
	"github.com/vyve/vyve-backend/internal/repository"
	"github.com/vyve/vyve-backend/internal/routes"
	"github.com/vyve/vyve-backend/internal/services"
	"github.com/vyve/vyve-backend/pkg/ai"
	"github.com/vyve/vyve-backend/pkg/analytics"
	"github.com/vyve/vyve-backend/pkg/cache"
	"github.com/vyve/vyve-backend/pkg/notifications"
	"github.com/vyve/vyve-backend/pkg/storage"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Version and BuildTime are set during build
var (
	Version   = "dev"
	BuildTime = "unknown"
)

// @title           Vyve API
// @version         1.0
// @description     Backend API for Vyve relationship management app
// @termsOfService  https://vyve.app/terms

// @contact.name   Vyve Support
// @contact.url    https://vyve.app/support
// @contact.email  support@vyve.app

// @license.name  Proprietary
// @license.url   https://vyve.app/license

// @host      api.vyve.app
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func initializeDatabase(cfg *config.Config) (*gorm.DB, *sql.DB, error) {
	dsn := cfg.GetDatabaseURL()
	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, nil, err
	}
	// Pool settings from config
	if cfg.Database.MaxConnections > 0 {
		sqlDB.SetMaxOpenConns(cfg.Database.MaxConnections)
	}
	if cfg.Database.MaxIdleConnections > 0 {
		sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConnections)
	}
	if cfg.Database.ConnectionMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.Database.ConnectionMaxLifetime)
	}
	return gdb, sqlDB, nil
}

func migrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Category{},
		&models.EnergyPattern{},
		&models.CommunicationMethod{},
		&models.RelationshipStatus{},
		&models.Intention{},
		&models.User{},
		&models.AuthProvider{},
		&models.Person{},
		&models.Interaction{},
		&models.Reflection{},
		&models.Nudge{},
		&models.Event{},
		&models.DailyMetric{},
		&models.PushToken{},
		&models.RelationshipAnalysis{},
		&models.AIAnalysisJob{},
	)
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Set up logging
	setupLogging(cfg)

	log.Printf("Starting Vyve API %s (built %s)", Version, BuildTime)
	log.Printf("Environment: %s", cfg.Env)

	// Initialize database and run migrations
	db, sqlDB, err := initializeDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
	}()
	// if err := migrateDatabase(db); err != nil {
	// 	log.Fatalf("Failed to run migrations: %v", err)
	// }

	// Initialize Redis cache
	redisClient, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize services
	storageService := initializeStorage(cfg)
	analyticsService := initializeAnalytics(cfg)
	notificationService := initializeNotifications(cfg)
	aiService := initializeAIService(cfg)

	// Initialize repositories
	repos := repository.NewRepositories(db)

	// Initialize services
	authService := services.NewAuthService(repos.User, redisClient, cfg.JWT, cfg, analyticsService)
	userService := services.NewUserService(repos.User, storageService, analyticsService)
	personService := services.NewPersonService(repos.Person, analyticsService)
	interactionService := services.NewInteractionService(repos.Interaction, repos.Person, analyticsService)
	reflectionService := services.NewReflectionService(repos.Reflection)
	nudgeService := services.NewNudgeService(repos.Nudge, notificationService, analyticsService)
	gdprService := services.NewGDPRService(repos, cfg.Encryption)
	dictionaryService := services.NewDictionaryService(db)
	analysisService := services.NewAnalysisService(aiService, repos.Analysis, repos.Person, repos.Interaction)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	onboardingHandler := handlers.NewOnboardingHandler(userService)
	personHandler := handlers.NewPersonHandler(personService)
	interactionHandler := handlers.NewInteractionHandler(interactionService)
	reflectionHandler := handlers.NewReflectionHandler(reflectionService)
	nudgeHandler := handlers.NewNudgeHandler(nudgeService)
	gdprHandler := handlers.NewGDPRHandler(gdprService)
	dictionaryHandler := handlers.NewDictionaryHandler(dictionaryService)
	analysisHandler := handlers.NewAnalysisHandler(analysisService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:               "Vyve API",
		ServerHeader:          "Vyve",
		DisableStartupMessage: cfg.Env == "production",
		ReadTimeout:           15 * time.Second,
		WriteTimeout:          15 * time.Second,
		IdleTimeout:           60 * time.Second,
		BodyLimit:             10 * 1024 * 1024, // 10MB
		Prefork:               false, // Disabled for containerized deployments (Railway, Docker, etc.)
	})

	// Global middleware
	setupMiddleware(app, cfg)

	// Initialize realtime hub
	hub := realtime.NewHub(redisClient)
	go hub.Run()

	// Setup routes
	routes.Setup(app, &routes.Handlers{
		Auth:        authHandler,
		User:        userHandler,
		Person:      personHandler, // Now handles both person and people operations
		Interaction: interactionHandler,
		Reflection:  reflectionHandler,
		Nudge:       nudgeHandler,
		GDPR:        gdprHandler,
		Realtime:    hub,
		Onboarding:  onboardingHandler,
		Dictionary:  dictionaryHandler,
		Analysis:    analysisHandler,
	}, authService, cfg)

	// Start background workers
	startBackgroundWorkers(cfg, repos, notificationService, analyticsService)

	// Graceful shutdown
	go gracefulShutdown(app)

	// Start server
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupLogging(cfg *config.Config) {
	if cfg.Env == "production" {
		// Production logging - JSON format
		log.SetFlags(0)
		// Could integrate with Sentry or other logging service here
	} else {
		// Development logging
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
}

func setupMiddleware(app *fiber.App, cfg *config.Config) {
	// Request ID
	app.Use(requestid.New())

	// Logger
	if cfg.Logging.Level == "debug" {
		app.Use(logger.New(logger.Config{
			Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
		}))
	}

	// Security headers
	app.Use(helmet.New(helmet.Config{
		XSSProtection:             "1; mode=block",
		ContentTypeNosniff:        "nosniff",
		XFrameOptions:             "SAMEORIGIN",
		HSTSMaxAge:                31536000, // 1 year
		HSTSExcludeSubdomains:     false,
		ContentSecurityPolicy:     "default-src 'self'",
		HSTSPreloadEnabled:        cfg.Env == "production",
	}))

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.Origins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Request-ID",
		ExposeHeaders:    "X-Request-ID",
		AllowCredentials: cfg.CORS.Credentials,
		MaxAge:           86400,
	}))

	// Rate limiting
	if cfg.Env == "production" {
		app.Use(limiter.New(limiter.Config{
			Max:        cfg.RateLimit.Max,
			Expiration: time.Duration(cfg.RateLimit.Window) * time.Second,
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(429).JSON(fiber.Map{
					"error": "Too many requests",
				})
			},
		}))
	}

	// Compression
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Panic recovery
	app.Use(recover.New(recover.Config{
		EnableStackTrace: cfg.Env != "production",
	}))
}

func initializeStorage(cfg *config.Config) storage.Storage {
	if cfg.Env == "production" {
		// Try S3 first
		if cfg.AWS.AccessKeyID != "" && cfg.AWS.SecretAccessKey != "" {
			s3Storage, err := storage.NewS3Storage(cfg.AWS)
			if err != nil {
				log.Printf("Warning: Failed to initialize S3 storage: %v", err)
				log.Println("Storage features will be disabled")
				return nil
			}
			return s3Storage
		}
		log.Println("Warning: S3 credentials not configured, storage features disabled")
		return nil
	}

	// Development mode - try MinIO
	minioStorage, err := storage.NewMinIOStorage(cfg.Storage)
	if err != nil {
		log.Printf("Warning: Failed to initialize MinIO storage: %v", err)
		log.Println("Storage features will be disabled in development")
		return nil
	}
	return minioStorage
}

func initializeAnalytics(cfg *config.Config) analytics.Analytics {
	if cfg.Analytics.AmplitudeKey != "" {
		return analytics.NewAmplitudeAnalytics(cfg.Analytics.AmplitudeKey)
	}
	// Fallback to database analytics
	return analytics.NewDatabaseAnalytics()
}

func initializeNotifications(cfg *config.Config) notifications.NotificationService {
	if cfg.FCM.Key != "" && cfg.FCM.ProjectID != "" {
		fcmService, err := notifications.NewFCMService(cfg.FCM)
		if err != nil {
			log.Printf("Failed to initialize FCM: %v", err)
			return notifications.NewMockNotificationService()
		}
		return fcmService
	}
	return notifications.NewMockNotificationService()
}

func initializeAIService(cfg *config.Config) *ai.Service {
	// Only initialize if AI features are enabled and API keys are configured
	if !cfg.Features.AIInsights {
		log.Println("AI insights feature is disabled")
		return nil
	}

	aiConfig := ai.Config{
		Provider:         cfg.AI.Provider,
		OpenAIKey:        cfg.AI.OpenAIKey,
		OpenAIModel:      cfg.AI.OpenAIModel,
		AnthropicKey:     cfg.AI.AnthropicKey,
		AnthropicModel:   cfg.AI.AnthropicModel,
		MaxTokens:        cfg.AI.MaxTokens,
		Temperature:      cfg.AI.Temperature,
		CacheEnabled:     cfg.AI.CacheEnabled,
		CacheTTL:         cfg.AI.CacheTTL,
		RateLimitPerUser: cfg.AI.RateLimitPerUser,
	}

	aiService, err := ai.NewService(aiConfig)
	if err != nil {
		log.Printf("Failed to initialize AI service: %v", err)
		return nil
	}

	log.Printf("AI service initialized with provider: %s", cfg.AI.Provider)
	return aiService
}

func startBackgroundWorkers(
	cfg *config.Config,
	repos *repository.Repositories,
	notificationService notifications.NotificationService,
	analyticsService analytics.Analytics,
) {
	// Daily reminder worker
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			services.SendDailyReminders(repos.User, notificationService)
		}
	}()

	// Nudge generator worker
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			services.GenerateNudges(repos, analyticsService, notificationService)
		}
	}()

	// Metrics aggregator worker
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			services.AggregateMetrics(repos, analyticsService)
		}
	}()
}

func gracefulShutdown(app *fiber.App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server shutdown complete")
}

func init() {
	// Set timezone to UTC
	os.Setenv("TZ", "UTC")
}
