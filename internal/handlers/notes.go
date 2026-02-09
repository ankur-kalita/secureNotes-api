package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/securenotes/securenotes-api/internal/models"
	"github.com/securenotes/securenotes-api/internal/repository"
)

// Version is the application version (set during build)
var Version = "1.0.1"

// NotesHandler handles HTTP requests for notes
type NotesHandler struct {
	repo repository.NotesRepository
}

// NewNotesHandler creates a new NotesHandler
func NewNotesHandler(repo repository.NotesRepository) *NotesHandler {
	return &NotesHandler{repo: repo}
}

// Health handles GET /health
func (h *NotesHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   Version,
	})
}

// GetAllNotes handles GET /api/v1/notes
func (h *NotesHandler) GetAllNotes(c *gin.Context) {
	notes := h.repo.GetAll()
	c.JSON(http.StatusOK, gin.H{
		"data":  notes,
		"count": len(notes),
	})
}

// GetNoteByID handles GET /api/v1/notes/:id
func (h *NotesHandler) GetNoteByID(c *gin.Context) {
	id := c.Param("id")

	note, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Note not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": note})
}

// CreateNote handles POST /api/v1/notes
func (h *NotesHandler) CreateNote(c *gin.Context) {
	var req models.CreateNoteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Title and content are required",
		})
		return
	}

	note := h.repo.Create(req)
	c.JSON(http.StatusCreated, gin.H{"data": note})
}

// UpdateNote handles PUT /api/v1/notes/:id
func (h *NotesHandler) UpdateNote(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request body",
		})
		return
	}

	note, err := h.repo.Update(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Note not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": note})
}

// DeleteNote handles DELETE /api/v1/notes/:id
func (h *NotesHandler) DeleteNote(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Note not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
