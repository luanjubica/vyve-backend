package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Env        string
	Server     ServerConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	JWT        JWTConfig
	AWS        AWSConfig
	Storage    StorageConfig
	OAuth      OAuthConfig
	Email      EmailConfig
	FCM        FCMConfig
	Analytics  AnalyticsConfig
	Encryption EncryptionConfig
	CORS       CORSConfig
	RateLimit  RateLimitConfig
	Logging    LoggingConfig
	Features   FeaturesConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host                  string
	Port                  int
	User                  string
	Password              string
	Name                  string
	SSLMode               string
	MaxConnections        int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
	URL                   string // Full connection URL if provided
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	URL      string // Full connection URL if provided
}

type JWTConfig struct {
	Secret               string
	Expiry               time.Duration
	RefreshTokenExpiry   time.Duration
	Issuer               string
}

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	S3Bucket        string
}

type StorageConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
	UseSSL    bool
}

type OAuthConfig struct {
	Google   OAuthProvider
	LinkedIn OAuthProvider
	Apple    AppleOAuthProvider
}

type OAuthProvider struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type AppleOAuthProvider struct {
	ClientID    string
	TeamID      string
	KeyID       string
	PrivateKey  string
	RedirectURL string
}

type EmailConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
	UseTLS   bool
}

type FCMConfig struct {
	Key       string
	ProjectID string
}

type AnalyticsConfig struct {
	AmplitudeKey string
	Enabled      bool
}

type EncryptionConfig struct {
	Enabled bool
	Key     string
}

type CORSConfig struct {
	Origins     string
	Credentials bool
}

type RateLimitConfig struct {
	Max    int
	Window int // seconds
}

type LoggingConfig struct {
	Level     string
	Format    string
	SentryDSN string
}

type FeaturesConfig struct {
	AIInsights         bool
	PushNotifications  bool
	SocialAuth         bool
	DataExport         bool
	EUDataResidency    bool
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Env: getEnv("ENV", "development"),
		
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  getDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		
		Database: DatabaseConfig{
			Host:                  getEnv("DB_HOST", "localhost"),
			Port:                  getEnvAsInt("DB_PORT", 5432),
			User:                  getEnv("DB_USER", "vyve"),
			Password:              getEnv("DB_PASSWORD", "vyve"),
			Name:                  getEnv("DB_NAME", "vyve_dev"),
			SSLMode:               getEnv("DB_SSL_MODE", "disable"),
			MaxConnections:        getEnvAsInt("DB_MAX_CONNECTIONS", 25),
			MaxIdleConnections:    getEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 5),
			ConnectionMaxLifetime: getDuration("DB_CONNECTION_MAX_LIFETIME", 5*time.Minute),
			URL:                   getEnv("DATABASE_URL", ""),
		},
		
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			URL:      getEnv("REDIS_URL", ""),
		},
		
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", "change-this-secret-in-production"),
			Expiry:             getDuration("JWT_EXPIRY", 24*time.Hour),
			RefreshTokenExpiry: getDuration("REFRESH_TOKEN_EXPIRY", 7*24*time.Hour),
			Issuer:             getEnv("JWT_ISSUER", "vyve.app"),
		},
		
		AWS: AWSConfig{
			Region:          getEnv("AWS_REGION", "eu-west-1"),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			S3Bucket:        getEnv("S3_BUCKET", "vyve-production"),
		},
		
		Storage: StorageConfig{
			Endpoint:  getEnv("S3_ENDPOINT", ""),
			AccessKey: getEnv("S3_ACCESS_KEY", ""),
			SecretKey: getEnv("S3_SECRET_KEY", ""),
			Bucket:    getEnv("S3_BUCKET", "vyve-dev"),
			Region:    getEnv("S3_REGION", "us-east-1"),
			UseSSL:    getEnvAsBool("S3_USE_SSL", true),
		},
		
		OAuth: OAuthConfig{
			Google: OAuthProvider{
				ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
			},
			LinkedIn: OAuthProvider{
				ClientID:     getEnv("LINKEDIN_CLIENT_ID", ""),
				ClientSecret: getEnv("LINKEDIN_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("LINKEDIN_REDIRECT_URL", ""),
			},
			Apple: AppleOAuthProvider{
				ClientID:    getEnv("APPLE_CLIENT_ID", ""),
				TeamID:      getEnv("APPLE_TEAM_ID", ""),
				KeyID:       getEnv("APPLE_KEY_ID", ""),
				PrivateKey:  getEnv("APPLE_PRIVATE_KEY", ""),
				RedirectURL: getEnv("APPLE_REDIRECT_URL", ""),
			},
		},
		
		Email: EmailConfig{
			Host:     getEnv("SMTP_HOST", "localhost"),
			Port:     getEnvAsInt("SMTP_PORT", 587),
			User:     getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@vyve.app"),
			UseTLS:   getEnvAsBool("SMTP_USE_TLS", true),
		},
		
		FCM: FCMConfig{
			Key:       getEnv("FCM_KEY", ""),
			ProjectID: getEnv("FCM_PROJECT_ID", ""),
		},
		
		Analytics: AnalyticsConfig{
			AmplitudeKey: getEnv("AMPLITUDE_KEY", ""),
			Enabled:      getEnvAsBool("ANALYTICS_ENABLED", true),
		},
		
		Encryption: EncryptionConfig{
			Enabled: getEnvAsBool("DB_ENCRYPTION", false),
			Key:     getEnv("ENCRYPTION_KEY", ""),
		},
		
		CORS: CORSConfig{
			Origins:     getEnv("CORS_ORIGINS", "*"),
			Credentials: getEnvAsBool("CORS_CREDENTIALS", true),
		},
		
		RateLimit: RateLimitConfig{
			Max:    getEnvAsInt("RATE_LIMIT", 100),
			Window: getEnvAsInt("RATE_LIMIT_WINDOW", 60),
		},
		
		Logging: LoggingConfig{
			Level:     getEnv("LOG_LEVEL", "info"),
			Format:    getEnv("LOG_FORMAT", "json"),
			SentryDSN: getEnv("SENTRY_DSN", ""),
		},
		
		Features: FeaturesConfig{
			AIInsights:        getEnvAsBool("FEATURE_AI_INSIGHTS", false),
			PushNotifications: getEnvAsBool("FEATURE_PUSH_NOTIFICATIONS", true),
			SocialAuth:        getEnvAsBool("FEATURE_SOCIAL_AUTH", true),
			DataExport:        getEnvAsBool("FEATURE_DATA_EXPORT", true),
			EUDataResidency:   getEnvAsBool("EU_DATA_RESIDENCY", false),
		},
	}
}

// GetDatabaseURL returns the database connection URL
func (c *Config) GetDatabaseURL() string {
	if c.Database.URL != "" {
		return c.Database.URL
	}
	
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// GetRedisURL returns the Redis connection URL
func (c *Config) GetRedisURL() string {
	if c.Redis.URL != "" {
		return c.Redis.URL
	}
	
	if c.Redis.Password != "" {
		return fmt.Sprintf(
			"redis://:%s@%s:%d/%d",
			c.Redis.Password,
			c.Redis.Host,
			c.Redis.Port,
			c.Redis.DB,
		)
	}
	
	return fmt.Sprintf(
		"redis://%s:%d/%d",
		c.Redis.Host,
		c.Redis.Port,
		c.Redis.DB,
	)
}

// IsProduction returns true if running in production
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

// IsDevelopment returns true if running in development
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// IsTest returns true if running in test
func (c *Config) IsTest() bool {
	return c.Env == "test"
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}