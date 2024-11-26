package storage

import (
	"sync"
	"time"

	"github.com/Codedude1/shorty/models"
)

type Storage struct {
	Mu         sync.RWMutex
	UrlMap     map[string]*models.URL
	LongURLMap map[string]string
	IdCounter  int64
}

func NewStorage() *Storage {
	return &Storage{
		UrlMap:     make(map[string]*models.URL),
		LongURLMap: make(map[string]string),
		IdCounter:  1,
	}
}

func (s *Storage) CleanupExpiredURLs() {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	now := time.Now()
	for shortCode, urlModel := range s.UrlMap {
		if !urlModel.ExpiresAt.IsZero() && now.After(urlModel.ExpiresAt) {
			delete(s.UrlMap, shortCode)
			delete(s.LongURLMap, urlModel.LongURL)
		}
	}
}
