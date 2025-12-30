package auth

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secretKey = []byte("your-secret-key") // You should store this in an env variable

// JWT Middleware to validate the token in request header
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization token required"})
			c.Abort()
			return
		}

		// The token comes as "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse the token and validate
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set the user ID and role to the context
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])

		c.Next()
	}
}
