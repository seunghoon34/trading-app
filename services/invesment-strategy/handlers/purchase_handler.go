package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seunghoon34/trading-app/services/investment-strategy/internal/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

// AccountDetails represents the account information from Alpaca
type AccountDetails struct {
	BuyingPower string `json:"buying_power"`
	AccountID   string `json:"id"`
	Status      string `json:"status"`
}

// OrderRequest represents the order request to trading service
type OrderRequest struct {
	AccountID string `json:"account_id"`
	Side      string `json:"side"`
	Symbol    string `json:"symbol"`
	Notional  string `json:"notional"`
}

// OrderResult represents the result of an individual order
type OrderResult struct {
	Symbol   string `json:"symbol"`
	Notional string `json:"notional"`
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
	OrderID  string `json:"order_id,omitempty"`
}

// PurchaseResult represents the overall purchase result
type PurchaseResult struct {
	TotalBuyingPower string        `json:"total_buying_power"`
	OrderResults     []OrderResult `json:"order_results"`
	SuccessCount     int           `json:"success_count"`
	FailureCount     int           `json:"failure_count"`
}

func makeAlpacaRequest(method, url string, payload io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	// Add authentication headers
	auth := os.Getenv("ALPACA_API_KEY") + ":" + os.Getenv("ALPACA_SECRET_KEY")
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", basicAuth)
	req.Header.Add("Accept", "application/json")

	if method == "POST" {
		req.Header.Add("Content-Type", "application/json")
	}

	return http.DefaultClient.Do(req)
}

func callTradingService(orderReq OrderRequest) (map[string]interface{}, error) {
	// Get trading service URL from environment variable
	tradingServiceURL := os.Getenv("TRADING_SERVICE_URL")
	if tradingServiceURL == "" {
		tradingServiceURL = "http://trading-engine:8083" // Default for local development
	}

	// Remove AccountID from the request body as it should be in header
	orderReqBody := map[string]string{
		"side":     orderReq.Side,
		"symbol":   orderReq.Symbol,
		"notional": orderReq.Notional,
	}

	jsonData, err := json.Marshal(orderReqBody)
	if err != nil {
		return nil, err
	}

	// Create request with proper headers
	req, err := http.NewRequest("POST", tradingServiceURL+"/orders", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Account-ID", orderReq.AccountID) // Set account ID in header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return result, fmt.Errorf("trading service error: %s", string(body))
	}

	return result, nil
}

func PurchasePortfolio(c *gin.Context) {
	// Get account_id from header
	accountID := c.GetHeader("X-Account-ID")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Account-ID header is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Step 1: Get buying power from Alpaca
	accountURL := fmt.Sprintf("https://broker-api.sandbox.alpaca.markets/v1/trading/accounts/%s/account", accountID)
	res, err := makeAlpacaRequest("GET", accountURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch account details"})
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read account response"})
		return
	}

	// Debug: Log the raw response
	fmt.Printf("Raw Alpaca response status: %d\n", res.StatusCode)
	fmt.Printf("Raw Alpaca response body: %s\n", string(body))

	// Check if the API call was successful
	if res.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       "Failed to fetch account details from Alpaca",
			"status_code": res.StatusCode,
			"response":    string(body),
		})
		return
	}

	var accountDetails AccountDetails
	if err := json.Unmarshal(body, &accountDetails); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":        "Failed to parse account details",
			"raw_response": string(body),
			"parse_error":  err.Error(),
		})
		return
	}

	// Debug: Log the buying power value
	fmt.Printf("Raw buying power from Alpaca: '%s'\n", accountDetails.BuyingPower)

	// Parse buying power to float
	buyingPower, err := strconv.ParseFloat(strings.TrimSpace(accountDetails.BuyingPower), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Invalid buying power format",
			"raw_value":   accountDetails.BuyingPower,
			"parse_error": err.Error(),
		})
		return
	}

	// Debug: Log the parsed buying power
	fmt.Printf("Parsed buying power: %f\n", buyingPower)

	if buyingPower <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient buying power"})
		return
	}

	// Step 2: Get portfolio positions from MongoDB
	var portfolio Portfolio
	err = mongo.PortfolioCollection.FindOne(ctx, bson.M{"alpaca_id": accountID}).Decode(&portfolio)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch portfolio"})
		return
	}

	if len(portfolio.Positions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Portfolio has no positions"})
		return
	}

	// Step 3: Calculate dollar amounts and execute orders
	var orderResults []OrderResult
	successCount := 0
	failureCount := 0

	for _, position := range portfolio.Positions {
		// Calculate dollar amount for this position (rounded down to no decimal)
		dollarAmount := math.Floor(buyingPower * position.Weight)

		if dollarAmount <= 0 {
			orderResults = append(orderResults, OrderResult{
				Symbol:   position.Symbol,
				Notional: "0",
				Success:  false,
				Error:    "Calculated amount is zero or negative",
			})
			failureCount++
			continue
		}

		// Prepare order request
		orderReq := OrderRequest{
			AccountID: accountID,
			Side:      "buy",
			Symbol:    strings.ToUpper(position.Symbol),
			Notional:  fmt.Sprintf("%.0f", dollarAmount), // No decimal places
		}

		// Call trading service
		result, err := callTradingService(orderReq)
		if err != nil {
			orderResults = append(orderResults, OrderResult{
				Symbol:   position.Symbol,
				Notional: orderReq.Notional,
				Success:  false,
				Error:    err.Error(),
			})
			failureCount++
			continue
		}

		// Extract order ID if available
		orderID := ""
		if id, exists := result["id"]; exists {
			if idStr, ok := id.(string); ok {
				orderID = idStr
			}
		}

		orderResults = append(orderResults, OrderResult{
			Symbol:   position.Symbol,
			Notional: orderReq.Notional,
			Success:  true,
			OrderID:  orderID,
		})
		successCount++
	}

	// Prepare final response
	purchaseResult := PurchaseResult{
		TotalBuyingPower: accountDetails.BuyingPower,
		OrderResults:     orderResults,
		SuccessCount:     successCount,
		FailureCount:     failureCount,
	}

	// Return appropriate status code
	if failureCount > 0 && successCount == 0 {
		// All orders failed
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "All orders failed",
			"result":  purchaseResult,
		})
	} else if failureCount > 0 {
		// Partial success
		c.JSON(http.StatusPartialContent, gin.H{
			"message": "Portfolio purchase completed with some failures",
			"result":  purchaseResult,
		})
	} else {
		// All orders succeeded
		c.JSON(http.StatusOK, gin.H{
			"message": "Portfolio purchased successfully",
			"result":  purchaseResult,
		})
	}
}
