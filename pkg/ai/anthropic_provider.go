package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AnthropicConfig represents Anthropic provider configuration
type AnthropicConfig struct {
	APIKey      string
	Model       string
	MaxTokens   int
	Temperature float64
}

// AnthropicProvider implements the Provider interface for Anthropic Claude
type AnthropicProvider struct {
	config     AnthropicConfig
	httpClient *http.Client
}

// NewAnthropicProvider creates a new Anthropic provider
func NewAnthropicProvider(config AnthropicConfig) (*AnthropicProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("Anthropic API key is required")
	}
	
	if config.Model == "" {
		config.Model = "claude-3-5-sonnet-20241022" // Default to Claude 3.5 Sonnet
	}
	
	if config.MaxTokens == 0 {
		config.MaxTokens = 2000
	}
	
	if config.Temperature == 0 {
		config.Temperature = 0.7
	}
	
	return &AnthropicProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// Analyze generates a relationship analysis using Anthropic Claude
func (p *AnthropicProvider) Analyze(ctx context.Context, req AnalysisRequest) (*AnalysisResponse, error) {
	startTime := time.Now()
	
	prompt := p.buildAnalysisPrompt(req)
	
	response, tokensUsed, err := p.callAnthropic(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("Anthropic API call failed: %w", err)
	}
	
	analysis, err := p.parseAnalysisResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse analysis response: %w", err)
	}
	
	analysis.TokensUsed = tokensUsed
	analysis.ProcessingTimeMs = int(time.Since(startTime).Milliseconds())
	
	return analysis, nil
}

// GenerateRecommendations generates action recommendations using Anthropic Claude
func (p *AnthropicProvider) GenerateRecommendations(ctx context.Context, req RecommendationRequest) (*RecommendationResponse, error) {
	prompt := p.buildRecommendationPrompt(req)
	
	response, tokensUsed, err := p.callAnthropic(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("Anthropic API call failed: %w", err)
	}
	
	recommendations, err := p.parseRecommendationResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse recommendation response: %w", err)
	}
	
	return &RecommendationResponse{
		Recommendations: recommendations,
		TokensUsed:      tokensUsed,
	}, nil
}

// GetProviderName returns the provider name
func (p *AnthropicProvider) GetProviderName() string {
	return "anthropic"
}

// GetModelName returns the model name
func (p *AnthropicProvider) GetModelName() string {
	return p.config.Model
}

// callAnthropic makes a request to the Anthropic API
func (p *AnthropicProvider) callAnthropic(ctx context.Context, prompt string) (string, int, error) {
	reqBody := map[string]interface{}{
		"model": p.config.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens":  p.config.MaxTokens,
		"temperature": p.config.Temperature,
		"system":      "You are an expert relationship analyst. Provide insightful, empathetic, and actionable analysis. Always respond with valid JSON.",
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, err
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 0, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.config.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	
	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("Anthropic API error (status %d): %s", resp.StatusCode, string(body))
	}
	
	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return "", 0, err
	}
	
	if len(result.Content) == 0 {
		return "", 0, fmt.Errorf("no response from Anthropic")
	}
	
	totalTokens := result.Usage.InputTokens + result.Usage.OutputTokens
	return result.Content[0].Text, totalTokens, nil
}

