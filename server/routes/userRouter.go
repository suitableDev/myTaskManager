package routes

import (
	"github.com/gin-gonic/gin"

	controller "task-manager/server/controllers"
	middleware "task-manager/server/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
}
