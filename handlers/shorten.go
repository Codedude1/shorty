// handlers/shorten.go

package handlers

import (
	"fmt"
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
		var request models.ShortenRequest

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

		// Hash the long URL
		hash := services.HashString(request.URL)

		// Generate the short code
		shortCode := services.EncodeHash(hash, 6) // Adjust length as desired

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

		counter := 1
		for {
			newHashInput := fmt.Sprintf("%s%d", request.URL, counter)
			hash = services.HashString(newHashInput)
			shortCode = services.EncodeHash(hash, 6)
			if _, exists := store.URLMap[shortCode]; !exists {
				break
			}
			counter++
		}

		// Set expiration time if provided
		var expiresAt time.Time
		if request.ExpiryInMins > 0 {
			expiresAt = time.Now().Add(time.Duration(request.ExpiryInMins) * time.Second)
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
		store.URLMap[shortCode] = urlModel
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
