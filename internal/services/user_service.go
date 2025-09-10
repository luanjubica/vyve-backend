package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vyve/vyve-backend/internal/models"
	"github.com/vyve/vyve-backend/internal/repository"
	"github.com/vyve/vyve-backend/pkg/storage"
)

// OnboardingStatus represents the onboarding status of a user
type OnboardingStatus struct {
	Completed      bool      `json:"completed"`
	CompletedSteps []string  `json:"completed_steps"`
	CurrentStep    string    `json:"current_step,omitempty"`
	LastUpdated    time.Time `json:"last_updated"`
}

// UserService handles user business logic
type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UploadAvatar(ctx context.Context, userID uuid.UUID, file []byte, filename string) (string, error)
	GetStats(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	GetSettings(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	UpdateSettings(ctx context.Context, userID uuid.UUID, settings map[string]interface{}) error
	RegisterPushToken(ctx context.Context, userID uuid.UUID, token string, platform string) error
	DeactivatePushToken(ctx context.Context, token string) error
	GetPushTokens(ctx context.Context, userID uuid.UUID) ([]*models.PushToken, error)

	// Onboarding related methods
	GetOnboardingStatus(ctx context.Context, userID uuid.UUID) (*OnboardingStatus, error)
	UpdateOnboardingStatus(ctx context.Context, userID uuid.UUID, completed bool, currentStep string) (*OnboardingStatus, error)
}

type userService struct {
	userRepo repository.UserRepository
	storage  storage.Storage
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, storage storage.Storage) UserService {
	return &userService{
		userRepo: userRepo,
		storage:  storage,
	}
}

// GetByID gets a user by ID
func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// Update updates a user
func (s *userService) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	// TODO: Add validation and proper field mapping
	if displayName, ok := updates["display_name"].(string); ok {
		user.DisplayName = displayName
	}
	if bio, ok := updates["bio"].(string); ok {
		user.Bio = bio
	}
	if timezone, ok := updates["timezone"].(string); ok {
		user.Timezone = timezone
	}
	if locale, ok := updates["locale"].(string); ok {
		user.Locale = locale
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Delete deletes a user
func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}

// UploadAvatar uploads user avatar
func (s *userService) UploadAvatar(ctx context.Context, userID uuid.UUID, file []byte, filename string) (string, error) {
	// Generate unique filename
	key := "avatars/" + userID.String() + "/" + uuid.New().String()

	// Upload to storage
	url, err := s.storage.Upload(ctx, key, file, "image/jpeg")
	if err != nil {
		return "", err
	}

	// Update user avatar URL
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", err
	}

	user.AvatarURL = url
	if err := s.userRepo.Update(ctx, user); err != nil {
		return "", err
	}

	return url, nil
}

// GetStats gets user statistics
func (s *userService) GetStats(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	return s.userRepo.GetUserStats(ctx, userID)
}

// UpdateSettings updates user settings
func (s *userService) UpdateSettings(ctx context.Context, userID uuid.UUID, settings map[string]interface{}) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Normalize payload: unwrap nested "data" wrappers and drop flags like "success"
	// Unwrap recursively in case the client sent nested { data: { data: {...} } }
	for {
		if raw, ok := settings["data"]; ok {
			if inner, ok := raw.(map[string]interface{}); ok {
				settings = inner
				continue
			}
		}
		break
	}
	delete(settings, "success")

	// Merge settings
	if user.Settings == nil {
		user.Settings = make(models.JSONB)
	}

	for key, value := range settings {
		user.Settings[key] = value
	}

	return s.userRepo.Update(ctx, user)
}

// GetSettings gets user settings
func (s *userService) GetSettings(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Return settings from user model
	if user.Settings == nil {
		return make(map[string]interface{}), nil
	}

	return user.Settings, nil
}

// RegisterPushToken registers a push notification token
func (s *userService) RegisterPushToken(ctx context.Context, userID uuid.UUID, token string, platform string) error {
	pushToken := &models.PushToken{
		UserID:   userID,
		Token:    token,
		Platform: platform,
		Active:   true,
	}

	return s.userRepo.SavePushToken(ctx, pushToken)
}

// DeactivatePushToken deactivates a push notification token
func (s *userService) DeactivatePushToken(ctx context.Context, token string) error {
	return s.userRepo.DeactivatePushToken(ctx, token)
}

// GetPushTokens gets user's push tokens
func (s *userService) GetPushTokens(ctx context.Context, userID uuid.UUID) ([]*models.PushToken, error) {
	return s.userRepo.GetUserPushTokens(ctx, userID) // Fixed method name
}

// GetOnboardingStatus gets the user's onboarding status
func (s *userService) GetOnboardingStatus(ctx context.Context, userID uuid.UUID) (*OnboardingStatus, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert OnboardingSteps to []string
	completedSteps := []string(user.OnboardingSteps)

	// Default onboarding status
	status := &OnboardingStatus{
		Completed:      user.OnboardingCompleted,
		CompletedSteps: completedSteps,
		LastUpdated:    time.Now(),
	}

	// If onboarding is marked as completed but no steps are recorded, add default steps
	if status.Completed && len(status.CompletedSteps) == 0 {
		status.CompletedSteps = []string{"welcome", "people_added", "first_interaction"}
	}

	return status, nil
}

// UpdateOnboardingStatus updates the user's onboarding status
func (s *userService) UpdateOnboardingStatus(ctx context.Context, userID uuid.UUID, completed bool, currentStep string) (*OnboardingStatus, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update completion status
	user.OnboardingCompleted = completed

	// If we're completing a step, add it to the completed steps
	if currentStep != "" {
		completedSteps := []string(user.OnboardingSteps)

		// Check if step is already completed
		stepExists := false
		for _, step := range completedSteps {
			if step == currentStep {
				stepExists = true
				break
			}
		}

		// Add the step if it doesn't exist
		if !stepExists {
			completedSteps = append(completedSteps, currentStep)
			user.OnboardingSteps = models.OnboardingSteps(completedSteps)
		}
	}

	// Update the user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Return the updated status
	return s.GetOnboardingStatus(ctx, userID)
}
