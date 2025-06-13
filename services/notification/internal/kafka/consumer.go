// internal/kafka/consumer.go
package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/seunghoon34/trading-app/services/notification/internal/models"
)

// Consumer handles consuming messages from Kafka
type Consumer struct {
	reader *kafka.Reader
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokerAddress, topic, consumerGroup string) *Consumer {
	// Create Kafka reader (consumer)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokerAddress},
		Topic:          topic,
		GroupID:        consumerGroup,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset, // Start from latest messages
	})

	log.Printf("📥 Kafka consumer created for topic: %s (group: %s)", topic, consumerGroup)

	return &Consumer{
		reader: reader,
	}
}

// ConsumeMessages starts consuming messages and calls the handler for each message
func (c *Consumer) ConsumeMessages(ctx context.Context, handler func(*models.KafkaTradeEvent) error) error {
	log.Printf("🎧 Starting to consume trade events...")

	for {
		select {
		case <-ctx.Done():
			log.Printf("🛑 Stopping Kafka consumer...")
			return ctx.Err()
		default:
		}

		// Read message with timeout
		message, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("❌ Error reading message: %v", err)
			continue
		}

		// Parse the trade event
		var tradeEvent models.KafkaTradeEvent
		if err := json.Unmarshal(message.Value, &tradeEvent); err != nil {
			log.Printf("❌ Error parsing trade event: %v", err)
			// Commit even if parsing fails to avoid getting stuck
			c.reader.CommitMessages(ctx, message)
			continue
		}

		// Log the received event
		log.Printf("📨 Received trade event: %s for %s (Order: %s)",
			tradeEvent.EventType, tradeEvent.Symbol, tradeEvent.OrderID)

		// Handle the event
		if err := handler(&tradeEvent); err != nil {
			log.Printf("❌ Error handling trade event: %v", err)
			// You might want to implement retry logic here
			// For now, we'll commit anyway to avoid getting stuck
		}

		// Commit the message
		if err := c.reader.CommitMessages(ctx, message); err != nil {
			log.Printf("❌ Error committing message: %v", err)
		} else {
			log.Printf("✅ Successfully processed and committed trade event")
		}
	}
}

// Close gracefully closes the Kafka consumer
func (c *Consumer) Close() error {
	log.Printf("🛑 Closing Kafka consumer...")
	return c.reader.Close()
}
