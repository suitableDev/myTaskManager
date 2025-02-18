package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	database "task-manager/server/database"
	helper "task-manager/server/helpers"
	routes "task-manager/server/routes"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	helper.SECRET_KEY = os.Getenv("SECRET_KEY")
	if helper.SECRET_KEY == "" {
		log.Panic("SECRET_KEY environment variable is not set")
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
	router.Use(gin.Logger())

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
