// services/event-listener/internal/alpaca/types_test.go
package alpaca

import (
	"encoding/json"
	"testing"
)

func TestAlpacaEventParsing(t *testing.T) {
	// Sample JSON from Alpaca documentation
	sampleJSON := `{
		"account_id": "529248ad-c4cc-4a50-bea4-6bfd2953f83a",
		"event": "new",
		"order": {
			"id": "edada91a-8b55-4916-a153-8c7a9817e708",
			"symbol": "TSLA",
			"side": "buy",
			"qty": "4",
			"status": "new"
		}
	}`

	var event AlpacaTradeEvent
	err := json.Unmarshal([]byte(sampleJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Test that we parsed correctly
	if event.AccountID != "529248ad-c4cc-4a50-bea4-6bfd2953f83a" {
		t.Errorf("Expected account ID to match")
	}

	if event.Order.Symbol != "TSLA" {
		t.Errorf("Expected symbol to be TSLA, got %s", event.Order.Symbol)
	}

	t.Logf("✅ Successfully parsed event: %s for %s", event.Event, event.Order.Symbol)
}
