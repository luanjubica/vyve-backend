package handlers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vyve/vyve-backend/internal/middleware"
	"github.com/vyve/vyve-backend/internal/repository"
	"github.com/vyve/vyve-backend/internal/services"
)

// AuthHandler defines authentication handler interface
type AuthHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	ForgotPassword(c *fiber.Ctx) error
	ResetPassword(c *fiber.Ctx) error
	VerifyEmail(c *fiber.Ctx) error

	// OAuth handlers
	GoogleAuth(c *fiber.Ctx) error
	GoogleCallback(c *fiber.Ctx) error
	LinkedInAuth(c *fiber.Ctx) error
	LinkedInCallback(c *fiber.Ctx) error
	AppleAuth(c *fiber.Ctx) error
	AppleCallback(c *fiber.Ctx) error

	// Token validation (used by middleware)
	ValidateToken(ctx context.Context, token string) (*services.Claims, error)
}

type authHandler struct {
	authService services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService services.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *authHandler) Register(c *fiber.Ctx) error {
	var req services.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validateRegisterRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Register user
	response, err := h.authService.Register(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// Login handles user login
func (h *authHandler) Login(c *fiber.Ctx) error {
	var req services.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Login user
	response, err := h.authService.Login(c.Context(), req)
	if err != nil {
		if err == repository.ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}
		// Log the actual error for debugging
		fmt.Printf("Login error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// Store user activity
	c.Locals("user_id", response.User.ID)

	return c.JSON(response)
}

// RefreshToken handles token refresh
func (h *authHandler) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Refresh token
	response, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return err
	}

	return c.JSON(response)
}

// Logout handles user logout
func (h *authHandler) Logout(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, err := middleware.GetUserID(c)
	if err != nil {
		// If no valid token, still return success (already logged out)
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Logged out successfully",
		})
	}

	// Get refresh token from body (optional)
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	c.BodyParser(&req)

	// Logout user
	if err := h.authService.Logout(c.Context(), userID, req.RefreshToken); err != nil {
		// Log error but don't fail logout
		fmt.Printf("Logout error: %v\n", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully",
	})
}

// ForgotPassword handles forgot password request
func (h *authHandler) ForgotPassword(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Send password reset email
	if err := h.authService.ForgotPassword(c.Context(), req.Email); err != nil {
		// Don't reveal if email exists or not
		return c.JSON(fiber.Map{
			"message": "If the email exists, a password reset link has been sent",
		})
	}

	return c.JSON(fiber.Map{
		"message": "If the email exists, a password reset link has been sent",
	})
}

// ResetPassword handles password reset
func (h *authHandler) ResetPassword(c *fiber.Ctx) error {
	var req struct {
		Token       string `json:"token" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=8"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Reset password
	if err := h.authService.ResetPassword(c.Context(), req.Token, req.NewPassword); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Password reset successfully",
	})
}

// VerifyEmail handles email verification
func (h *authHandler) VerifyEmail(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid verification token",
		})
	}

	// Verify email
	if err := h.authService.VerifyEmail(c.Context(), token); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Email verified successfully",
	})
}

// OAuth handlers

// GoogleAuth initiates Google OAuth flow
func (h *authHandler) GoogleAuth(c *fiber.Ctx) error {
	// TODO: Implement Google OAuth flow initiation
	return c.Redirect("https://accounts.google.com/oauth/authorize?...")
}

// GoogleCallback handles Google OAuth callback
func (h *authHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing authorization code",
		})
	}

	// Handle Google authentication
	response, err := h.authService.HandleGoogleAuth(c.Context(), code)
	if err != nil {
		return err
	}

	return c.JSON(response)
}

// LinkedInAuth initiates LinkedIn OAuth flow
func (h *authHandler) LinkedInAuth(c *fiber.Ctx) error {
	// TODO: Implement LinkedIn OAuth flow initiation
	return c.Redirect("https://www.linkedin.com/oauth/v2/authorization?...")
}

// LinkedInCallback handles LinkedIn OAuth callback
func (h *authHandler) LinkedInCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing authorization code",
		})
	}

	// Handle LinkedIn authentication
	response, err := h.authService.HandleLinkedInAuth(c.Context(), code)
	if err != nil {
		return err
	}

	return c.JSON(response)
}

// AppleAuth initiates Apple OAuth flow
func (h *authHandler) AppleAuth(c *fiber.Ctx) error {
	// TODO: Implement Apple OAuth flow initiation
	return c.Redirect("https://appleid.apple.com/auth/authorize?...")
}

// AppleCallback handles Apple OAuth callback
func (h *authHandler) AppleCallback(c *fiber.Ctx) error {
	// Apple sends a POST request with form data
	code := c.FormValue("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing authorization code",
		})
	}

	// Handle Apple authentication
	response, err := h.authService.HandleAppleAuth(c.Context(), code)
	if err != nil {
		return err
	}

	return c.JSON(response)
}

// ValidateToken validates JWT token (used by middleware)
func (h *authHandler) ValidateToken(ctx context.Context, token string) (*services.Claims, error) {
	return h.authService.ValidateToken(ctx, token)
}

// Helper functions

func validateRegisterRequest(req services.RegisterRequest) error {
	// TODO: Add validation logic
	return nil
}
