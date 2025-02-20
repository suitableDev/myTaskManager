package routes

import (
	"github.com/gin-gonic/gin"

	controller "task-manager/server/controllers"
	middleware "task-manager/server/middleware"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	signupRateLimiter := middleware.RateLimitMiddleware(middleware.SignupLimiter)
	loginRateLimiter := middleware.RateLimitMiddleware(middleware.LoginLimiter)

	incomingRoutes.POST("users/signup", signupRateLimiter, controller.Signup())
	incomingRoutes.POST("users/login", loginRateLimiter, controller.Login())
}
