// File: services/user-management/handlers/auth.go
package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seunghoon34/trading-app/services/user-management/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context, db *pgxpool.Pool) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	var existingID int

	err = db.QueryRow(context.Background(),
		"SELECT id FROM users WHERE email=$1", user.Email).Scan(&existingID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	} else if err != pgx.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	// If we reach here, it means the user does not exist, and we can proceed to insert
	// Insert user into database
	query := `
        INSERT INTO users (first_name, last_name, email, password, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id, created_at, updated_at`

	var id int
	var createdAt, updatedAt time.Time

	err = db.QueryRow(context.Background(), query,
		user.FirstName, user.LastName, user.Email, string(hashedPassword)).
		Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Set the returned values
	user.ID = id
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt
	user.Password = "" // Don't return password in response

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

func Login(c *gin.Context, db *pgxpool.Pool) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// for testing purposes

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var password string
	err := db.QueryRow(context.Background(), "SELECT password FROM users WHERE email=$1", loginData.Email).Scan(&password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect password or email"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect password or email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User Login successfully",
	})
}
