package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
	"users-crud/api"
	"users-crud/methods"

	"github.com/google/uuid"
)

func main() {
	err := run()

	if err != nil {
		slog.Error("failed to execute code", "error", err)
		os.Exit(1)
		return
	}

	slog.Info("all system off")
}

func run() error {
	id := uuid.New().String()

	db := methods.UserDatabase{
		id: {
			Id:        id,
			FirstName: "First User",
			LastName:  "Admin",
			Biography: "This is the owner of the project!",
		},
	}

	handler := api.ApiHandler(db)

	s := http.Server{
		Addr:         ":3000",
		Handler:      handler,
		IdleTimeout:  time.Minute,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	err := s.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
