package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/securenotes/securenotes-api/internal/models"
)

var (
	ErrNoteNotFound = errors.New("note not found")
)

// NotesRepository defines the interface for notes data operations
type NotesRepository interface {
	GetAll() []models.Note
	GetByID(id string) (*models.Note, error)
	Create(req models.CreateNoteRequest) *models.Note
	Update(id string, req models.UpdateNoteRequest) (*models.Note, error)
	Delete(id string) error
}

// InMemoryNotesRepository implements NotesRepository with in-memory storage
type InMemoryNotesRepository struct {
	notes map[string]models.Note
	mu    sync.RWMutex
}

// NewInMemoryNotesRepository creates a new in-memory repository
func NewInMemoryNotesRepository() *InMemoryNotesRepository {
	return &InMemoryNotesRepository{
		notes: make(map[string]models.Note),
	}
}

// GetAll returns all notes
func (r *InMemoryNotesRepository) GetAll() []models.Note {
	r.mu.RLock()
	defer r.mu.RUnlock()

	notes := make([]models.Note, 0, len(r.notes))
	for _, note := range r.notes {
		notes = append(notes, note)
	}
	return notes
}

// GetByID returns a note by its ID
func (r *InMemoryNotesRepository) GetByID(id string) (*models.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	note, exists := r.notes[id]
	if !exists {
		return nil, ErrNoteNotFound
	}
	return &note, nil
}

// Create creates a new note
func (r *InMemoryNotesRepository) Create(req models.CreateNoteRequest) *models.Note {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	note := models.Note{
		ID:        uuid.New().String(),
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	r.notes[note.ID] = note
	return &note
}

// Update updates an existing note
func (r *InMemoryNotesRepository) Update(id string, req models.UpdateNoteRequest) (*models.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	note, exists := r.notes[id]
	if !exists {
		return nil, ErrNoteNotFound
	}

	if req.Title != "" {
		note.Title = req.Title
	}
	if req.Content != "" {
		note.Content = req.Content
	}
	note.UpdatedAt = time.Now()

	r.notes[id] = note
	return &note, nil
}

// Delete removes a note by its ID
func (r *InMemoryNotesRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.notes[id]; !exists {
		return ErrNoteNotFound
	}

	delete(r.notes, id)
	return nil
}
