package main

import (
	"github.com/gin-gonic/gin"
)

// Utility function for error response
func respondWithError(context *gin.Context, code int, message string, details string) {
	context.JSON(code, gin.H{
		"error":   message,
		"details": details,
	})
}
