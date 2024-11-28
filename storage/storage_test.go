package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddURL(t *testing.T) {
	store := NewStorage()

	// Test case: Add a single URL without expiration
	url1 := "https://www.example.com"
	shortCode1 := "exmpl1"
	store.AddURL(url1, shortCode1, time.Time{})

	// Verify that URLMap has the shortCode1
	urlModel, exists := store.GetURL(shortCode1)
	assert.True(t, exists, "Short code should exist in URLMap")
	assert.Equal(t, url1, urlModel.LongURL, "Long URL should match the input URL")
	assert.Equal(t, 0, urlModel.AccessCount, "Access count should be initialized to 0")
	assert.True(t, urlModel.ExpiresAt.IsZero(), "ExpiresAt should be zero for no expiration")

	// Verify that LongURLMap has the url1
	retrievedShortCode, exists := store.GetShortCode(url1)
	assert.True(t, exists, "Long URL should exist in LongURLMap")
	assert.Equal(t, shortCode1, retrievedShortCode, "Retrieved short code should match the input short code")

	// Test case: Add another URL with expiration
	url2 := "https://www.google.com"
	shortCode2 := "googl1"
	expiration := time.Now().Add(24 * time.Hour)
	store.AddURL(url2, shortCode2, expiration)

	// Verify that URLMap has the shortCode2
	urlModel2, exists := store.GetURL(shortCode2)
	assert.True(t, exists, "Short code2 should exist in URLMap")
	assert.Equal(t, url2, urlModel2.LongURL, "Long URL2 should match the input URL")
	assert.Equal(t, 0, urlModel2.AccessCount, "Access count2 should be initialized to 0")
	assert.False(t, urlModel2.ExpiresAt.IsZero(), "ExpiresAt should be set for URL2")
	assert.WithinDuration(t, expiration, urlModel2.ExpiresAt, time.Second, "ExpiresAt should be correctly set for URL2")

	// Verify that LongURLMap has the url2
	retrievedShortCode2, exists := store.GetShortCode(url2)
	assert.True(t, exists, "Long URL2 should exist in LongURLMap")
	assert.Equal(t, shortCode2, retrievedShortCode2, "Retrieved short code2 should match the input short code2")
}

func TestGetURL(t *testing.T) {
	store := NewStorage()

	// Add a URL to storage
	url := "https://www.openai.com"
	shortCode := "openai1"
	store.AddURL(url, shortCode, time.Time{})

	// Test case: Retrieve existing URL
	urlModel, exists := store.GetURL(shortCode)
	assert.True(t, exists, "Short code should exist in storage")
	assert.Equal(t, url, urlModel.LongURL, "Long URL should match")
	assert.Equal(t, 0, urlModel.AccessCount, "Access count should be 0")
	assert.True(t, urlModel.ExpiresAt.IsZero(), "ExpiresAt should be zero")

	// Test case: Retrieve non-existent URL
	_, exists = store.GetURL("nonexist")
	assert.False(t, exists, "Non-existent short code should not exist in storage")
}

func TestGetShortCode(t *testing.T) {
	store := NewStorage()

	// Add URLs to storage
	url1 := "https://www.github.com"
	shortCode1 := "ghub1"
	store.AddURL(url1, shortCode1, time.Time{})

	url2 := "https://www.stackoverflow.com"
	shortCode2 := "so1"
	store.AddURL(url2, shortCode2, time.Time{})

	// Test case: Retrieve existing short codes by long URLs
	retrievedShortCode1, exists := store.GetShortCode(url1)
	assert.True(t, exists, "Long URL1 should exist in storage")
	assert.Equal(t, shortCode1, retrievedShortCode1, "Retrieved short code1 should match")

	retrievedShortCode2, exists := store.GetShortCode(url2)
	assert.True(t, exists, "Long URL2 should exist in storage")
	assert.Equal(t, shortCode2, retrievedShortCode2, "Retrieved short code2 should match")

	// Test case: Retrieve non-existent short code by long URL
	_, exists = store.GetShortCode("https://www.nonexistent.com")
	assert.False(t, exists, "Non-existent long URL should not exist in storage")
}

func TestDeleteURL(t *testing.T) {
	store := NewStorage()

	// Add URLs to storage
	url1 := "https://www.reddit.com"
	shortCode1 := "reddit1"
	store.AddURL(url1, shortCode1, time.Time{})

	url2 := "https://www.medium.com"
	shortCode2 := "medium1"
	store.AddURL(url2, shortCode2, time.Time{})

	// Test case: Delete existing URL
	store.DeleteURL(shortCode1)

	// Verify that URLMap no longer has shortCode1
	_, exists := store.GetURL(shortCode1)
	assert.False(t, exists, "Short code1 should be deleted from URLMap")

	// Verify that LongURLMap no longer has url1
	_, exists = store.GetShortCode(url1)
	assert.False(t, exists, "Long URL1 should be deleted from LongURLMap")

	// Test case: Delete non-existent URL
	store.DeleteURL("nonexist") // Should not panic or cause errors

	// Ensure existing URL2 is still present
	urlModel2, exists := store.GetURL(shortCode2)
	assert.True(t, exists, "Short code2 should still exist in storage")
	assert.Equal(t, url2, urlModel2.LongURL, "Long URL2 should match")
}

