package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"

	"github.com/seunghoon34/trading-app/services/payment/handlers"
	rdb "github.com/seunghoon34/trading-app/services/payment/redis"
)

// TestHelper provides common test utilities
type TestHelper struct {
	MockRedis   redismock.ClientMock
	TestServer  *httptest.Server
	OriginalEnv map[string]string
}

// NewTestHelper creates a new test helper with mocked dependencies
func NewTestHelper() *TestHelper {
	// Setup mock Redis
	db, mock := redismock.NewClientMock()
	rdb.Client = db

	// Save original environment variables
	originalEnv := map[string]string{
		"ALPACA_API_KEY":    os.Getenv("ALPACA_API_KEY"),
		"ALPACA_SECRET_KEY": os.Getenv("ALPACA_SECRET_KEY"),
	}

	// Set test environment variables
	os.Setenv("ALPACA_API_KEY", "test_key")
	os.Setenv("ALPACA_SECRET_KEY", "test_secret")

	return &TestHelper{
		MockRedis:   mock,
		OriginalEnv: originalEnv,
	}
}

// Cleanup restores original state
func (th *TestHelper) Cleanup() {
	// Restore environment variables
	for key, value := range th.OriginalEnv {
		if value == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, value)
		}
	}

	// Close test server if it exists
	if th.TestServer != nil {
		th.TestServer.Close()
	}
}

// CreateMockAlpacaServer creates a mock Alpaca API server
func (th *TestHelper) CreateMockAlpacaServer(accountID string, responses map[string]interface{}) {
	th.TestServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case fmt.Sprintf("/v1/accounts/%s/ach_relationships", accountID):
			if r.Method == "GET" {
				if response, ok := responses["get_ach"]; ok {
					data, _ := json.Marshal(response)
					w.WriteHeader(http.StatusOK)
					w.Write(data)
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("[]"))
				}
			} else if r.Method == "POST" {
				if response, ok := responses["create_ach"]; ok {
					data, _ := json.Marshal(response)
					w.WriteHeader(http.StatusOK)
					w.Write(data)
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"id": "test-ach-id"}`))
				}
			}
		case fmt.Sprintf("/v1/accounts/%s/transfers", accountID):
			if r.Method == "POST" {
				if response, ok := responses["transfer"]; ok {
					data, _ := json.Marshal(response)
					w.WriteHeader(http.StatusCreated)
					w.Write(data)
				} else {
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte(`{"id": "test-transfer-id"}`))
				}
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

// SetupMockRedisCache sets up Redis cache expectations
func (th *TestHelper) SetupMockRedisCache(accountID string, hit bool, achDetails *handlers.ACHDetails) {
	cacheKey := fmt.Sprintf("ach_details:%s", accountID)

	if hit && achDetails != nil {
		achJSON, _ := json.Marshal(achDetails)
		th.MockRedis.ExpectGet(cacheKey).SetVal(string(achJSON))
	} else {
		th.MockRedis.ExpectGet(cacheKey).RedisNil()
		if achDetails != nil {
			achJSON, _ := json.Marshal(achDetails)
			th.MockRedis.ExpectSet(cacheKey, string(achJSON), time.Hour).SetVal("OK")
		}
	}
}

// CreateTestRouter creates a Gin router for testing
func CreateTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Payment health endpoint",
		})
	})

	r.POST("/deposit/:amount", handlers.DepositFunds)

	return r
}

// MockRedisClient provides a mock Redis client for testing
type MockRedisClient struct {
	Data map[string]string
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		Data: make(map[string]string),
	}
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx, "get", key)
	if value, exists := m.Data[key]; exists {
		cmd.SetVal(value)
	} else {
		cmd.SetErr(redis.Nil)
	}
	return cmd
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx, "set", key, value)
	m.Data[key] = fmt.Sprintf("%v", value)
	cmd.SetVal("OK")
	return cmd
}

// TestEnvironment manages test environment setup
type TestEnvironment struct {
	Router    *gin.Engine
	Helper    *TestHelper
	AccountID string
}

func NewTestEnvironment() *TestEnvironment {
	gin.SetMode(gin.TestMode)

	return &TestEnvironment{
		Router:    CreateTestRouter(),
		Helper:    NewTestHelper(),
		AccountID: "test-account-123",
	}
}

func (te *TestEnvironment) Cleanup() {
	te.Helper.Cleanup()
}
