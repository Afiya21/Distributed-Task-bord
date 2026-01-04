package main

import (
	"Task-service/db"
	"Task-service/rabbitmq"
	"Task-service/routes"
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

	// Initialize RabbitMQ
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}
	rabbitmq.InitRabbitMQ(rabbitURL)

	// Initialize MongoDB
	if err := db.InitDB(); err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Define routes for the Task Management Service
	r.POST("/tasks", routes.CreateTask)
	r.GET("/tasks", routes.GetTasks)
	r.PUT("/tasks/:id/status", routes.UpdateTaskStatus)
	r.DELETE("/tasks/:id", routes.DeleteTask)

	// Run the service on port 8081
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	if err := r.Run(":" + port); err != nil {
		panic(err)
	}
}
