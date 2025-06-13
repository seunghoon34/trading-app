package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/seunghoon34/trading-app/services/event-listener/internal/alpaca"
)

// Producer handles sending events to Kafka
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a new Kafka producer
func NewProducer(brokerAddress, topic string) *Producer {
	// Create Kafka writer (producer)
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokerAddress), // Kafka server address
		Topic:        topic,                    // Which topic to write to
		Balancer:     &kafka.LeastBytes{},      // How to distribute messages across partitions
		WriteTimeout: 10 * time.Second,         // Timeout for writes
		ReadTimeout:  10 * time.Second,         // Timeout for reads
	}

	log.Printf("📤 Kafka producer created for topic: %s", topic)

	return &Producer{
		writer: writer,
	}
}

// PublishTradeEvent converts Alpaca event to Kafka event and publishes it
func (p *Producer) PublishTradeEvent(alpacaEvent alpaca.AlpacaTradeEvent) error {
	// Convert Alpaca event to our Kafka format
	kafkaEvent := p.convertToKafkaEvent(alpacaEvent)

	// Convert to JSON
	eventJSON, err := json.Marshal(kafkaEvent)
	if err != nil {
		log.Printf("❌ Failed to marshal event to JSON: %v", err)
		return err
	}

	// Create Kafka message
	message := kafka.Message{
		Key:   []byte(kafkaEvent.AccountID), // Use account_id as key for partitioning
		Value: eventJSON,                    // The actual event data
		Time:  time.Now(),                   // When we're sending it
	}

	// Send to Kafka
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = p.writer.WriteMessages(ctx, message)
	if err != nil {
		log.Printf("❌ Failed to write message to Kafka: %v", err)
		return err
	}

	log.Printf("✅ Published %s event for %s (Order: %s)",
		kafkaEvent.EventType, kafkaEvent.Symbol, kafkaEvent.OrderID)

	return nil
}

// convertToKafkaEvent transforms Alpaca format to our internal format
func (p *Producer) convertToKafkaEvent(alpacaEvent alpaca.AlpacaTradeEvent) alpaca.KafkaTradeEvent {
	// Map Alpaca event types to our standardized types
	eventType := p.mapEventType(alpacaEvent.Event)

	// Extract price information (only available for fills)
	price := ""
	if alpacaEvent.Order.FilledAvgPrice != nil {
		price = *alpacaEvent.Order.FilledAvgPrice
	}

	return alpaca.KafkaTradeEvent{
		EventType:     eventType,
		AccountID:     alpacaEvent.AccountID,
		OrderID:       alpacaEvent.Order.ID,
		Symbol:        alpacaEvent.Order.Symbol,
		Side:          alpacaEvent.Order.Side,
		Quantity:      alpacaEvent.Order.Quantity,
		Price:         price,
		Status:        alpacaEvent.Order.Status,
		Timestamp:     alpacaEvent.Timestamp,
		OriginalEvent: alpacaEvent.Event,
	}
}

// mapEventType converts Alpaca event types to our standardized naming
func (p *Producer) mapEventType(alpacaEvent string) string {
	switch alpacaEvent {
	case "accepted":
		return "ORDER_ACCEPTED"
	case "new":
		return "ORDER_NEW"
	case "fill":
		return "ORDER_FILLED"
	case "partial_fill":
		return "ORDER_PARTIAL_FILLED"
	case "canceled":
		return "ORDER_CANCELED"
	case "rejected":
		return "ORDER_REJECTED"
	case "expired":
		return "ORDER_EXPIRED"
	default:
		return "ORDER_" + alpacaEvent // Fallback: ORDER_whatever
	}
}

// Close gracefully shuts down the producer
func (p *Producer) Close() error {
	log.Printf("🛑 Closing Kafka producer...")
	return p.writer.Close()
}
