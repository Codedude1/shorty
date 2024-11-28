package services

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidURL(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		inputURL string
		expected bool
	}{
		// Valid URLs
		{
			name:     "Valid HTTP URL",
			inputURL: "http://www.example.com",
			expected: true,
		},
		{
			name:     "Valid HTTPS URL with Path and Query",
			inputURL: "https://example.com/path?query=123",
			expected: true,
		},
		{
			name:     "Valid HTTPS URL with Fragment",
			inputURL: "https://example.com/path#section",
			expected: true,
		},
		{
			name:     "Valid URL with IP Address",
			inputURL: "http://127.0.0.1",
			expected: true,
		},
		{
			name:     "Valid URL with IPv6 Address",
			inputURL: "https://[::1]/",
			expected: true,
		},
		{
			name:     "Valid URL with Port",
			inputURL: "http://example.com:8080",
			expected: true,
		},
		{
			name:     "Valid URL with User Info",
			inputURL: "https://user:pass@example.com",
			expected: true,
		},
		// Invalid URLs
		{
			name:     "Invalid Scheme",
			inputURL: "ftp://www.example.com",
			expected: false,
		},
		{
			name:     "Unsupported Scheme",
			inputURL: "mailto:user@example.com",
			expected: false,
		},
		{
			name:     "Missing Scheme",
			inputURL: "www.example.com",
			expected: false,
		},
		{
			name:     "Empty String",
			inputURL: "",
			expected: false,
		},
		{
			name:     "Malformed URL",
			inputURL: "htp:/example.com",
			expected: false,
		},
		{
			name:     "URL with Empty Host",
			inputURL: "http:///path",
			expected: false,
		},
		{
			name:     "URL with Spaces",
			inputURL: "http://example .com",
			expected: false,
		},
		{
			name:     "URL with Non-ASCII Characters",
			inputURL: "http://例子.测试",
			expected: true, // Parsed correctly, host is not empty
		},
		{
			name:     "URL with Only Path",
			inputURL: "/just/a/path",
			expected: false,
		},
		{
			name:     "Extremely Long URL",
			inputURL: "http://www.example.com/" + strings.Repeat("a", 2000),
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidURL(tt.inputURL)
			assert.Equal(t, tt.expected, result, "IsValidURL(%s) should be %v", tt.inputURL, tt.expected)
		})
	}
}
