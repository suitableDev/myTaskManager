package helper

import (
	"log"

	"github.com/gin-gonic/gin"

	model "task-manager/server/models"
)

// RespondWithError - Logs the error and sends an HTTP error response
func RespondWithError(c *gin.Context, code int, message string, details string) {
	log.Printf("%s: %v", message, details)

	c.IndentedJSON(code, gin.H{
		"error":   message,
		"details": details,
	})
}

// RespondWithSuccess - Sends an HTTP success response
func RespondWithSuccess(c *gin.Context, code int, message string, data interface{}) {
	c.IndentedJSON(code, model.SuccessResponse{
		Message: message,
		Data:    data,
	})
}
