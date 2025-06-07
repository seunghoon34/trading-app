// File: services/market-data/handlers/market.go

package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/seunghoon34/trading-app/services/market-data/config"
)

func GetLatestQuote(c *gin.Context) {
	symbol := c.Param("symbol")

	// Use the same endpoint that worked in curl
	endpoint := fmt.Sprintf("/v2/stocks/%s/quotes/latest", symbol)
	data, err := config.MakeMarketDataRequest(endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to fetch quote for %s: %v", symbol, err),
		})
		return
	}

	// Return the raw JSON data (you can parse it later)
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(data))
}

func GetLatestBar(c *gin.Context) {
	symbol := c.Param("symbol")

	// Use the same endpoint that worked in curl
	endpoint := fmt.Sprintf("/v2/stocks/%s/bars/latest", symbol)
	data, err := config.MakeMarketDataRequest(endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to fetch bar for %s: %v", symbol, err),
		})
		return
	}

	// Return the raw JSON data (you can parse it later)
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(data))
}

func GetLatestQuoteMulti(c *gin.Context) {
	symbols := c.QueryArray("symbols")

	if len(symbols) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No symbols provided. Use ?symbols=AAPL&symbols=TSLA or ?symbols=AAPL,TSLA",
		})
		return
	}

	// Join symbols with comma and URL encode
	symbolsParam := url.QueryEscape(strings.Join(symbols, ","))

	endpoint := fmt.Sprintf("/v2/stocks/quotes/latest?symbols=%s", symbolsParam)
	data, err := config.MakeMarketDataRequest(endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to fetch quotes for %v: %v", symbols, err),
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(data))
}

func GetLatestBarMulti(c *gin.Context) {
	symbols := c.QueryArray("symbols")

	if len(symbols) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No symbols provided. Use ?symbols=AAPL&symbols=TSLA or ?symbols=AAPL,TSLA",
		})
		return
	}

	// Join symbols with comma and URL encode
	symbolsParam := url.QueryEscape(strings.Join(symbols, ","))

	endpoint := fmt.Sprintf("/v2/stocks/bars/latest?symbols=%s", symbolsParam)
	data, err := config.MakeMarketDataRequest(endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to fetch bars for %v: %v", symbols, err),
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(data))
}
