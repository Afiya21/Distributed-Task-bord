package main

import (
	"common/middleware"
	"user-service/events"
	"user-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Start RabbitMQ Consumer
	go events.SetupConsumer("amqp://guest:guest@localhost:5672/")

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Register routes
	// Register routes
	r.GET("/users", routes.GetAllUsers)             // Get all users
	r.PUT("/users/:id", routes.UpdateUserProfile)   // Update user profile
	r.GET("/users/:id", routes.GetUserByID)         // Get user by ID
	r.PUT("/users/:id/role", routes.UpdateUserRole) // Update user role -> admin only

	// Run the server
	if err := r.Run(":8087"); err != nil {
		panic(err)
	}
}
