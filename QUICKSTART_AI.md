# Quick Start Guide: AI-Powered Relationship Analysis

Get your AI analysis system up and running in 5 minutes!

## Step 1: Get API Keys (2 minutes)

### Option A: OpenAI (Recommended)
1. Go to https://platform.openai.com/api-keys
2. Sign in or create an account
3. Click "Create new secret key"
4. Copy the key (starts with `sk-proj-` or `sk-`)
5. Add billing information if needed

### Option B: Anthropic
1. Go to https://console.anthropic.com/
2. Sign in or create an account
3. Navigate to API Keys
4. Create a new key
5. Copy the key (starts with `sk-ant-`)

## Step 2: Configure Environment (1 minute)

Add to your `.env` file:

```bash
# Enable AI Features
FEATURE_AI_INSIGHTS=true

# Choose Provider (openai or anthropic)
AI_PROVIDER=openai

# Add Your API Key
OPENAI_API_KEY=sk-proj-YOUR-KEY-HERE
# OR
# ANTHROPIC_API_KEY=sk-ant-YOUR-KEY-HERE

# Optional: Customize Settings
AI_MAX_TOKENS=2000
AI_TEMPERATURE=0.7
AI_CACHE_ENABLED=true
AI_CACHE_TTL=24h
AI_RATE_LIMIT_PER_USER=10
```

## Step 3: Run Database Migration (1 minute)

### Option A: Using migrate CLI
```bash
migrate -path migrations -database "postgresql://user:pass@localhost:5432/vyve_dev?sslmode=disable" up
```

### Option B: Using the application
Uncomment the migration line in `cmd/api/main.go`:
```go
if err := migrateDatabase(db); err != nil {
    log.Fatalf("Failed to run migrations: %v", err)
}
```

Then restart the application.

## Step 4: Start the Server (30 seconds)

```bash
# Development
make dev

# Or directly
go run cmd/api/main.go
```

Look for this log message:
```
AI service initialized with provider: openai
```

## Step 5: Test It! (30 seconds)

### Get Authentication Token
```bash
# Login first
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "your@email.com",
    "password": "yourpassword"
  }'

# Save the token from response
export TOKEN="your-jwt-token-here"
```

### Test Analysis Endpoint
```bash
# Replace {person_id} with an actual person ID from your database
curl -X GET http://localhost:8080/api/v1/people/{person_id}/analysis \
  -H "Authorization: Bearer $TOKEN"
```

### Expected Response
```json
{
  "analysis": {
    "id": "uuid",
    "overall_score": 78.5,
    "connection_strength": 82.0,
    "summary": "Strong, consistent relationship...",
    "key_insights": [
      "Regular weekly check-ins maintain connection",
      "Shared interests create engaging conversations"
    ],
    "recommendations": [...]
  }
}
```

## Troubleshooting

### "AI insights feature is disabled"
**Solution**: Set `FEATURE_AI_INSIGHTS=true` in `.env` and restart

### "Failed to initialize AI service"
**Solution**: Check your API key is correct and has credits

### "OpenAI API error (status 401)"
**Solution**: Your API key is invalid. Generate a new one

### "No person found"
**Solution**: Create a person first:
```bash
curl -X POST http://localhost:8080/api/v1/people \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "relationship": "friend"
  }'
```

### "No interactions to analyze"
**Solution**: Add some interactions first:
```bash
curl -X POST http://localhost:8080/api/v1/interactions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "person_id": "{person_id}",
    "energy_impact": "energizing",
    "quality": 4,
    "duration": 30
  }'
```

## Quick Test Script

Save as `test_ai.sh`:
```bash
#!/bin/bash

# Configuration
API_URL="http://localhost:8080/api/v1"
EMAIL="test@example.com"
PASSWORD="password123"

# Login
echo "Logging in..."
TOKEN=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" | jq -r '.token')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
  echo "Login failed!"
  exit 1
fi

echo "Token: $TOKEN"

# Create a person
echo -e "\nCreating test person..."
PERSON_ID=$(curl -s -X POST "$API_URL/people" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Person","relationship":"friend"}' | jq -r '.person.id')

echo "Person ID: $PERSON_ID"

# Add some interactions
echo -e "\nAdding interactions..."
for i in {1..3}; do
  curl -s -X POST "$API_URL/interactions" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"person_id\":\"$PERSON_ID\",\"energy_impact\":\"energizing\",\"quality\":4}" > /dev/null
  echo "Added interaction $i"
done

# Get AI analysis
echo -e "\nGetting AI analysis..."
curl -s -X GET "$API_URL/people/$PERSON_ID/analysis" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo -e "\nDone!"
```

Run it:
```bash
chmod +x test_ai.sh
./test_ai.sh
```

## What's Next?

### Explore Features
1. **Get Recommendations**: `GET /api/v1/people/{id}/recommendations`
2. **Batch Analyze**: `POST /api/v1/analytics/batch-analyze`
3. **View Insights**: `GET /api/v1/analytics/insights`
4. **Update Status**: `POST /api/v1/recommendations/{id}/status`

### Monitor Usage
```bash
# Check analysis count
curl -X GET http://localhost:8080/api/v1/analytics/insights \
  -H "Authorization: Bearer $TOKEN"

# View job status
curl -X GET http://localhost:8080/api/v1/analytics/jobs/{job_id} \
  -H "Authorization: Bearer $TOKEN"
```

### Optimize Costs
1. Enable caching: `AI_CACHE_ENABLED=true`
2. Set token limits: `AI_MAX_TOKENS=2000`
3. Use rate limiting: `AI_RATE_LIMIT_PER_USER=10`
4. Monitor usage in provider dashboard

### Production Deployment
1. Move API keys to secrets manager
2. Set up monitoring and alerts
3. Configure billing limits
4. Enable audit logging
5. Update privacy policy
6. Test error scenarios

## Common Use Cases

### 1. Analyze All Relationships
```bash
# Get all people
PEOPLE=$(curl -s -X GET "$API_URL/people" \
  -H "Authorization: Bearer $TOKEN" | jq -r '.people[].id')

# Batch analyze
curl -X POST "$API_URL/analytics/batch-analyze" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"person_ids\": [$(echo $PEOPLE | tr ' ' ',')]}"
```

### 2. Get Daily Recommendations
```bash
# Get active recommendations
curl -X GET "$API_URL/recommendations" \
  -H "Authorization: Bearer $TOKEN" | jq '.recommendations[] | select(.timing == "today")'
```

### 3. Track Recommendation Completion
```bash
# Mark as completed
curl -X POST "$API_URL/recommendations/{id}/status" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'
```

## Need Help?

- **Documentation**: See `AI_ANALYSIS_README.md`
- **Configuration**: See `AI_ENV_VARS.md`
- **Summary**: See `IMPLEMENTATION_SUMMARY.md`
- **Issues**: Open a GitHub issue
- **Support**: support@vyve.app

## Success! 🎉

You now have AI-powered relationship analysis running! The system will:
- ✅ Analyze relationships automatically
- ✅ Generate smart recommendations
- ✅ Track relationship health over time
- ✅ Provide actionable insights
- ✅ Help users maintain meaningful connections

Happy analyzing! 🚀
