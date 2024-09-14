package api

import (
	"github.com/dgrijalva/jwt-go"
	"rest-notes/internal/app/models"
	"rest-notes/internal/app/repository"
)

// Authorization defines interface for user authentication and authorization logic
type Authorization interface {
	CreateUser(user models.User) error
	GenerateToken(user models.User) (string, error)
	IsTokenValid(tokenString string) (bool, jwt.MapClaims, error)
}

// Note defines interface for note-related operations
type Note interface {
	CreateNote(note models.Note) (models.Note, error)
	GetNoteList(userID int) ([]models.Note, error)
}

// Service aggregates Authorization and Note interfaces
// It combines business logic for user authentication and note management
type Service struct {
	Authorization
	Note
}

// New returns new instance of Service, initializing dependencies
// It takes repository that holds database access logic and initializes speller service
func New(repo *repository.Repository) *Service {
	spellerService := NewSpellerService()

	return &Service{
		Authorization: NewAuthService(repo.UserRepo),
		Note:          NewNoteService(repo.NoteRepo, spellerService),
	}
}
