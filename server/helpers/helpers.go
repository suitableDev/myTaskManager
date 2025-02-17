package helpers

import (
	"github.com/gin-gonic/gin"

	"task-manager/server/models"
)

// respondWithError - Helper function for HTTP error responses
func RespondWithError(context *gin.Context, code int, message string, details string) {
	context.IndentedJSON(code, gin.H{
		"error":   message,
		"details": details,
	})
}

// respondWithSuccess - Helper function for HTTP success responses
func RespondWithSuccess(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.IndentedJSON(code, models.SuccessResponse{
		Message: message,
		Data:    data,
	})
}
