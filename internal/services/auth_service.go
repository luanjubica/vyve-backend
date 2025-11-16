package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	oauth2 "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleoauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"

	"github.com/vyve/vyve-backend/internal/config"
	"github.com/vyve/vyve-backend/internal/models"
	"github.com/vyve/vyve-backend/internal/repository"
	"github.com/vyve/vyve-backend/pkg/analytics"
	"github.com/vyve/vyve-backend/pkg/cache"
	"github.com/vyve/vyve-backend/pkg/utils"
)

// AuthService handles authentication logic
type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)
	Logout(ctx context.Context, userID uuid.UUID, token string) error
	LogoutAll(ctx context.Context, userID uuid.UUID) error
	VerifyEmail(ctx context.Context, token string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error

	// OAuth methods
	GetGoogleAuthURL(state string) string
	HandleGoogleAuth(ctx context.Context, code string) (*AuthResponse, error)
	HandleLinkedInAuth(ctx context.Context, code string) (*AuthResponse, error)
	HandleAppleAuth(ctx context.Context, code string) (*AuthResponse, error)
	LinkOAuthAccount(ctx context.Context, userID uuid.UUID, provider string, code string) error
	UnlinkOAuthAccount(ctx context.Context, userID uuid.UUID, provider string) error

	// Token methods
	GenerateTokenPair(ctx context.Context, user *models.User) (*TokenPair, error)
	ValidateToken(ctx context.Context, tokenString string) (*Claims, error)
	RevokeToken(ctx context.Context, token string) error

	// Session management
	CreateSession(ctx context.Context, userID uuid.UUID, metadata SessionMetadata) (*Session, error)
	GetSession(ctx context.Context, sessionID string) (*Session, error)
	EndSession(ctx context.Context, sessionID string) error
}

type authService struct {
	userRepo  repository.UserRepository
	cache     cache.Cache
	jwtCfg    config.JWTConfig
	cfg       *config.Config
	analytics analytics.Analytics
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, cache cache.Cache, jwtCfg config.JWTConfig, cfg *config.Config, analyticsService analytics.Analytics) AuthService {
	return &authService{
		userRepo:  userRepo,
		cache:     cache,
		jwtCfg:    jwtCfg,
		cfg:       cfg,
		analytics: analyticsService,
	}
}

// Request/Response DTOs
type RegisterRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=50"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	DisplayName string `json:"display_name" validate:"omitempty,min=2,max=50"`
}

type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User         *UserDTO `json:"user"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"`
	TokenType    string   `json:"token_type"`
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	SessionID string    `json:"session_id"`
	jwt.RegisteredClaims
}

