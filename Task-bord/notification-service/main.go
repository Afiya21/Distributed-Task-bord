package main

import (
	"common/middleware"
	"notification-service/events"
	"notification-service/routes"
	"notification-service/websockets"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize WebSocket Hub
	hub := websockets.NewHub()
	go hub.Run()

	// Start RabbitMQ Consumer
	go events.SetupConsumer("amqp://guest:guest@localhost:5672/", hub)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Register notification routes
	r.POST("/notifications", routes.CreateNotification)       // Send a notification
	r.GET("/notifications/:user_id", routes.GetNotifications) // Get all notifications for a user
	r.PUT("/notifications/:id", routes.MarkAsRead)            // Mark a notification as read

	// WebSocket Route
	r.GET("/ws", func(c *gin.Context) {
		userId := c.Query("userId")
		if userId == "" {
			c.JSON(400, gin.H{"error": "userId required"})
			return
		}
		websockets.ServeWs(hub, c.Writer, c.Request, userId)
	})

	// Start the server
	if err := r.Run(":8083"); err != nil {
		panic(err)
	}
}
