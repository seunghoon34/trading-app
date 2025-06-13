package kafka

import (
	"testing"
	"time"

	"github.com/seunghoon34/trading-app/services/event-listener/internal/alpaca"
)

func TestEventTypeMapping(t *testing.T) {
	producer := &Producer{}

	// Test our event type mapping
	tests := []struct {
		alpacaEvent string
		expected    string
	}{
		{"accepted", "ORDER_ACCEPTED"},
		{"fill", "ORDER_FILLED"},
		{"canceled", "ORDER_CANCELED"},
		{"unknown", "ORDER_unknown"},
	}

	for _, test := range tests {
		result := producer.mapEventType(test.alpacaEvent)
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}

	t.Logf("✅ Event type mapping working correctly")
}

func TestKafkaEventConversion(t *testing.T) {
	producer := &Producer{}

	// Create sample Alpaca event
	alpacaEvent := alpaca.AlpacaTradeEvent{
		AccountID: "test-account-123",
		Event:     "fill",
		Timestamp: time.Now(),
		Order: alpaca.Order{
			ID:       "order-123",
			Symbol:   "AAPL",
			Side:     "buy",
			Quantity: "10",
			Status:   "filled",
		},
	}

	// Convert to Kafka event
	kafkaEvent := producer.convertToKafkaEvent(alpacaEvent)

	// Test conversion
	if kafkaEvent.EventType != "ORDER_FILLED" {
		t.Errorf("Expected ORDER_FILLED, got %s", kafkaEvent.EventType)
	}

	if kafkaEvent.Symbol != "AAPL" {
		t.Errorf("Expected AAPL, got %s", kafkaEvent.Symbol)
	}

	t.Logf("✅ Kafka event conversion working correctly")
}
