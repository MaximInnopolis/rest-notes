package api

import (
	"errors"
	"fmt"
	"log"

	"rest-notes/internal/app/models"
	"rest-notes/internal/app/repository"
)

var ErrSpell = errors.New("обнаружены орфографические ошибки")

// NoteService represents service for handling notes and checking spelling
type NoteService struct {
	repo           repository.NoteRepo
	spellerService *SpellerService
}

// NewNoteService creates new instance of NoteService with repository and spelling service
func NewNoteService(repo repository.NoteRepo, spellerService *SpellerService) *NoteService {
	return &NoteService{
		repo:           repo,
		spellerService: spellerService,
	}
}

// CreateNote creates new note using repository and returns created note
func (n *NoteService) CreateNote(note models.Note) (models.Note, error) {

	// Check spelling of note's description using SpellerService
	spellingResult, err := n.spellerService.CheckText(note.Description, "rus", 0)
	if err != nil {
		log.Printf("Ошибка при проверке орфографии: %v", err)
		return models.Note{}, err
	}

	// If there are spelling errors, format error message and return it
	if len(spellingResult) > 0 {
		log.Printf("Обнаружены орфографические ошибки в тексте: %v", spellingResult)

		errorDetails := formatSpellingErrors(spellingResult)

		// Return formatted error with details
		return models.Note{}, fmt.Errorf("%w: %s", ErrSpell, errorDetails)
	}

	// If no mistakes found, create note in repository
	createdNote, err := n.repo.Create(note)
	if err != nil {
		log.Printf("Ошибка при создании заметки: %v", err)
		return models.Note{}, err
	}

	log.Printf("Заметка успешно создана: %v", createdNote)
	return createdNote, nil
}

// GetNoteList retrieves list of all notes from repository
func (n *NoteService) GetNoteList(userID int) ([]models.Note, error) {
	notes, err := n.repo.GetAll(userID)
	if err != nil {
		log.Printf("Ошибка при получении списка заметок: %v", err)
		return nil, err
	}

	log.Printf("Список заметок успешно получен для пользователя ID: %d", userID)
	return notes, nil
}

// formatSpellingErrors formats the list of spelling errors
// Takes spelling check result and returns string with errors and suggestions
func formatSpellingErrors(spellingResult []map[string]interface{}) string {
	var errorDetails string
	for _, err := range spellingResult {
		word := err["word"].(string)
		suggestions := err["s"].([]interface{})
		errorDetails += fmt.Sprintf("Слово: %s, предложения: %v\n", word, suggestions)
	}
	return errorDetails
}
