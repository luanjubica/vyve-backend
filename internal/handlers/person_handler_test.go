package handlers_test

import (
	"context"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vyve/vyve-backend/internal/handlers"
	"github.com/vyve/vyve-backend/internal/models"
	"github.com/vyve/vyve-backend/internal/services"
)

func TestPersonHandler_CountPeople(t *testing.T) {
	// Setup test
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock person service
	mockService := services.NewMockPersonService(ctrl)
	handler := handlers.NewPersonHandler(mockService)

	// Create a test Fiber app
	app := fiber.New()
	app.Get("/people/count", handler.CountPeople)

	t.Run("Success", func(t *testing.T) {
		// Mock data
		userID := uuid.New()
		expectedCount := int64(5)

		// Setup expectations
		mockService.EXPECT().
			CountPeople(gomock.Any(), userID).
			Return(expectedCount, nil)

		// Create request
	req := httptest.NewRequest("GET", "/people/count", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userID", userID.String()))

	// Execute request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result struct {
		Count int64 `json:"count"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, result.Count)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/people/count", nil)
		// No user ID in context

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/people/count", nil)
		req = req.WithContext(context.WithValue(req.Context(), "userID", "invalid-uuid"))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		userID := uuid.New()

		mockService.EXPECT().
			CountPeople(gomock.Any(), userID).
			Return(int64(0), errors.New("database error"))

		req := httptest.NewRequest("GET", "/people/count", nil)
		req = req.WithContext(context.WithValue(req.Context(), "userID", userID.String()))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}
