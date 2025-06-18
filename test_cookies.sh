#!/bin/bash

# Test Script for Investment Strategy API Changes
# Uses cookies for authentication
# Test user: test4@example.com / password123

BASE_URL="http://localhost:3000/api/v1"
EMAIL="test4@example.com"
PASSWORD="password123"
COOKIE_JAR="cookies.txt"

echo "🧪 Testing Investment Strategy API Changes (Cookie-based Auth)"
echo "=============================================================="

# Clean up any existing cookie file
rm -f $COOKIE_JAR

# Step 1: Login and save cookies
echo "📝 Step 1: Logging in..."
LOGIN_RESPONSE=$(curl -s -c $COOKIE_JAR -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

echo "Login Response: $LOGIN_RESPONSE"

ACCOUNT_ID=$(echo $LOGIN_RESPONSE | jq -r '.account_id')
if [ "$ACCOUNT_ID" = "null" ] || [ -z "$ACCOUNT_ID" ]; then
    echo "❌ Login failed - no account_id found!"
    exit 1
fi

echo "✅ Login successful!"
echo "👤 Account ID: $ACCOUNT_ID"

# Step 2: Test Investment Strategy - Create Risk Profile
echo -e "\n📝 Step 2: Creating Risk Profile..."
RISK_PROFILE_RESPONSE=$(curl -s -b $COOKIE_JAR -X POST "$BASE_URL/investment/risk-profile" \
  -H "Content-Type: application/json" \
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
GET_RISK_RESPONSE=$(curl -s -b $COOKIE_JAR -X GET "$BASE_URL/investment/risk-profile")

echo "Response: $GET_RISK_RESPONSE"

# Step 4: Test Investment Strategy - Create Portfolio
echo -e "\n📝 Step 4: Creating Portfolio..."
PORTFOLIO_RESPONSE=$(curl -s -b $COOKIE_JAR -X POST "$BASE_URL/investment/portfolio" \
  -H "Content-Type: application/json" \
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
GET_PORTFOLIO_RESPONSE=$(curl -s -b $COOKIE_JAR -X GET "$BASE_URL/investment/portfolio")

echo "Response: $GET_PORTFOLIO_RESPONSE"

# Step 6: Test Trading Engine Integration - Purchase Portfolio
echo -e "\n📝 Step 6: Testing Trading Engine Integration (Purchase Portfolio)..."
PURCHASE_RESPONSE=$(curl -s -b $COOKIE_JAR -X POST "$BASE_URL/investment/portfolio/purchase")

echo "Response: $PURCHASE_RESPONSE"

# Step 7: Test Trading Engine - Get Orders
echo -e "\n📝 Step 7: Getting Trading Orders..."
ORDERS_RESPONSE=$(curl -s -b $COOKIE_JAR -X GET "$BASE_URL/trading/orders")

echo "Response: $ORDERS_RESPONSE"

# Cleanup
rm -f $COOKIE_JAR

echo -e "\n🏁 Test completed!"
echo "=============================================================="
echo "Summary:"
echo "✅ Login: $(echo $LOGIN_RESPONSE | jq -r '.message // "Success"')"
echo "✅ Risk Profile: $(echo $RISK_PROFILE_RESPONSE | jq -r '.message // "Check response"')"
echo "✅ Portfolio: $(echo $PORTFOLIO_RESPONSE | jq -r '.message // "Check response"')"
echo "✅ Purchase: $(echo $PURCHASE_RESPONSE | jq -r '.message // "Check response"')" 