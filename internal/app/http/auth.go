package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"rest-notes/internal/app/api"
	"rest-notes/internal/app/models"
)

// RegisterUserHandler handles user registration requests
// It parses request body to get username and password, creates new user
// and responds with appropriate HTTP status code based on result.
func (h *Handler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"username"`
		Password string `json:"password"`
	}

	// Decode request body into input struct
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неправильный формат данных", http.StatusBadRequest)
		return
	}

	user := models.User{
		Name:     input.Name,
		Password: input.Name,
	}

	// Attempt to create user using service
	err := h.service.Authorization.CreateUser(user)
	if err != nil {
		if errors.Is(err, api.ErrUserAlreadyExists) {
			http.Error(w, "Пользователь с этим именем уже существует", http.StatusConflict)
			return
		}

		http.Error(w, "Проблема на сервере", http.StatusInternalServerError)
		return
	}

	// Respond with created user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// LoginUserHandler handles user login requests
// It parses request body to get username and password, generates token
// and responds with token if successful
func (h *Handler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"username"`
		Password string `json:"password"`
	}

	// Decode request body into input struct
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неправильный формат данных", http.StatusBadRequest)
		return
	}

	user := models.User{
		Name:     input.Name,
		Password: input.Name,
	}

	// Attempt to generate token using service
	token, err := h.service.Authorization.GenerateToken(user)
	if err != nil {
		http.Error(w, "Проблема на сервере", http.StatusInternalServerError)
		return
	}
	// Respond with created user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}
