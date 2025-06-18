#!/bin/bash

# Test Script for Investment Strategy API Changes
# Test user: test4@example.com / password123

BASE_URL="http://localhost:3000/api/v1"
EMAIL="test4@example.com"
PASSWORD="password123"

echo "🧪 Testing Investment Strategy API Changes"
echo "=========================================="

# Step 1: Login and get JWT token
echo "📝 Step 1: Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
ACCOUNT_ID=$(echo $LOGIN_RESPONSE | jq -r '.user.alpaca_account_id')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo "❌ Login failed!"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

echo "✅ Login successful!"
echo "🔑 Token: ${TOKEN:0:20}..."
echo "👤 Account ID: $ACCOUNT_ID"

# Step 2: Test Investment Strategy - Create Risk Profile
echo -e "\n📝 Step 2: Creating Risk Profile..."
RISK_PROFILE_RESPONSE=$(curl -s -X POST "$BASE_URL/investment/risk-profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "risk_tolerance": "moderate",
    "investment_timeline": "long_term",
    "financial_goals": ["retirement", "wealth_building"],
    "age_bracket": "26-35",
    "annual_income_bracket": "50000-75000",
    "investment_experience": "intermediate",
    "risk_capacity": "medium"
  }')

echo "Response: $RISK_PROFILE_RESPONSE"

# Step 3: Test Investment Strategy - Get Risk Profile
echo -e "\n📝 Step 3: Getting Risk Profile..."
GET_RISK_RESPONSE=$(curl -s -X GET "$BASE_URL/investment/risk-profile" \
  -H "Authorization: Bearer $TOKEN")

echo "Response: $GET_RISK_RESPONSE"

# Step 4: Test Investment Strategy - Create Portfolio
echo -e "\n📝 Step 4: Creating Portfolio..."
PORTFOLIO_RESPONSE=$(curl -s -X POST "$BASE_URL/investment/portfolio" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "positions": [
      {"symbol": "AAPL", "weight": 0.4},
      {"symbol": "GOOGL", "weight": 0.3},
      {"symbol": "MSFT", "weight": 0.3}
    ]
  }')

echo "Response: $PORTFOLIO_RESPONSE"

# Step 5: Test Investment Strategy - Get Portfolio
echo -e "\n📝 Step 5: Getting Portfolio..."
GET_PORTFOLIO_RESPONSE=$(curl -s -X GET "$BASE_URL/investment/portfolio" \
  -H "Authorization: Bearer $TOKEN")

echo "Response: $GET_PORTFOLIO_RESPONSE"

# Step 6: Test Trading Engine Integration - Purchase Portfolio
echo -e "\n📝 Step 6: Testing Trading Engine Integration (Purchase Portfolio)..."
PURCHASE_RESPONSE=$(curl -s -X POST "$BASE_URL/investment/purchase-portfolio" \
  -H "Authorization: Bearer $TOKEN")

echo "Response: $PURCHASE_RESPONSE"

# Step 7: Test Trading Engine - Get Orders
echo -e "\n📝 Step 7: Getting Trading Orders..."
ORDERS_RESPONSE=$(curl -s -X GET "$BASE_URL/trading/orders" \
  -H "Authorization: Bearer $TOKEN")

echo "Response: $ORDERS_RESPONSE"

echo -e "\n🏁 Test completed!"
echo "=========================================="
echo "Summary:"
echo "✅ Login: $(echo $LOGIN_RESPONSE | jq -r '.message // "Success"')"
echo "✅ Risk Profile: $(echo $RISK_PROFILE_RESPONSE | jq -r '.message // "Check response"')"
echo "✅ Portfolio: $(echo $PORTFOLIO_RESPONSE | jq -r '.message // "Check response"')"
echo "✅ Purchase: $(echo $PURCHASE_RESPONSE | jq -r '.message // "Check response"')" 