type Session struct {
	ID        string          `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	CreatedAt time.Time       `json:"created_at"`
	ExpiresAt time.Time       `json:"expires_at"`
	Metadata  SessionMetadata `json:"metadata"`
}

type SessionMetadata struct {
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	DeviceID  string `json:"device_id,omitempty"`
}

type UserDTO struct {
	ID          uuid.UUID  `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	DisplayName string     `json:"display_name"`
	AvatarURL   string     `json:"avatar_url"`
	Bio         string     `json:"bio"`
	Timezone    string     `json:"timezone"`
	Locale      string     `json:"locale"`
	StreakCount int        `json:"streak_count"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// Register registers a new user
func (s *authService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	// Check if username exists
	exists, err := s.userRepo.CheckUsernameExists(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already taken")
	}

	// Check if email exists
	exists, err = s.userRepo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Set display name to username if not provided
	displayName := req.DisplayName
	if displayName == "" {
		displayName = req.Username
	}

	// Create user
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		DisplayName:  displayName,
		Timezone:     "UTC",
		Locale:       "en",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate tokens
	tokenPair, err := s.GenerateTokenPair(ctx, user)
	if err != nil {
		return nil, err
	}

	// Send verification email
	go s.sendVerificationEmail(user.Email, user.ID.String())

	// Track sign up event
	go s.analytics.Track(ctx, analytics.Event{
		UserID:    user.ID.String(),
		EventType: analytics.EventUserSignUp,
		Properties: map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
		},
		Timestamp: time.Now(),
	})

	return &AuthResponse{
		User:         s.mapUserToDTO(user),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		TokenType:    "Bearer",
	}, nil
}

// Login authenticates a user
func (s *authService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	log.Printf("Login attempt for email: %s", req.Email)
	var user *models.User
	var err error

	// Find user by email or username
	if req.Email != "" {
		log.Printf("Looking up user by email: %s", req.Email)
		user, err = s.userRepo.FindByEmail(ctx, req.Email)
	} else if req.Username != "" {
		log.Printf("Looking up user by username: %s", req.Username)
		user, err = s.userRepo.FindByUsername(ctx, req.Username)
	} else {
		log.Println("No email or username provided")
		return nil, errors.New("email or username required")
	}

	if err != nil {
		log.Printf("Error finding user: %v", err)
		if repository.IsNotFound(err) {
			log.Println("User not found")
			return nil, repository.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	log.Printf("User found: ID=%s, Email=%s, Username=%s, EmailVerified=%v",
		user.ID, user.Email, user.Username, user.EmailVerified)

	// Check password
	log.Printf("user hash %v", user.PasswordHash)
	log.Printf("user password %v", req.Password)

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		log.Printf("Password comparison failed: %v", err)
		log.Printf("Stored hash length: %d, starts with: %s...",
			len(user.PasswordHash), safeSubstring(user.PasswordHash, 0, 10))
		return nil, repository.ErrInvalidCredentials
	}

	log.Println("Password check passed")

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update last login: %v\n", err)
	}

	// Generate tokens
	tokenPair, err := s.GenerateTokenPair(ctx, user)
	if err != nil {
		return nil, err
	}

	// Track login event
	go s.analytics.Track(ctx, analytics.Event{
		UserID:    user.ID.String(),
		EventType: analytics.EventUserLogin,
		Properties: map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
		},
		Timestamp: time.Now(),
	})

	return &AuthResponse{
		User:         s.mapUserToDTO(user),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		TokenType:    "Bearer",
	}, nil
}

// RefreshToken refreshes access token using refresh token
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	// Find refresh token
	token, err := s.userRepo.FindRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, repository.ErrTokenInvalid
	}

	// Check if expired
	if time.Now().After(token.ExpiresAt) {
		return nil, repository.ErrTokenExpired
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, token.UserID)
	if err != nil {
		return nil, err
	}

	// Revoke old refresh token
	if err := s.userRepo.RevokeRefreshToken(ctx, refreshToken); err != nil {
		// Log error but continue
		fmt.Printf("Failed to revoke old refresh token: %v\n", err)
	}

	// Generate new tokens
	tokenPair, err := s.GenerateTokenPair(ctx, user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         s.mapUserToDTO(user),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		TokenType:    "Bearer",
	}, nil
}

// Logout logs out a user
func (s *authService) Logout(ctx context.Context, userID uuid.UUID, token string) error {
	// Track logout event
	go s.analytics.Track(ctx, analytics.Event{
		UserID:    userID.String(),
		EventType: analytics.EventUserLogout,
		Timestamp: time.Now(),
	})

	// Revoke refresh token if provided
	if token != "" {
		if err := s.userRepo.RevokeRefreshToken(ctx, token); err != nil {
			// Log error but don't fail logout
			fmt.Printf("Failed to revoke refresh token: %v\n", err)
		}
	}

	// Clear session from cache
	sessionKey := fmt.Sprintf("session:%s", userID.String())
	return s.cache.Delete(ctx, sessionKey)
}

// LogoutAll logs out user from all devices
func (s *authService) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	// Revoke all refresh tokens
	if err := s.userRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		return err
	}

	// Clear all sessions from cache
	sessionPattern := fmt.Sprintf("session:%s:*", userID.String())
	return s.cache.DeletePattern(ctx, sessionPattern)
}

// GenerateTokenPair generates access and refresh tokens
func (s *authService) GenerateTokenPair(ctx context.Context, user *models.User) (*TokenPair, error) {
	// Create session
	sessionID := uuid.New().String()
	sessionKey := fmt.Sprintf("session:%s:%s", user.ID.String(), sessionID)

	// Generate access token
	now := time.Now()
	expiresAt := now.Add(s.jwtCfg.Expiry)

	claims := &Claims{
		UserID:    user.ID,
		Email:     user.Email,
		Username:  user.Username,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.jwtCfg.Issuer,
			Subject:   user.ID.String(),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(s.jwtCfg.Secret))
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     utils.GenerateRandomString(64),
		ExpiresAt: now.Add(s.jwtCfg.RefreshTokenExpiry),
	}

	if err := s.userRepo.SaveRefreshToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	// Cache session
	session := &Session{
		ID:        sessionID,
		UserID:    user.ID,
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}

	if err := s.cache.Set(ctx, sessionKey, session, s.jwtCfg.Expiry); err != nil {
		// Log error but don't fail token generation
		fmt.Printf("Failed to cache session: %v\n", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
		ExpiresIn:    int(s.jwtCfg.Expiry.Seconds()),
	}, nil
}

// ValidateToken validates JWT token
func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtCfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, repository.ErrTokenInvalid
	}

	// Check if session exists in cache
	sessionKey := fmt.Sprintf("session:%s:%s", claims.UserID.String(), claims.SessionID)
	exists, err := s.cache.Exists(ctx, sessionKey)
	if err != nil || !exists {
		return nil, repository.ErrTokenInvalid
	}

	return claims, nil
}

// Helper methods

// safeSubstring returns a substring of s from start to start+length, or the entire string if out of bounds
func safeSubstring(s string, start, length int) string {
	runes := []rune(s)
	if start >= len(runes) {
		return ""
	}
	end := start + length
	if end > len(runes) {
		end = len(runes)
	}
	return string(runes[start:end])
}

func (s *authService) mapUserToDTO(user *models.User) *UserDTO {
	return &UserDTO{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Bio:         user.Bio,
		Timezone:    user.Timezone,
		Locale:      user.Locale,
		StreakCount: user.StreakCount,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
	}
}

func (s *authService) sendVerificationEmail(email, userID string) {
	// TODO: Implement email sending
	fmt.Printf("Sending verification email to %s for user %s\n", email, userID)
}

// Stub implementations for other methods
func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	// TODO: Implement
	return nil
}

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	// TODO: Implement
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// TODO: Implement
	return nil
}

func (s *authService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	// TODO: Implement
	return nil
}

func (s *authService) GetGoogleAuthURL(state string) string {
	if s.cfg == nil {
		return ""
	}
	conf := &oauth2.Config{
		ClientID:     s.cfg.OAuth.Google.ClientID,
		ClientSecret: s.cfg.OAuth.Google.ClientSecret,
		RedirectURL:  s.cfg.OAuth.Google.RedirectURL,
		Scopes: []string{
			googleoauth.UserinfoEmailScope,
			googleoauth.UserinfoProfileScope,
			"openid",
		},
		Endpoint: google.Endpoint,
	}
	return conf.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func (s *authService) HandleGoogleAuth(ctx context.Context, code string) (*AuthResponse, error) {
	if s.cfg == nil {
		return nil, errors.New("oauth config not initialized")
	}

	if s.cfg.OAuth.Google.ClientID == "" || s.cfg.OAuth.Google.ClientSecret == "" || s.cfg.OAuth.Google.RedirectURL == "" {
		return nil, errors.New("google oauth is not configured")
	}

	// Build OAuth2 config
	conf := &oauth2.Config{
		ClientID:     s.cfg.OAuth.Google.ClientID,
		ClientSecret: s.cfg.OAuth.Google.ClientSecret,
		RedirectURL:  s.cfg.OAuth.Google.RedirectURL,
		Scopes: []string{
			googleoauth.UserinfoEmailScope,
			googleoauth.UserinfoProfileScope,
			"openid",
		},
		Endpoint: google.Endpoint,
	}

	// Exchange code for token
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Fetch user info
	httpClient := conf.Client(ctx, tok)
	oauth2Service, err := googleoauth.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("failed to create google oauth2 service: %w", err)
	}
	userInfo, err := oauth2Service.Userinfo.V2.Me.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// userInfo contains Id, Email, VerifiedEmail, Name, GivenName, FamilyName, Picture
	provider := "google"
	providerID := userInfo.Id

	// Try to find by provider
	var user *models.User
	user, err = s.userRepo.FindByAuthProvider(ctx, provider, providerID)
	if err != nil {
		if !repository.IsNotFound(err) {
			return nil, err
		}
	}

	if user == nil {
		// Try find by email
		if userInfo.Email != "" {
			u, err := s.userRepo.FindByEmail(ctx, userInfo.Email)
			if err == nil {
				user = u
			} else if !repository.IsNotFound(err) {
				return nil, err
			}
		}
	}

	if user == nil {
		// Create new user
		username := userInfo.Email
		if idx := strings.Index(username, "@"); idx > 0 {
			username = username[:idx]
		}
		displayName := userInfo.Name
		if displayName == "" {
			displayName = username
		}
		user = &models.User{
			Username:      username,
			Email:         userInfo.Email,
			EmailVerified: userInfo.VerifiedEmail != nil && *userInfo.VerifiedEmail,
			DisplayName:   displayName,
			AvatarURL:     userInfo.Picture,
			Timezone:      "UTC",
			Locale:        "en",
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	// Link provider (idempotent create)
	var expiresAt *time.Time
	if tok.Expiry.After(time.Now()) {
		t := tok.Expiry
		expiresAt = &t
	}
	ap := &models.AuthProvider{
		UserID:       user.ID,
		Provider:     provider,
		ProviderID:   providerID,
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		ExpiresAt:    expiresAt,
		RawData:      models.JSONB{"provider": provider},
	}
	// Best-effort: try to create, ignore if exists
	if err := s.userRepo.LinkAuthProvider(ctx, ap); err != nil {
		// continue on unique violation or similar; repository layer not exposing, so just log
		log.Printf("LinkAuthProvider error (continuing): %v", err)
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	// Generate tokens
	tokenPair, err := s.GenerateTokenPair(ctx, user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         s.mapUserToDTO(user),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		TokenType:    "Bearer",
	}, nil
}

func (s *authService) HandleLinkedInAuth(ctx context.Context, code string) (*AuthResponse, error) {
	// LinkedIn OAuth implementation placeholder
	// In production, this would:
	// 1. Exchange code for access token
	// 2. Fetch user info from LinkedIn API
	// 3. Create or update user account
	// 4. Generate JWT tokens
	return nil, errors.New("LinkedIn OAuth not configured - please add LinkedIn OAuth credentials to enable this feature")
}

func (s *authService) HandleAppleAuth(ctx context.Context, code string) (*AuthResponse, error) {
	// Apple Sign In implementation placeholder
	// In production, this would:
	// 1. Verify Apple ID token
	// 2. Extract user info
	// 3. Create or update user account
	// 4. Generate JWT tokens
	return nil, errors.New("Apple Sign In not configured - please add Apple OAuth credentials to enable this feature")
}

func (s *authService) LinkOAuthAccount(ctx context.Context, userID uuid.UUID, provider string, code string) error {
	// Link an OAuth account to an existing user
	// This allows users to add additional login methods
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// In production, exchange code for provider tokens and link account
	// For now, return placeholder
	_ = user
	return errors.New("OAuth account linking not fully implemented - provider: " + provider)
}

func (s *authService) UnlinkOAuthAccount(ctx context.Context, userID uuid.UUID, provider string) error {
	return s.userRepo.UnlinkAuthProvider(ctx, userID, provider)
}

func (s *authService) RevokeToken(ctx context.Context, token string) error {
	return s.userRepo.RevokeRefreshToken(ctx, token)
}

func (s *authService) CreateSession(ctx context.Context, userID uuid.UUID, metadata SessionMetadata) (*Session, error) {
	sessionID := uuid.New().String()
	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Metadata:  metadata,
	}

	// Store session in cache
	sessionKey := fmt.Sprintf("session:%s:%s", userID.String(), sessionID)
	if err := s.cache.Set(ctx, sessionKey, session, 24*time.Hour); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *authService) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	// Parse session ID to extract user ID (format: userID:sessionID)
	// For now, we'll search for the session in cache
	// In production, you'd want a more efficient lookup mechanism

	// Placeholder implementation
	return nil, errors.New("session not found - get session by ID needs session index")
}

func (s *authService) EndSession(ctx context.Context, sessionID string) error {
	// Find and delete session from cache
	// Placeholder implementation - needs session index for efficient lookup
	return errors.New("end session requires session indexing - not fully implemented")
}

// (duplicate removed)
