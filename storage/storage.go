package storage

import (
	"sync"
	"time"

	"github.com/Codedude1/shorty/models"
)

// Storage defines the in-memory storage structure.
type Storage struct {
	Mu         sync.RWMutex
	URLMap     map[string]*models.URL
	LongURLMap map[string]string
}

// NewStorage initializes and returns a new Storage instance.
func NewStorage() *Storage {
	return &Storage{
		URLMap:     make(map[string]*models.URL),
		LongURLMap: make(map[string]string),
	}
}

// AddURL adds a new URL mapping to the storage.
func (s *Storage) AddURL(url string, shortCode string, expiresAt time.Time) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.URLMap[shortCode] = &models.URL{
		BaseURL: models.BaseURL{
			LongURL:     url,
			AccessCount: 0,
			CreatedAt:   time.Now(),
			ExpiresAt:   expiresAt,
		},
		ShortCode: shortCode,
	}
	s.LongURLMap[url] = shortCode
}

// GetURL retrieves a URL model by its short code.
func (s *Storage) GetURL(shortCode string) (*models.URL, bool) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	urlModel, exists := s.URLMap[shortCode]
	return urlModel, exists
}

// GetShortCode retrieves the short code for a given long URL.
func (s *Storage) GetShortCode(url string) (string, bool) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	shortCode, exists := s.LongURLMap[url]
	return shortCode, exists
}

// DeleteURL removes a URL mapping from the storage.
func (s *Storage) DeleteURL(shortCode string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if urlModel, exists := s.URLMap[shortCode]; exists {
		delete(s.URLMap, shortCode)
		delete(s.LongURLMap, urlModel.LongURL)
	}
}

// IncrementAccessCount increments the access count for a given short code.
func (s *Storage) IncrementAccessCount(shortCode string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if urlModel, exists := s.URLMap[shortCode]; exists {
		urlModel.AccessCount++
	}
}

// CleanupExpiredURLs removes expired URLs from the storage.
func (s *Storage) CleanupExpiredURLs() {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	now := time.Now()
	for shortCode, urlModel := range s.URLMap {
		if !urlModel.ExpiresAt.IsZero() && now.After(urlModel.ExpiresAt) {
			delete(s.URLMap, shortCode)
			delete(s.LongURLMap, urlModel.LongURL)
		}
	}
}
