package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv" // Added for string conversion
	"testing"
	"time"

	"github.com/Codedude1/shorty/models"
	"github.com/Codedude1/shorty/storage"
	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

// StatsTestCase defines the structure for each test case in TestStatsHandler
type StatsTestCase struct {
	name               string
	shortCode          string
	expectedStatusCode int
	expectError        bool
	expectedAccess     int
	expectExpiresAt    bool // New field to indicate expectation
}

func TestStatsHandler(t *testing.T) {
	// Initialize Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a new storage instance
	store := storage.NewStorage()

	// Add a valid URL to storage (without expiration)
	validShortCode := "stat123"
	validLongURL := "https://www.stats.com"
	store.AddURL(validLongURL, validShortCode, time.Time{})

	// Add an expired URL to storage (will be removed by the handler)
	expiredShortCode := "statexpired"
	expiredLongURL := "https://www.statexpired.com"
	expiredExpiresAt := time.Now().Add(-1 * time.Hour)
	store.AddURL(expiredLongURL, expiredShortCode, expiredExpiresAt)

	// Initialize the router with the handler
	router := gin.Default()
	router.GET("/stats/:shortCode", StatsHandler(store))

	// Define test cases
	tests := []StatsTestCase{
		{
			name:               "Valid Short Code without Expiry",
			shortCode:          validShortCode,
			expectedStatusCode: http.StatusOK,
			expectError:        false,
			expectedAccess:     0,
			expectExpiresAt:    false, // ExpiresAt is zero
		},
		{
			name:               "Expired Short Code",
			shortCode:          expiredShortCode,
			expectedStatusCode: http.StatusNotFound,
			expectError:        true,
			expectedAccess:     0,
			expectExpiresAt:    false, // ExpiresAt is zero (since it's expired and removed)
		},
		{
			name:               "Non-Existent Short Code",
			shortCode:          "statnonexistent",
			expectedStatusCode: http.StatusNotFound,
			expectError:        true,
			expectedAccess:     0,
			expectExpiresAt:    false, // Does not exist
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new HTTP request
			req, err := http.NewRequest(http.MethodGet, "/stats/"+tt.shortCode, nil)
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
				_, exists := response["error"]
				assert.True(t, exists, "Expected error message in response")
			} else {
				// Parse the success response
				var response models.StatsResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, validLongURL, response.LongURL, "Long URL should match")
				assert.Equal(t, tt.expectedAccess, response.AccessCount, "Access count should match")

				if tt.expectExpiresAt {
					assert.False(t, response.ExpiresAt.IsZero(), "ExpiresAt should not be zero")
				} else {
					assert.True(t, response.ExpiresAt.IsZero(), "ExpiresAt should be zero")
				}
			}
		})
	}
}

func TestStatsHandler_AfterAccess(t *testing.T) {
	// Initialize Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a new storage instance
	store := storage.NewStorage()

	// Add a URL to storage (without expiration)
	shortCode := "stataccess"
	longURL := "https://www.stataccess.com"
	store.AddURL(longURL, shortCode, time.Time{})

	// Initialize the router with the handlers
	router := gin.Default()
	router.GET("/:shortCode", RedirectHandler(store))
	router.GET("/stats/:shortCode", StatsHandler(store))

	// Define the number of times to access the URL
	accessCount := 3
	for i := 1; i <= accessCount; i++ {
		currentCount := i // Capture the current value of i
		t.Run("AccessCount_"+strconv.Itoa(currentCount), func(t *testing.T) {
			// Create a new HTTP request for redirection
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
			assert.Equal(t, longURL, location)

			// Verify access count
			urlModel, exists := store.GetURL(shortCode)
			assert.True(t, exists, "Short code should exist in storage")
			assert.Equal(t, currentCount, urlModel.AccessCount, "Access count should be incremented correctly")
		})
	}

	// Now, retrieve the stats and verify the access count
	t.Run("Retrieve Stats", func(t *testing.T) {
		// Create a new HTTP request for stats
		req, err := http.NewRequest(http.MethodGet, "/stats/"+shortCode, nil)
		assert.NoError(t, err)

		// Create a ResponseRecorder to record the response
		w := httptest.NewRecorder()

		// Serve the HTTP request
		router.ServeHTTP(w, req)

		// Assert the HTTP status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Parse the success response
		var response models.StatsResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, longURL, response.LongURL, "Long URL should match")
		assert.Equal(t, accessCount, response.AccessCount, "Access count should match the number of accesses")
		assert.True(t, response.CreatedAt.Before(time.Now()), "CreatedAt should be in the past")
		assert.True(t, response.ExpiresAt.IsZero(), "ExpiresAt should be zero if not set")
	})
}
