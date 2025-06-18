package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Performance struct {
	Timestamp     []int64   `json:"timestamp"`
	Equity        []float64 `json:"equity"`
	ProfitLoss    []float64 `json:"profit_loss"`
	ProfitLossPct []float64 `json:"profit_loss_pct"`
}

type BuyingPower struct {
	BuyingPower string `json:"buying_power"`
}

type Position struct {
	Symbol                 string `json:"symbol"`
	Quantity               string `json:"qty"`
	AvgEntryPrice          string `json:"avg_entry_price"`
	CurrentPrice           string `json:"current_price"`
	MarketValue            string `json:"market_value"`
	CostBasis              string `json:"cost_basis"`
	UnrealizedPL           string `json:"unrealized_pl"`
	UnrealizedPLPC         string `json:"unrealized_plpc"`
	UnrealizedIntradayPL   string `json:"unrealized_intraday_pl"`
	UnrealizedIntradayPLPC string `json:"unrealized_intraday_plpc"`
	Side                   string `json:"side"`
}

type PortfolioResponse struct {
	Positions []Position `json:"positions"`
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

// func sendAlpacaResponse(c *gin.Context, res *http.Response) {
// 	defer res.Body.Close()

// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
// 		return
// 	}

// 	c.Header("Content-Type", "application/json")
// 	c.String(res.StatusCode, string(body))
// }

func getPositionsHelper(account_id string) ([]Position, error) {
	url := fmt.Sprintf("https://broker-api.sandbox.alpaca.markets/v1/trading/accounts/%s/positions", account_id)
	fmt.Printf("Making request to: %s\n", url)

	res, err := makeAlpacaRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Response body: %s\n", string(body))

	var positions []Position
	if err := json.Unmarshal(body, &positions); err != nil {
		return nil, err
	}

	return positions, nil
}
func GetPosition(c *gin.Context) {
	symbol := c.Param("symbol")

	accountID := c.GetHeader("X-Account-ID") // ← Get from header
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID header missing"})
		return
	}

	url := fmt.Sprintf("https://broker-api.sandbox.alpaca.markets/v1/trading/accounts/%s/positions/%s", accountID, symbol)

	res, err := makeAlpacaRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute request"})
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Unmarshal directly into your struct - Go will ignore extra fields!
	var positions Position
	if err := json.Unmarshal(body, &positions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	// Return clean response with only the fields you defined
	c.JSON(http.StatusOK, gin.H{
		"account_id": accountID,
		"positions":  positions,
	})
}

func GetPositions(c *gin.Context) {
	accountID := c.GetHeader("X-Account-ID") // ← CORRECT: Get from header
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID header missing"})
		return
	}

	positions, err := getPositionsHelper(accountID) // Use your helper!
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get positions" + err.Error()})
		return
	}

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

	// Check if the API call was successful
	if res.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       "Failed to fetch account details from Alpaca",
			"status_code": res.StatusCode,
			"response":    string(body),
		})
		return
	}

	var buyingPower BuyingPower
	if err := json.Unmarshal(body, &buyingPower); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":        "Failed to parse account details",
			"raw_response": string(body),
			"parse_error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account_id": accountID,
		"positions":  positions,
		"Cash":       buyingPower.BuyingPower,
	})
}

func GetPortfolioWorth(c *gin.Context) {
	accountID := c.GetHeader("X-Account-ID") // ← Get from header
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID header missing"})
		return
	}

	positions, err := getPositionsHelper(accountID) // Get both data and error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get positions"})
		return
	}

	var worth float64
	for _, position := range positions { // Directly iterate over positions
		marketValue, _ := strconv.ParseFloat(position.MarketValue, 64)
		worth += marketValue
	}

	c.JSON(http.StatusOK, gin.H{
		"account_id":  accountID,
		"total_worth": worth,
	})
}

