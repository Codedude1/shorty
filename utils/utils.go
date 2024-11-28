package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

// RespondWithJSON sends a JSON response with the specified HTTP status code and payload.
func RespondWithJSON(c *gin.Context, code int, payload interface{}) {
	if code >= 400 {
		log.Printf("[WARN] Responding with error %d: %+v", code, payload)
	} else {
		log.Printf("[INFO] Responding with status %d: %+v", code, payload)
	}
	c.JSON(code, payload)
}

// RespondWithError sends a JSON error response with the specified HTTP status code and message.
func RespondWithError(c *gin.Context, code int, message string) {
	RespondWithJSON(c, code, gin.H{"error": message})
}
