// File: services/market-data/main.go
package main

import (
	"net/http"

	"github.com/seunghoon34/trading-app/services/market-data/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "health endpoint",
		})
	})
	r.GET("/quotes/:symbol", func(c *gin.Context) {
		handlers.GetLatestQuote(c)
	})
	r.GET("/bars/:symbol", func(c *gin.Context) {
		handlers.GetLatestBar(c)
	})

	r.GET("/quotes", func(c *gin.Context) {
		handlers.GetLatestQuoteMulti(c)
	})

	r.GET("/bars", func(c *gin.Context) {
		handlers.GetLatestBarMulti(c)
	})

	r.Run(":8082") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
