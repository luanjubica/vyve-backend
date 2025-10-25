// File: internal/repository/person_repository.go
package repository

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/vyve/vyve-backend/internal/models"
)

// PersonRepository defines person data access interface
type PersonRepository interface {
	Create(ctx context.Context, person *models.Person) error
	Update(ctx context.Context, person *models.Person) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Person, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Person, error)
	List(ctx context.Context, opts FilterOptions) ([]*models.Person, *PaginationResult, error)
	Search(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*models.Person, error)
	GetCategories(ctx context.Context, userID uuid.UUID) ([]string, error)
	GetRecentInteractions(ctx context.Context, personID uuid.UUID, limit int) ([]*models.Interaction, error)
	UpdateHealthScore(ctx context.Context, personID uuid.UUID, score float64) error
	GetByCategory(ctx context.Context, userID uuid.UUID, category string) ([]*models.Person, error)
	GetPeopleNeedingAttention(ctx context.Context, userID uuid.UUID) ([]*models.Person, error)
	GetPeopleForReminders(ctx context.Context, userID uuid.UUID) ([]*models.Person, error)
	IncrementInteractionCount(ctx context.Context, personID uuid.UUID) error
	UpdateLastInteraction(ctx context.Context, personID uuid.UUID) error
}

type personRepository struct {
	BaseRepository
}

// NewPersonRepository creates a new person repository
func NewPersonRepository(db *gorm.DB) PersonRepository {
	return &personRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create creates a new person
func (r *personRepository) Create(ctx context.Context, person *models.Person) error {
	return r.db.WithContext(ctx).Create(person).Error
}

// Update updates a person
func (r *personRepository) Update(ctx context.Context, person *models.Person) error {
	return r.db.WithContext(ctx).Save(person).Error
}

// Delete soft deletes a person
func (r *personRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Person{}, "id = ?", id).Error
}

// FindByID finds a person by ID
func (r *personRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Person, error) {
	var person models.Person
	err := r.db.WithContext(ctx).First(&person, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPersonNotFound
		}
		return nil, err
	}
	return &person, nil
}

