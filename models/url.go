package models

import "time"

// BaseURL contains fields common to multiple responses.
type BaseURL struct {
	LongURL     string    `json:"long_url"`
	CreatedAt   time.Time `json:"created_at"`
	AccessCount int       `json:"access_count"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
}

// URL represents the internal storage model for a shortened URL.
type URL struct {
	BaseURL
	ShortCode string `json:"short_code"`
}

// StatsResponse represents the API response for URL statistics.
type StatsResponse struct {
	BaseURL
}
