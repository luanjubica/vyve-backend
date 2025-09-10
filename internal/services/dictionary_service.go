package services

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/vyve/vyve-backend/internal/models"
)

// DictionaryService provides read-only access to dictionary tables
// Categories are user-scoped; other dictionaries are global.
type DictionaryService interface {
	ListCategories(ctx context.Context, userID uuid.UUID) ([]*models.Category, error)
	ListCommunicationMethods(ctx context.Context) ([]*models.CommunicationMethod, error)
	ListRelationshipStatuses(ctx context.Context) ([]*models.RelationshipStatus, error)
	ListIntentions(ctx context.Context) ([]*models.Intention, error)
	ListEnergyPatterns(ctx context.Context) ([]*models.EnergyPattern, error)
}

type dictionaryService struct {
	db *gorm.DB
}

// NewDictionaryService creates a new DictionaryService
func NewDictionaryService(db *gorm.DB) DictionaryService {
	return &dictionaryService{db: db}
}

func (s *dictionaryService) ListCategories(ctx context.Context, userID uuid.UUID) ([]*models.Category, error) {
	var items []*models.Category
	q := s.db.WithContext(ctx).Model(&models.Category{}).
		Where("user_id = ?", userID).
		Order("name ASC")
	if err := q.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *dictionaryService) ListCommunicationMethods(ctx context.Context) ([]*models.CommunicationMethod, error) {
	var items []*models.CommunicationMethod
	if err := s.db.WithContext(ctx).Model(&models.CommunicationMethod{}).Order("name ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *dictionaryService) ListRelationshipStatuses(ctx context.Context) ([]*models.RelationshipStatus, error) {
	var items []*models.RelationshipStatus
	if err := s.db.WithContext(ctx).Model(&models.RelationshipStatus{}).Order("name ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *dictionaryService) ListIntentions(ctx context.Context) ([]*models.Intention, error) {
	var items []*models.Intention
	if err := s.db.WithContext(ctx).Model(&models.Intention{}).Order("name ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *dictionaryService) ListEnergyPatterns(ctx context.Context) ([]*models.EnergyPattern, error) {
	var items []*models.EnergyPattern
	if err := s.db.WithContext(ctx).Model(&models.EnergyPattern{}).Order("name ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
