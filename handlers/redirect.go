package handlers

import (
	"net/http"
	"time"

	"github.com/Codedude1/shorty/storage"
	"github.com/Codedude1/shorty/utils"
	"github.com/gin-gonic/gin"
)


func RedirectHandler(store *storage.Storage) gin.HandlerFunc {
    return func(c *gin.Context) {
        shortCode := c.Param("shortCode")

        store.Mu.RLock()
        urlModel, exists := store.UrlMap[shortCode]
        store.Mu.RUnlock()

        if !exists {
            utils.RespondWithError(c, http.StatusNotFound, "Short URL not found")
            return
        }

        // Checking for expiration
        if !urlModel.ExpiresAt.IsZero() && time.Now().After(urlModel.ExpiresAt) {
            // Removing expired URL from storage
            store.Mu.Lock()
            delete(store.UrlMap, shortCode)
            delete(store.LongURLMap, urlModel.LongURL)
            store.Mu.Unlock()

            utils.RespondWithError(c, http.StatusGone, "Short URL has expired")
            return
        }

        // Increment access count
        store.Mu.Lock()
        urlModel.AccessCount++
        store.Mu.Unlock()

        // Redirect to the original long URL
        c.Redirect(http.StatusFound, urlModel.LongURL)
    }
}
