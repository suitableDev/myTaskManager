package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func ConnectToMongoDB(ctx context.Context) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGO_URI")).SetServerAPIOptions(serverAPI)

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI is not set in the environment variables")
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB successfully")

	MongoClient = client
	return client, nil
}

// GetTaskCollection retrieves the "tasks" collection from the database
func GetTaskCollection() *mongo.Collection {
	if MongoClient == nil {
		log.Println("MongoDB client is not initialized. Please ensure ConnectToMongoDB is successful.TASK_COLLECTION")
		return nil
	}
	return MongoClient.Database("task_manager").Collection("tasks")
}

// GetUserCollection retrieves the "users" collection from the database
func GetUserCollection() *mongo.Collection {
	if MongoClient == nil {
		log.Println("MongoDB client is not initialized. Please ensure ConnectToMongoDB is successful.USER_COLLECTION")
		return nil
	}
	return MongoClient.Database("task_manager").Collection("users")
}
