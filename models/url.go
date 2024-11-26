package models

import "time"

type URL struct {
	LongURL     string    `json:"long_url"`
	ShortCode   string    `json:"short_code"`
	CreatedAt   time.Time `json:"created_at"`
	AccessCount int       `json:"access_count"`
	ExpiresAt   time.Time `json:"expires_at"`
}
