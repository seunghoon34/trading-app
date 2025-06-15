package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/seunghoon34/trading-app/services/notification/internal/models"
	"github.com/seunghoon34/trading-app/services/notification/internal/service"
	"go.mongodb.org/mongo-driver/bson"
)

func GetNotification(n *service.NotificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID := c.Param("account_id")
		filter := bson.M{"account_id": accountID}

		cursor, err := n.MongoClient.Collection.Find(context.TODO(), filter)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(context.TODO())

		var notifications []models.TradeEvent
		if err = cursor.All(context.TODO(), &notifications); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"notifications": notifications,
			"count":         len(notifications),
		})
	}
}
