package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/Codedude1/shorty/storage"
	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func TestRedirectHandler(t *testing.T) {
	// Initialize Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a new storage instance
	store := storage.NewStorage()

	// Add a valid URL to storage (without expiration)
	validShortCode := "abc123"
	validLongURL := "https://www.example.com"
	store.AddURL(validLongURL, validShortCode, time.Time{})

	// Add an expired URL to storage (expires in the past)
	expiredShortCode := "expired123"
	expiredLongURL := "https://www.expired.com"
	expiredAt := time.Now().Add(-1 * time.Hour)
	store.AddURL(expiredLongURL, expiredShortCode, expiredAt)

	// Initialize the router with the RedirectHandler
	router := gin.Default()
	router.GET("/:shortCode", RedirectHandler(store))

	// Define test cases
	tests := []struct {
		name               string
		shortCode          string
		expectedStatusCode int
		expectedLocation   string
		expectError        bool
	}{
		{
			name:               "Valid Short Code",
			shortCode:          validShortCode,
			expectedStatusCode: http.StatusFound, // 302
			expectedLocation:   validLongURL,
			expectError:        false,
		},
		{
			name:               "Non-Existent Short Code",
			shortCode:          "nonexistent",
			expectedStatusCode: http.StatusNotFound, // 404
			expectedLocation:   "",
			expectError:        true,
		},
		{
			name:               "Expired Short Code",
			shortCode:          expiredShortCode,
			expectedStatusCode: http.StatusGone, // 410
			expectedLocation:   "",
			expectError:        true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create a new HTTP GET request
			req, err := http.NewRequest(http.MethodGet, "/"+tt.shortCode, nil)
			assert.NoError(t, err)

			// Create a ResponseRecorder to record the response
			w := httptest.NewRecorder()

			// Serve the HTTP request
			router.ServeHTTP(w, req)

			// Assert the HTTP status code
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if tt.expectError {
				// Parse the error response
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Check for the "error" key in the response
				errorMessage, exists := response["error"]
				assert.True(t, exists, "Expected error message in response")
				assert.NotEmpty(t, errorMessage, "Error message should not be empty")
			} else {
				// For valid redirects, check the "Location" header
				location := w.Header().Get("Location")
				assert.Equal(t, tt.expectedLocation, location, "Redirect location should match the long URL")

				// Verify that the access count has been incremented
				urlModel, exists := store.GetURL(tt.shortCode)
				assert.True(t, exists, "Short code should exist in storage")
				assert.Equal(t, 1, urlModel.AccessCount, "Access count should be incremented to 1")
			}
		})
	}
}

func TestRedirectHandler_AccessCount(t *testing.T) {
	// Initialize Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a new storage instance
	store := storage.NewStorage()

	// Add a URL to storage
	shortCode := "access123"
	longURL := "https://www.access.com"
	store.AddURL(longURL, shortCode, time.Time{})

	// Initialize the router with the RedirectHandler
	router := gin.Default()
	router.GET("/:shortCode", RedirectHandler(store))

	// Define the number of times to access the URL
	accessCount := 5

	for i := 1; i <= accessCount; i++ {
		currentCount := i // Capture the current value
		t.Run("Access_"+strconv.Itoa(currentCount), func(t *testing.T) {
			// Create a new HTTP GET request
			req, err := http.NewRequest(http.MethodGet, "/"+shortCode, nil)
			assert.NoError(t, err)

			// Create a ResponseRecorder to record the response
			w := httptest.NewRecorder()

			// Serve the HTTP request
			router.ServeHTTP(w, req)

			// Assert the HTTP status code
			assert.Equal(t, http.StatusFound, w.Code)

			// Assert the Location header for redirection
			location := w.Header().Get("Location")
			assert.Equal(t, longURL, location, "Redirect location should match the long URL")

			// Verify that the access count has been incremented correctly
			urlModel, exists := store.GetURL(shortCode)
			assert.True(t, exists, "Short code should exist in storage")
			assert.Equal(t, currentCount, urlModel.AccessCount, "Access count should be incremented correctly")
		})
	}
}

func TestRedirectHandler_ExpiredURLCleanup(t *testing.T) {
	// Initialize Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a new storage instance
	store := storage.NewStorage()

	// Add an expired URL to storage
	expiredShortCode := "cleanup123"
	expiredLongURL := "https://www.cleanup.com"
	expiredAt := time.Now().Add(-30 * time.Minute) // Expired 30 minutes ago
	store.AddURL(expiredLongURL, expiredShortCode, expiredAt)

	// Initialize the router with the RedirectHandler
	router := gin.Default()
	router.GET("/:shortCode", RedirectHandler(store))

	// Create a new HTTP GET request for the expired short code
	req, err := http.NewRequest(http.MethodGet, "/"+expiredShortCode, nil)
	assert.NoError(t, err)

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(w, req)

	// Assert the HTTP status code
	assert.Equal(t, http.StatusGone, w.Code)

	// Parse the error response
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check for the "error" key in the response
	errorMessage, exists := response["error"]
	assert.True(t, exists, "Expected error message in response")
	assert.Equal(t, "Short URL has expired", errorMessage, "Error message should indicate expiration")

	// Verify that the expired short code has been removed from storage
	_, exists = store.GetURL(expiredShortCode)
	assert.False(t, exists, "Expired short code should be removed from storage")
}
