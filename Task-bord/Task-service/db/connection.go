package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB connects to MongoDB Atlas
func ConnectDB() (*mongo.Client, context.Context, context.CancelFunc, error) {
	// MongoDB Atlas connection string
	clientOptions := options.Client().ApplyURI("mongodb+srv://nebyatahmed21_db_user:Uv8qT79Qi3OIl7PA@task-service.vj20nn5.mongodb.net/?appName=Task-service")

	// Connect to MongoDB Atlas
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, nil, nil, err
	}

	// Ping the database to confirm connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Ping(ctx, nil)
	if err != nil {
		defer cancel()
		return nil, nil, nil, err
	}

	log.Println("Successfully connected to MongoDB Atlas")

	return client, ctx, cancel, nil
}
