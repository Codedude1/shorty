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

		// Check if the long URL already exists using encapsulated method
		if existingShortCode, exists := store.GetShortCode(request.URL); exists {
			shortURL := constructShortURL(c, existingShortCode)
			response := gin.H{"short_url": shortURL}
			utils.RespondWithJSON(c, http.StatusOK, response)
			return
		}

		// Hash the long URL
		hash := services.HashString(request.URL)

		// Generate the short code
		shortCode, err := services.EncodeHash(hash, 6) // Adjust length as desired
		if err != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, "Error generating short code")
			return
		}

		// Handle potential collisions by appending a counter
		counter := 1
		originalURL := request.URL // Keep the original URL unchanged
		for {
			// Check if the short code already exists
			if _, exists := store.GetURL(shortCode); !exists {
				break // Unique short code found
			}

			// Collision detected, generate a new hash with a counter
			newHashInput := fmt.Sprintf("%s%d", originalURL, counter)
			hash = services.HashString(newHashInput)
			shortCode, err = services.EncodeHash(hash, 6)
			if err != nil {
				utils.RespondWithError(c, http.StatusInternalServerError, "Error generating short code")
				return
			}
			counter++
		}

		// Set expiration time if provided
		var expiresAt time.Time
		if request.ExpiryInMins > 0 {
			expiresAt = time.Now().Add(time.Duration(request.ExpiryInMins) * time.Minute)
		}

		// Store the mapping in the storage using encapsulated method
		store.AddURL(request.URL, shortCode, expiresAt)

		// Construct the short URL with scheme
		shortURL := constructShortURL(c, shortCode)

		// Prepare the response
		response := gin.H{
			"short_url": shortURL,
		}

		// Respond with the short URL
		utils.RespondWithJSON(c, http.StatusOK, response)
	}
}

// constructShortURL constructs the full short URL based on the request context and short code.
func constructShortURL(c *gin.Context, shortCode string) string {
	// Determine the scheme based on TLS
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	// Construct the full short URL
	return fmt.Sprintf("%s://%s/%s", scheme, c.Request.Host, shortCode)
}
