package events

import (
	"common/rabbitmq"
	"log"
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
		if event.Type == "UserRegistered" {
			handleUserRegistered(event.Payload)
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
	client, ctx, cancel, err := db.ConnectDB()
	if err != nil {
		log.Printf("Failed to connect to DB: %v", err)
		return
	}
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting: %v", err)
		}
	}()

	collection := client.Database("user-management-service").Collection("users")

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