// buildAnalysisPrompt builds the prompt for relationship analysis
func (p *AnthropicProvider) buildAnalysisPrompt(req AnalysisRequest) string {
	prompt := fmt.Sprintf(`Analyze the following relationship and provide a comprehensive assessment.

**Person Details:**
- Name: %s
- Relationship Type: %s
- Total Interactions: %d
- Current Health Score: %.1f/100
- Energy Pattern: %s
`, req.PersonName, req.Relationship, req.InteractionCount, req.HealthScore, req.EnergyPattern)
	
	if req.LastInteraction != nil {
		daysSince := int(time.Since(*req.LastInteraction).Hours() / 24)
		prompt += fmt.Sprintf("- Last Interaction: %d days ago\n", daysSince)
	}
	
	if len(req.Context) > 0 {
		prompt += fmt.Sprintf("- Context: %v\n", req.Context)
	}
	
	prompt += "\n**Recent Interactions:**\n"
	for i, interaction := range req.RecentInteractions {
		if i >= 10 {
			break // Limit to 10 most recent
		}
		daysSince := int(time.Since(interaction.Date).Hours() / 24)
		prompt += fmt.Sprintf("- %d days ago: %s energy, quality %d/5", daysSince, interaction.EnergyImpact, interaction.Quality)
		if interaction.Duration > 0 {
			prompt += fmt.Sprintf(", %d minutes", interaction.Duration)
		}
		if len(interaction.Context) > 0 {
			prompt += fmt.Sprintf(", context: %v", interaction.Context)
		}
		prompt += "\n"
	}
	
	prompt += `

**Task:** Provide a detailed analysis in JSON format with the following structure:
{
  "connection_strength": <0-100>,
  "engagement_quality": <0-100>,
  "communication_balance": <0-100>,
  "energy_alignment": <0-100>,
  "relationship_health": <0-100>,
  "overall_score": <0-100>,
  "summary": "<2-3 sentence overview>",
  "key_insights": ["<insight 1>", "<insight 2>", "<insight 3>"],
  "patterns": ["<pattern 1>", "<pattern 2>"],
  "strengths": ["<strength 1>", "<strength 2>"],
  "concerns": ["<concern 1>", "<concern 2>"],
  "trend_direction": "<improving|stable|declining>"
}

Provide actionable, empathetic insights. Focus on patterns, not individual interactions.`
	
	return prompt
}

// buildRecommendationPrompt builds the prompt for generating recommendations
func (p *AnthropicProvider) buildRecommendationPrompt(req RecommendationRequest) string {
	prompt := fmt.Sprintf(`Based on the relationship analysis, generate 2-4 actionable recommendations.

**Person:** %s (%s)
**Overall Score:** %.1f/100
**Trend:** %s

**Analysis Summary:**
%s

**Key Concerns:**
%v

**Recent Interaction Pattern:**
`, req.PersonName, req.Relationship, req.Analysis.OverallScore, req.Analysis.TrendDirection, req.Analysis.Summary, req.Analysis.Concerns)
	
	for i, interaction := range req.RecentInteractions {
		if i >= 5 {
			break
		}
		daysSince := int(time.Since(interaction.Date).Hours() / 24)
		prompt += fmt.Sprintf("- %d days ago: %s energy\n", daysSince, interaction.EnergyImpact)
	}
	
	prompt += `

**Task:** Generate recommendations in JSON format:
{
  "recommendations": [
    {
      "type": "<reach_out|schedule_call|set_boundary|celebrate|check_in>",
      "priority": "<high|medium|low>",
      "title": "<short title>",
      "description": "<detailed description>",
      "reasoning": "<why this matters>",
      "suggested_actions": ["<action 1>", "<action 2>"],
      "conversation_starters": ["<starter 1>", "<starter 2>"],
      "timing": "<now|today|this_week|this_month>",
      "estimated_impact": "<high|medium|low>"
    }
  ]
}

Focus on practical, specific actions. Prioritize based on urgency and impact.`
	
	return prompt
}

// parseAnalysisResponse parses the Anthropic response into AnalysisResponse
func (p *AnthropicProvider) parseAnalysisResponse(response string) (*AnalysisResponse, error) {
	// Try to extract JSON from markdown code blocks if present
	response = extractJSON(response)
	
	var result AnalysisResponse
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	
	return &result, nil
}

// parseRecommendationResponse parses the Anthropic response into recommendations
func (p *AnthropicProvider) parseRecommendationResponse(response string) ([]Recommendation, error) {
	// Try to extract JSON from markdown code blocks if present
	response = extractJSON(response)
	
	var result struct {
		Recommendations []Recommendation `json:"recommendations"`
	}
	
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	
	return result.Recommendations, nil
}
