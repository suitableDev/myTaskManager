package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"task-manager/server/database"
	"task-manager/server/routes"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ctx := context.Background()
	client, err := database.ConnectToMongoDB(ctx)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}
	database.MongoClient = client
	log.Println("Connected to MongoDB successfully")
}

func main() {
	router := gin.Default()

	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := router.Run(fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
