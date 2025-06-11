package handlers

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/seunghoon34/trading-app/services/trading-engine/internal/logger"
)

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

func sendAlpacaResponse(c *gin.Context, res *http.Response) {
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	c.Header("Content-Type", "application/json")
	c.String(res.StatusCode, string(body))
}

func CreateOrder(c *gin.Context) {
	var OrderData struct {
		AccountID string `json:"account_id"`
		Side      string `json:"side"`
		Symbol    string `json:"symbol"`
		Qty       string `json:"qty"`
		Notional  string `json:"notional"`
	}

	if err := c.ShouldBindJSON(&OrderData); err != nil {
		logger.WithFields(map[string]interface{}{
			"account_id": OrderData.AccountID,
			"symbol":     OrderData.Symbol,
			"error":      err.Error(),
			"action":     "order_create_failed",
		}).Error("Order creation failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if OrderData.Side != "buy" && OrderData.Side != "sell" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Side must be 'buy' or 'sell'"})
		return
	}

	url := fmt.Sprintf("https://broker-api.sandbox.alpaca.markets/v1/trading/accounts/%s/orders", OrderData.AccountID)

	var payload *strings.Reader
	if OrderData.Qty != "" {
		payload = strings.NewReader(fmt.Sprintf(`{"type":"market","time_in_force":"day","side":"%s","symbol":"%s","qty":"%s"}`, OrderData.Side, OrderData.Symbol, OrderData.Qty))
	} else {
		payload = strings.NewReader(fmt.Sprintf(`{"type":"market","time_in_force":"day","side":"%s","symbol":"%s","notional":"%s"}`, OrderData.Side, OrderData.Symbol, OrderData.Notional))
	}

	logger.WithFields(map[string]interface{}{
		"account_id": OrderData.AccountID,
		"symbol":     OrderData.Symbol,
		"side":       OrderData.Side,
		"qty":        OrderData.Qty,
		"action":     "order_create_attempt",
	}).Info("Order creation started")

	res, err := makeAlpacaRequest("POST", url, payload)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute order"})
		logger.WithFields(map[string]interface{}{
			"account_id": OrderData.AccountID,
			"symbol":     OrderData.Symbol,
			"error":      err.Error(),
			"action":     "order_create_failed",
		}).Error("Order creation failed")
		return
	}

	sendAlpacaResponse(c, res)
}

func GetOrder(c *gin.Context) {
	orderID := c.Param("order_id")
	accountID := c.Param("account_id")

	url := fmt.Sprintf("https://broker-api.sandbox.alpaca.markets/v1/trading/accounts/%s/orders/%s", accountID, orderID)

	logger.WithFields(map[string]interface{}{
		"account_id": accountID,
		"order_id":   orderID,
		"action":     "get_order_attempt",
	}).Info("Get order started")

	res, err := makeAlpacaRequest("GET", url, nil)
	if err != nil {

		logger.WithFields(map[string]interface{}{
			"account_id": accountID,
			"order_id":   orderID,
			"error":      err.Error(),
			"action":     "get_order_failed",
		}).Error("Get order failed")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute request"})
		return
	}

	sendAlpacaResponse(c, res)
}

func GetOrders(c *gin.Context) {
	accountID := c.Param("account_id")

	url := fmt.Sprintf("https://broker-api.sandbox.alpaca.markets/v1/trading/accounts/%s/orders", accountID)

	logger.WithFields(map[string]interface{}{
		"account_id": accountID,
		"action":     "get_all_orders_attempt",
	}).Info("Get all orders started")

	res, err := makeAlpacaRequest("GET", url, nil)
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"account_id": accountID,
			"error":      err.Error(),
			"action":     "get_all_orders_failed",
		}).Error("Get all orders failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute request"})
		return
	}

	sendAlpacaResponse(c, res)
}

func DeleteOrder(c *gin.Context) {
	orderID := c.Param("order_id")
	accountID := c.Param("account_id")

	url := fmt.Sprintf("https://broker-api.sandbox.alpaca.markets/v1/trading/accounts/%s/orders/%s", accountID, orderID)

	logger.WithFields(map[string]interface{}{
		"account_id": accountID,
		"order_id":   orderID,
		"action":     "delete_order_attempt",
	}).Info("Delete order started")

	res, err := makeAlpacaRequest("DELETE", url, nil)
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"account_id": accountID,
			"order_id":   orderID,
			"error":      err.Error(),
			"action":     "delete_order_failed",
		}).Error("Delete order failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute request"})
		return
	}

	sendAlpacaResponse(c, res)
}

func DeleteAllOrders(c *gin.Context) {
	accountID := c.Param("account_id")

	url := fmt.Sprintf("https://broker-api.sandbox.alpaca.markets/v1/trading/accounts/%s/orders", accountID)

	logger.WithFields(map[string]interface{}{
		"account_id": accountID,
		"action":     "delete_all_orders_attempt",
	}).Info("Delete all orders started")
	res, err := makeAlpacaRequest("DELETE", url, nil)
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"account_id": accountID,
			"error":      err.Error(),
			"action":     "delete_all_orders_failed",
		}).Error("Delete all orders failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute request"})
		return
	}

	sendAlpacaResponse(c, res)
}
