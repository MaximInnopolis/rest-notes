package main

import (
	"fmt"
	"os"

	"rest-notes/internal/app/api"
	"rest-notes/internal/app/config"
	httpHandler "rest-notes/internal/app/http"
	"rest-notes/internal/app/repository"
	"rest-notes/internal/app/repository/database"
)

func main() {

	// Create config
	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create a new connection pool to database
	pool, err := database.NewPool(cfg.DbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer pool.Close()

	// Create a new Database with connection pool
	db := database.NewDatabase(pool)

	// Create a new repo with Database
	repo := repository.New(*db)

	// Create a new service
	service := api.New(repo)

	// Create Http handler
	handler := httpHandler.New(*service)

	// Start server
	handler.StartServer(cfg.HttpPort)
}
