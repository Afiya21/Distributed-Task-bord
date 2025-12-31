package main

import (
	"Task-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Define routes for the Task Management Service
	r.POST("/tasks", routes.CreateTask)
	r.GET("/tasks", routes.GetTasks)
	r.PUT("/tasks/:id", routes.UpdateTask)
	r.DELETE("/tasks/:id", routes.DeleteTask)

	// Run the service on port 8081
	r.Run(":8081")
}
