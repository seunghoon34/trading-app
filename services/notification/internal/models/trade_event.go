// internal/models/trade_event.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TradeEvent represents a trade event stored in MongoDB
type TradeEvent struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	EventType     string             `bson:"event_type" json:"event_type"`
	AccountID     string             `bson:"account_id" json:"account_id"`
	OrderID       string             `bson:"order_id" json:"order_id"`
	Symbol        string             `bson:"symbol" json:"symbol"`
	Side          string             `bson:"side" json:"side"`
	Quantity      string             `bson:"quantity" json:"quantity"`
	Price         string             `bson:"price,omitempty" json:"price,omitempty"`
	Status        string             `bson:"status" json:"status"`
	Timestamp     time.Time          `bson:"timestamp" json:"timestamp"`
	OriginalEvent string             `bson:"original_event" json:"original_event"`
	// MongoDB metadata
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	ProcessedAt time.Time `bson:"processed_at" json:"processed_at"`
}

// KafkaTradeEvent represents the event structure from Kafka (matches your producer)
type KafkaTradeEvent struct {
	EventType     string    `json:"event_type"`
	AccountID     string    `json:"account_id"`
	OrderID       string    `json:"order_id"`
	Symbol        string    `json:"symbol"`
	Side          string    `json:"side"`
	Quantity      string    `json:"quantity"`
	Price         string    `json:"price,omitempty"`
	Status        string    `json:"status"`
	Timestamp     time.Time `json:"timestamp"`
	OriginalEvent string    `json:"original_event"`
}

// ToTradeEvent converts KafkaTradeEvent to MongoDB TradeEvent
func (k *KafkaTradeEvent) ToTradeEvent() *TradeEvent {
	now := time.Now()
	return &TradeEvent{
		EventType:     k.EventType,
		AccountID:     k.AccountID,
		OrderID:       k.OrderID,
		Symbol:        k.Symbol,
		Side:          k.Side,
		Quantity:      k.Quantity,
		Price:         k.Price,
		Status:        k.Status,
		Timestamp:     k.Timestamp,
		OriginalEvent: k.OriginalEvent,
		CreatedAt:     k.Timestamp,
		ProcessedAt:   now,
	}
}
