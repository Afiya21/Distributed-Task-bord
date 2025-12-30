package main

import (
	"auth-service/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routes.RegisterRoutes(r) // Register all routes in auth_routes.go

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
