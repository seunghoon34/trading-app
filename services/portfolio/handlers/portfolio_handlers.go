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
	accountID := c.Param("account_id")

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
	accountID := c.Param("account_id")

	positions, err := getPositionsHelper(accountID) // Use your helper!
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get positions" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account_id": accountID,
		"positions":  positions,
	})
}

func GetPortfolioWorth(c *gin.Context) {
	accountID := c.Param("account_id")

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
	accountID := c.Param("account_id")

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
