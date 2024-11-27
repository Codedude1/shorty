// handlers/stats.go

package handlers

import (
	"net/http"
	"time"

	"github.com/Codedude1/shorty/models"
	"github.com/Codedude1/shorty/storage"
	"github.com/Codedude1/shorty/utils"
	"github.com/gin-gonic/gin"
)

func StatsHandler(store *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortCode := c.Param("shortCode")

		store.Mu.RLock()
		urlModel, exists := store.URLMap[shortCode]
		store.Mu.RUnlock()

		if !exists {
			utils.RespondWithError(c, http.StatusNotFound, "Short URL not found")
			return
		}

		// Check for expiration
		if !urlModel.ExpiresAt.IsZero() && time.Now().After(urlModel.ExpiresAt) {
			utils.RespondWithError(c, http.StatusNotFound, "Short URL not found or has expired")
			return
		}

		// Prepare the response
		response := models.StatsResponse{
			LongURL:     urlModel.LongURL,
			AccessCount: urlModel.AccessCount,
			CreatedAt:   urlModel.CreatedAt,
			ExpiresAt:   urlModel.ExpiresAt,
		}

		utils.RespondWithJSON(c, http.StatusOK, response)
	}
}
