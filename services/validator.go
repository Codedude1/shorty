package services

import (
	"net/url"
)

func IsValidURL(inputURL string) bool {
	parsedURL, err := url.ParseRequestURI(inputURL)
	if err != nil {
		return false
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}
	if parsedURL.Host == "" {
		return false
	}
	return true
}
