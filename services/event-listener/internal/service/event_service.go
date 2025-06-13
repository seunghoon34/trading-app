// internal/service/event_service.go
package service

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/seunghoon34/trading-app/services/event-listener/internal/alpaca"
	"github.com/seunghoon34/trading-app/services/event-listener/internal/kafka"
)

// EventService orchestrates the SSE client and Kafka producer for trade events
type EventService struct {
	sseClient     *alpaca.SSEClient
	kafkaProducer *kafka.Producer
	accountID     string
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

// NewEventService creates a new event service
func NewEventService() *EventService {
	ctx, cancel := context.WithCancel(context.Background())

	// Get account ID from environment
	accountID := os.Getenv("ALPACA_ACCOUNT_ID")
	if accountID == "" {
		log.Fatal("❌ ALPACA_ACCOUNT_ID environment variable is required")
	}

	// Create SSE client
	sseClient := alpaca.NewSSEClient()

	// Create Kafka producer (using your existing implementation)
	kafkaProducer := kafka.NewProducer("kafka:29092", "trade-events")

	return &EventService{
		sseClient:     sseClient,
		kafkaProducer: kafkaProducer,
		accountID:     accountID,
		ctx:           ctx,
		cancel:        cancel,
	}
}

// Start begins the event listening and processing
func (s *EventService) Start() error {
	log.Printf("🚀 Starting Trade Event Service for account: %s", s.accountID)

	// Connect to Alpaca SSE
	if err := s.sseClient.Connect(s.accountID); err != nil {
		return err
	}

	// Start event processing goroutine
	s.wg.Add(1)
	go s.processTradeEvents()

	// Start error handling goroutine
	s.wg.Add(1)
	go s.handleErrors()

	log.Printf("✅ Trade Event Service started successfully")
	return nil
}

// processTradeEvents handles incoming trade events from SSE and sends them to Kafka
func (s *EventService) processTradeEvents() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			log.Printf("🛑 Stopping trade event processing...")
			return

		case event := <-s.sseClient.EventChannel:
			// Process the trade event using your existing Kafka producer
			if err := s.kafkaProducer.PublishTradeEvent(event); err != nil {
				log.Printf("❌ Failed to publish trade event to Kafka: %v", err)
				// You might want to implement retry logic here
				continue
			}

		case <-time.After(30 * time.Second):
			// Heartbeat - log that we're still alive
			log.Printf("💓 Trade event service heartbeat - waiting for events...")
		}
	}
}

// handleErrors processes errors from the SSE client
func (s *EventService) handleErrors() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return

		case err := <-s.sseClient.ErrorChannel:
			log.Printf("⚠️ SSE Client Error: %v", err)

			// Implement reconnection logic
			s.handleReconnection(err)
		}
	}
}

// handleReconnection attempts to reconnect to the SSE stream
func (s *EventService) handleReconnection(err error) {
	log.Printf("🔄 Attempting to reconnect due to error: %v", err)

	// Wait before reconnecting
	time.Sleep(5 * time.Second)

	// Try to reconnect
	if reconnectErr := s.sseClient.Connect(s.accountID); reconnectErr != nil {
		log.Printf("❌ Reconnection failed: %v", reconnectErr)

		// Wait longer before next attempt
		time.Sleep(30 * time.Second)

		// You might want to implement exponential backoff here
		go s.handleReconnection(reconnectErr)
	} else {
		log.Printf("✅ Successfully reconnected to trade events SSE stream")
	}
}

// GetStats returns statistics about the service
func (s *EventService) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"event_channel_size": len(s.sseClient.EventChannel),
		"error_channel_size": len(s.sseClient.ErrorChannel),
		"account_id":         s.accountID,
		"service_type":       "trade-events",
		"status":             "running",
		"timestamp":          time.Now(),
	}
}

// Stop gracefully shuts down the event service
func (s *EventService) Stop() {
	log.Printf("🛑 Stopping Trade Event Service...")

	// Cancel context to stop all goroutines
	s.cancel()

	// Close SSE client
	s.sseClient.Close()

	// Close Kafka producer (using your existing method)
	if err := s.kafkaProducer.Close(); err != nil {
		log.Printf("⚠️ Error closing Kafka producer: %v", err)
	}

	// Wait for all goroutines to finish
	s.wg.Wait()

	log.Printf("✅ Trade Event Service stopped")
}
