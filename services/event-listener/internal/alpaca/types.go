package alpaca

import "time"

// AlpacaTradeEvent represents the event structure from Alpaca SSE
type AlpacaTradeEvent struct {
	AccountID string    `json:"account_id"`
	Event     string    `json:"event"` // "new", "fill", "canceled", etc.
	EventID   string    `json:"event_id"`
	Timestamp time.Time `json:"timestamp"`
	Order     Order     `json:"order"`
}

// Order represents the order details within the event
type Order struct {
	ID             string  `json:"id"`
	ClientOrderID  string  `json:"client_order_id"`
	Symbol         string  `json:"symbol"`
	Side           string  `json:"side"` // "buy" or "sell"
	Quantity       string  `json:"qty"`
	FilledQuantity string  `json:"filled_qty"`
	Status         string  `json:"status"`
	OrderType      string  `json:"order_type"`  // "market", "limit"
	LimitPrice     *string `json:"limit_price"` // Pointer because it can be null
	CreatedAt      string  `json:"created_at"`
	FilledAt       *string `json:"filled_at"` // Pointer because it can be null
	FilledAvgPrice *string `json:"filled_avg_price"`
}

// KafkaTradeEvent represents what we'll send to Kafka (simplified)
type KafkaTradeEvent struct {
	EventType     string    `json:"event_type"` // "ORDER_ACCEPTED", "ORDER_FILLED", etc.
	AccountID     string    `json:"account_id"`
	OrderID       string    `json:"order_id"`
	Symbol        string    `json:"symbol"`
	Side          string    `json:"side"`
	Quantity      string    `json:"quantity"`
	Price         string    `json:"price,omitempty"` // Only for fills
	Status        string    `json:"status"`
	Timestamp     time.Time `json:"timestamp"`
	OriginalEvent string    `json:"original_event"` // Store the raw Alpaca event type
}
