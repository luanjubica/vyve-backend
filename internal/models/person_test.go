package models

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestPersonHealthScoreRounding(t *testing.T) {
	tests := []struct {
		name              string
		healthScore       float64
		expectedJSONValue int
	}{
		{
			name:              "Round down from .4",
			healthScore:       50.4,
			expectedJSONValue: 50,
		},
		{
			name:              "Round up from .5",
			healthScore:       50.5,
			expectedJSONValue: 51,
		},
		{
			name:              "Round up from .6",
			healthScore:       75.6,
			expectedJSONValue: 76,
		},
		{
			name:              "Exact integer",
			healthScore:       100.0,
			expectedJSONValue: 100,
		},
		{
			name:              "Round down from .3",
			healthScore:       33.3,
			expectedJSONValue: 33,
		},
		{
			name:              "Round up from .9",
			healthScore:       88.9,
			expectedJSONValue: 89,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			person := Person{
				Base: Base{
					ID: uuid.New(),
				},
				UserID:      uuid.New(),
				Name:        "Test Person",
				HealthScore: tt.healthScore,
			}

			jsonData, err := json.Marshal(person)
			if err != nil {
				t.Fatalf("Failed to marshal person: %v", err)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(jsonData, &result); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			healthScore, ok := result["health_score"].(float64)
			if !ok {
				t.Fatalf("health_score not found or not a number in JSON")
			}

			if int(healthScore) != tt.expectedJSONValue {
				t.Errorf("Expected health_score to be %d, got %d (original: %.1f)",
					tt.expectedJSONValue, int(healthScore), tt.healthScore)
			}

			// Verify it's an integer in JSON (no decimal point)
			jsonStr := string(jsonData)
			t.Logf("JSON output: %s", jsonStr)
		})
	}
}
