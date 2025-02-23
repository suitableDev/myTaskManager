package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"task-manager/server/routes"
)

func main() {
	router := gin.New()
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(fmt.Sprintf("0.0.0.0:%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
