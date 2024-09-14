package postgresql

import (
	"context"

	"rest-notes/internal/app/models"
	"rest-notes/internal/app/repository/database"
)

// NotePostgres is repository implementation for managing notes in PostgreSQL database
type NotePostgres struct {
	db database.Database
}

// NewNotePostgres creates new NotePostgres instance with given database connection
func NewNotePostgres(db database.Database) *NotePostgres {
	return &NotePostgres{db: db}
}

// Create inserts new note into database and returns created note with its ID and timestamps
func (n *NotePostgres) Create(note models.Note) (models.Note, error) {
	query := `INSERT INTO notes (user_id, title, description, due_date, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id, created_at, updated_at`
	ctx := context.Background()

	// Execute query and scan returned ID, created_at, and updated_at into note object
	err := n.db.GetPool().QueryRow(ctx, query, note.UserID, note.Title, note.Description, note.DueDate).
		Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return models.Note{}, err
	}
	return note, nil
}

// GetAll retrieves all notes for a given user from the database
func (n *NotePostgres) GetAll(userID int) ([]models.Note, error) {
	query := `SELECT id, user_id, title, description, due_date, created_at, updated_at FROM notes WHERE user_id = $1`
	var notes []models.Note
	ctx := context.Background()

	// Execute query and iterate over result rows
	rows, err := n.db.GetPool().Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan each row into Note object and append to notes slice
	for rows.Next() {
		var note models.Note
		err = rows.Scan(
			&note.ID, &note.UserID, &note.Title, &note.Description, &note.DueDate, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	// Check for any error that occurred during iteration over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}
