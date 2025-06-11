// File: services/trading engine/main.go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seunghoon34/trading-app/services/trading-engine/handlers"
)

func main() {

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "health endpoint",
		})
	})
	r.POST("/orders", handlers.CreateOrder)

	r.DELETE("/orders", handlers.DeleteAllOrders)

	r.GET("/:order_id", handlers.GetOrder)

	r.GET("/orders", handlers.GetOrders)

	r.DELETE("/:order_id", handlers.DeleteOrder)

	r.Run(":8083") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
