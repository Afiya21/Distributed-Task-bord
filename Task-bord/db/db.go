package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// ConnectDB initializes the MongoDB connection
func ConnectDB() (*mongo.Client, context.Context, context.CancelFunc) {
	// Replace with your MongoDB Atlas connection string
	clientOptions := options.Client().ApplyURI("mongodb+srv://nebyatahmed21_db_user:6hiwhcfPOwA3yXzV@task-bord.fcgflno.mongodb.net/?appName=Task-bord") // Atlas URI here
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client, ctx, cancel
}

// DisconnectDB disconnects from MongoDB
func DisconnectDB() {
	if err := client.Disconnect(context.Background()); err != nil {
		log.Fatal(err)
	}
}
