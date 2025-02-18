package routes

import (
	"github.com/gin-gonic/gin"

	controller "task-manager/server/controllers"
)

// SetupRoutes configures the routes and binds them to the corresponding handler functions
func SetupRoutes(router *gin.Engine) {
	router.GET("/tasks", controller.GetTasks)
	router.GET("/tasks/:id", controller.GetTaskByID)
	router.POST("/tasks", controller.PostTask)
	router.PUT("/tasks/:id", controller.UpdateTask)
	router.DELETE("/tasks/:id", controller.DeleteTask)
	router.DELETE("/tasks/", controller.DeleteAllTasks)
}
