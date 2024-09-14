package repository

import (
	"rest-notes/internal/app/models"
	"rest-notes/internal/app/repository/database"
	"rest-notes/internal/app/repository/postgresql"
)

// UserRepo defines interface for user-related database operations
type UserRepo interface {
	Create(user models.User) error
	Get(user models.User) (models.User, error)
}

// NoteRepo defines interface for note-related database operations
type NoteRepo interface {
	Create(note models.Note) (models.Note, error)
	GetAll(userID int) ([]models.Note, error)
}

// Repository combines UserRepo and NoteRepo interfaces into single struct
type Repository struct {
	UserRepo
	NoteRepo
}

// New initializes and returns new Repository instance with PostgreSQL implementations for UserRepo and NoteRepo
func New(db database.Database) *Repository {
	return &Repository{
		UserRepo: postgresql.NewUserPostgres(db),
		NoteRepo: postgresql.NewNotePostgres(db),
	}
}
