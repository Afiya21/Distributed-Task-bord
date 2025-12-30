package routes

import (
	"net/http"
	"task-board/auth"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all user-related routes
func RegisterRoutes(router *gin.Engine) {
	router.POST("/register", func(c *gin.Context) {
		var userInput struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}

		if err := c.BindJSON(&userInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		user, err := auth.RegisterUser(userInput.Email, userInput.Password, userInput.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User registered", "user": user})
	})
}
