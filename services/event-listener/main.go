package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seunghoon34/trading-app/services/event-listener/internal/service"
)

func main() {
	// Set up logging for our service
	log.Printf("🚀 Starting Trade Event Listener Service...")

	// Validate required environment variables
	requiredEnvVars := []string{
		"ALPACA_API_KEY",
		"ALPACA_SECRET_KEY",
		"ALPACA_ACCOUNT_ID",
	}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("❌ Required environment variable %s is not set", envVar)
		}
	}

	// Create and start the trade event service
	eventService := service.NewEventService()

	if err := eventService.Start(); err != nil {
		log.Fatalf("❌ Failed to start trade event service: %v", err)
	}

	// Create Gin router for health checks and monitoring
	r := gin.Default()

	// Health check endpoint (important for Docker health checks)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"service":   "trade-event-listener",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Event statistics endpoint (useful for monitoring)
	r.GET("/stats", func(c *gin.Context) {
		stats := eventService.GetStats()
		c.JSON(200, stats)
	})

	// Start HTTP server in a goroutine (non-blocking)
	go func() {
		log.Printf("🌐 Health check server starting on :8085")
		log.Printf("📊 Stats endpoint available at: http://localhost:8085/stats")
		if err := r.Run(":8085"); err != nil {
			log.Printf("❌ Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("✅ Trade Event Listener Service is running:")
	log.Printf("   📈 Listening for trade events from Alpaca SSE")
	log.Printf("   📤 Publishing to trade-events Kafka topic")
	log.Printf("   🏥 Health check: http://localhost:8085/health")
	log.Printf("   📊 Statistics: http://localhost:8085/stats")
	log.Printf("   Press Ctrl+C to stop.")

	<-quit // This blocks until we receive a signal

	log.Printf("🛑 Shutting down Trade Event Listener Service...")

	// Gracefully stop the event service
	eventService.Stop()

	log.Printf("👋 Trade Event Listener Service shutdown complete")
}
