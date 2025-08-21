package services_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vyve/vyve-backend/internal/models"
	"github.com/vyve/vyve-backend/internal/repository"
	svc "github.com/vyve/vyve-backend/internal/services"
)

func TestPersonService_CountPeople(t *testing.T) {
	// Setup test
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock person repository
	mockRepo := repository.NewMockPersonRepository(ctrl)
	service := svc.NewPersonService(mockRepo)

	// Test data
	userID := uuid.New()
	t.Run("Success", func(t *testing.T) {
		// Mock data
		people := []*models.Person{
			{ID: uuid.New(), UserID: userID, Name: "Person 1"},
			{ID: uuid.New(), UserID: userID, Name: "Person 2"},
		}

		// Setup expectations
		mockRepo.EXPECT().
			FindByUserID(gomock.Any(), userID).
			Return(people, nil)

		// Call the method
		count, err := service.CountPeople(context.Background(), userID)

		// Assert results
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
	})

	t.Run("EmptyResult", func(t *testing.T) {
		// Mock empty result
		mockRepo.EXPECT().
			FindByUserID(gomock.Any(), userID).
			Return([]*models.Person{}, nil)

		// Call the method
		count, err := service.CountPeople(context.Background(), userID)

		// Assert results
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		// Mock error
		expectedErr := assert.AnError
		mockRepo.EXPECT().
			FindByUserID(gomock.Any(), userID).
			Return(nil, expectedErr)

		// Call the method
		count, err := service.CountPeople(context.Background(), userID)

		// Assert results
		assert.ErrorIs(t, err, expectedErr)
		assert.Equal(t, int64(0), count)
	})
}
