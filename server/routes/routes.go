package routes

import (
	"github.com/gin-gonic/gin"

	controller "task-manager/server/controllers"
	middleware "task-manager/server/middleware"
)

func SetupRoutes(router *gin.Engine) {
	// Health check route
	HealthLimiter := middleware.RateLimitMiddleware(middleware.HealthLimiter)

	router.GET("/health", HealthLimiter, controller.HealthCheck)

	// Authentication routes
	signupRateLimiter := middleware.RateLimitMiddleware(middleware.SignupLimiter)
	loginRateLimiter := middleware.RateLimitMiddleware(middleware.LoginLimiter)

	router.POST("users/signup", signupRateLimiter, controller.Signup())
	router.POST("users/login", loginRateLimiter, controller.Login())

	// User routes (protected with authentication middleware)
	router.Use(middleware.Authenticate())
	router.GET("/users", controller.GetUsers())
	router.GET("/users/:userid", controller.GetUser())

	// Task routes
	taskCreateLimiter := middleware.RateLimitMiddleware(middleware.CreateLimiter)
	taskDeleteLimiter := middleware.RateLimitMiddleware(middleware.DeleteLimiter)

	router.GET("/tasks", controller.GetTasks)
	router.GET("/tasks/:id", controller.GetTaskByID)
	router.POST("/tasks", taskCreateLimiter, controller.PostTask)
	router.PUT("/tasks/:id", controller.UpdateTask)
	router.DELETE("/tasks/:id", controller.DeleteTask)
	router.DELETE("/tasks", taskDeleteLimiter, controller.DeleteAllTasks)
}
