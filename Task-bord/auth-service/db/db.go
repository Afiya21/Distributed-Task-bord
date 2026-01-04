package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

// InitDB - Connect to MongoDB and initialize the global client
func InitDB() error {
	// Define the MongoDB URI (connection string)
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return fmt.Errorf("MONGO_URI not set in environment")
	}

	// Set a timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	// Ping MongoDB to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("MongoDB connection failed: %v", err)
	}

	Client = client
	fmt.Println("Connected to MongoDB!")
	return nil
}

// GetCollection returns a handle to a MongoDB collection
func GetCollection(collectionName string) *mongo.Collection {
	if Client == nil {
		log.Fatal("Database not initialized")
	}
	return Client.Database("taskboard").Collection(collectionName)
}

// DisconnectDB - Close the connection to MongoDB
func DisconnectDB() {
	if Client == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := Client.Disconnect(ctx); err != nil {
		log.Fatal("Error disconnecting from MongoDB:", err)
	}
	fmt.Println("Disconnected from MongoDB")
}
