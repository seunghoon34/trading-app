package alpaca

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type AlpacaAccount struct {
	Id string `json:"id"`
}

type AlpacaError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func generateUniqueTaxId() string {
	// Generate a unique tax ID for sandbox (format: XXX-XX-XXXX)
	return fmt.Sprintf("%03d-%02d-%04d",
		rand.Intn(900)+100,   // 100-999
		rand.Intn(90)+10,     // 10-99
		rand.Intn(9000)+1000) // 1000-9999
}

func generateUniquePhone() string {
	// Generate a unique phone number for sandbox
	return fmt.Sprintf("+1555%07d", rand.Intn(10000000))
}

func CreateAlpacaAccount(email string, firstName string, lastName string) (*AlpacaAccount, error) {
	url := "https://broker-api.sandbox.alpaca.markets/v1/accounts"

	// Generate unique identifiers for sandbox
	taxId := generateUniqueTaxId()
	phoneNumber := generateUniquePhone()

	// Create payload with unique data
	payloadData := map[string]interface{}{
		"contact": map[string]interface{}{
			"email_address":  email,
			"phone_number":   phoneNumber,
			"city":           "San Mateo",
			"postal_code":    "94401",
			"street_address": []string{"20 N San Mateo Dr"},
			"country":        "USA",
			"state":          "CA",
		},
		"identity": map[string]interface{}{
			"tax_id_type":              "NOT_SPECIFIED",
			"given_name":               firstName,
			"family_name":              lastName,
			"date_of_birth":            "1990-01-01",
			"tax_id":                   taxId,
			"funding_source":           []string{"employment_income"},
			"country_of_tax_residence": "USA",
			"country_of_citizenship":   "USA",
			"country_of_birth":         "USA",
		},
		"disclosures": map[string]interface{}{
			"is_control_person":               false,
			"is_affiliated_exchange_or_finra": false,
			"is_politically_exposed":          false,
			"immediate_family_exposed":        false,
		},
		"agreements": []map[string]interface{}{
			{
				"agreement":  "account_agreement",
				"signed_at":  "2019-09-11T18:09:33Z",
				"ip_address": "185.13.21.99",
			},
			{
				"agreement":  "margin_agreement",
				"signed_at":  "2019-09-11T18:09:33Z",
				"ip_address": "185.13.21.99",
			},
			{
				"agreement":  "customer_agreement",
				"signed_at":  "2019-09-11T18:09:33Z",
				"ip_address": "185.13.21.99",
			},
		},
		"account_type": "trading",
	}

	// Convert to JSON
	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Basic Q0tKSzM2UzNFRTZPTDRBWTU1SDE6THJkYXFpaUQyZzRGeHV4cG1rOE1yWmxibE5qYTJVOU5HR0lsQ0tGMg==")

	// Execute request
	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer res.Body.Close()

	// Read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if res.StatusCode >= 400 {
		var alpacaErr AlpacaError
		if err := json.Unmarshal(body, &alpacaErr); err == nil {
			return nil, fmt.Errorf("alpaca API error (%d): %s", res.StatusCode, alpacaErr.Message)
		}
		return nil, fmt.Errorf("alpaca API error (%d): %s", res.StatusCode, string(body))
	}

	// Parse successful response
	var account AlpacaAccount
	if err := json.Unmarshal(body, &account); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w, response: %s", err, string(body))
	}

	return &account, nil
}
