# Merging AI Recommendations with Nudges

## Rationale

The original implementation created a separate `AIRecommendation` model alongside the existing `Nudge` model. However, these serve essentially the same purpose:

### Original Separation:
- **Nudges**: System-generated reminders (rule-based)
- **AI Recommendations**: AI-generated suggestions (AI-powered)

### Problem:
- Duplicate functionality
- Two separate tables for similar data
- Two separate APIs for recommendations
- Confusing for frontend to manage both
- Harder to prioritize mixed recommendations

### Solution: Unified Model
Enhance the existing `Nudge` model to support both AI-generated and rule-based recommendations.

## Benefits

1. **Single Source of Truth**: One table, one API for all recommendations
2. **Simpler Frontend**: One endpoint to fetch all suggestions
3. **Better Prioritization**: Mix AI and rule-based recommendations naturally
4. **Unified Notifications**: Single notification system
5. **Easier Analytics**: Track all recommendation engagement in one place
6. **Backward Compatible**: Existing nudge functionality preserved

## Changes Made

### 1. Enhanced Nudge Model

**Added Fields**:
- `AnalysisID` - Links to AI analysis (if AI-generated)
- `Source` - 'ai' or 'system' to distinguish origin
- `Reasoning` - Why this recommendation (AI-specific)
- `SuggestedActions` - Specific action steps
- `ConversationStarters` - Opening lines for conversations
- `Timing` - When to act (now, today, this_week, this_month)
- `EstimatedImpact` - Expected relationship improvement
- `Status` - Unified status tracking (pending, seen, accepted, completed, dismissed)
- `AcceptedAt`, `CompletedAt`, `DismissedAt` - Detailed status timestamps
- `Provider`, `Model` - AI metadata for cost tracking

**Preserved Fields**:
- All existing nudge fields remain
- Legacy `Seen`/`ActedOn` fields kept for backward compatibility
- `Action`/`ActionData` for system-generated nudges

### 2. Database Migration Needed

**New Migration** (`000005_merge_recommendations_into_nudges.up.sql`):
```sql
-- Add new columns to nudges table
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS analysis_id UUID REFERENCES relationship_analyses(id) ON DELETE SET NULL;
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS source VARCHAR(20) DEFAULT 'system';
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS reasoning TEXT;
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS suggested_actions TEXT[];
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS conversation_starters TEXT[];
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS timing VARCHAR(50);
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS estimated_impact VARCHAR(20);
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'pending';
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS accepted_at TIMESTAMP;
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS completed_at TIMESTAMP;
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS dismissed_at TIMESTAMP;
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS provider VARCHAR(50);
ALTER TABLE nudges ADD COLUMN IF NOT EXISTS model VARCHAR(100);

-- Update message column to TEXT if not already
ALTER TABLE nudges ALTER COLUMN message TYPE TEXT;

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_nudges_analysis_id ON nudges(analysis_id);
CREATE INDEX IF NOT EXISTS idx_nudges_source ON nudges(source);
CREATE INDEX IF NOT EXISTS idx_nudges_status ON nudges(status);
CREATE INDEX IF NOT EXISTS idx_nudges_timing ON nudges(timing);

-- Drop ai_recommendations table if it exists
DROP TABLE IF EXISTS ai_recommendations;
```

### 3. API Changes

**No Breaking Changes** - Existing nudge endpoints continue to work:
- `GET /api/v1/nudges` - Now returns both system and AI nudges
- `POST /api/v1/nudges/:id/seen` - Works for all nudges
- `POST /api/v1/nudges/:id/acted` - Works for all nudges

**New/Updated Endpoints**:
- `GET /api/v1/nudges?source=ai` - Filter by source
- `GET /api/v1/nudges?status=pending` - Filter by status
- `POST /api/v1/nudges/:id/status` - Update status (accepted, completed, dismissed)
- `GET /api/v1/people/:id/nudges` - Get nudges for specific person (both types)

**Removed Endpoints** (were in AI implementation):
- `/api/v1/recommendations/*` - Now use `/api/v1/nudges/*`

### 4. Service Layer Updates

**NudgeService** now handles:
- System-generated nudges (existing logic)
- AI-generated nudges (from analysis service)
- Unified status management
- Mixed prioritization

**AnalysisService** now:
- Creates `Nudge` records instead of `AIRecommendation`
- Sets `source='ai'` and `analysis_id`
- Populates AI-specific fields

## Migration Guide

### For Existing Nudges
No action needed - existing nudges continue to work with `source='system'`

### For New AI Recommendations
```go
// Old way (removed)
recommendation := &models.AIRecommendation{
    UserID: userID,
    PersonID: &personID,
    Type: "reach_out",
    Title: "Check in with John",
    // ...
}

// New way (unified)
nudge := &models.Nudge{
    UserID: userID,
    PersonID: &personID,
    AnalysisID: &analysisID,
    Source: "ai",  // Mark as AI-generated
    Type: "reach_out",
    Title: "Check in with John",
    Reasoning: "You haven't connected in 2 weeks...",
    SuggestedActions: []string{"Send a text", "Schedule a call"},
    ConversationStarters: []string{"Hey! How's the project going?"},
    Timing: "this_week",
    EstimatedImpact: "high",
    Provider: "openai",
    Model: "gpt-4o",
    // ...
}
```

