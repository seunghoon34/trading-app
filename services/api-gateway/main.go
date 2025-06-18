// File: services/api-gateway/main.go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seunghoon34/trading-app/services/api-gateway/handlers"
	"github.com/seunghoon34/trading-app/services/api-gateway/middleware"
)

func main() {

	r := gin.Default()
	r.Use(middleware.ErrorHandlingMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestLoggingMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	// routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "health endpoint",
		})
	})
	r.Any("/api/v1/auth/*path", handlers.ForwardToAuthService)

	// Protected routes (JWT required)
	protected := r.Group("/api/v1")
	protected.Use(middleware.JWTMiddleware()) // ← Add JWT middleware
	{
		protected.Any("/user/*path", handlers.ForwardToAuthService)
		protected.Any("/market/*path", handlers.ForwardToMarketDataService)
		protected.Any("/trading/*path", handlers.ForwardToTradingService)
		protected.Any("/portfolio/*path", handlers.ForwardToPortfolioService)
		protected.Any("/investment-strategy/*path", handlers.ForwardToInvestmentStrategy)
		protected.Any("/crewai-portfolio/*path", handlers.ForwardToCrewAIPortfolio)
		protected.Any("/payment/*path", handlers.ForwardToPaymentService)
	}

	r.Run(":3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
