package mongo

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var PortfolioCollection *mongo.Collection
var RiskProfileCollection *mongo.Collection

// initMongoDB initializes the MongoDB connection and creates indexes

func InitMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get MongoDB credentials from environment variables
	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")

	var uri string
	if mongoUser != "" && mongoPassword != "" {
		// Connect to admin database since these are root credentials
		uri = fmt.Sprintf("mongodb://%s:%s@localhost:27017/admin", mongoUser, mongoPassword)
	} else {
		uri = "mongodb://localhost:27017"
	}

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	MongoClient = client
	PortfolioCollection = client.Database("trading").Collection("portfolios")
	RiskProfileCollection = client.Database("trading").Collection("risk_profile")

	// Create unique index on alpaca_id for fast lookups and prevent duplicates
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "alpaca_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err = PortfolioCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("Warning: Failed to create index: %v", err)
	}

	riskProfileIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "alpaca_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err = RiskProfileCollection.Indexes().CreateOne(ctx, riskProfileIndex)
	if err != nil {
		log.Printf("Warning: Failed to create risk profile index: %v", err)
	}

	log.Println("Connected to MongoDB and created indexes!")
}

func DisconnectMongoDB() error {
	if MongoClient == nil {
		return nil
	}
	return MongoClient.Disconnect(context.Background())
}
