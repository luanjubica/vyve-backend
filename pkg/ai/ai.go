package ai

import (
	"context"
	"time"
)

// Provider represents an AI provider interface
type Provider interface {
	// Analyze generates a relationship analysis
	Analyze(ctx context.Context, req AnalysisRequest) (*AnalysisResponse, error)
	
	// GenerateRecommendations generates action recommendations
	GenerateRecommendations(ctx context.Context, req RecommendationRequest) (*RecommendationResponse, error)
	
	// GetProviderName returns the provider name
	GetProviderName() string
	
	// GetModelName returns the model name being used
	GetModelName() string
}

// AnalysisRequest represents a request for relationship analysis
type AnalysisRequest struct {
	PersonName        string
	Relationship      string
	InteractionCount  int
	RecentInteractions []InteractionData
	EnergyPattern     string
	HealthScore       float64
	LastInteraction   *time.Time
	Context           []string
	Notes             string
}

// InteractionData represents interaction data for analysis
type InteractionData struct {
	Date         time.Time
	EnergyImpact string
	Quality      int
	Duration     int
	Context      []string
	Notes        string
}

// AnalysisResponse represents the AI analysis response
type AnalysisResponse struct {
	// Scores (0-100)
	ConnectionStrength   float64
	EngagementQuality    float64
	CommunicationBalance float64
	EnergyAlignment      float64
	RelationshipHealth   float64
	OverallScore         float64
	
	// Analysis content
	Summary        string
	KeyInsights    []string
	Patterns       []string
	Strengths      []string
	Concerns       []string
	TrendDirection string // improving, stable, declining
	
	// Metadata
	TokensUsed       int
	ProcessingTimeMs int
}

// RecommendationRequest represents a request for recommendations
type RecommendationRequest struct {
	PersonName         string
	Relationship       string
	Analysis           *AnalysisResponse
	RecentInteractions []InteractionData
	LastInteraction    *time.Time
	Context            []string
}

// RecommendationResponse represents AI-generated recommendations
type RecommendationResponse struct {
	Recommendations []Recommendation
	TokensUsed      int
}

// Recommendation represents a single recommendation
type Recommendation struct {
	Type                 string   // reach_out, schedule_call, set_boundary, celebrate, check_in
	Priority             string   // high, medium, low
	Title                string
	Description          string
	Reasoning            string
	SuggestedActions     []string
	ConversationStarters []string
	Timing               string // now, today, this_week, this_month
	EstimatedImpact      string // high, medium, low
	ExpiresIn            *time.Duration
}

// Config represents AI service configuration
type Config struct {
	Provider         string // openai, anthropic
	OpenAIKey        string
	OpenAIModel      string
	AnthropicKey     string
	AnthropicModel   string
	MaxTokens        int
	Temperature      float64
	CacheEnabled     bool
	CacheTTL         time.Duration
	RateLimitPerUser int
}

// Service represents the AI service
type Service struct {
	provider Provider
	config   Config
}

// NewService creates a new AI service
func NewService(config Config) (*Service, error) {
	var provider Provider
	var err error
	
	switch config.Provider {
	case "openai":
		provider, err = NewOpenAIProvider(OpenAIConfig{
			APIKey:      config.OpenAIKey,
			Model:       config.OpenAIModel,
			MaxTokens:   config.MaxTokens,
			Temperature: config.Temperature,
		})
	case "anthropic":
		provider, err = NewAnthropicProvider(AnthropicConfig{
			APIKey:      config.AnthropicKey,
			Model:       config.AnthropicModel,
			MaxTokens:   config.MaxTokens,
			Temperature: config.Temperature,
		})
	default:
		// Default to OpenAI
		provider, err = NewOpenAIProvider(OpenAIConfig{
			APIKey:      config.OpenAIKey,
			Model:       config.OpenAIModel,
			MaxTokens:   config.MaxTokens,
			Temperature: config.Temperature,
		})
	}
	
	if err != nil {
		return nil, err
	}
	
	return &Service{
		provider: provider,
		config:   config,
	}, nil
}

// Analyze performs relationship analysis
func (s *Service) Analyze(ctx context.Context, req AnalysisRequest) (*AnalysisResponse, error) {
	return s.provider.Analyze(ctx, req)
}

// GenerateRecommendations generates action recommendations
func (s *Service) GenerateRecommendations(ctx context.Context, req RecommendationRequest) (*RecommendationResponse, error) {
	return s.provider.GenerateRecommendations(ctx, req)
}

// GetProviderName returns the current provider name
func (s *Service) GetProviderName() string {
	return s.provider.GetProviderName()
}

// GetModelName returns the current model name
func (s *Service) GetModelName() string {
	return s.provider.GetModelName()
}
