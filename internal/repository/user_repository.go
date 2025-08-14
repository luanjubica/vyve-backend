package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/vyve/vyve-backend/internal/models"
)

// UserRepository defines user data access interface
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindByAuthProvider(ctx context.Context, provider, providerID string) (*models.User, error)
	List(ctx context.Context, opts FilterOptions) ([]*models.User, *PaginationResult, error)
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	UpdateLastActivity(ctx context.Context, id uuid.UUID) error
	UpdateStreak(ctx context.Context, id uuid.UUID, count int) error
	GetActiveUsers(ctx context.Context, since time.Time) ([]*models.User, error)
	GetUsersForReminders(ctx context.Context, hour int) ([]*models.User, error)
	SaveRefreshToken(ctx context.Context, token *models.RefreshToken) error
	FindRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error
	SavePushToken(ctx context.Context, token *models.PushToken) error
	GetUserPushTokens(ctx context.Context, userID uuid.UUID) ([]*models.PushToken, error)
	DeactivatePushToken(ctx context.Context, token string) error
	LinkAuthProvider(ctx context.Context, provider *models.AuthProvider) error
	UnlinkAuthProvider(ctx context.Context, userID uuid.UUID, provider string) error
	GetAuthProviders(ctx context.Context, userID uuid.UUID) ([]*models.AuthProvider, error)
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	SearchUsers(ctx context.Context, query string, limit int) ([]*models.User, error)
	GetUserStats(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
}

type userRepository struct {
	BaseRepository
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete soft deletes a user
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}

// FindByID finds a user by ID
func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("AuthProviders").
		First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("AuthProviders").
		First(&user, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername finds a user by username
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("AuthProviders").
		First(&user, "username = ?", username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByAuthProvider finds a user by auth provider
func (r *userRepository) FindByAuthProvider(ctx context.Context, provider, providerID string) (*models.User, error) {
	var authProvider models.AuthProvider
	err := r.db.WithContext(ctx).
		Where("provider = ? AND provider_id = ?", provider, providerID).
		First(&authProvider).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return r.FindByID(ctx, authProvider.UserID)
}

// List lists users with pagination
func (r *userRepository) List(ctx context.Context, opts FilterOptions) ([]*models.User, *PaginationResult, error) {
	query := r.db.WithContext(ctx).Model(&models.User{})
	query = ApplyFilters(query, opts)

	var users []*models.User
	result, err := Paginate(ctx, query, opts.Page, opts.Limit, &users)
	if err != nil {
		return nil, nil, err
	}

	return users, result, nil
}

// UpdateLastLogin updates user's last login time
func (r *userRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_login_at":    now,
			"last_activity_at": now,
		}).Error
}

// UpdateLastActivity updates user's last activity time
func (r *userRepository) UpdateLastActivity(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("last_activity_at", time.Now()).Error
}

// UpdateStreak updates user's streak count
func (r *userRepository) UpdateStreak(ctx context.Context, id uuid.UUID, count int) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"streak_count":        count,
			"last_reflection_at": time.Now(),
		}).Error
}

// GetActiveUsers gets users active since a given time
func (r *userRepository) GetActiveUsers(ctx context.Context, since time.Time) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).
		Where("last_activity_at >= ?", since).
		Find(&users).Error
	return users, err
}

// GetUsersForReminders gets users for reminder notifications
func (r *userRepository) GetUsersForReminders(ctx context.Context, hour int) ([]*models.User, error) {
	var users []*models.User
	
	// Get users based on their timezone where it's the specified hour
	err := r.db.WithContext(ctx).
		Joins("JOIN push_tokens ON push_tokens.user_id = users.id").
		Where("push_tokens.active = ?", true).
		Where("DATE_PART('hour', NOW() AT TIME ZONE users.timezone) = ?", hour).
		Where("users.last_reflection_at < CURRENT_DATE OR users.last_reflection_at IS NULL").
		Group("users.id").
		Find(&users).Error
		
	return users, err
}

// SaveRefreshToken saves a refresh token
func (r *userRepository) SaveRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// FindRefreshToken finds a refresh token
func (r *userRepository) FindRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token = ? AND revoked = ? AND expires_at > ?", token, false, time.Now()).
		First(&refreshToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}
	return &refreshToken, nil
}

