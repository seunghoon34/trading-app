// File: services/user-management/handlers/auth.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seunghoon34/trading-app/services/user-management/models"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}
