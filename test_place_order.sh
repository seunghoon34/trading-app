#!/bin/bash

# Script to place a test order to populate the portfolio
API_BASE="http://localhost:3000"

echo "🔥 Placing Test Order to Populate Portfolio"
echo "============================================"

# Place a small market buy order for Apple stock
echo "Placing order: Buy 5 shares of AAPL..."
curl -X POST "$API_BASE/api/v1/trading/orders" \
  -b cookies.txt \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "AAPL",
    "qty": "5",
    "side": "buy", 
    "type": "market",
    "time_in_force": "day"
  }' \
  -w "\nStatus: %{http_code}\n" \
  -s | jq '.' 2>/dev/null || echo "Response received"

echo ""
echo "Order placed! Wait a few seconds then check:"
echo "1. Portfolio positions: curl -b cookies.txt $API_BASE/api/v1/portfolio/positions"
echo "2. Portfolio performance: curl -b cookies.txt $API_BASE/api/v1/portfolio/performance"
echo "3. Refresh your frontend dashboard"

echo ""
echo "Note: In sandbox mode, orders are executed immediately"
echo "      Your portfolio should now show the AAPL position" 