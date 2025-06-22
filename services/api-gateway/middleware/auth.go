package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JWT Claims structure (must match user-management service)
type Claims struct {
	AccountID string `json:"account_id"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

// JWT Middleware function
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get token from cookie first
		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			// Fallback to Authorization header for backwards compatibility
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
				c.Abort()
				return
			}

			// Check Bearer format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
				c.Abort()
				return
			}

			tokenString = parts[1]
		}

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Get JWT secret (same as user-management service)
			jwtSecret := os.Getenv("JWT_SECRET")
			if jwtSecret == "" {
				log.Fatal("JWT_SECRET environment variable is required")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			// Store account_id in context for forwarding handlers
			c.Set("account_id", claims.AccountID)
			c.Set("email", claims.Email)

			// Continue to next handler
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}
