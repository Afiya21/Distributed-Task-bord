package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client instance
var Client *mongo.Client

// InitDB initializes the MongoDB connection
func InitDB() error {
	clientOptions := options.Client().ApplyURI("mongodb+srv://nebyatahmed21_db_user:Uv8qT79Qi3OIl7PA@task-service.vj20nn5.mongodb.net/?appName=Task-service")

	var err error
	Client, err = mongo.Connect(context.Background(), clientOptions)
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

	log.Println("Successfully connected to MongoDB Atlas")
	return nil
}

// GetCollection returns a MongoDB collection
func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database("taskboard").Collection(collectionName)
}
