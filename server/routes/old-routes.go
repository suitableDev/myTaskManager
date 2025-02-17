package routes

import (
	"github.com/gin-gonic/gin"

	"task-manager/server/controllers"
)

// SetupRoutes configures the routes and binds them to the corresponding handler functions
func SetupRoutes(router *gin.Engine) {
	router.GET("/tasks", controllers.GetTasks)
	router.GET("/tasks/:id", controllers.GetTaskByID)
	router.POST("/tasks", controllers.PostTask)
	router.PUT("/tasks/:id", controllers.UpdateTask)
	router.DELETE("/tasks/:id", controllers.DeleteTask)
	router.DELETE("/tasks/", controllers.DeleteAllTasks)
}
