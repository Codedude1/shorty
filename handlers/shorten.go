package handlers

import (
	"net/http"
	"time"

	"github.com/Codedude1/shorty/models"
	"github.com/Codedude1/shorty/services"
	"github.com/Codedude1/shorty/storage"
	"github.com/Codedude1/shorty/utils"
	"github.com/gin-gonic/gin"
)

func ShortenURLHandler(store *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			URL          string `json:"url" binding:"required"`
			ExpiryInMins int    `json:"expiry_in_mins"` // Optional TTL parameter
		}

		// Bind the JSON request body to the request struct
		if err := c.ShouldBindJSON(&request); err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Validate the URL format
		if !services.IsValidURL(request.URL) {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid URL")
			return
		}

		// Lock the storage for writing
		store.Mu.Lock()
		defer store.Mu.Unlock()

		// Check if the long URL already exists
		if shortCode, exists := store.LongURLMap[request.URL]; exists {
			shortURL := c.Request.Host + "/" + shortCode
			response := gin.H{"short_url": shortURL}
			utils.RespondWithJSON(c, http.StatusOK, response)
			return
		}

		// Generate a unique ID and encode it to create the short code
		id := services.GenerateID(&store.IdCounter)
		shortCode := services.Encode(id)

		// Set expiration time if provided
		var expiresAt time.Time
		if request.ExpiryInMins > 0 {
			expiresAt = time.Now().Add(time.Duration(request.ExpiryInMins) * time.Minute)
		}

		// Create the URL model
		urlModel := &models.URL{
			LongURL:     request.URL,
			ShortCode:   shortCode,
			CreatedAt:   time.Now(),
			AccessCount: 0,
			ExpiresAt:   expiresAt,
		}

		// Store the mapping in the storage
		store.UrlMap[shortCode] = urlModel
		store.LongURLMap[request.URL] = shortCode

		// Construct the short URL
		shortURL := c.Request.Host + "/" + shortCode

		// Prepare the response
		response := gin.H{
			"short_url": shortURL,
		}

		// Respond with the short URL
		utils.RespondWithJSON(c, http.StatusOK, response)
	}
}
