# AI Analysis Environment Variables

Add these to your `.env` file to enable AI-powered relationship analysis:

```bash
# ============================================
# AI Analysis Configuration
# ============================================

# AI Provider Selection
# Options: 'openai' or 'anthropic'
AI_PROVIDER=openai

# OpenAI Configuration
OPENAI_API_KEY=sk-proj-your-key-here
OPENAI_MODEL=gpt-4o

# Anthropic Configuration (alternative to OpenAI)
ANTHROPIC_API_KEY=sk-ant-your-key-here
ANTHROPIC_MODEL=claude-3-5-sonnet-20241022

# AI Service Settings
AI_MAX_TOKENS=2000
AI_TEMPERATURE=0.7
AI_CACHE_ENABLED=true
AI_CACHE_TTL=24h
AI_RATE_LIMIT_PER_USER=10

# Feature Flag
FEATURE_AI_INSIGHTS=true
```

## Configuration Details

### AI_PROVIDER
- **Default**: `openai`
- **Options**: `openai`, `anthropic`
- **Description**: Which AI provider to use for analysis

### OPENAI_API_KEY
- **Required**: Yes (if using OpenAI)
- **Format**: `sk-proj-...` or `sk-...`
- **Get Key**: https://platform.openai.com/api-keys

### OPENAI_MODEL
- **Default**: `gpt-4o`
- **Options**: `gpt-4o`, `gpt-4-turbo`, `gpt-4`, `gpt-3.5-turbo`
- **Recommended**: `gpt-4o` for best quality/cost balance

### ANTHROPIC_API_KEY
- **Required**: Yes (if using Anthropic)
- **Format**: `sk-ant-...`
- **Get Key**: https://console.anthropic.com/

### ANTHROPIC_MODEL
- **Default**: `claude-3-5-sonnet-20241022`
- **Options**: `claude-3-5-sonnet-20241022`, `claude-3-opus-20240229`, `claude-3-sonnet-20240229`
- **Recommended**: `claude-3-5-sonnet-20241022` for best performance

### AI_MAX_TOKENS
- **Default**: `2000`
- **Range**: `500-4000`
- **Description**: Maximum tokens per AI request (affects cost and response length)

### AI_TEMPERATURE
- **Default**: `0.7`
- **Range**: `0.0-1.0`
- **Description**: Controls response creativity (0=deterministic, 1=creative)

### AI_CACHE_ENABLED
- **Default**: `true`
- **Description**: Enable caching of AI responses to reduce costs

### AI_CACHE_TTL
- **Default**: `24h`
- **Format**: Duration string (e.g., `1h`, `30m`, `24h`)
- **Description**: How long to cache AI responses

### AI_RATE_LIMIT_PER_USER
- **Default**: `10`
- **Description**: Maximum AI requests per user per hour

### FEATURE_AI_INSIGHTS
- **Default**: `false`
- **Description**: Master switch to enable/disable AI features

## Cost Estimates

### OpenAI Pricing (as of 2024)
- **GPT-4o**: ~$2.50 per 1M input tokens, ~$10 per 1M output tokens
- **GPT-4-turbo**: ~$10 per 1M input tokens, ~$30 per 1M output tokens
- **Estimated per analysis**: $0.01-0.05 depending on interaction history

### Anthropic Pricing (as of 2024)
- **Claude 3.5 Sonnet**: ~$3 per 1M input tokens, ~$15 per 1M output tokens
- **Estimated per analysis**: $0.01-0.05 depending on interaction history

### Cost Optimization Tips
1. Enable caching (`AI_CACHE_ENABLED=true`)
2. Set reasonable token limits (`AI_MAX_TOKENS=2000`)
3. Use rate limiting (`AI_RATE_LIMIT_PER_USER=10`)
4. Choose cost-effective models (GPT-4o or Claude 3.5 Sonnet)
5. Batch analyze when possible

## Quick Start

1. Copy these variables to your `.env` file
2. Get an API key from OpenAI or Anthropic
3. Set `FEATURE_AI_INSIGHTS=true`
4. Restart your application
5. Test with: `GET /api/v1/people/{id}/analysis`

## Troubleshooting

### "AI insights feature is disabled"
- Set `FEATURE_AI_INSIGHTS=true`
- Restart the application

### "Failed to initialize AI service"
- Check your API key is valid
- Verify the key format is correct
- Ensure you have credits/quota available

### "OpenAI API error (status 401)"
- API key is invalid or expired
- Generate a new key from the provider dashboard

### "Rate limit exceeded"
- Increase `AI_RATE_LIMIT_PER_USER`
- Wait for the rate limit window to reset
- Consider upgrading your API plan

## Security Best Practices

1. **Never commit `.env` files** to version control
2. **Rotate API keys regularly** (every 90 days)
3. **Use environment-specific keys** (dev, staging, prod)
4. **Monitor API usage** to detect anomalies
5. **Set up billing alerts** with your AI provider
6. **Restrict API key permissions** to minimum required
7. **Use secrets management** in production (AWS Secrets Manager, etc.)

## Production Checklist

- [ ] API keys stored in secure secrets manager
- [ ] Rate limiting configured appropriately
- [ ] Caching enabled to reduce costs
- [ ] Monitoring and alerting set up
- [ ] Billing alerts configured with provider
- [ ] Error handling and fallbacks tested
- [ ] Privacy policy updated to mention AI usage
- [ ] User consent flow implemented
- [ ] Cost tracking dashboard created
- [ ] Backup provider configured (optional)
