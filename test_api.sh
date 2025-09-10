
#!/bin/bash

# Vyve API Test Script
# Usage: ./test_api.sh

BASE_URL="http://localhost:8080/api/v1"
EMAIL="alice@example.com"
PASSWORD="password123"
TOKEN=""

echo "🧪 Testing Vyve API Endpoints"
echo "========================================="

# 1. Login
echo "📍 1. Testing Login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "'$EMAIL'",
    "password": "'$PASSWORD'"
  }')

echo "Login Response: $LOGIN_RESPONSE"

# Extract token from response
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "❌ Login failed - no token received"
    exit 1
fi

echo "✅ Login successful! Token: ${TOKEN:0:20}..."
echo ""

# 2. Get Profile
echo "📍 2. Testing Get Profile..."
PROFILE_RESPONSE=$(curl -s -X GET "$BASE_URL/users/me" \
  -H "Authorization: Bearer $TOKEN")

echo "Profile Response: $PROFILE_RESPONSE"
echo ""

# 3. Get Settings
echo "📍 3. Testing Get Settings..."
SETTINGS_RESPONSE=$(curl -s -X GET "$BASE_URL/users/me/settings" \
  -H "Authorization: Bearer $TOKEN")

echo "Settings Response: $SETTINGS_RESPONSE"
echo ""

# 4. Update Settings
echo "📍 4. Testing Update Settings..."
UPDATE_SETTINGS_RESPONSE=$(curl -s -X PUT "$BASE_URL/users/me/settings" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "theme": "dark",
    "notifications": true,
    "language": "en"
  }')

echo "Update Settings Response: $UPDATE_SETTINGS_RESPONSE"
echo ""

# 5. Get People Count
echo "📍 5. Testing Get People Count..."
PEOPLE_COUNT_RESPONSE=$(curl -s -X GET "$BASE_URL/people/count" \
  -H "Authorization: Bearer $TOKEN")

echo "People Count Response: $PEOPLE_COUNT_RESPONSE"
echo ""

# 6. List People
echo "📍 6. Testing List People..."
PEOPLE_LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/people" \
  -H "Authorization: Bearer $TOKEN")

echo "People List Response: $PEOPLE_LIST_RESPONSE"
echo ""

# 7. Create Person
echo "📍 7. Testing Create Person..."
CREATE_PERSON_RESPONSE=$(curl -s -X POST "$BASE_URL/people" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Person",
    "relationship": "friend",
    "notes": "Created via API test"
  }')

echo "Create Person Response: $CREATE_PERSON_RESPONSE"

# Extract person ID for further tests
PERSON_ID=$(echo $CREATE_PERSON_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo ""

if [ ! -z "$PERSON_ID" ]; then
    # 8. Get Person
    echo "📍 8. Testing Get Person..."
    PERSON_RESPONSE=$(curl -s -X GET "$BASE_URL/people/$PERSON_ID" \
      -H "Authorization: Bearer $TOKEN")

    echo "Get Person Response: $PERSON_RESPONSE"
    echo ""

    # 9. Update Person
    echo "📍 9. Testing Update Person..."
    UPDATE_PERSON_RESPONSE=$(curl -s -X PUT "$BASE_URL/people/$PERSON_ID" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "name": "Updated Test Person",
        "relationship": "close friend"
      }')

    echo "Update Person Response: $UPDATE_PERSON_RESPONSE"
    echo ""
fi

# 10. Get User Stats
echo "📍 10. Testing Get User Stats..."
STATS_RESPONSE=$(curl -s -X GET "$BASE_URL/users/me/stats" \
  -H "Authorization: Bearer $TOKEN")

echo "Stats Response: $STATS_RESPONSE"
echo ""

# 11. Logout
echo "📍 11. Testing Logout..."
LOGOUT_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/logout" \
  -H "Authorization: Bearer $TOKEN")

echo "Logout Response: $LOGOUT_RESPONSE"
echo ""

echo "🎉 API Test Complete!"
echo "========================================="