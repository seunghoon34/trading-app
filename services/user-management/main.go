// File: services/user-management/main.go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seunghoon34/trading-app/services/user-management/handlers"
)

func main() {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "health endpoint",
		})
	})
	r.POST("/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "login endpoint",
		})
	})
	r.POST("/register", handlers.Register)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
