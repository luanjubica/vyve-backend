package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// OpenAIConfig represents OpenAI provider configuration
type OpenAIConfig struct {
	APIKey      string
	Model       string
	MaxTokens   int
	Temperature float64
}

// OpenAIProvider implements the Provider interface for OpenAI
type OpenAIProvider struct {
	config     OpenAIConfig
	httpClient *http.Client
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(config OpenAIConfig) (*OpenAIProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}
	
	if config.Model == "" {
		config.Model = "gpt-4o" // Default to GPT-4o
	}
	
	if config.MaxTokens == 0 {
		config.MaxTokens = 2000
	}
	
	if config.Temperature == 0 {
		config.Temperature = 0.7
	}
	
	return &OpenAIProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// Analyze generates a relationship analysis using OpenAI
func (p *OpenAIProvider) Analyze(ctx context.Context, req AnalysisRequest) (*AnalysisResponse, error) {
	startTime := time.Now()
	log.Printf("[OPENAI_PROVIDER] Starting analysis for person: %s", req.PersonName)
	
	prompt := p.buildAnalysisPrompt(req)
	log.Printf("[OPENAI_PROVIDER] Built prompt, calling OpenAI API with model=%s", p.config.Model)
	
	response, tokensUsed, err := p.callOpenAI(ctx, prompt)
	if err != nil {
		log.Printf("[OPENAI_PROVIDER] ❌ OpenAI API call failed: %v", err)
		return nil, fmt.Errorf("OpenAI API call failed: %w", err)
	}
	
	log.Printf("[OPENAI_PROVIDER] Received response, tokens used: %d", tokensUsed)
	
	analysis, err := p.parseAnalysisResponse(response)
	if err != nil {
		log.Printf("[OPENAI_PROVIDER] ❌ Failed to parse response: %v", err)
		return nil, fmt.Errorf("failed to parse analysis response: %w", err)
	}
	
	analysis.TokensUsed = tokensUsed
	analysis.ProcessingTimeMs = int(time.Since(startTime).Milliseconds())
	
	log.Printf("[OPENAI_PROVIDER] ✅ Analysis completed in %dms", analysis.ProcessingTimeMs)
	return analysis, nil
}

// GenerateRecommendations generates action recommendations using OpenAI
func (p *OpenAIProvider) GenerateRecommendations(ctx context.Context, req RecommendationRequest) (*RecommendationResponse, error) {
	prompt := p.buildRecommendationPrompt(req)
	
	response, tokensUsed, err := p.callOpenAI(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API call failed: %w", err)
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
func (p *OpenAIProvider) GetProviderName() string {
	return "openai"
}

// GetModelName returns the model name
func (p *OpenAIProvider) GetModelName() string {
	return p.config.Model
}

// callOpenAI makes a request to the OpenAI API
func (p *OpenAIProvider) callOpenAI(ctx context.Context, prompt string) (string, int, error) {
	log.Printf("[OPENAI_API] Preparing request to OpenAI...")
	reqBody := map[string]interface{}{
		"model": p.config.Model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are an expert relationship analyst. Provide insightful, empathetic, and actionable analysis. Always respond with valid JSON.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens":  p.config.MaxTokens,
		"temperature": p.config.Temperature,
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("[OPENAI_API] ❌ Failed to marshal request: %v", err)
		return "", 0, err
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[OPENAI_API] ❌ Failed to create request: %v", err)
		return "", 0, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	
	// Log API key for debugging (first/last chars only for security)
	keyLen := len(p.config.APIKey)
	if keyLen > 20 {
		log.Printf("[OPENAI_API] Using API key: %s...%s (length: %d)", p.config.APIKey[:10], p.config.APIKey[keyLen-4:], keyLen)
	} else {
		log.Printf("[OPENAI_API] Using API key length: %d", keyLen)
	}
	
	log.Printf("[OPENAI_API] Sending request to OpenAI...")
	startTime := time.Now()
	resp, err := p.httpClient.Do(req)
	if err != nil {
		log.Printf("[OPENAI_API] ❌ HTTP request failed after %v: %v", time.Since(startTime), err)
		return "", 0, err
	}
	defer resp.Body.Close()
	
	log.Printf("[OPENAI_API] Received response with status %d in %v", resp.StatusCode, time.Since(startTime))
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[OPENAI_API] ❌ Failed to read response body: %v", err)
		return "", 0, err
	}
	
	if resp.StatusCode != http.StatusOK {
		log.Printf("[OPENAI_API] ❌ Non-OK status %d, response: %s", resp.StatusCode, string(body))
		return "", 0, fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
	}
	
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"usage"`
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return "", 0, err
	}
	
	if len(result.Choices) == 0 {
		return "", 0, fmt.Errorf("no response from OpenAI")
	}
	
	return result.Choices[0].Message.Content, result.Usage.TotalTokens, nil
}

// buildAnalysisPrompt builds the prompt for relationship analysis
func (p *OpenAIProvider) buildAnalysisPrompt(req AnalysisRequest) string {
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
func (p *OpenAIProvider) buildRecommendationPrompt(req RecommendationRequest) string {
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

// parseAnalysisResponse parses the OpenAI response into AnalysisResponse
func (p *OpenAIProvider) parseAnalysisResponse(response string) (*AnalysisResponse, error) {
	log.Printf("[OPENAI_PROVIDER] Raw response length: %d, first 100 chars: %s", len(response), truncate(response, 100))
	
	// Try to extract JSON from markdown code blocks if present
	cleaned := extractJSON(response)
	log.Printf("[OPENAI_PROVIDER] After extraction, length: %d, first 100 chars: %s", len(cleaned), truncate(cleaned, 100))
	
	var result AnalysisResponse
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		log.Printf("[OPENAI_PROVIDER] ❌ JSON parse error. Full response: %s", response)
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	
	return &result, nil
}

// parseRecommendationResponse parses the OpenAI response into recommendations
func (p *OpenAIProvider) parseRecommendationResponse(response string) ([]Recommendation, error) {
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

// extractJSON extracts JSON from markdown code blocks
func extractJSON(text string) string {
	// Remove markdown code blocks if present
	if len(text) > 7 && text[:3] == "```" {
		// Find the first newline after ```
		start := 3
		for start < len(text) && text[start] != '\n' {
			start++
		}
		start++ // Skip the newline
		
		// Find the closing ```
		end := len(text) - 3
		if end > start && text[end:] == "```" {
			return text[start:end]
		}
	}
	return text
}
