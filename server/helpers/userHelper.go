package helper

import "github.com/gin-gonic/gin"

// GetUserDetails - Helper function to extract UID and Username
func GetUserDetails(c *gin.Context) (string, string, bool) {
	UserID, exists := c.Get("uid")
	if !exists {
		return "", "", false
	}

	username, exists := c.Get("username")
	if !exists {
		return "", "", false
	}

	return UserID.(string), username.(string), true
}
