package events

import (
	"common/rabbitmq"
	"context"
	"log"
	"time"
	"user-service/db"
	"user-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SetupConsumer starts listening for user events
func SetupConsumer(url string) {
	client, err := rabbitmq.Connect(url)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	// Do not close client immediately, it needs to stay open for consuming
	err = client.Consume("user_queue", func(event rabbitmq.Event) {
		log.Printf("Received event: %s", event.Type)
		switch event.Type {
		case "UserRegistered":
			handleUserRegistered(event.Payload)
		case "UserRoleUpdated":
			handleUserRoleUpdated(event.Payload)
		}
	})
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}
}

func handleUserRegistered(payload interface{}) {
	// Payload comes as interface{}, need to marshal/unmarshal or type assert map
	// Since JSON unmarshal to interface{} produces map[string]interface{}
	data, ok := payload.(map[string]interface{})
	if !ok {
		log.Printf("Invalid payload type: %T", payload)
		return
	}

	// Extract fields
	email, _ := data["email"].(string)
	role, _ := data["role"].(string)
	idStr, _ := data["userId"].(string)
	username, _ := data["username"].(string) // Extract username

	// In a real app we might want to store more info or fetch it using the ID
	// valid mongo ID?
	objID, _ := primitive.ObjectIDFromHex(idStr)

	log.Printf("Syncing user: %s (%s)", email, role)

	// Save to DB
	collection := db.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user exists (idempotency)
	existingCount, err := collection.CountDocuments(ctx, bson.M{"_id": objID})
	if err != nil {
		log.Printf("Error checking existing user: %v", err)
		return
	}
	if existingCount > 0 {
		log.Printf("User %s already exists, skipping sync.", email)
		return
	}

	newUser := models.User{
		ID:       objID,
		Email:    email,
		Role:     role,
		Username: username, // Save username
		Theme:    "light",  // Default theme
	}

	_, err = collection.InsertOne(ctx, newUser)
	if err != nil {
		log.Printf("Failed to insert synced user: %v", err)
	} else {
		log.Printf("User %s (ID: %s) synced successfully to User Service DB", email, idStr)
	}
}

func handleUserRoleUpdated(payload interface{}) {
	data, ok := payload.(map[string]interface{})
	if !ok {
		log.Printf("Invalid payload type for role update: %T", payload)
		return
	}

	userID, _ := data["userId"].(string)
	newRole, _ := data["role"].(string)

	if userID == "" || newRole == "" {
		log.Printf("Invalid role update data: %v", data)
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Invalid user ID in role update: %v", err)
		return
	}

	collection := db.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"role": newRole}}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Failed to update user role in sync: %v", err)
		return
	}

	if result.ModifiedCount > 0 {
		log.Printf("Synced role update for user %s to %s", userID, newRole)
	} else {
		log.Printf("No user found or role unchanged for user %s", userID)
	}
}
