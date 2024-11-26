package services

import "net/url"

func IsValidURL(inputURL string) bool {
	parsedURL, err := url.ParseRequestURI(inputURL)
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""

}
