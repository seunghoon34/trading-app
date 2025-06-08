// File: services/api-gateway/main.go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seunghoon34/trading-app/services/api-gateway/handlers"
)

func main() {

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "health endpoint",
		})
	})
	r.Any("/api/v1/auth/*path", handlers.ForwardToAuthService)
	r.Any("/api/v1/market/*path", handlers.ForwardToMarketDataService)

	r.Run(":3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
