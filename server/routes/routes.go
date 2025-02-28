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
	router.POST("users/signup", middleware.RateLimitMiddleware(0.1, 1), controller.Signup())
	router.POST("users/login", middleware.RateLimitMiddleware(0.1, 1), controller.Login())
	router.POST("/users/logout", controller.Logout())
	router.POST("/refresh", middleware.RateLimitMiddleware(0.5, 2), controller.RefreshAccessToken())
	router.GET("/verify", controller.VerifyEmail())
	router.POST("/users/forgot-password", middleware.RateLimitMiddleware(0.1, 1), controller.ForgotPassword())
	router.POST("/users/reset-password", middleware.RateLimitMiddleware(0.1, 1), controller.ResetPassword())

	// Authenticate
	router.Use(middleware.Authenticate())

	// User Routes
	router.GET("/users", middleware.RateLimitMiddleware(3, 6), controller.GetUsers())
	router.GET("/users/:userid", middleware.RateLimitMiddleware(3, 5), controller.GetUser())

	// Task Routes
	router.GET("/tasks", middleware.RateLimitMiddleware(10, 20), controller.GetTasks())
	router.GET("/tasks/:id", middleware.RateLimitMiddleware(3, 6), controller.GetTaskByID())
	router.POST("/tasks", middleware.RateLimitMiddleware(1, 3), controller.PostTask())
	router.PUT("/tasks/:id", middleware.RateLimitMiddleware(2, 5), controller.UpdateTask())
	router.DELETE("/tasks/:id", middleware.RateLimitMiddleware(0.5, 1), controller.DeleteTask())
	router.DELETE("/tasks/all", middleware.RateLimitMiddleware(0.5, 1), controller.DeleteAllTasks())
}
