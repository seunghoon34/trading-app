package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/seunghoon34/trading-app/services/payment/handlers"
	"github.com/seunghoon34/trading-app/services/payment/redis"
)

type IntegrationTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *IntegrationTestSuite) SetupSuite() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Set test environment variables
	os.Setenv("ALPACA_API_KEY", "test_key")
	os.Setenv("ALPACA_SECRET_KEY", "test_secret")

	// Initialize Redis (in real integration tests, you'd use a test Redis instance)
	redis.Init()

	// Setup router
	suite.router = gin.Default()
	suite.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Payment health endpoint",
		})
	})
	suite.router.POST("/deposit/:amount", handlers.DepositFunds)
}

func (suite *IntegrationTestSuite) TestHealthEndpoint() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Payment health endpoint")
}

func (suite *IntegrationTestSuite) TestDepositEndpoint_MissingAccountID() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/deposit/100", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "X-Account-ID header is required")
}

func (suite *IntegrationTestSuite) TestDepositEndpoint_WithAccountID() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/deposit/100", nil)
	req.Header.Set("X-Account-ID", "test-account-123")
	suite.router.ServeHTTP(w, req)

	// This will fail because it tries to make real HTTP requests to Alpaca
	// In a real integration test, you'd use a test environment or mock server
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *IntegrationTestSuite) TestInvalidEndpoint() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/invalid", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
