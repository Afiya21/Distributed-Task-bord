package auth

import (
	"auth-service/db"
	"auth-service/models"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Hash password before saving
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// Compare password with the hashed password
func comparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// Generate JWT token
func generateJWT(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte("your-secret-key") // Store securely in environment variables
	tokenString, err := token.SignedString(secretKey)
	return tokenString, err
}

// Register a new user
func RegisterUser(email, password, role string) (string, error) {
	client, ctx, cancel := db.ConnectDB()
	defer cancel()

	collection := client.Database("taskboard").Collection("users")
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return "", err
	}

	user := models.User{
		Email:    email,
		Password: hashedPassword,
		Role:     role,
	}

	// Insert user into the database
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	// Generate JWT token for the user
	userID := result.InsertedID.(primitive.ObjectID).Hex()
	token, err := generateJWT(userID, role)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Login a user
func LoginUser(email, password string) (string, error) {
	client, ctx, cancel := db.ConnectDB()
	defer cancel()

	collection := client.Database("taskboard").Collection("users")
	var user models.User

	// Find the user by email
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return "", fmt.Errorf("no user found with this email")
	}

	// Compare the password with the stored hashed password
	if !comparePassword(user.Password, password) {
		return "", fmt.Errorf("invalid password")
	}

	// Generate JWT token for the user
	token, err := generateJWT(user.ID.Hex(), user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Logout - Invalidate the session
func LogoutUser(c *gin.Context) {
	// In JWT, logout is simply deleting the token on the client-side
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out. Please delete the token on the client-side."})
}
