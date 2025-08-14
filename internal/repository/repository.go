package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repositories holds all repository instances
type Repositories struct {
	User        UserRepository
	Person      PersonRepository
	Interaction InteractionRepository
	Reflection  ReflectionRepository
	Nudge       NudgeRepository
	Event       EventRepository
	Consent     ConsentRepository
	AuditLog    AuditLogRepository
	DataExport  DataExportRepository
}

// NewRepositories creates new repository instances
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:        NewUserRepository(db),
		Person:      NewPersonRepository(db),
		Interaction: NewInteractionRepository(db),
		Reflection:  NewReflectionRepository(db),
		Nudge:       NewNudgeRepository(db),
		Event:       NewEventRepository(db),
		Consent:     NewConsentRepository(db),
		AuditLog:    NewAuditLogRepository(db),
		DataExport:  NewDataExportRepository(db),
	}
}

// BaseRepository provides common database operations
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *gorm.DB) BaseRepository {
	return BaseRepository{db: db}
}

// Create creates a new record
func (r *BaseRepository) Create(ctx context.Context, entity interface{}) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// Update updates a record
func (r *BaseRepository) Update(ctx context.Context, entity interface{}) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete soft deletes a record
func (r *BaseRepository) Delete(ctx context.Context, entity interface{}) error {
	return r.db.WithContext(ctx).Delete(entity).Error
}

// HardDelete permanently deletes a record
func (r *BaseRepository) HardDelete(ctx context.Context, entity interface{}) error {
	return r.db.WithContext(ctx).Unscoped().Delete(entity).Error
}

// FindByID finds a record by ID
func (r *BaseRepository) FindByID(ctx context.Context, id uuid.UUID, entity interface{}) error {
	return r.db.WithContext(ctx).First(entity, "id = ?", id).Error
}

// Transaction executes a function within a database transaction
func (r *BaseRepository) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

// Paginate adds pagination to query
func (r *BaseRepository) Paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if limit <= 0 {
			limit = 10
		}
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

// OrderBy adds ordering to query
func (r *BaseRepository) OrderBy(field string, desc bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if desc {
			return db.Order(field + " DESC")
		}
		return db.Order(field + " ASC")
	}
}

// WithPreload adds preloading to query
func (r *BaseRepository) WithPreload(associations ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, association := range associations {
			db = db.Preload(association)
		}
		return db
	}
}

// PaginationResult holds pagination metadata
type PaginationResult struct {
	Total       int64       `json:"total"`
	Page        int         `json:"page"`
	Limit       int         `json:"limit"`
	TotalPages  int         `json:"total_pages"`
	HasNext     bool        `json:"has_next"`
	HasPrevious bool        `json:"has_previous"`
	Items       interface{} `json:"items"`
}

// Paginate executes a paginated query
func Paginate(ctx context.Context, db *gorm.DB, page, limit int, result interface{}) (*PaginationResult, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	var total int64
	db.Count(&total)

	offset := (page - 1) * limit
	if err := db.WithContext(ctx).Offset(offset).Limit(limit).Find(result).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &PaginationResult{
		Total:       total,
		Page:        page,
		Limit:       limit,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
		Items:       result,
	}, nil
}

// FilterOptions represents common filter options
type FilterOptions struct {
	UserID    uuid.UUID
	StartDate *time.Time
	EndDate   *time.Time
	Search    string
	Status    string
	Type      string
	Category  string
	OrderBy   string
	Desc      bool
	Page      int
	Limit     int
}

// ApplyFilters applies common filters to a query
func ApplyFilters(db *gorm.DB, opts FilterOptions) *gorm.DB {
	if opts.UserID != uuid.Nil {
		db = db.Where("user_id = ?", opts.UserID)
	}
	
	if opts.StartDate != nil {
		db = db.Where("created_at >= ?", opts.StartDate)
	}
	
	if opts.EndDate != nil {
		db = db.Where("created_at <= ?", opts.EndDate)
	}
	
	if opts.Search != "" {
		db = db.Where("name ILIKE ? OR email ILIKE ?", "%"+opts.Search+"%", "%"+opts.Search+"%")
	}
	
	if opts.Status != "" {
		db = db.Where("status = ?", opts.Status)
	}
	
	if opts.Type != "" {
		db = db.Where("type = ?", opts.Type)
	}
	
	if opts.Category != "" {
		db = db.Where("category = ?", opts.Category)
	}
	
	if opts.OrderBy != "" {
		if opts.Desc {
			db = db.Order(opts.OrderBy + " DESC")
		} else {
			db = db.Order(opts.OrderBy + " ASC")
		}
	}
	
	return db
}

// BatchCreate creates multiple records in a batch
func BatchCreate(ctx context.Context, db *gorm.DB, entities interface{}) error {
	return db.WithContext(ctx).CreateInBatches(entities, 100).Error
}

// BulkUpdate performs bulk update
func BulkUpdate(ctx context.Context, db *gorm.DB, model interface{}, updates map[string]interface{}, conditions ...interface{}) error {
	return db.WithContext(ctx).Model(model).Where(conditions[0], conditions[1:]...).Updates(updates).Error
}

// SoftDelete performs soft delete
func SoftDelete(ctx context.Context, db *gorm.DB, model interface{}, id uuid.UUID) error {
	return db.WithContext(ctx).Model(model).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

// Restore restores soft deleted record
func Restore(ctx context.Context, db *gorm.DB, model interface{}, id uuid.UUID) error {
	return db.WithContext(ctx).Model(model).Unscoped().Where("id = ?", id).Update("deleted_at", nil).Error
}

// Upsert performs insert or update
func Upsert(ctx context.Context, db *gorm.DB, entity interface{}, conflictColumns []string, updateColumns []string) error {
	return db.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: conflictColumns[0]}},
			DoUpdates: clause.AssignmentColumns(updateColumns),
		},
	).Create(entity).Error
}

// Count returns count of records
func Count(ctx context.Context, db *gorm.DB, model interface{}, conditions ...interface{}) (int64, error) {
	var count int64
	err := db.WithContext(ctx).Model(model).Where(conditions[0], conditions[1:]...).Count(&count).Error
	return count, err
}

// Exists checks if record exists
func Exists(ctx context.Context, db *gorm.DB, model interface{}, conditions ...interface{}) (bool, error) {
	var count int64
	err := db.WithContext(ctx).Model(model).Where(conditions[0], conditions[1:]...).Count(&count).Error
	return count > 0, err
}