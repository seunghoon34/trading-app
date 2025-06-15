package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Investment-strategy health endpoint",
		})
	})
	r.Run(":8089") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
