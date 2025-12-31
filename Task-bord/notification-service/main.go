package main

import (
	"notification-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Register notification routes
	r.POST("/notifications", routes.CreateNotification)       // Send a notification
	r.GET("/notifications/:user_id", routes.GetNotifications) // Get all notifications for a user
	r.PUT("/notifications/:id", routes.MarkAsRead)            // Mark a notification as read

	// Start the server
	r.Run(":8083") // Port 8083 for Notification Service
}
