package models

type ShortenRequest struct {
	URL          string `json:"url" binding:"required"`
	ExpiryInMins int    `json:"expiry_in_mins"` // Optional TTL parameter
}
