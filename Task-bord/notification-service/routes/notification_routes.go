package routes

import (
	"net/http"
	"notification-service/db"
	"notification-service/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// Create a notification
func CreateNotification(c *gin.Context) {
	var notification models.Notification
	if err := c.BindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set created time
	notification.CreatedAt = time.Now().Format(time.RFC3339)

	client, ctx, cancel := db.ConnectDB()
	defer cancel()

	collection := client.Database("taskboard").Collection("notifications")

	// Insert the notification into the database
	_, err := collection.InsertOne(ctx, notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification sent successfully"})
}

// Get all notifications for a user
func GetNotifications(c *gin.Context) {
	userID := c.Param("user_id")
	client, ctx, cancel := db.ConnectDB()
	defer cancel()

	collection := client.Database("taskboard").Collection("notifications")

	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	var notifications []models.Notification
	if err = cursor.All(ctx, &notifications); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// Mark notification as read
func MarkAsRead(c *gin.Context) {
	notificationID := c.Param("id")
	client, ctx, cancel := db.ConnectDB()
	defer cancel()

	collection := client.Database("taskboard").Collection("notifications")

	// Update the notification's "is_read" field to true
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": notificationID},
		bson.M{"$set": bson.M{"is_read": true}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}
