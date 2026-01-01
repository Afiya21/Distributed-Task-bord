package main

import (
	"user-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Register routes
	r.POST("/register", routes.RegisterUser) // User Registration
	r.POST("/login", routes.LoginUser)       // User Login

	// Run the server
	r.Run(":8087") // Run on port 8087
}