func TestIncrementAccessCount(t *testing.T) {
	store := NewStorage()

	// Add a URL to storage
	url := "https://www.twitter.com"
	shortCode := "twit1"
	store.AddURL(url, shortCode, time.Time{})

	// Test case: Increment access count multiple times
	for i := 1; i <= 5; i++ {
		store.IncrementAccessCount(shortCode)
		urlModel, exists := store.GetURL(shortCode)
		assert.True(t, exists, "Short code should exist in storage")
		assert.Equal(t, i, urlModel.AccessCount, "Access count should increment correctly")
	}

	// Test case: Increment access count for non-existent short code
	store.IncrementAccessCount("nonexist") // Should not panic or cause errors
	// No assertion needed as it should silently fail
}

func TestCleanupExpiredURLs(t *testing.T) {
	store := NewStorage()

	// Current time
	now := time.Now()

	// Add URLs with and without expiration
	url1 := "https://www.unexpired.com"
	shortCode1 := "unexp1"
	store.AddURL(url1, shortCode1, now.Add(2*time.Hour)) // Expires in 2 hours

	url2 := "https://www.expired.com"
	shortCode2 := "expd1"
	store.AddURL(url2, shortCode2, now.Add(-1*time.Hour)) // Expired 1 hour ago

	url3 := "https://www.noexpiry.com"
	shortCode3 := "noexp1"
	store.AddURL(url3, shortCode3, time.Time{}) // No expiration

	// Perform cleanup
	store.CleanupExpiredURLs()

	// Verify that expired URL is removed
	_, exists := store.GetURL(shortCode2)
	assert.False(t, exists, "Expired short code should be removed from storage")

	// Verify that unexpired URL still exists
	urlModel1, exists := store.GetURL(shortCode1)
	assert.True(t, exists, "Unexpired short code should still exist in storage")
	assert.Equal(t, url1, urlModel1.LongURL, "Long URL1 should match")

	// Verify that no expiration is set for URL3
	urlModel3, exists := store.GetURL(shortCode3)
	assert.True(t, exists, "No-expiry short code should still exist in storage")
	assert.Equal(t, url3, urlModel3.LongURL, "Long URL3 should match")
	assert.True(t, urlModel3.ExpiresAt.IsZero(), "ExpiresAt should be zero for no-expiry URL")
}

func TestStorageConcurrency(t *testing.T) {
	store := NewStorage()

	// Define a set of URLs and short codes
	urls := []string{
		"https://www.site1.com",
		"https://www.site2.com",
		"https://www.site3.com",
		"https://www.site4.com",
		"https://www.site5.com",
	}

	shortCodes := []string{
		"s1",
		"s2",
		"s3",
		"s4",
		"s5",
	}

	// Add URLs to storage concurrently
	done := make(chan bool)
	for i := 0; i < len(urls); i++ {
		go func(i int) {
			store.AddURL(urls[i], shortCodes[i], time.Time{})
			done <- true
		}(i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < len(urls); i++ {
		<-done
	}

	// Verify all URLs are added correctly
	for i := 0; i < len(urls); i++ {
		urlModel, exists := store.GetURL(shortCodes[i])
		assert.True(t, exists, "Short code %s should exist in storage", shortCodes[i])
		assert.Equal(t, urls[i], urlModel.LongURL, "Long URL should match for short code %s", shortCodes[i])

		retrievedShortCode, exists := store.GetShortCode(urls[i])
		assert.True(t, exists, "Long URL %s should exist in LongURLMap", urls[i])
		assert.Equal(t, shortCodes[i], retrievedShortCode, "Retrieved short code should match for URL %s", urls[i])
	}

	// Increment access counts concurrently
	for i := 0; i < len(shortCodes); i++ {
		go func(i int) {
			store.IncrementAccessCount(shortCodes[i])
			done <- true
		}(i)
	}

	// Wait for all increments
	for i := 0; i < len(shortCodes); i++ {
		<-done
	}

	// Verify access counts
	for i := 0; i < len(shortCodes); i++ {
		urlModel, exists := store.GetURL(shortCodes[i])
		assert.True(t, exists, "Short code %s should exist in storage", shortCodes[i])
		assert.Equal(t, 1, urlModel.AccessCount, "Access count should be 1 for short code %s", shortCodes[i])
	}
}
