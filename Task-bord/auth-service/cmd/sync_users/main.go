package main

import (
	"auth-service/db"
	"auth-service/models"
	"common/rabbitmq"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	log.Println("Starting User Sync...")

	// 1. Connect to DB
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.DisconnectDB()

	// 2. Connect to RabbitMQ
	rabbitClient, err := rabbitmq.Connect("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitClient.Close()

	// 3. Get all users
	collection := db.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to fetch users: %v", err)
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		log.Fatalf("Failed to decode users: %v", err)
	}

	log.Printf("Found %d users. Starting sync...", len(users))

	// 4. Publish Event for each user
	for _, user := range users {
		eventPayload := map[string]string{
			"userId": user.ID.Hex(),
			"email":  user.Email,
			"role":   user.Role,
		}

		err = rabbitClient.Publish("", "user_queue", "UserRegistered", eventPayload)
		if err != nil {
			log.Printf("Failed to sync user %s: %v", user.Email, err)
		} else {
			fmt.Printf("Parsed sync event for: %s (%s)\n", user.Email, user.Role)
		}
		// Small delay to avoid flooding if many users
		time.Sleep(10 * time.Millisecond)
	}

	log.Println("Sync complete!")
}
