package main

import (
	"Task-service/rabbitmq"
	"Task-service/routes"
	"common/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize RabbitMQ
	rabbitmq.InitRabbitMQ("amqp://guest:guest@localhost:5672/")

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Define routes for the Task Management Service
	r.POST("/tasks", routes.CreateTask)
	r.GET("/tasks", routes.GetTasks)
	r.PUT("/tasks/:id/status", routes.UpdateTaskStatus)
	r.DELETE("/tasks/:id", routes.DeleteTask)

	// Run the service on port 8081
	if err := r.Run(":8081"); err != nil {
		panic(err)
	}
}
