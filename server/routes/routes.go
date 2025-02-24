package routes

import (
	"github.com/gin-gonic/gin"

	controller "task-manager/server/controllers"
	middleware "task-manager/server/middleware"
)

func SetupRoutes(router *gin.Engine) {

	// Health check route
	router.GET("/health", controller.HealthCheck())

	// Authentication routes
	router.POST("users/signup", controller.Signup())
	router.POST("users/login", controller.Login())
	router.POST("/users/logout", controller.Logout())
	router.POST("/refresh", controller.RefreshAccessToken())

	// Protected user routes
	router.Use(middleware.Authenticate())
	router.GET("/users", controller.GetUsers())
	router.GET("/users/:userid", controller.GetUser())

	// Task routes
	router.GET("/tasks", controller.GetTasks())
	router.GET("/tasks/:id", controller.GetTaskByID())
	router.POST("/tasks", controller.PostTask())
	router.PUT("/tasks/:id", controller.UpdateTask())
	router.DELETE("/tasks/:id", controller.DeleteTask())
	router.DELETE("/tasks/all", controller.DeleteAllTasks())
}
