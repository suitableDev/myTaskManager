package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	helper "task-manager/server/helpers"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			helper.RespondWithError(c, http.StatusUnauthorized, "No Authorization header provided or invalid format", "")
			c.Abort()
			return
		}

		clientToken := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			helper.RespondWithError(c, http.StatusInternalServerError, "Invalid token", err)
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("username", claims.Username)
		c.Set("uid", claims.Uid)
		c.Set("user_type", claims.UserType)

		c.Next()
	}
}
