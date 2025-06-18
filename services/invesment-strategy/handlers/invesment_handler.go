package handlers

// services/invesment-strategy/handlers/invesment_handler.go
import (
	"context"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seunghoon34/trading-app/services/investment-strategy/internal/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo" // Alias to avoid naming conflict
)

// Position represents a single stock position
type Position struct {
	Symbol string  `json:"symbol" bson:"symbol" binding:"required"`
	Weight float64 `json:"weight" bson:"weight" binding:"required,min=0,max=1"`
}

// Portfolio represents a user's portfolio
type Portfolio struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AccountID string             `json:"account_id" bson:"alpaca_id"`
	Positions []Position         `json:"positions" bson:"positions"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// PortfolioRequest represents the request body for creating/updating portfolios
type PortfolioRequest struct {
	Positions []Position `json:"positions" binding:"required,dive"`
}

// createPortfolio creates a new portfolio (prevents duplicates)
func CreatePortfolio(c *gin.Context) {
	// Get account_id from header
	accountID := c.GetHeader("X-Account-ID")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Account-ID header is required"})
		return
	}

	var req PortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate weights sum to 1.0
	totalWeight := 0.0
	for _, pos := range req.Positions {
		totalWeight += pos.Weight
	}

	if math.Abs(totalWeight-1.0) > 0.001 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "weights must sum to 1.0"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if portfolio already exists
	var existingPortfolio Portfolio
	err := mongo.PortfolioCollection.FindOne(ctx, bson.M{"alpaca_id": accountID}).Decode(&existingPortfolio)

	if err == nil {
		// Portfolio exists
		c.JSON(http.StatusConflict, gin.H{
			"error":                 "Portfolio already exists for this user",
			"existing_portfolio_id": existingPortfolio.ID,
			"message":               "Use update endpoint to modify existing portfolio",
		})
		return
	} else if err != mongodriver.ErrNoDocuments {
		// Some other error occurred
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing portfolio " + err.Error()})
		return
	}

	// Create new portfolio (no existing portfolio found)
	portfolio := Portfolio{
		AccountID: accountID,
		Positions: req.Positions,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := mongo.PortfolioCollection.InsertOne(ctx, portfolio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create portfolio"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "portfolio created successfully",
		"account_id":   accountID,
		"portfolio_id": result.InsertedID,
	})
}

// updatePortfolio updates an existing portfolio
func UpdatePortfolio(c *gin.Context) {
	accountID := c.GetHeader("X-Account-ID")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Account-ID header is required"})
		return
	}

	var req PortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate weights sum to 1.0
	totalWeight := 0.0
	for _, pos := range req.Positions {
		totalWeight += pos.Weight
	}

	if math.Abs(totalWeight-1.0) > 0.001 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "weights must sum to 1.0"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Update portfolio (keep original created_at, update updated_at)
	update := bson.M{
		"$set": bson.M{
			"positions":  req.Positions,
			"updated_at": time.Now(),
		},
	}

	result, err := mongo.PortfolioCollection.UpdateOne(
		ctx,
		bson.M{"alpaca_id": accountID},
		update,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update portfolio"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "portfolio updated successfully",
		"account_id": accountID,
	})
}

// getPortfolio retrieves a portfolio by account_id
func GetPortfolio(c *gin.Context) {
	accountID := c.GetHeader("X-Account-ID")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Account-ID header is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var portfolio Portfolio
	err := mongo.PortfolioCollection.FindOne(ctx, bson.M{"alpaca_id": accountID}).Decode(&portfolio)

	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch portfolio"})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

// getAllPortfolios retrieves all portfolios (admin function)
func GetAllPortfolios(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := mongo.PortfolioCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch portfolios"})
		return
	}
	defer cursor.Close(ctx)

	var portfolios []Portfolio
	if err = cursor.All(ctx, &portfolios); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode portfolios"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"portfolios": portfolios,
		"count":      len(portfolios),
	})
}

// deletePortfolio deletes a portfolio by account_id
func DeletePortfolio(c *gin.Context) {
	accountID := c.GetHeader("X-Account-ID")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Account-ID header is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := mongo.PortfolioCollection.DeleteOne(ctx, bson.M{"alpaca_id": accountID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete portfolio"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "portfolio deleted successfully",
		"account_id": accountID,
	})
}
