package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	SetupRoutes(router)
	router.Run("localhost:8080")
}
