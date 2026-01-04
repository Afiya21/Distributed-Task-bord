package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

// InitDB initializes the MongoDB connection
func InitDB() error {
	// Define the MongoDB URI (connection string)
	uri := "mongodb+srv://nebyatahmed21_db_user:zEzew7TtvTHmJAcY@user-management-service.neiqyrn.mongodb.net/?appName=user-management-service"

	var err error
	Client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	// Ping the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = Client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	fmt.Println("Connected to MongoDB!")
	return nil
}

// GetCollection returns a MongoDB collection
func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database("taskboard").Collection(collectionName)
}
