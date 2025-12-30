package main

import (
	"task-board/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Register routes for user
	routes.RegisterRoutes(r)

	// Run the server
	r.Run(":8080")
}
