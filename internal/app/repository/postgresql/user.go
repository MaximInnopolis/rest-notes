package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"rest-notes/internal/app/models"
	"rest-notes/internal/app/repository/database"
)

var ErrNotFound = errors.New("user not found")

// UserPostgres implements the UserRepo interface for PostgreSQL database operations related to users
type UserPostgres struct {
	db database.Database
}

// NewUserPostgres creates new UserPostgres instance with provided database connection
func NewUserPostgres(db database.Database) *UserPostgres {
	return &UserPostgres{db: db}
}

// Create inserts new user into the users table and returns error if operation fails
func (up *UserPostgres) Create(user models.User) error {
	query := `INSERT INTO users (username, password, created_at) 
	          VALUES ($1, $2, NOW()) RETURNING id, created_at`
	ctx := context.Background()

	// Execute query and scan returned ID, created_at, and updated_at into note object
	err := up.db.GetPool().QueryRow(ctx, query, user.Name, user.Password).
		Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Get retrieves user from the users table by username and verifies provided password
// Returns user and error if user is not found or if password is incorrect
func (up *UserPostgres) Get(user models.User) (models.User, error) {
	query := `SELECT id, username, password, created_at FROM users WHERE username = $1`
	var dbUser models.User

	ctx := context.Background()

	err := up.db.GetPool().QueryRow(ctx, query, user.Name).
		Scan(&dbUser.ID, &dbUser.Name, &dbUser.Password, &dbUser.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, err
	}

	// Compare provided password with hashed password stored in database
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return models.User{}, errors.New("invalid password")
	}

	return dbUser, nil
}
