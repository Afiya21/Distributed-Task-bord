package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() (*mongo.Client, context.Context, context.CancelFunc) {
	// Define the MongoDB URI (connection string)
	uri := "mongodb+srv://nebyatahmed21_db_user:zEzew7TtvTHmJAcY@user-management-service.neiqyrn.mongodb.net/?appName=user-management-service"

	// Set a timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// Ping MongoDB to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	fmt.Println("Connected to MongoDB!")

	// Return the client, context, and cancel function
	return client, ctx, cancel
}

// DisconnectDB - Close the connection to MongoDB
func DisconnectDB(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	err := client.Disconnect(ctx)
	if err != nil {
		log.Fatal("Error disconnecting from MongoDB:", err)
	}
	cancel()
}
