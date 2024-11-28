// main.go

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Codedude1/shorty/handlers"
	"github.com/Codedude1/shorty/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Initialize the Gin router
	router := gin.Default()

	// Initialize the in-memory storage
	store := storage.NewStorage()

	// Register routes
	router.POST("/shorten", handlers.ShortenURLHandler(store))
	router.GET("/stats/:shortCode", handlers.StatsHandler(store))
	router.GET("/:shortCode", handlers.RedirectHandler(store))

	// Determine server port from environment variable or default to 8081
	port := getEnv("PORT", "8081")

	// Create an HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start the cleanup goroutine to remove expired URLs periodically
	go func() {
		cleanupInterval := getEnvAsDuration("CLEANUP_INTERVAL", time.Hour)
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		for {
			<-ticker.C
			store.CleanupExpiredURLs()
			log.Println("[INFO] Cleanup of expired URLs completed.")
		}
	}()

	// Start the server in a separate goroutine
	go func() {
		log.Printf("[INFO] Server is running on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR] ListenAndServe(): %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("[INFO] Shutdown signal received.")

	// Create a deadline to wait for ongoing requests to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[ERROR] Server forced to shutdown: %v", err)
	}

	log.Println("[INFO] Server exiting.")
}

// getEnv retrieves the value of the environment variable named by the key.
// It returns the defaultValue if the variable is not present.
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsDuration retrieves the value of the environment variable named by the key
// and parses it as a time.Duration. It returns the defaultValue if the variable is not present or invalid.
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := time.ParseDuration(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
