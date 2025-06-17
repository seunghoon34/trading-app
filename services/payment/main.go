// File: services/payment/main.go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seunghoon34/trading-app/services/payment/handlers"
	"github.com/seunghoon34/trading-app/services/payment/redis"
)

func main() {

	redis.Init()
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Payment health endpoint",
		})
	})

	r.POST("/deposit/:account_id/:amount", handlers.DepositFunds)

	r.Run(":8090") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