func (r *personRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Person, error) {
	log.Printf("🎯 PersonRepository.FindByUserID called with userID: %s", userID)

	log.Printf("🔍 Step 1: Declaring people slice...")
	var people []*models.Person
	log.Printf("✅ Step 1: People slice declared")

	log.Printf("🔍 Step 2: Getting database connection with context...")
	//dbWithCtx := r.db.WithContext(ctx)
	log.Printf("✅ Step 2: Database connection with context obtained")

	log.Printf("🔍 Step 3: Building WHERE clause...")
	//queryWithWhere := dbWithCtx.Where("user_id = ?", userID)
	log.Printf("✅ Step 3: WHERE clause built")

	log.Printf("🔍 Step 4: Adding ORDER BY clause...")
	//queryWithOrder := queryWithWhere.Order("name ASC")
	log.Printf("✅ Step 4: ORDER BY clause added")

	log.Printf("🔍 Step 5: About to execute Find()...")
	log.Printf("🚨 THIS IS THE CRITICAL MOMENT - if logs stop here, Find() is hanging")

	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			*
		FROM people 
		WHERE user_id = $1 AND (deleted_at IS NULL OR deleted_at = '0001-01-01 00:00:00+00')
		ORDER BY name ASC
	`, userID).Scan(&people).Error

	log.Printf("✅ Step 5: Find() completed! Found %d people", len(people))

	if err != nil {
		log.Printf("❌ Find() returned error: %v", err)
		return nil, err
	}

	log.Printf("✅ FindByUserID completed successfully, returning %d people", len(people))
	return people, nil
}

// List lists people with pagination
func (r *personRepository) List(ctx context.Context, opts FilterOptions) ([]*models.Person, *PaginationResult, error) {
	query := r.db.WithContext(ctx).Model(&models.Person{})

	// Apply filters
	if opts.UserID != uuid.Nil {
		query = query.Where("user_id = ?", opts.UserID)
	}

	if opts.Category != "" {
		query = query.Where("category_id = ?", opts.Category)
	}

	if opts.Search != "" {
		query = query.Where("name ILIKE ?", "%"+opts.Search+"%")
	}

	// Apply ordering
	if opts.OrderBy != "" {
		// Handle cases where orderBy might be in format "field:direction"
		orderBy := opts.OrderBy
		direction := "ASC"
		if opts.Desc {
			direction = "DESC"
		}
		
		// If orderBy contains a colon, split it into field and direction
		if parts := strings.Split(orderBy, ":"); len(parts) == 2 {
			orderBy = parts[0]
			direction = strings.ToUpper(parts[1])
			if direction != "ASC" && direction != "DESC" {
				direction = "ASC" // Default to ASC if invalid direction
			}
		}
		
		// Ensure the field is properly quoted to prevent SQL injection
		query = query.Order(orderBy + " " + direction)
	} else {
		query = query.Order("name ASC")
	}

	var people []*models.Person
	result, err := Paginate(ctx, query, opts.Page, opts.Limit, &people)
	if err != nil {
		return nil, nil, err
	}

	return people, result, nil
}

// Search searches people by name
func (r *personRepository) Search(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*models.Person, error) {
	var people []*models.Person
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND name ILIKE ?", userID, "%"+query+"%").
		Limit(limit).
		Find(&people).Error
	return people, err
}

// GetCategories gets unique categories for a user
func (r *personRepository) GetCategories(ctx context.Context, userID uuid.UUID) ([]string, error) {
	type row struct{ Name string }
	var names []string
	err := r.db.WithContext(ctx).
		Table("categories").
		Where("user_id = ?", userID).
		Order("name ASC").
		Pluck("name", &names).Error
	return names, err
}

// GetRecentInteractions gets recent interactions for a person
func (r *personRepository) GetRecentInteractions(ctx context.Context, personID uuid.UUID, limit int) ([]*models.Interaction, error) {
	var interactions []*models.Interaction
	err := r.db.WithContext(ctx).
		Where("person_id = ?", personID).
		Order("interaction_at DESC").
		Limit(limit).
		Find(&interactions).Error
	return interactions, err
}

// UpdateHealthScore updates a person's health score
func (r *personRepository) UpdateHealthScore(ctx context.Context, personID uuid.UUID, score float64) error {
	return r.db.WithContext(ctx).
		Model(&models.Person{}).
		Where("id = ?", personID).
		Update("health_score", score).Error
}

// GetByCategory gets people by category
func (r *personRepository) GetByCategory(ctx context.Context, userID uuid.UUID, category string) ([]*models.Person, error) {
	var catID uuid.UUID
	if err := r.db.WithContext(ctx).
		Table("categories").
		Select("id").
		Where("user_id = ? AND name = ?", userID, category).
		Limit(1).
		Scan(&catID).Error; err != nil {
		return nil, err
	}
	var people []*models.Person
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND category_id = ?", userID, catID).
		Find(&people).Error
	return people, err
}

// GetPeopleNeedingAttention gets people with low health scores or draining energy patterns
func (r *personRepository) GetPeopleNeedingAttention(ctx context.Context, userID uuid.UUID) ([]*models.Person, error) {
	var people []*models.Person
	sub := r.db.Table("energy_patterns").Select("id").Where("name = ?", "draining")
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND (health_score < ? OR energy_pattern_id IN (?))", userID, 60.0, sub).
		Order("health_score ASC").
		Find(&people).Error
	return people, err
}

// GetPeopleForReminders gets people due for reminders
func (r *personRepository) GetPeopleForReminders(ctx context.Context, userID uuid.UUID) ([]*models.Person, error) {
	var people []*models.Person
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND next_reminder_at <= NOW()", userID).
		Find(&people).Error
	return people, err
}

// IncrementInteractionCount increments interaction count for a person
func (r *personRepository) IncrementInteractionCount(ctx context.Context, personID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Person{}).
		Where("id = ?", personID).
		UpdateColumn("interaction_count", gorm.Expr("interaction_count + ?", 1)).
		Error
}

// UpdateLastInteraction updates last interaction timestamp
func (r *personRepository) UpdateLastInteraction(ctx context.Context, personID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Person{}).
		Where("id = ?", personID).
		UpdateColumn("last_interaction_at", gorm.Expr("NOW()")).
		Error
}
