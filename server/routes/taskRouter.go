package routes

import (
	"github.com/gin-gonic/gin"

	controller "task-manager/server/controllers"
	middleware "task-manager/server/middleware"
)

// SetupRoutes configures the routes and binds them to the corresponding handler functions
func SetupRoutes(router *gin.Engine) {
	taskCreateLimiter := middleware.RateLimitMiddleware(middleware.CreateLimiter)
	taskDeleteLimiter := middleware.RateLimitMiddleware(middleware.DeleteLimiter)

	router.GET("/tasks", controller.GetTasks)
	router.GET("/tasks/:id", controller.GetTaskByID)
	router.POST("/tasks", taskCreateLimiter, controller.PostTask)
	router.PUT("/tasks/:id", controller.UpdateTask)
	router.DELETE("/tasks/:id", controller.DeleteTask)
	router.DELETE("/tasks", taskDeleteLimiter, controller.DeleteAllTasks)
}
