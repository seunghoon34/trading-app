// internal/mongodb/client.go
package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/seunghoon34/trading-app/services/notification/internal/models"
)

// Client handles MongoDB operations
type Client struct {
	client     *mongo.Client
	database   *mongo.Database
	Collection *mongo.Collection
}

// NewClient creates a new MongoDB client
func NewClient() (*Client, error) {
	// Get MongoDB connection string from environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://mongodb:27017" // Default for Docker
	}

	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "trading_notifications" // Default database name
	}

	log.Printf("🔌 Connecting to MongoDB: %s", mongoURI)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test the connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Printf("✅ Connected to MongoDB successfully")

	// Get database and collection
	database := client.Database(dbName)
	collection := database.Collection("trade_events")

	return &Client{
		client:     client,
		database:   database,
		Collection: collection,
	}, nil
}

// StoreTradeEvent stores a trade event in MongoDB
func (c *Client) StoreTradeEvent(event *models.TradeEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert the event
	result, err := c.Collection.InsertOne(ctx, event)
	if err != nil {
		return fmt.Errorf("failed to store trade event: %w", err)
	}

	log.Printf("✅ Stored trade event in MongoDB: %s for %s (Order: %s)",
		event.EventType, event.Symbol, event.OrderID)
	log.Printf("   📄 MongoDB Document ID: %s", result.InsertedID)

	return nil
}

// GetTradeEventsByAccount retrieves trade events for a specific account
func (c *Client) GetTradeEventsByAccount(accountID string, limit int64) ([]*models.TradeEvent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create filter
	filter := map[string]interface{}{"account_id": accountID}

	// Create options with limit and sort by timestamp descending
	opts := options.Find().SetLimit(limit).SetSort(map[string]interface{}{"timestamp": -1})

	// Find events
	cursor, err := c.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find trade events: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode results
	var events []*models.TradeEvent
	if err := cursor.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("failed to decode trade events: %w", err)
	}

	return events, nil
}

// GetTradeEventsByOrder retrieves trade events for a specific order
func (c *Client) GetTradeEventsByOrder(orderID string) ([]*models.TradeEvent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create filter
	filter := map[string]interface{}{"order_id": orderID}

	// Create options to sort by timestamp ascending (order lifecycle)
	opts := options.Find().SetSort(map[string]interface{}{"timestamp": 1})

	// Find events
	cursor, err := c.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find trade events for order: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode results
	var events []*models.TradeEvent
	if err := cursor.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("failed to decode trade events: %w", err)
	}

	return events, nil
}

// Close gracefully closes the MongoDB connection
func (c *Client) Close() error {
	log.Printf("🛑 Closing MongoDB connection...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.client.Disconnect(ctx)
}
