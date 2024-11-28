package models

// ShortenRequest contains fields from the incoming request.
type ShortenRequest struct {
	URL          string `json:"url" binding:"required"`
	ExpiryInMins int    `json:"expiry_in_mins"` // Optional TTL parameter
}