### Frontend Changes Needed

**Before** (two separate calls):
```typescript
// Get system nudges
const nudges = await api.get('/nudges');

// Get AI recommendations
const recommendations = await api.get('/recommendations');

// Merge and sort manually
const all = [...nudges, ...recommendations].sort(...);
```

**After** (single call):
```typescript
// Get all nudges (system + AI)
const nudges = await api.get('/nudges');

// Filter if needed
const aiNudges = nudges.filter(n => n.source === 'ai');
const systemNudges = nudges.filter(n => n.source === 'system');

// Or filter by status
const pending = nudges.filter(n => n.status === 'pending');
```

## Response Format

### Unified Nudge Response
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "person_id": "uuid",
  "analysis_id": "uuid",
  "source": "ai",
  "type": "reach_out",
  "title": "Check in with Sarah",
  "message": "It's been 2 weeks since your last interaction...",
  "priority": "high",
  "reasoning": "Interaction frequency has decreased by 30%...",
  "suggested_actions": [
    "Send a text asking about her recent project",
    "Propose a specific time for a call this week"
  ],
  "conversation_starters": [
    "Hey! I was thinking about you. How did that presentation go?",
    "It's been too long! Want to catch up over coffee this week?"
  ],
  "timing": "this_week",
  "estimated_impact": "high",
  "status": "pending",
  "provider": "openai",
  "model": "gpt-4o",
  "expires_at": "2025-10-29T00:00:00Z",
  "created_at": "2025-10-22T01:30:00Z"
}
```

### System Nudge (backward compatible)
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "person_id": "uuid",
  "source": "system",
  "type": "reconnect",
  "title": "Time to reconnect with John",
  "message": "You haven't interacted with John in 14 days",
  "priority": "medium",
  "action": "create_interaction",
  "action_data": {"person_id": "uuid"},
  "status": "pending",
  "seen": false,
  "acted_on": false,
  "created_at": "2025-10-22T01:30:00Z"
}
```

## Advantages of Unified Approach

### 1. Natural Prioritization
```go
// Get all pending nudges, sorted by priority and timing
nudges := nudgeService.GetPending(userID)
// Returns mixed AI and system nudges, properly prioritized
```

### 2. Unified Analytics
```sql
-- Track all recommendation engagement
SELECT 
  source,
  COUNT(*) as total,
  SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed,
  AVG(CASE WHEN completed_at IS NOT NULL 
      THEN EXTRACT(EPOCH FROM (completed_at - created_at)) 
      END) as avg_time_to_complete
FROM nudges
GROUP BY source;
```

### 3. Simpler Notification System
```go
// One notification handler for all nudges
func SendNudgeNotification(nudge *models.Nudge) {
    // Works for both AI and system nudges
    if nudge.Source == "ai" {
        // Include conversation starters
    }
    // Send notification
}
```

### 4. Better User Experience
- Single "Suggestions" or "Nudges" section in UI
- Natural mixing of AI and rule-based suggestions
- Consistent interaction model (accept, dismiss, complete)
- Unified notification preferences

## Implementation Checklist

- [x] Update Nudge model with new fields
- [x] Remove AIRecommendation model
- [ ] Create migration to add columns to nudges table
- [ ] Update NudgeRepository to handle new fields
- [ ] Update AnalysisService to create Nudges instead of AIRecommendations
- [ ] Update NudgeService to handle both sources
- [ ] Update NudgeHandler to support new status updates
- [ ] Update routes (remove /recommendations endpoints)
- [ ] Update main.go (remove AIRecommendation from migrations)
- [ ] Update tests
- [ ] Update API documentation
- [ ] Update frontend to use unified endpoint

## Rollout Strategy

### Phase 1: Backend Changes (This PR)
1. Update models and migrations
2. Update services and repositories
3. Keep both endpoints temporarily for backward compatibility

### Phase 2: Frontend Updates
1. Update frontend to use `/nudges` endpoint
2. Add filtering by source if needed
3. Update UI to show AI-specific fields

### Phase 3: Cleanup
1. Remove deprecated `/recommendations` endpoints
2. Remove old migration files
3. Update documentation

## Questions & Answers

**Q: What happens to existing nudges?**  
A: They continue to work unchanged. They'll have `source='system'` by default.

**Q: Can we still distinguish AI vs system nudges?**  
A: Yes, via the `source` field. Filter with `?source=ai` or `?source=system`.

**Q: What about the separate recommendations table?**  
A: It's removed. The migration drops `ai_recommendations` table.

**Q: Is this a breaking change?**  
A: No for existing nudge functionality. Yes for the new AI recommendations API (but it wasn't in production yet).

**Q: Can we add more sources later (e.g., 'calendar', 'email')?**  
A: Yes! The `source` field supports any value. Easy to extend.

## Conclusion

This merge simplifies the architecture, improves user experience, and makes the system more maintainable while preserving all functionality. The unified model is more flexible and easier to extend in the future.
