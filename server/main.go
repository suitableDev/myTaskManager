package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ctx := context.Background()
	client, err := connectToMongoDB(ctx)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}
	mongoClient = client
	log.Println("Connected to MongoDB successfully")
}

func main() {
	router := gin.Default()
	SetupRoutes(router)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := router.Run(fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
