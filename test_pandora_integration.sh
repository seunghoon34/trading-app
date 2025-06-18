#!/bin/bash

# Test script for Pandora Integration
echo "🔮 Testing Pandora Integration..."

API_BASE="http://localhost:3000/api/v1"

# Test 1: Check if services are running
echo -e "\n1. Checking service health..."
echo "API Gateway:"
curl -s -o /dev/null -w "%{http_code}" "$API_BASE/../health" || echo "FAILED"

echo -e "\nInvestment Strategy Service:"
curl -s -o /dev/null -w "%{http_code}" "http://localhost:8089/health" || echo "FAILED"

echo -e "\nCrewAI Portfolio Service:"
curl -s -o /dev/null -w "%{http_code}" "http://localhost:8000/health" || echo "FAILED"

# Test 2: Test API Gateway routing
echo -e "\n\n2. Testing API Gateway routes..."

echo "Testing investment-strategy route:"
curl -s -o /dev/null -w "%{http_code}" "$API_BASE/investment-strategy/health" || echo "FAILED"

echo -e "\nTesting crewai-portfolio route:"
curl -s -o /dev/null -w "%{http_code}" "$API_BASE/crewai-portfolio/health" || echo "FAILED"

# Test 3: Test with authentication (requires valid JWT)
echo -e "\n\n3. Testing authenticated endpoints (requires login)..."
echo "Note: These will fail with 401 if not logged in, which is expected"

echo "Testing risk profile endpoint:"
curl -s -o /dev/null -w "%{http_code}" "$API_BASE/investment-strategy/risk-profile" || echo "FAILED"

echo -e "\nTesting portfolio generation endpoint:"
curl -s -o /dev/null -w "%{http_code}" -X POST "$API_BASE/crewai-portfolio/generate-portfolio" \
  -H "Content-Type: application/json" \
  -d '{"risk_tolerance": "moderate", "investment_timeline": "medium_term", "financial_goals": ["wealth_building"]}' || echo "FAILED"

echo -e "\n\n✅ Test complete! Check the HTTP status codes above."
echo "200 = Success, 401 = Unauthorized (expected for protected routes), 404 = Not Found" 