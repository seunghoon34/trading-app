package config

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
)

func MakeMarketDataRequest(endpoint string) ([]byte, error) {
	// Build full URL
	baseURL := os.Getenv("ALPACA_MARKET_DATA_URL")
	url := fmt.Sprintf("%s%s", baseURL, endpoint)

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add Basic Auth (same as your working curl command)
	auth := os.Getenv("ALPACA_API_KEY") + ":" + os.Getenv("ALPACA_SECRET_KEY")
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", basicAuth)
	req.Header.Add("Accept", "application/json")

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
