// File: services/user-management/main.go
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seunghoon34/trading-app/services/user-management/config"
	"github.com/seunghoon34/trading-app/services/user-management/handlers"
	"github.com/seunghoon34/trading-app/services/user-management/middleware"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := config.CreateUsersTable(db); err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "health endpoint",
		})
	})
	r.POST("/login", func(c *gin.Context) {
		handlers.Login(c, db)
	})
	r.POST("/register", func(c *gin.Context) {
		handlers.Register(c, db)
	})
	r.POST("/logout", func(c *gin.Context) {
		handlers.Logout(c)
	})

	// Protected route
	r.GET("/me", middleware.JWTMiddleware(), func(c *gin.Context) {
		handlers.GetCurrentUser(c, db)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
