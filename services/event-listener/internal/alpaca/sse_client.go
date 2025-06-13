// internal/alpaca/sse_client.go
package alpaca

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// SSEClient handles the connection to Alpaca's trade events SSE endpoint
type SSEClient struct {
	apiKey    string
	secretKey string
	baseURL   string
	client    *http.Client
	// Channel to send parsed trade events
	EventChannel chan AlpacaTradeEvent
	// Channel to handle errors
	ErrorChannel chan error
	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc
}

// NewSSEClient creates a new SSE client for trade events
func NewSSEClient() *SSEClient {
	ctx, cancel := context.WithCancel(context.Background())

	return &SSEClient{
		apiKey:       os.Getenv("ALPACA_API_KEY"),
		secretKey:    os.Getenv("ALPACA_SECRET_KEY"),
		baseURL:      "https://broker-api.sandbox.alpaca.markets", // Use sandbox for testing
		client:       &http.Client{Timeout: 0},                    // No timeout for SSE connections
		EventChannel: make(chan AlpacaTradeEvent, 100),            // Buffered channel
		ErrorChannel: make(chan error, 10),
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Connect establishes connection to Alpaca trade events SSE endpoint
func (s *SSEClient) Connect(accountID string) error {
	// Build the SSE URL for trade events
	url := "https://broker-api.sandbox.alpaca.markets/v2/events/trades"

	log.Printf("🔌 Connecting to Alpaca Trade Events SSE: %s", url)

	// Create request
	req, err := http.NewRequestWithContext(s.ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication (same pattern as your other services)
	auth := s.apiKey + ":" + s.secretKey
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", basicAuth)
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	// Make the request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to SSE: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("SSE connection failed with status: %d", resp.StatusCode)
	}

	log.Printf("✅ Connected to Alpaca Trade Events SSE successfully")

	// Start reading the stream in a goroutine
	go s.readStream(resp.Body)

	return nil
}

// readStream reads and processes the SSE stream
func (s *SSEClient) readStream(body io.ReadCloser) {
	defer body.Close()

	scanner := bufio.NewScanner(body)
	var eventData strings.Builder

	for scanner.Scan() {
		select {
		case <-s.ctx.Done():
			log.Printf("🛑 SSE stream reading stopped")
			return
		default:
		}

		line := scanner.Text()

		// SSE format parsing
		if strings.HasPrefix(line, "data: ") {
			// Extract the JSON data
			data := strings.TrimPrefix(line, "data: ")
			eventData.WriteString(data)
		} else if line == "" {
			// Empty line indicates end of event
			if eventData.Len() > 0 {
				s.processEvent(eventData.String())
				eventData.Reset()
			}
		}
		// Ignore other SSE fields like "event:", "id:", etc.
	}

	if err := scanner.Err(); err != nil {
		select {
		case s.ErrorChannel <- fmt.Errorf("error reading SSE stream: %w", err):
		case <-s.ctx.Done():
		}
	}
}

// processEvent parses and sends the trade event to the channel
func (s *SSEClient) processEvent(data string) {
	// Skip heartbeat or empty events
	if strings.TrimSpace(data) == "" || data == "heartbeat" {
		return
	}

	log.Printf("📨 Received trade event data: %s", data)

	var event AlpacaTradeEvent
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		select {
		case s.ErrorChannel <- fmt.Errorf("failed to parse trade event: %w", err):
		case <-s.ctx.Done():
		}
		return
	}

	// Send the parsed event to the channel
	select {
	case s.EventChannel <- event:
		log.Printf("✅ Processed %s trade event for order %s (%s)",
			event.Event, event.Order.ID, event.Order.Symbol)
	case <-s.ctx.Done():
		return
	default:
		// Channel is full, log warning
		log.Printf("⚠️ Event channel is full, dropping trade event")
	}
}

// Close gracefully shuts down the SSE client
func (s *SSEClient) Close() {
	log.Printf("🛑 Closing SSE client...")
	s.cancel()
	close(s.EventChannel)
	close(s.ErrorChannel)
}
