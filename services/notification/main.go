// main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/seunghoon34/trading-app/services/notification/handlers"
	"github.com/seunghoon34/trading-app/services/notification/internal/service"
)

func main() {
	// Set up logging for our service
	log.Printf("🚀 Starting Notification Service...")

	// Create and start the notification service
	notificationService, err := service.NewNotificationService()
	if err != nil {
		log.Fatalf("❌ Failed to create notification service: %v", err)
	}

	if err := notificationService.Start(); err != nil {
		log.Fatalf("❌ Failed to start notification service: %v", err)
	}

	// Create Gin router for health checks and monitoring
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"service":   "notification-service",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	r.GET("/notifications/:account_id", handlers.GetNotification(notificationService))

	// Service statistics endpoint
	r.GET("/stats", func(c *gin.Context) {
		stats := notificationService.GetStats()
		c.JSON(200, stats)
	})

	// Start HTTP server in a goroutine (non-blocking)
	go func() {
		log.Printf("🌐 Health check server starting on :8087")
		log.Printf("📊 Stats endpoint available at: http://localhost:8087/stats")
		if err := r.Run(":8087"); err != nil {
			log.Printf("❌ Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("✅ Notification Service is running:")
	log.Printf("   📥 Consuming from trade-events Kafka topic")
	log.Printf("   💾 Storing events in MongoDB")
	log.Printf("   📝 Logging detailed notifications")
	log.Printf("   🏥 Health check: http://localhost:8087/health")
	log.Printf("   📊 Statistics: http://localhost:8087/stats")
	log.Printf("   Press Ctrl+C to stop.")

	<-quit // This blocks until we receive a signal

	log.Printf("🛑 Shutting down Notification Service...")

	// Gracefully stop the notification service
	notificationService.Stop()

	log.Printf("👋 Notification Service shutdown complete")
}
