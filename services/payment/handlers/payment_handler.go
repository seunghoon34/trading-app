package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	rdb "github.com/seunghoon34/trading-app/services/payment/redis"
)

type ACHDetails struct {
	Id string `json:"id"`
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

func createACHDetails(account_id string) (*ACHDetails, error) {
	url := fmt.Sprintf("https://broker-api.sandbox.alpaca.markets/v1/accounts/%s/ach_relationships", account_id)
	payload := strings.NewReader("{\"bank_account_type\":\"CHECKING\",\"account_owner_name\":\"seunghoon han\",\"bank_account_number\":\"32131231ab\",\"bank_routing_number\":\"123103716\"}")
	res, err := makeAlpacaRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check HTTP status code
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Response body: %s\n", string(body))

	var ach_details ACHDetails
	if err := json.Unmarshal(body, &ach_details); err != nil {
		return nil, err
	}
	return &ach_details, nil

}

func retrieveACHDetails(account_id string) (*ACHDetails, error) {

	// Try cache first
	cacheKey := fmt.Sprintf("ach_details:%s", account_id)

	cached, err := rdb.Client.Get(context.Background(), cacheKey).Result()
	if err == nil {
		// Cache hit - unmarshal and return
		var ach_details ACHDetails
		if err := json.Unmarshal([]byte(cached), &ach_details); err == nil {
			fmt.Println("Cache hit for ACH details")
			return &ach_details, nil
		}
	}

	// Cache miss - fetch from API
	fmt.Println("Cache miss - fetching from API")

	url := fmt.Sprintf("https://broker-api.sandbox.alpaca.markets/v1/accounts/%s/ach_relationships", account_id)
	fmt.Printf("Making request to: %s\n", url)
	res, err := makeAlpacaRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check HTTP status code
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Response body: %s\n", string(body))

	var ach_relationships []ACHDetails
	if err := json.Unmarshal(body, &ach_relationships); err != nil {
		return nil, err
	}

	// Check if we have at least one relationship
	if len(ach_relationships) == 0 {
		ach_relationship, err := createACHDetails(account_id)
		if err != nil {
			return nil, fmt.Errorf("failed to create ACH details: %w", err)
		}
		// Cache the result for 1 hour
		achJSON, _ := json.Marshal(ach_relationship)
		rdb.Client.Set(context.Background(), cacheKey, achJSON, time.Hour)
		return ach_relationship, nil
	}

	// // Cache the result for 1 hour
	achJSON, _ := json.Marshal(ach_relationships[0])
	rdb.Client.Set(context.Background(), cacheKey, achJSON, time.Hour)

	// Return the first one (or you could add logic to pick a specific one)
	return &ach_relationships[0], nil

}

func DepositFunds(c *gin.Context) {
	amount := c.Param("amount")
	accountID := c.GetHeader("X-Account-ID")

	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Account-ID header is required"})
		return
	}

	fmt.Printf("Received account_id: '%s'\n", accountID)

	ach_details, err := retrieveACHDetails(accountID)
	if err != nil {
		fmt.Printf("Error retrieving ACH details: %v\n", err) // Add this debug line
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ACH details"})
		return
	}
	url := fmt.Sprintf("https://broker-api.sandbox.alpaca.markets/v1/accounts/%s/transfers", accountID)
	payload := strings.NewReader(fmt.Sprintf("{\"transfer_type\":\"ach\",\"direction\":\"INCOMING\",\"timing\":\"immediate\",\"relationship_id\":\"%s\",\"amount\":\"%s\"}", ach_details.Id, amount))
	res, err := makeAlpacaRequest("POST", url, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute request"})
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(res.Body)
		fmt.Printf("Transfer failed with status %d: %s\n", res.StatusCode, string(body))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transfer failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account_id": accountID,
		"amount":     amount,
		"message":    "Funds deposited successfully",
	})

}
