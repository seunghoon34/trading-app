// File: services/user-management/handlers/auth.go
package handlers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seunghoon34/trading-app/services/user-management/internal/alpaca"
	"github.com/seunghoon34/trading-app/services/user-management/models"
	"golang.org/x/crypto/bcrypt"
)

// JWT Claims structure
type Claims struct {
	AccountID string `json:"account_id"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

// Generate JWT token
func generateJWT(accountID, email string) (string, error) {
	// Create claims with user data
	claims := Claims{
		AccountID: accountID,
		Email:     email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hour expiration
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Get JWT secret from environment (we'll set this up)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key" // Default for development
	}

	// Sign and return the token
	return token.SignedString([]byte(jwtSecret))
}

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

	var exists bool
	err = db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", user.Email).Scan(&exists)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}
	// If we reach here, it means the user does not exist, and we can proceed to insert
	// Insert user into database
	alpacaAccount, err := alpaca.CreateAlpacaAccount(user.Email, user.FirstName, user.LastName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Alpaca account: " + err.Error()})
		return
	}
	user.AlpacaAccountID = alpacaAccount.Id

	query := `
        INSERT INTO users (first_name, alpaca_account_id, last_name, email, password, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
        RETURNING id, created_at, updated_at`

	var id int
	var createdAt, updatedAt time.Time

	err = db.QueryRow(context.Background(), query,
		user.FirstName, user.AlpacaAccountID, user.LastName, user.Email, string(hashedPassword)).
		Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create use: " + err.Error()})
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

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var password string
	var alpacaAccountID string
	err := db.QueryRow(context.Background(),
		"SELECT password, alpaca_account_id FROM users WHERE email=$1",
		loginData.Email).Scan(&password, &alpacaAccountID)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect password or email"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect password or email"})
		return
	}

	accountID := alpacaAccountID
	email := loginData.Email // Get email from login request

	// Generate JWT token
	token, err := generateJWT(accountID, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return token and user info
	c.JSON(http.StatusOK, gin.H{
		"message":    "User login successful",
		"token":      token,
		"account_id": accountID,
		"email":      email,
	})
}
