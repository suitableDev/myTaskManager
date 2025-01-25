package main

import (
	"github.com/gin-gonic/gin"
)

// respondWithError - Helper function for HTTP error responses
func respondWithError(context *gin.Context, code int, message string, details string) {
	context.IndentedJSON(code, gin.H{
		"error":   message,
		"details": details,
	})
}

// respondWithSuccess - Helper function for HTTP success responses
func respondWithSuccess(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.IndentedJSON(code, SuccessResponse{
		Message: message,
		Data:    data,
	})
}