func GetPortfolioPerformance(c *gin.Context) {
	accountID := c.GetHeader("X-Account-ID") // ← Get from header
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID header missing"})
		return
	}

	positions, err := getPositionsHelper(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get positions"})
		return
	}

	var dailyPL, totalPL, totalCostBasis, totalMarketValue float64

	for _, position := range positions {
		unrealizedIntradayPL, _ := strconv.ParseFloat(position.UnrealizedIntradayPL, 64)
		unrealizedPL, _ := strconv.ParseFloat(position.UnrealizedPL, 64)
		costBasis, _ := strconv.ParseFloat(position.CostBasis, 64)
		marketValue, _ := strconv.ParseFloat(position.MarketValue, 64)

		dailyPL += unrealizedIntradayPL
		totalPL += unrealizedPL
		totalCostBasis += costBasis
		totalMarketValue += marketValue
	}

	// Calculate portfolio-level percentages
	var dailyPLPC, totalPLPC float64
	if totalMarketValue > 0 {
		totalPLPC = (totalPL / totalCostBasis) * 100               // Portfolio total return %
		dailyPLPC = (dailyPL / (totalMarketValue - dailyPL)) * 100 // Portfolio daily return %
	}

	c.JSON(http.StatusOK, gin.H{
		"account_id":         accountID,
		"daily_pl":           dailyPL,
		"daily_plpc":         dailyPLPC,
		"total_pl":           totalPL,
		"total_plpc":         totalPLPC,
		"total_market_value": totalMarketValue,
		"total_cost_basis":   totalCostBasis,
	})
}

func GetMultiTimeFramePerformance(c *gin.Context) {
	accountID := c.GetHeader("X-Account-ID")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID header missing"})
		return
	}

	// Portfolio history API works independently of positions
	// Even accounts with no positions can have cash equity history

	// Define the different periods and their parameters
	// Valid timeframes: 1Min, 5Min, 15Min, 1H, 1D
	periods := map[string]string{
		"1D": "period=1D&timeframe=1H",
		"1W": "period=1W&timeframe=1D",
		"1M": "period=1M&timeframe=1D",
		"1Y": "period=1A&timeframe=1D",
	}

	// Channel to collect results
	type result struct {
		period      string
		performance Performance
		err         error
	}

	results := make(chan result, len(periods))

	// Launch concurrent requests
	for period, params := range periods {
		go func(p, prms string) {
			url := fmt.Sprintf("https://broker-api.sandbox.alpaca.markets/v1/trading/accounts/%s/account/portfolio/history?%s&intraday_reporting=market_hours&pnl_reset=per_day&cashflow_types=NONE", accountID, prms)

			res, err := makeAlpacaRequest("GET", url, nil)
			if err != nil {
				results <- result{period: p, err: fmt.Errorf("failed to fetch %s performance: %w", p, err)}
				return
			}
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				results <- result{period: p, err: fmt.Errorf("failed to read %s response: %w", p, err)}
				return
			}

			if res.StatusCode != http.StatusOK {
				fmt.Printf("Alpaca API error for %s: Status %d, Body: %s\n", p, res.StatusCode, string(body))
				results <- result{period: p, err: fmt.Errorf("failed to fetch %s performance from Alpaca: status %d, response: %s", p, res.StatusCode, string(body))}
				return
			}

			var performance Performance
			if err := json.Unmarshal(body, &performance); err != nil {
				results <- result{period: p, err: fmt.Errorf("failed to parse %s performance: %w", p, err)}
				return
			}

			results <- result{period: p, performance: performance, err: nil}
		}(period, params)
	}

	// Collect all results
	performanceData := make(map[string]Performance)
	for i := 0; i < len(periods); i++ {
		result := <-results
		if result.err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to fetch performance data",
				"details": result.err.Error(),
			})
			return
		}
		performanceData[result.period] = result.performance
	}

	// Return successful response
	c.JSON(http.StatusOK, gin.H{
		"account_id": accountID,
		"1D":         performanceData["1D"],
		"1W":         performanceData["1W"],
		"1M":         performanceData["1M"],
		"1Y":         performanceData["1Y"],
	})
}
