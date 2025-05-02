package main

import (
	"chatgpt-wrapper/config"
	"chatgpt-wrapper/handlers"
	"chatgpt-wrapper/middleware"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Load configuration
	cfg := config.NewConfig()

	// Load banned words configuration
	config.LoadBannedWords()

	// Initialize router
	r := gin.Default()

	// Configure CORS
	r.Use(middleware.CORSMiddleware())

	// Initialize handlers with configuration
	handlers.InitHandlers(cfg)

	// Routes
	r.POST("/api/chat/stream", handlers.HandleStream)
	r.POST("/api/chat/title", handlers.HandleTitle)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
