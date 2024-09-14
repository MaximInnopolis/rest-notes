package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"rest-notes/internal/app/api"
)

// Handler struct holds service used for handling requests
type Handler struct {
	service api.Service
}

// New creates new Handler instance with provided service
func New(service api.Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers HTTP routes
func (h *Handler) RegisterRoutes(r *mux.Router) {
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", h.RegisterUserHandler).Methods("POST")
	authRouter.HandleFunc("/login", h.LoginUserHandler).Methods("POST")

	noteRouter := r.PathPrefix("/notes").Subrouter()
	createNoteRouter := http.HandlerFunc(h.CreateNoteHandler)
	noteRouter.Handle("/new", h.RequireValidTokenMiddleware(createNoteRouter)).Methods("POST")

	getNotesRouter := http.HandlerFunc(h.GetNoteListHandler)
	noteRouter.Handle("/list", h.RequireValidTokenMiddleware(getNotesRouter)).Methods("GET")
}

// StartServer initializes and starts HTTP server on given port
func (h *Handler) StartServer(port string) {
	router := mux.NewRouter()

	// Middleware for processing request ID
	router.Use(h.RequestIDMiddleware)
	h.RegisterRoutes(router)

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
