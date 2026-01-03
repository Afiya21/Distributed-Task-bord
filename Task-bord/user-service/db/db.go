package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB connects to MongoDB
// ConnectDB connects to MongoDB
func ConnectDB() (*mongo.Client, context.Context, context.CancelFunc, error) {
	// Define the MongoDB URI (connection string) - Specific to User Service
	uri := "mongodb+srv://nebyatahmed21_db_user:zEzew7TtvTHmJAcY@user-management-service.neiqyrn.mongodb.net/?appName=user-management-service"

	// Connect to MongoDB
	// Use context.Background() for the client connection so it doesn't timeout the client itself unexpectedly
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, nil, err
	}

	// Ping MongoDB to verify connection
	// Use a timeout specifically for the ping check
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Ping(ctx, nil)
	if err != nil {
		cancel() // Cancel the context if ping fails
		return nil, nil, nil, err
	}

	fmt.Println("Connected to MongoDB!")

	// Return the client, context, and cancel function
	return client, ctx, cancel, nil
}

// DisconnectDB - Close the connection to MongoDB
func DisconnectDB(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	if client != nil {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}
	cancel()
}
