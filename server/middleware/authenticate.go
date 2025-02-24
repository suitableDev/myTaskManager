package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	helper "task-manager/server/helpers"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken, err := c.Cookie("access_token")

		if err != nil {
			helper.RespondWithError(c, http.StatusUnauthorized, "No Authorization cookie provided", "")
			c.Abort()
			return
		}

		claims, msg := helper.ValidateToken(clientToken)
		if msg != "" {
			helper.RespondWithError(c, http.StatusUnauthorized, "Invalid token", msg)
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
