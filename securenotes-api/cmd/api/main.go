package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/securenotes/securenotes-api/internal/handlers"
	"github.com/securenotes/securenotes-api/internal/middleware"
	"github.com/securenotes/securenotes-api/internal/repository"
)

func main() {
	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set Gin mode
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "release"
	}
	gin.SetMode(ginMode)

	// Initialize repository
	repo := repository.NewInMemoryNotesRepository()

	// Initialize handlers
	notesHandler := handlers.NewNotesHandler(repo)

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	// Health endpoint
	router.GET("/health", notesHandler.Health)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		notes := v1.Group("/notes")
		{
			notes.GET("", notesHandler.GetAllNotes)
			notes.GET("/:id", notesHandler.GetNoteByID)
			notes.POST("", notesHandler.CreateNote)
			notes.PUT("/:id", notesHandler.UpdateNote)
			notes.DELETE("/:id", notesHandler.DeleteNote)
		}
	}

	// Start server
	log.Printf("Starting SecureNotes API on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
