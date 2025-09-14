package middleware

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/vyve/vyve-backend/internal/services"
)

// ContextKey type for context keys
type ContextKey string

const (
	// UserContextKey is the key for user in context
	UserContextKey ContextKey = "user"
	// ClaimsContextKey is the key for claims in context
	ClaimsContextKey ContextKey = "claims"
	// RequestIDContextKey is the key for request ID in context
	RequestIDContextKey ContextKey = "request_id"
)

// TokenValidator defines the function signature for token validation
type TokenValidator func(ctx context.Context, token string) (*services.Claims, error)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(validateToken TokenValidator) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// Extract token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		token := parts[1]

		// Validate token
		claims, err := validateToken(c.Context(), token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Invalid or expired token",
				"details": err.Error(),
			})
		}

		log.Printf("[AUTH] Token validated successfully for user: %s", claims.UserID)

		// Set user context
		c.Locals("user_id", claims.UserID)
		c.Locals("claims", claims)

		return c.Next()
	}
}

// OptionalAuthMiddleware allows both authenticated and unauthenticated requests
func OptionalAuthMiddleware(validateToken TokenValidator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		// Extract token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		token := parts[1]

		// Validate token
		claims, err := validateToken(c.Context(), token)
		if err == nil {
			// Set user context if valid
			c.Locals("user_id", claims.UserID)
			c.Locals("claims", claims)
		}

		return c.Next()
	}
}

// RequireRole checks if user has required role
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("claims").(*services.Claims)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied",
			})
		}

		// Check if user has any of the required roles
		// TODO: Add role checking logic based on your requirements
		// For now, we'll skip this check
		_ = claims
		_ = roles

		return c.Next()
	}
}

// RequireEmailVerified checks if user's email is verified
func RequireEmailVerified() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Check if user's email is verified
		// For now, we'll skip this check
		return c.Next()
	}
}

// GetUserID gets user ID from context
func GetUserID(c *fiber.Ctx) (uuid.UUID, error) {

	// Try to get user_id from locals
	userIDValue := c.Locals("user_id")

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	return userID, nil
}

// GetClaims gets JWT claims from context
func GetClaims(c *fiber.Ctx) (*services.Claims, error) {
	claims, ok := c.Locals("claims").(*services.Claims)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid claims")
	}
	return claims, nil
}

// WebSocketUpgrade handles WebSocket upgrade
func WebSocketUpgrade() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Get("Upgrade") != "websocket" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "WebSocket upgrade required",
			})
		}
		return c.Next()
	}
}

// SetUserContext sets user context for downstream handlers
func SetUserContext(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, UserContextKey, userID)
}

// GetUserFromContext gets user ID from context
func GetUserFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserContextKey).(uuid.UUID)
	return userID, ok
}
