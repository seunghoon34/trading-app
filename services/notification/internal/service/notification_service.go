// internal/service/notification_service.go
package service

import (
	"context"
	"log"
	"sync"

	"github.com/seunghoon34/trading-app/services/notification/internal/kafka"
	"github.com/seunghoon34/trading-app/services/notification/internal/models"
	"github.com/seunghoon34/trading-app/services/notification/internal/mongodb"
)

// NotificationService handles consuming trade events and storing them
type NotificationService struct {
	kafkaConsumer *kafka.Consumer
	MongoClient   *mongodb.Client
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

// NewNotificationService creates a new notification service
func NewNotificationService() (*NotificationService, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create Kafka consumer
	kafkaConsumer := kafka.NewConsumer(
		"kafka:29092",          // Broker address
		"trade-events",         // Topic
		"notification-service", // Consumer group
	)

	// Create MongoDB client
	mongoClient, err := mongodb.NewClient()
	if err != nil {
		cancel()
		return nil, err
	}

	return &NotificationService{
		kafkaConsumer: kafkaConsumer,
		MongoClient:   mongoClient,
		ctx:           ctx,
		cancel:        cancel,
	}, nil
}

// Start begins consuming trade events and processing them
func (n *NotificationService) Start() error {
	log.Printf("🚀 Starting Notification Service...")

	// Start consuming messages in a goroutine
	n.wg.Add(1)
	go func() {
		defer n.wg.Done()

		if err := n.kafkaConsumer.ConsumeMessages(n.ctx, n.handleTradeEvent); err != nil {
			if err != context.Canceled {
				log.Printf("❌ Error consuming messages: %v", err)
			}
		}
	}()

	log.Printf("✅ Notification Service started successfully")
	return nil
}

// handleTradeEvent processes a single trade event
func (n *NotificationService) handleTradeEvent(kafkaEvent *models.KafkaTradeEvent) error {
	// Log the event processing
	log.Printf("🔔 Processing notification for %s: %s %s %s (Order: %s)",
		kafkaEvent.EventType,
		kafkaEvent.Side,
		kafkaEvent.Quantity,
		kafkaEvent.Symbol,
		kafkaEvent.OrderID)

	// Convert to MongoDB model
	tradeEvent := kafkaEvent.ToTradeEvent()

	// Store in MongoDB
	if err := n.MongoClient.StoreTradeEvent(tradeEvent); err != nil {
		log.Printf("❌ Failed to store trade event in MongoDB: %v", err)
		return err
	}

	// Log detailed notification
	n.logNotification(kafkaEvent)

	return nil
}

// logNotification creates detailed log notifications
func (n *NotificationService) logNotification(event *models.KafkaTradeEvent) {
	switch event.EventType {
	case "ORDER_ACCEPTED":
		log.Printf("🎯 NOTIFICATION: Order accepted - %s %s %s",
			event.Side, event.Quantity, event.Symbol)

	case "ORDER_NEW":
		log.Printf("📝 NOTIFICATION: New order placed - %s %s %s (Order: %s)",
			event.Side, event.Quantity, event.Symbol, event.OrderID)

	case "ORDER_FILLED":
		price := event.Price
		if price == "" {
			price = "market price"
		}
		log.Printf("💰 NOTIFICATION: Order filled - %s %s %s at %s (Order: %s)",
			event.Side, event.Quantity, event.Symbol, price, event.OrderID)

	case "ORDER_PARTIAL_FILLED":
		log.Printf("📊 NOTIFICATION: Order partially filled - %s %s %s (Order: %s)",
			event.Side, event.Quantity, event.Symbol, event.OrderID)

	case "ORDER_CANCELED":
		log.Printf("❌ NOTIFICATION: Order canceled - %s %s %s (Order: %s)",
			event.Side, event.Quantity, event.Symbol, event.OrderID)

	case "ORDER_REJECTED":
		log.Printf("🚫 NOTIFICATION: Order rejected - %s %s %s (Order: %s)",
			event.Side, event.Quantity, event.Symbol, event.OrderID)

	default:
		log.Printf("📢 NOTIFICATION: %s - %s %s %s (Order: %s)",
			event.EventType, event.Side, event.Quantity, event.Symbol, event.OrderID)
	}
}

// GetStats returns statistics about the notification service
func (n *NotificationService) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"service":   "notification-service",
		"status":    "running",
		"consuming": "trade-events",
		"storing":   "MongoDB",
		"group_id":  "notification-service",
	}
}

// Stop gracefully shuts down the notification service
func (n *NotificationService) Stop() {
	log.Printf("🛑 Stopping Notification Service...")

	// Cancel context to stop consuming
	n.cancel()

	// Wait for goroutines to finish
	n.wg.Wait()

	// Close Kafka consumer
	if err := n.kafkaConsumer.Close(); err != nil {
		log.Printf("⚠️ Error closing Kafka consumer: %v", err)
	}

	// Close MongoDB client
	if err := n.MongoClient.Close(); err != nil {
		log.Printf("⚠️ Error closing MongoDB client: %v", err)
	}

	log.Printf("✅ Notification Service stopped")
}
