// File: services/user-management/handlers/auth.go
package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seunghoon34/trading-app/services/user-management/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

func Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var userEmails = []string{"test@example.com", "test2@example.com"}

	// for testing purposes
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	emailFound := false
	for _, email := range userEmails {
		if email == loginData.Email {
			emailFound = true
			break
		}
	}

	// If email not found, return error
	if !emailFound {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect password or email"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "incorrect password or email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User Login successfully",
	})
}
