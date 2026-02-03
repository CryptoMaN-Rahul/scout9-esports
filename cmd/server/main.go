package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"scout9/internal/api"
	"scout9/pkg/cache"
	"scout9/pkg/grid"
	"scout9/pkg/llm"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize Redis cache
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	var gridCache grid.Cache
	var cacheClient *cache.RedisCache
	var err error
	cacheClient, err = cache.NewRedisCache(redisURL)
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v. Running without cache.", err)
		cacheClient = nil
	} else {
		gridCache = cacheClient
	}

	// Initialize GRID API client
	gridAPIKey := os.Getenv("GRID_API_KEY")
	if gridAPIKey == "" {
		log.Fatal("GRID_API_KEY environment variable is required")
	}
	gridClient := grid.NewClient(gridAPIKey, gridCache)

	// Initialize LLM service
	llmAPIKey := os.Getenv("LLM_API_KEY")
	var llmService llm.Service
	if llmAPIKey != "" {
		llmService = llm.NewOpenAIService(llmAPIKey)
	} else {
		log.Println("LLM_API_KEY not set, using template-based generation")
		llmService = llm.NewTemplateService()
	}

	// Create API router
	router := api.NewRouter(gridClient, llmService, cacheClient)

	// Configure server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second, // Longer for report generation
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 SCOUT9 server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
