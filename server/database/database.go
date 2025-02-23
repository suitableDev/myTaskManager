package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ConnectToMongoDB(ctx)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}

	MongoClient = client
}

func ConnectToMongoDB(ctx context.Context) (*mongo.Client, error) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI is not set in environment variables")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}

// GetTaskCollection retrieves the "tasks" collection from the database.
func GetTaskCollection() *mongo.Collection {
	if MongoClient == nil {
		log.Fatal("MongoDB client is not initialized. Ensure ConnectToMongoDB() is successful.")
	}
	return MongoClient.Database("task_manager").Collection("tasks")
}

// GetUserCollection retrieves the "users" collection from the database.
func GetUserCollection() *mongo.Collection {
	if MongoClient == nil {
		log.Fatal("MongoDB client is not initialized. Ensure ConnectToMongoDB() is successful.")
	}
	return MongoClient.Database("task_manager").Collection("users")
}
