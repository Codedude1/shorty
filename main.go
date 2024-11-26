// main.go

package main

import (
	"time"

	"github.com/Codedude1/shorty/handlers"
	"github.com/Codedude1/shorty/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the Gin router
	router := gin.Default()

	// Initialize the in-memory storage
	store := storage.NewStorage()

	// Start the cleanup goroutine to remove expired URLs periodically
	go func() {
		ticker := time.NewTicker(1 * time.Hour) // Adjust the duration as needed
		defer ticker.Stop()
		for {
			<-ticker.C
			store.CleanupExpiredURLs()
		}
	}()

	// Register routes
	router.POST("/shorten", handlers.ShortenURLHandler(store))
	router.GET("/stats/:shortCode", handlers.StatsHandler(store))
	router.GET("/:shortCode", handlers.RedirectHandler(store))

	// Start the server on port 8080
	router.Run(":8080")
}
