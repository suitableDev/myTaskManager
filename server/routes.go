package main

import "github.com/gin-gonic/gin"

// SetupRoutes configures the routes and binds them to the corresponding handler functions
func SetupRoutes(router *gin.Engine) {
	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", getTaskByID)
	router.POST("/tasks", postTask)
	router.PUT("/tasks/:id", updateTask)
	router.DELETE("/tasks/:id", deleteTask)
	router.DELETE("/tasks/", deleteAllTasks)
}
