# AI-Powered Relationship Analysis System

## Overview

This implementation provides a comprehensive AI-powered relationship analysis system that evaluates user interactions with people and generates actionable recommendations. The system is designed to be flexible, scalable, and privacy-first.

## Architecture

### Core Components

1. **AI Service Layer** (`pkg/ai/`)
   - Provider abstraction supporting OpenAI and Anthropic
   - Structured prompt templates for consistency
   - JSON response parsing
   - Cost tracking and token usage monitoring

2. **Database Models** (`internal/models/models.go`)
   - `RelationshipAnalysis`: Stores AI-generated analysis with scores and insights
   - `AIRecommendation`: Action suggestions with priority and timing
   - `AIAnalysisJob`: Background job tracking for batch operations

3. **Repository Layer** (`internal/repository/analysis_repository.go`)
   - CRUD operations for analyses and recommendations
   - Query optimization with pagination
   - Job status tracking

4. **Service Layer** (`internal/services/analysis_service.go`)
   - Orchestrates AI analysis workflow
   - Manages recommendation generation
   - Handles batch processing
   - Caching strategy for cost optimization

5. **API Handlers** (`internal/handlers/analysis_handler.go`)
   - RESTful endpoints for analysis operations
   - Request validation and error handling
   - Response formatting

## API Endpoints

### Individual Person Analysis

```
GET    /api/v1/people/:id/analysis              # Get latest analysis
POST   /api/v1/people/:id/analysis/refresh      # Force new analysis
GET    /api/v1/people/:id/analysis/history      # Get analysis history
GET    /api/v1/people/:id/recommendations       # Get recommendations
```

### Global Recommendations

```
GET    /api/v1/recommendations                  # Get all active recommendations
POST   /api/v1/recommendations/:id/status       # Update recommendation status
```

### Analytics & Batch Operations

```
GET    /api/v1/analytics/insights               # Overall relationship insights
POST   /api/v1/analytics/batch-analyze          # Batch analyze multiple people
GET    /api/v1/analytics/jobs/:id               # Get batch job status
```

## Configuration

Add the following environment variables to your `.env` file:

```bash
# AI Provider Configuration
AI_PROVIDER=openai                              # or 'anthropic'
OPENAI_API_KEY=sk-...                          # Your OpenAI API key
OPENAI_MODEL=gpt-4o                            # Model to use
ANTHROPIC_API_KEY=sk-ant-...                   # Your Anthropic API key
ANTHROPIC_MODEL=claude-3-5-sonnet-20241022     # Model to use

# AI Service Settings
AI_MAX_TOKENS=2000                             # Max tokens per request
AI_TEMPERATURE=0.7                             # Response creativity (0-1)
AI_CACHE_ENABLED=true                          # Enable response caching
AI_CACHE_TTL=24h                               # Cache duration
AI_RATE_LIMIT_PER_USER=10                      # Requests per user per hour

# Feature Flag
FEATURE_AI_INSIGHTS=true                       # Enable AI features
```

## Database Migration

Run the migration to create the necessary tables:

```bash
# Using migrate CLI
migrate -path migrations -database "postgresql://..." up

# Or using the application's auto-migrate (uncomment in main.go)
# The migration includes:
# - relationship_analyses table
# - ai_recommendations table
# - ai_analysis_jobs table
```

## Features

### 1. Relationship Scoring (0-100)

Each analysis provides multiple dimensions:
- **Connection Strength**: Overall bond quality
- **Engagement Quality**: Interaction depth and meaningfulness
- **Communication Balance**: Two-way communication health
- **Energy Alignment**: Mutual energy impact
- **Relationship Health**: Overall trajectory
- **Overall Score**: Weighted composite score

### 2. AI-Generated Insights

- **Summary**: 2-3 sentence overview of the relationship
- **Key Insights**: 3-5 actionable observations
- **Patterns**: Recurring interaction patterns
- **Strengths**: Positive aspects to maintain
- **Concerns**: Areas needing attention
- **Trend Direction**: improving, stable, or declining

### 3. Smart Recommendations

Recommendations include:
- **Type**: reach_out, schedule_call, set_boundary, celebrate, check_in
- **Priority**: high, medium, low
- **Timing**: now, today, this_week, this_month
- **Suggested Actions**: Specific steps to take
- **Conversation Starters**: Context-aware opening lines
- **Estimated Impact**: Expected relationship improvement

### 4. Background Processing

- Automatic analysis refresh after 24 hours
- Batch processing for multiple relationships
- Job status tracking with progress indicators
- Cost tracking per job

## Usage Examples

### 1. Get Analysis for a Person

```bash
curl -X GET \
  https://api.vyve.app/api/v1/people/{person_id}/analysis \
  -H "Authorization: Bearer {token}"
```

Response:
```json
{
  "analysis": {
    "id": "uuid",
    "person_id": "uuid",
    "overall_score": 78.5,
    "connection_strength": 82.0,
    "engagement_quality": 75.0,
    "summary": "Strong, consistent relationship with regular meaningful interactions...",
    "key_insights": [
      "Regular weekly check-ins maintain connection",
      "Shared interests in technology create engaging conversations",
      "Recent decrease in interaction frequency noted"
    ],
    "patterns": [
      "Primarily communicates on weekends",
      "Energizing interactions correlate with shared activities"
    ],
    "strengths": [
      "Mutual support during challenging times",
      "Consistent communication rhythm"
    ],
    "concerns": [
      "Interaction frequency has decreased by 30% this month",
      "May benefit from more proactive outreach"
    ],
    "trend_direction": "stable",
    "analyzed_at": "2025-10-22T01:30:00Z"
  }
}
```

