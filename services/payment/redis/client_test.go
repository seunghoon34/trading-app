package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	// Test that Init() creates a client
	Init()

	assert.NotNil(t, Client)
	assert.Equal(t, "localhost:6379", Client.Options().Addr)
}

func TestInit_ClientConfiguration(t *testing.T) {
	// Reset client to nil to test initialization
	Client = nil

	Init()

	assert.NotNil(t, Client)

	// Test client options
	options := Client.Options()
	assert.Equal(t, "localhost:6379", options.Addr)
	assert.Equal(t, "", options.Password) // Default should be empty
	assert.Equal(t, 0, options.DB)        // Default should be 0
}

func TestClient_NotNilAfterInit(t *testing.T) {
	// Ensure client is not nil after initialization
	originalClient := Client
	Client = nil

	Init()

	assert.NotNil(t, Client)

	// Restore original client
	Client = originalClient
}
