package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestRespondWithJSON tests the RespondWithJSON function
func TestRespondWithJSON(t *testing.T) {
	// Capture log output
	var logBuffer bytes.Buffer
	originalOutput := log.Writer()      // Store the original log output
	log.SetOutput(&logBuffer)           // Redirect logs to logBuffer
	defer log.SetOutput(originalOutput) // Restore original log output after the test

	// Initialize Gin in test mode without Logger middleware
	router := gin.New()

	// Define test routes that call RespondWithJSON with different scenarios
	router.GET("/testjson", func(c *gin.Context) {
		payload := gin.H{"message": "success"}
		RespondWithJSON(c, http.StatusOK, payload)
	})

	router.GET("/testjsonerror", func(c *gin.Context) {
		payload := gin.H{"message": "error occurred"}
		RespondWithJSON(c, http.StatusBadRequest, payload)
	})

	// Define test cases
	tests := []struct {
		name               string
		endpoint           string
		expectedStatusCode int
		expectedBody       map[string]string
		expectedLogPrefix  string
	}{
		{
			name:               "Successful JSON Response",
			endpoint:           "/testjson",
			expectedStatusCode: http.StatusOK,
			expectedBody:       map[string]string{"message": "success"},
			expectedLogPrefix:  "[INFO]",
		},
		{
			name:               "Error JSON Response",
			endpoint:           "/testjsonerror",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       map[string]string{"message": "error occurred"},
			expectedLogPrefix:  "[WARN]",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable to avoid closure issues
		t.Run(tt.name, func(t *testing.T) {
			// Reset log buffer before each test
			logBuffer.Reset()

			// Create a new HTTP GET request
			req, err := http.NewRequest(http.MethodGet, tt.endpoint, nil)
			assert.NoError(t, err)

			// Create a ResponseRecorder to record the response
			w := httptest.NewRecorder()

			// Serve the HTTP request
			router.ServeHTTP(w, req)

			// Assert the HTTP status code
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			// Assert the response body
			var response map[string]string
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			// Assert the log output contains the expected prefix
			logOutput := logBuffer.String()
			assert.Contains(t, logOutput, tt.expectedLogPrefix, "Log should contain %s", tt.expectedLogPrefix)
			assert.Contains(t, logOutput, tt.expectedBody["message"], "Log should contain the message")
		})
	}
}

// TestRespondWithError tests the RespondWithError function
func TestRespondWithError(t *testing.T) {
	// Capture log output
	var logBuffer bytes.Buffer
	originalOutput := log.Writer()      // Store the original log output
	log.SetOutput(&logBuffer)           // Redirect logs to logBuffer
	defer log.SetOutput(originalOutput) // Restore original log output after the test

	// Initialize Gin in test mode without Logger middleware
	router := gin.New()

	// Define test routes that call RespondWithError with different scenarios
	router.GET("/testerror", func(c *gin.Context) {
		RespondWithError(c, http.StatusNotFound, "Resource not found")
	})

	router.GET("/testinternalerror", func(c *gin.Context) {
		RespondWithError(c, http.StatusInternalServerError, "Internal server error")
	})

	// Define test cases
	tests := []struct {
		name               string
		endpoint           string
		expectedStatusCode int
		expectedBody       map[string]string
		expectedLogPrefix  string
	}{
		{
			name:               "Not Found Error",
			endpoint:           "/testerror",
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       map[string]string{"error": "Resource not found"},
			expectedLogPrefix:  "[WARN]",
		},
		{
			name:               "Internal Server Error",
			endpoint:           "/testinternalerror",
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       map[string]string{"error": "Internal server error"},
			expectedLogPrefix:  "[WARN]",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable to avoid closure issues
		t.Run(tt.name, func(t *testing.T) {
			// Reset log buffer before each test
			logBuffer.Reset()

			// Create a new HTTP GET request
			req, err := http.NewRequest(http.MethodGet, tt.endpoint, nil)
			assert.NoError(t, err)

			// Create a ResponseRecorder to record the response
			w := httptest.NewRecorder()

			// Serve the HTTP request
			router.ServeHTTP(w, req)

			// Assert the HTTP status code
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			// Assert the response body
			var response map[string]string
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			// Assert the log output contains the expected prefix
			logOutput := logBuffer.String()
			assert.Contains(t, logOutput, tt.expectedLogPrefix, "Log should contain %s", tt.expectedLogPrefix)
			assert.Contains(t, logOutput, tt.expectedBody["error"], "Log should contain the error message")
		})
	}
}
