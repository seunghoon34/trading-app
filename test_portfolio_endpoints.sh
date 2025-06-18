#!/bin/bash

# Test script for portfolio performance endpoints
# Make sure you're logged in and have a valid session cookie

API_BASE="http://localhost:3000"

echo "🧪 Testing Portfolio Performance Endpoints"
echo "==========================================="

echo ""
echo "1. Testing current portfolio performance..."
curl -X GET "$API_BASE/api/v1/portfolio/performance" \
  -b cookies.txt \
  -H "Content-Type: application/json" \
  -w "\nStatus: %{http_code}\n" \
  -s | jq '.' 2>/dev/null || echo "Response received (jq not available for formatting)"

echo ""
echo "2. Testing multi-timeframe performance..."
curl -X GET "$API_BASE/api/v1/portfolio/performance/all" \
  -b cookies.txt \
  -H "Content-Type: application/json" \
  -w "\nStatus: %{http_code}\n" \
  -s | jq '.' 2>/dev/null || echo "Response received (jq not available for formatting)"

echo ""
echo "3. Testing portfolio positions..."
curl -X GET "$API_BASE/api/v1/portfolio/positions" \
  -b cookies.txt \
  -H "Content-Type: application/json" \
  -w "\nStatus: %{http_code}\n" \
  -s | jq '.' 2>/dev/null || echo "Response received (jq not available for formatting)"

echo ""
echo "4. Testing portfolio value..."
curl -X GET "$API_BASE/api/v1/portfolio/value" \
  -b cookies.txt \
  -H "Content-Type: application/json" \
  -w "\nStatus: %{http_code}\n" \
  -s | jq '.' 2>/dev/null || echo "Response received (jq not available for formatting)"

echo ""
echo "✅ Portfolio endpoint testing complete!"
echo ""
echo "Note: Make sure you have:"
echo "  - API Gateway running on port 3000"
echo "  - Portfolio service running on port 8084"
echo "  - Valid authentication cookies in cookies.txt"
echo "  - Alpaca API credentials configured" 