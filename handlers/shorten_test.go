package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Codedude1/shorty/models"
	"github.com/Codedude1/shorty/storage"
	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

// ShortenTestCase defines the structure for each test case in TestShortenURLHandler
type ShortenTestCase struct {
	name           string
	requestBody    models.ShortenRequest
	expectedStatus int
	expectError    bool
	expectedURL    string // Expected short URL prefix
}

// TestShortenURLHandler tests the ShortenURLHandler with various scenarios
func TestShortenURLHandler(t *testing.T) {
	// Initialize Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a new storage instance
	store := storage.NewStorage()

	// Initialize the router with the handler
	router := gin.Default()
	router.POST("/shorten", ShortenURLHandler(store))

	// Define test cases
	tests := []ShortenTestCase{
		{
			name: "Valid URL without Expiry",
			requestBody: models.ShortenRequest{
				URL: "https://www.example.com",
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedURL:    "http://",
		},
		{
			name: "Valid URL with Expiry",
			requestBody: models.ShortenRequest{
				URL:          "https://www.google.com",
				ExpiryInMins: 60,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedURL:    "http://",
		},
		{
			name: "Invalid URL",
			requestBody: models.ShortenRequest{
				URL: "not-a-valid-url",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
			expectedURL:    "",
		},
		{
			name: "Empty URL",
			requestBody: models.ShortenRequest{
				URL: "",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
			expectedURL:    "",
		},
		{
			name: "Unsupported Scheme",
			requestBody: models.ShortenRequest{
				URL: "ftp://www.example.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
			expectedURL:    "",
		},
		{
			name: "URL with Special Characters",
			requestBody: models.ShortenRequest{
				URL: "https://www.example.com/path?query=param&another=param2",
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedURL:    "http://",
		},
		{
			name: "Extremely Long URL",
			requestBody: models.ShortenRequest{
				URL: "https://www.example.com/" + strings.Repeat("a", 1000),
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedURL:    "http://",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the request body to JSON
			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			// Create a new HTTP request
			req, err := http.NewRequest(http.MethodPost, "/shorten", strings.NewReader(string(body)))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Create a ResponseRecorder to record the response
			w := httptest.NewRecorder()

			// Serve the HTTP request
			router.ServeHTTP(w, req)

			// Assert the HTTP status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectError {
				// Parse the error response
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				_, exists := response["error"]
				assert.True(t, exists, "Expected error message in response")
			} else {
				// Parse the success response
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				shortURL, exists := response["short_url"]
				assert.True(t, exists, "Expected short_url in response")
				assert.NotEmpty(t, shortURL, "short_url should not be empty")

				// Validate the short URL format
				assert.True(t, strings.HasPrefix(shortURL, tt.expectedURL) || strings.HasPrefix(shortURL, "https://") || strings.HasPrefix(shortURL, "http://"), "Short URL should start with http:// or https://")

				// Validate that the short code exists in storage
				parts := strings.Split(shortURL, "/")
				shortCode := parts[len(parts)-1]
				urlModel, exists := store.GetURL(shortCode)
				assert.True(t, exists, "Short code should exist in storage")
				assert.Equal(t, tt.requestBody.URL, urlModel.LongURL, "Long URL should match the input URL")

				// If ExpiryInMins is set, verify ExpiresAt
				if tt.requestBody.ExpiryInMins > 0 {
					expectedExpiry := time.Now().Add(time.Duration(tt.requestBody.ExpiryInMins) * time.Minute)
					assert.WithinDuration(t, expectedExpiry, urlModel.ExpiresAt, time.Minute, "ExpiresAt should be set correctly")
				} else {
					assert.True(t, urlModel.ExpiresAt.IsZero(), "ExpiresAt should not be set when ExpiryInMins is not provided")
				}
			}
		})
	}
}
func TestShortenURLHandler_DuplicateURL(t *testing.T) {
	// Initialize Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a new storage instance
	store := storage.NewStorage()

	// Add a URL to storage
	existingURL := "https://www.duplicate.com"
	existingShortCode := "dup123"
	store.AddURL(existingURL, existingShortCode, time.Time{})

	// Initialize the router with the handler
	router := gin.Default()
	router.POST("/shorten", ShortenURLHandler(store))

	// Define the request with the same URL
	requestBody := models.ShortenRequest{
		URL: existingURL,
	}

	// Marshal the request body to JSON
	body, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodPost, "/shorten", strings.NewReader(string(body)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(w, req)

	// Assert the HTTP status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the success response
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	shortURL, exists := response["short_url"]
	assert.True(t, exists, "Expected short_url in response")
	assert.NotEmpty(t, shortURL, "short_url should not be empty")

	// Validate that the short code matches the existing one
	parts := strings.Split(shortURL, "/")
	shortCode := parts[len(parts)-1]
	assert.Equal(t, existingShortCode, shortCode, "Short code should match the existing one")
}