### 2. Get Recommendations

```bash
curl -X GET \
  https://api.vyve.app/api/v1/people/{person_id}/recommendations \
  -H "Authorization: Bearer {token}"
```

Response:
```json
{
  "recommendations": [
    {
      "id": "uuid",
      "type": "reach_out",
      "priority": "high",
      "title": "Schedule a catch-up call",
      "description": "It's been 2 weeks since your last interaction. A call would help maintain your connection.",
      "reasoning": "Interaction frequency has decreased, and this relationship benefits from regular contact.",
      "suggested_actions": [
        "Send a text asking about their recent project",
        "Propose a specific time for a call this week"
      ],
      "conversation_starters": [
        "Hey! I was thinking about you. How did that presentation go?",
        "It's been too long! Want to catch up over coffee this week?"
      ],
      "timing": "this_week",
      "estimated_impact": "high",
      "status": "pending"
    }
  ]
}
```

### 3. Batch Analyze Multiple People

```bash
curl -X POST \
  https://api.vyve.app/api/v1/analytics/batch-analyze \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "person_ids": ["uuid1", "uuid2", "uuid3"]
  }'
```

Response:
```json
{
  "job": {
    "id": "uuid",
    "status": "pending",
    "total_items": 3,
    "processed_items": 0,
    "progress": 0
  },
  "message": "Batch analysis job created"
}
```

### 4. Update Recommendation Status

```bash
curl -X POST \
  https://api.vyve.app/api/v1/recommendations/{recommendation_id}/status \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed"
  }'
```

## Cost Management

### Optimization Strategies

1. **Caching**: Analyses are cached for 24 hours by default
2. **Smart Refresh**: Only refresh when new interactions occur
3. **Batch Processing**: Group multiple analyses to reduce overhead
4. **Token Limits**: Configurable max tokens per request
5. **Rate Limiting**: Per-user rate limits prevent abuse

### Cost Tracking

Each analysis and job tracks:
- Total tokens used
- Estimated cost (calculated based on provider pricing)
- Processing time

Monitor costs via the database:
```sql
SELECT 
  DATE(analyzed_at) as date,
  COUNT(*) as analyses,
  SUM(tokens_used) as total_tokens,
  AVG(processing_time_ms) as avg_time_ms
FROM relationship_analyses
GROUP BY DATE(analyzed_at)
ORDER BY date DESC;
```

## Privacy & Security

### Data Protection

1. **Encryption**: Sensitive notes are encrypted before AI processing
2. **User Consent**: AI features require explicit user opt-in
3. **Data Minimization**: Only necessary data sent to AI providers
4. **Retention Policies**: Configurable analysis history retention
5. **Audit Logging**: All AI operations are logged

### Compliance

- GDPR compliant with data export and deletion
- User can disable AI features at any time
- Transparent about data usage in privacy policy

## Extensibility

### Adding New AI Providers

1. Implement the `Provider` interface in `pkg/ai/`
2. Add configuration in `config.go`
3. Update `NewService()` to support the new provider

Example:
```go
type CustomProvider struct {
    config CustomConfig
}

func (p *CustomProvider) Analyze(ctx context.Context, req AnalysisRequest) (*AnalysisResponse, error) {
    // Implementation
}

func (p *CustomProvider) GenerateRecommendations(ctx context.Context, req RecommendationRequest) (*RecommendationResponse, error) {
    // Implementation
}
```

### Custom Analysis Templates

Modify prompt templates in provider files to customize:
- Analysis focus areas
- Recommendation types
- Response format
- Language and tone

## Monitoring & Debugging

### Logging

Enable debug logging:
```bash
LOG_LEVEL=debug
```

### Metrics to Monitor

1. **Analysis Success Rate**: Track failed vs successful analyses
2. **Token Usage**: Monitor costs and optimize prompts
3. **Processing Time**: Identify performance bottlenecks
4. **User Engagement**: Track recommendation acceptance rates
5. **Cache Hit Rate**: Optimize caching strategy

### Common Issues

**Issue**: AI service returns errors
- Check API keys are valid
- Verify rate limits not exceeded
- Check network connectivity

**Issue**: Analysis quality is poor
- Adjust temperature setting
- Modify prompt templates
- Ensure sufficient interaction data

**Issue**: High costs
- Enable caching
- Reduce max_tokens
- Implement stricter rate limiting

## Testing

### Unit Tests

```bash
go test ./internal/services/analysis_service_test.go
go test ./pkg/ai/...
```

### Integration Tests

```bash
# Test with real AI providers (requires API keys)
AI_PROVIDER=openai OPENAI_API_KEY=sk-... go test -tags=integration ./...
```

### Manual Testing

Use the provided test script:
```bash
./test_ai_analysis.sh {person_id}
```

## Future Enhancements

1. **ML Model Training**: Train custom models on aggregated data
2. **Predictive Analytics**: Forecast relationship health trends
3. **Multi-Language Support**: Analyze relationships in different languages
4. **Integration with Calendar**: Suggest meeting times
5. **Voice Analysis**: Analyze tone and sentiment from calls
6. **Relationship Network**: Analyze connections between people
7. **Custom Metrics**: User-defined scoring dimensions

## Support

For issues or questions:
- GitHub Issues: https://github.com/vyve/vyve-backend/issues
- Email: support@vyve.app
- Documentation: https://docs.vyve.app

## License

Proprietary - See LICENSE file for details
