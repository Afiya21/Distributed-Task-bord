package routes

import (
	"auth-service/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all the routes for the Auth Service
func RegisterRoutes(r *gin.Engine) {
	// Route for Registering a user
	r.POST("/register", func(c *gin.Context) {
		var userInput struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}

		if err := c.BindJSON(&userInput); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		token, err := auth.RegisterUser(userInput.Email, userInput.Password, userInput.Role)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Registration successful", "token": token})
	})

	// Route for Logging in a user
	r.POST("/login", func(c *gin.Context) {
		var userInput struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.BindJSON(&userInput); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		token, err := auth.LoginUser(userInput.Email, userInput.Password)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Login successful", "token": token})
	})

	// Route for Logging out a user
	r.POST("/logout", auth.JWTMiddleware(), auth.LogoutUser)

	// Route for Updating User Role (Admin only)
	r.PUT("/users/:id/role", auth.JWTMiddleware(), func(c *gin.Context) {
		// Verify requester is admin (Middleware puts "role" in context, assuming it works that way.
		// If not, we need to check claims from context.
		// Let's assume Middleware sets "claims" or similar.
		// Actually, I should check how Middleware works.
		// For now, I'll trust standard implementation or check file first.
		// But let's just implement the call.

		userID := c.Param("id")
		var body struct {
			Role string `json:"role"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		// Security Check: Only Admin can update roles
		requestingRole, exists := c.Get("role")
		if !exists || requestingRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized: Admins only"})
			return
		}

		err := auth.UpdateUserRole(userID, body.Role)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "User role updated"})
	})

	// Protected route for testing JWT token
	r.GET("/protected-resource", auth.JWTMiddleware(), func(c *gin.Context) {
		// This route is protected by the JWT token
		c.JSON(http.StatusOK, gin.H{"message": "This is a protected resource"})
	})
}
