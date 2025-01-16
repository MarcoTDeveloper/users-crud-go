package api

import (
	"net/http"
	"users-crud/methods"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func ApiHandler(db methods.UserDatabase) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Post("/api/users", handleCreateUser(db))
	r.Get("/api/users", handleGetAllUsers(db))
	r.Get("/api/users/{id}", handleGetUser(db))
	r.Delete("/api/users/{id}", handleDeleteUser(db))
	r.Put("/api/users/{id}", handleUpdateUser(db))

	return r
}
