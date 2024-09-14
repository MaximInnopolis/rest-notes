package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"rest-notes/internal/app/api"
	"rest-notes/internal/app/models"
)

// CreateNoteHandler handles HTTP POST request to create new note
// It parses request body, validates data, and calls service to create note
func (h *Handler) CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("UserID").(int)
	if !ok {
		http.Error(w, "Ошибка аутентификации", http.StatusUnauthorized)
		return
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"`
	}

	// Decode request body into input struct
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неправильный формат данных", http.StatusBadRequest)
		return
	}

	// Parse due date from string to time.Time format
	dueDate, err := time.Parse(time.RFC3339, input.DueDate)
	if err != nil {
		http.Error(w, "Неправильный формат даты", http.StatusBadRequest)
		return
	}

	// Create new note object
	note := models.Note{
		//
		UserID: userID,
		//
		Title:       input.Title,
		Description: input.Description,
		DueDate:     dueDate,
	}

	// Call service to create note
	createdNote, err := h.service.CreateNote(note)
	if err != nil {
		if strings.Contains(err.Error(), api.ErrSpell.Error()) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Проблема на сервере", http.StatusInternalServerError)
		return
	}

	// Respond with created note
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdNote)
}

// GetNoteListHandler handles HTTP GET request to retrieve all notes
// It calls service to get list of notes and returns them in response
func (h *Handler) GetNoteListHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("UserID").(int)
	if !ok {
		http.Error(w, "Ошибка аутентификации", http.StatusUnauthorized)
		return
	}

	notes, err := h.service.GetNoteList(userID)
	if err != nil {
		http.Error(w, "Проблема на сервере", http.StatusInternalServerError)
		return
	}

	// Respond with list of notes
	w.Header().Set("Content-Type", "application/json")
	if len(notes) == 0 {
		json.NewEncoder(w).Encode("Список заметок пуст")
		return
	}
	json.NewEncoder(w).Encode(notes)
}
