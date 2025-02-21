package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time" // Import time

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"task-manager/server/database"
	"task-manager/server/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := database.ConnectToMongoDB(ctx)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	database.MongoClient = client
	log.Println("Connected to MongoDB successfully")

	router := gin.Default()

	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err = router.Run(fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