// RevokeRefreshToken revokes a refresh token
func (r *userRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("token = ?", token).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": time.Now(),
		}).Error
}

// RevokeAllUserTokens revokes all user's refresh tokens
func (r *userRepository) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": time.Now(),
		}).Error
}

// SavePushToken saves a push notification token
func (r *userRepository) SavePushToken(ctx context.Context, token *models.PushToken) error {
	// Deactivate existing tokens for the same device
	if token.DeviceID != "" {
		r.db.WithContext(ctx).
			Model(&models.PushToken{}).
			Where("user_id = ? AND device_id = ?", token.UserID, token.DeviceID).
			Update("active", false)
	}
	
	// Upsert the new token
	return r.db.WithContext(ctx).
		Where("token = ?", token.Token).
		Assign(token).
		FirstOrCreate(token).Error
}

// GetUserPushTokens gets user's active push tokens
func (r *userRepository) GetUserPushTokens(ctx context.Context, userID uuid.UUID) ([]*models.PushToken, error) {
	var tokens []*models.PushToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND active = ?", userID, true).
		Find(&tokens).Error
	return tokens, err
}

// DeactivatePushToken deactivates a push token
func (r *userRepository) DeactivatePushToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).
		Model(&models.PushToken{}).
		Where("token = ?", token).
		Update("active", false).Error
}

// LinkAuthProvider links an auth provider to a user
func (r *userRepository) LinkAuthProvider(ctx context.Context, provider *models.AuthProvider) error {
	return r.db.WithContext(ctx).Create(provider).Error
}

// UnlinkAuthProvider unlinks an auth provider from a user
func (r *userRepository) UnlinkAuthProvider(ctx context.Context, userID uuid.UUID, provider string) error {
	// Check if user has other auth methods
	var count int64
	r.db.WithContext(ctx).
		Model(&models.AuthProvider{}).
		Where("user_id = ?", userID).
		Count(&count)
	
	// Also check if user has password
	var user models.User
	r.db.WithContext(ctx).First(&user, userID)
	
	if count <= 1 && user.PasswordHash == "" {
		return ErrLastAuthMethod
	}
	
	return r.db.WithContext(ctx).
		Where("user_id = ? AND provider = ?", userID, provider).
		Delete(&models.AuthProvider{}).Error
}

// GetAuthProviders gets user's auth providers
func (r *userRepository) GetAuthProviders(ctx context.Context, userID uuid.UUID) ([]*models.AuthProvider, error) {
	var providers []*models.AuthProvider
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&providers).Error
	return providers, err
}

// CheckUsernameExists checks if username exists
func (r *userRepository) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("username = ?", username).
		Count(&count).Error
	return count > 0, err
}

// CheckEmailExists checks if email exists
func (r *userRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).Error
	return count > 0, err
}

// SearchUsers searches users by query
func (r *userRepository) SearchUsers(ctx context.Context, query string, limit int) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).
		Where("username ILIKE ? OR email ILIKE ? OR display_name ILIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%").
		Limit(limit).
		Find(&users).Error
	return users, err
}

// GetUserStats gets user statistics
func (r *userRepository) GetUserStats(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Get various counts
	var peopleCount int64
	r.db.WithContext(ctx).Model(&models.Person{}).Where("user_id = ?", userID).Count(&peopleCount)
	stats["people_count"] = peopleCount
	
	var interactionsCount int64
	r.db.WithContext(ctx).Model(&models.Interaction{}).Where("user_id = ?", userID).Count(&interactionsCount)
	stats["interactions_count"] = interactionsCount
	
	var reflectionsCount int64
	r.db.WithContext(ctx).Model(&models.Reflection{}).Where("user_id = ?", userID).Count(&reflectionsCount)
	stats["reflections_count"] = reflectionsCount
	
	// Get average health score
	var avgHealthScore float64
	r.db.WithContext(ctx).
		Model(&models.Person{}).
		Where("user_id = ?", userID).
		Pluck("AVG(health_score)", &avgHealthScore)
	stats["avg_health_score"] = avgHealthScore
	
	// Get energy distribution
	var energyStats []struct {
		EnergyImpact string
		Count        int
	}
	r.db.WithContext(ctx).
		Model(&models.Interaction{}).
		Select("energy_impact, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("energy_impact").
		Scan(&energyStats)
	stats["energy_distribution"] = energyStats
	
	return stats, nil
}