package routes

import (
	"user-service/db"
	"user-service/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterUser handles user registration

// GetAllUsers fetches all users (for Admin selection)
func GetAllUsers(c *gin.Context) {
	client, ctx, cancel, err := db.ConnectDB()
	if err != nil {
		c.JSON(500, gin.H{"error": "Database connection failed: " + err.Error()})
		return
	}
	defer db.DisconnectDB(client, ctx, cancel)

	collection := client.Database("user-management-service").Collection("users")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer cursor.Close(ctx)

	users := []models.User{}
	if err = cursor.All(ctx, &users); err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse users"})
		return
	}

	// Filter sensitive data if needed, or just return relevant fields.
	// For simplicity, returning full user struct (password is hashed anyway).
	c.JSON(200, users)
}

// UpdateUserProfile updates user profile (Name/Username)
func UpdateUserProfile(c *gin.Context) {
	userId := c.Param("id")
	var body struct {
		Username string `json:"username"`
		Theme    string `json:"theme"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	client, ctx, cancel, err := db.ConnectDB()
	if err != nil {
		c.JSON(500, gin.H{"error": "Database connection failed: " + err.Error()})
		return
	}
	defer db.DisconnectDB(client, ctx, cancel)
	collection := client.Database("user-management-service").Collection("users")

	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid User ID"})
		return
	}

	updateFields := bson.M{}
	if body.Username != "" {
		updateFields["username"] = body.Username
	}
	if body.Theme != "" {
		updateFields["theme"] = body.Theme
	}

	update := bson.M{"$set": updateFields}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(200, gin.H{"message": "Profile updated"})
}

// GetUserByID fetches a single user by ID
func GetUserByID(c *gin.Context) {
	userId := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid User ID"})
		return
	}

	client, ctx, cancel, err := db.ConnectDB()
	if err != nil {
		c.JSON(500, gin.H{"error": "Database connection failed: " + err.Error()})
		return
	}
	defer db.DisconnectDB(client, ctx, cancel)

	collection := client.Database("user-management-service").Collection("users")

	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(200, user)
}

// UpdateUserRole updates a user's role (Admin only)
func UpdateUserRole(c *gin.Context) {
	userId := c.Param("id")
	var body struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if body.Role != "user" && body.Role != "admin" {
		c.JSON(400, gin.H{"error": "Invalid role. Must be 'user' or 'admin'"})
		return
	}

	client, ctx, cancel, err := db.ConnectDB()
	if err != nil {
		c.JSON(500, gin.H{"error": "Database connection failed: " + err.Error()})
		return
	}
	defer db.DisconnectDB(client, ctx, cancel)

	collection := client.Database("user-management-service").Collection("users")

	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid User ID"})
		return
	}

	update := bson.M{"$set": bson.M{"role": body.Role}}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update role"})
		return
	}

	c.JSON(200, gin.H{"message": "User role updated successfully"})
}
