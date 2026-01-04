package main

import (
	"common/middleware"
	"user-service/db"
	"user-service/events"
	"user-service/routes"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Start RabbitMQ Consumer
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}
	go events.SetupConsumer(rabbitURL)

	// Initialize MongoDB
	if err := db.InitDB(); err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Register routes
	r.GET("/users", routes.GetAllUsers)             // Get all users
	r.PUT("/users/:id", routes.UpdateUserProfile)   // Update user profile
	r.GET("/users/:id", routes.GetUserByID)         // Get user by ID
	r.PUT("/users/:id/role", routes.UpdateUserRole) // Update user role -> admin only

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}
	if err := r.Run(":" + port); err != nil {
		panic(err)
	}
}
