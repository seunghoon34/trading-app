// File: services/portfolio/main.go
package main

import (
	"net/http"

	"github.com/seunghoon34/trading-app/services/portfolio/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Portfolio health endpoint",
		})
	})
	r.GET("/positions", handlers.GetPositions)
	r.GET("/value", handlers.GetPortfolioWorth)

	r.GET("/performance", handlers.GetPortfolioPerformance)

	r.GET("/:symbol", handlers.GetPosition)

	r.Run(":8084") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
