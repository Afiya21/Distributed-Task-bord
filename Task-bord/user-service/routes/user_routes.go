package routes

import (
	"time"
	"user-service/db"
	"user-service/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// GenerateJWT creates a JWT token
func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	secretKey := []byte("your-secret-key")
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// RegisterUser handles user registration
func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error hashing password"})
		return
	}
	user.Password = string(hashedPassword)

	client, ctx, cancel := db.ConnectDB()
	defer cancel()

	collection := client.Database("taskboard").Collection("users")

	// Check if user already exists
	var existingUser models.User
	err = collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err != mongo.ErrNoDocuments {
		c.JSON(400, gin.H{"error": "User already exists"})
		return
	}

	// Insert the user into the database
	user.ID = primitive.NewObjectID()
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to insert user"})
		return
	}

	c.JSON(200, gin.H{"message": "User registered successfully"})
}

// LoginUser handles user login
func LoginUser(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	client, ctx, cancel := db.ConnectDB()
	defer cancel()

	collection := client.Database("taskboard").Collection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": loginData.Email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	// Check if password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(400, gin.H{"error": "Invalid password"})
		return
	}

	// Generate JWT token
	token, err := GenerateJWT(user.ID.Hex())
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(200, gin.H{"message": "Login successful", "token": token})
}
