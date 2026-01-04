package main

import (
	"auth-service/auth"
	"auth-service/db"
	"auth-service/routes"
	"common/middleware"
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

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Initialize RabbitMQ
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}
	if err := auth.InitRabbitMQ(rabbitURL); err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
	} else {
		defer auth.RabbitClient.Close()
	}

	// Initialize Database
	if err := db.InitDB(); err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}
	defer db.DisconnectDB()

	routes.RegisterRoutes(r) // Register all routes in auth_routes.go

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
