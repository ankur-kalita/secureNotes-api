package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/securenotes/securenotes-api/internal/handlers"
	"github.com/securenotes/securenotes-api/internal/models"
	"github.com/securenotes/securenotes-api/internal/repository"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	repo := repository.NewInMemoryNotesRepository()
	handler := handlers.NewNotesHandler(repo)

	router.GET("/health", handler.Health)
	router.GET("/api/v1/notes", handler.GetAllNotes)
	router.GET("/api/v1/notes/:id", handler.GetNoteByID)
	router.POST("/api/v1/notes", handler.CreateNote)
	router.PUT("/api/v1/notes/:id", handler.UpdateNote)
	router.DELETE("/api/v1/notes/:id", handler.DeleteNote)

	return router
}

func TestHealthEndpoint(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.HealthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response.Status)
	}
}

func TestGetAllNotesEmpty(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notes", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	count := int(response["count"].(float64))
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

func TestCreateNote(t *testing.T) {
	router := setupRouter()

	noteData := models.CreateNoteRequest{
		Title:   "Test Note",
		Content: "This is a test note content",
	}
	jsonData, _ := json.Marshal(noteData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notes", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	data := response["data"].(map[string]interface{})
	if data["title"] != noteData.Title {
		t.Errorf("Expected title '%s', got '%s'", noteData.Title, data["title"])
	}
	if data["content"] != noteData.Content {
		t.Errorf("Expected content '%s', got '%s'", noteData.Content, data["content"])
	}
}

func TestCreateNoteValidation(t *testing.T) {
	router := setupRouter()

	// Test with missing title
	noteData := map[string]string{
		"content": "Content without title",
	}
	jsonData, _ := json.Marshal(noteData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notes", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetNoteNotFound(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notes/nonexistent-id", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUpdateNote(t *testing.T) {
	router := setupRouter()

	// First create a note
	createData := models.CreateNoteRequest{
		Title:   "Original Title",
		Content: "Original Content",
	}
	jsonData, _ := json.Marshal(createData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notes", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	noteID := createResponse["data"].(map[string]interface{})["id"].(string)

	// Now update the note
	updateData := models.UpdateNoteRequest{
		Title:   "Updated Title",
		Content: "Updated Content",
	}
	jsonData, _ = json.Marshal(updateData)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/notes/"+noteID, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var updateResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &updateResponse)
	data := updateResponse["data"].(map[string]interface{})

	if data["title"] != updateData.Title {
		t.Errorf("Expected title '%s', got '%s'", updateData.Title, data["title"])
	}
}

func TestDeleteNote(t *testing.T) {
	router := setupRouter()

	// First create a note
	createData := models.CreateNoteRequest{
		Title:   "Note to Delete",
		Content: "This will be deleted",
	}
	jsonData, _ := json.Marshal(createData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notes", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	noteID := createResponse["data"].(map[string]interface{})["id"].(string)

	// Now delete the note
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/notes/"+noteID, nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify it's deleted
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/notes/"+noteID, nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d after deletion, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteNoteNotFound(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/notes/nonexistent-id", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
