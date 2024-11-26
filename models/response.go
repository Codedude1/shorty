package models

import "time"

type StatsResponse struct {
	LongURL     string    `json:"long_url"`
	AccessCount int       `json:"access_count"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
}
