// Investment Strategy Service main.go
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seunghoon34/trading-app/services/investment-strategy/handlers"
	"github.com/seunghoon34/trading-app/services/investment-strategy/internal/mongo"
)

func main() {

	mongo.InitMongoDB()
	defer func() {
		if err := mongo.DisconnectMongoDB(); err != nil {
			log.Printf("Warning: Failed to disconnect from MongoDB: %v", err) // Changed from log.Fatal
		}
	}()
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Investment-strategy health endpoint",
		})
	})

	// Portfolio routes (API gateway handles auth and sets X-Alpaca-ID header)
	r.POST("/portfolio", handlers.CreatePortfolio)   // Create new portfolio
	r.PUT("/portfolio", handlers.UpdatePortfolio)    // Update existing portfolio
	r.GET("/portfolio", handlers.GetPortfolio)       // Get portfolio by alpaca_id
	r.DELETE("/portfolio", handlers.DeletePortfolio) // Delete portfolio

	r.POST("/risk-profile", handlers.CreateRiskProfile) // Create new risk profile
	r.PUT("/risk-profile", handlers.UpdateRiskProfile)  // Update existing risk profile
	r.GET("/risk-profile", handlers.GetRiskProfile)     // Get risk profile by alpaca_id

	// Admin route
	r.GET("/api/portfolios", handlers.GetAllPortfolios) // Get all portfolios

	r.POST("/portfolio/purchase", handlers.PurchasePortfolio)

	log.Println("Investment Strategy Service starting on :8089")
	r.Run(":8089") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
