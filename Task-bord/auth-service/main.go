package main

import (
	"auth-service/auth"
	"auth-service/db"
	"auth-service/routes"
	"common/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Initialize RabbitMQ
	if err := auth.InitRabbitMQ("amqp://guest:guest@localhost:5672/"); err != nil {
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
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
