package handlers

import (
	"net/http"
	"time"

	"github.com/Codedude1/shorty/storage"
	"github.com/Codedude1/shorty/utils"
	"github.com/gin-gonic/gin"
)

func RedirectHandler(store *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortCode := c.Param("shortCode")

		// Retrieve URL from storage using encapsulated method
		urlModel, exists := store.GetURL(shortCode)

		if !exists {
			utils.RespondWithError(c, http.StatusNotFound, "Short URL not found")
			return
		}

		// Check for expiration
		if !urlModel.ExpiresAt.IsZero() && time.Now().After(urlModel.ExpiresAt) {
			// Remove expired URL from storage
			store.DeleteURL(shortCode)
			utils.RespondWithError(c, http.StatusGone, "Short URL has expired")
			return
		}

		// Increment access count using encapsulated method
		store.IncrementAccessCount(shortCode)

		// Redirect to the original long URL
		c.Redirect(http.StatusFound, urlModel.LongURL)
	}
}
