package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"users-crud/methods"

	"github.com/go-chi/chi/v5"
)

type UserMethods interface {
	FindAll() []methods.User
	FindById(uuid string) (methods.User, error)
	Insert(firstName, lastName, bio string) (methods.User, error)
	Update(uuid, firstName, lastName, bio string) (methods.User, error)
	Delete(uuid string) (methods.User, error)
}

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func sendJSON(w http.ResponseWriter, res Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(res)
	if err != nil {
		slog.Error("error in json marshal", "error", err)
		sendJSON(w, Response{Error: "something went wrong"}, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	_, writeErr := w.Write(data)
	if writeErr != nil {
		slog.Error("error to send response", "error", err)
		return
	}
}

type CreateUserBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
}

func handleCreateUser(db methods.UserDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUserBody CreateUserBody
		err := json.NewDecoder(r.Body).Decode(&newUserBody)

		if err != nil {
			sendJSON(w, Response{Error: "Invalid body"}, http.StatusUnprocessableEntity)
			return
		}

		if newUserBody.FirstName == "" {
			sendJSON(w, Response{Error: "Please provide first_name for the user"}, http.StatusBadRequest)
			return
		}

		if newUserBody.LastName == "" {
			sendJSON(w, Response{Error: "Please provide last_name for the user"}, http.StatusBadRequest)
			return
		}

		if newUserBody.Biography == "" {
			sendJSON(w, Response{Error: "Please provide biography for the user"}, http.StatusBadRequest)
			return
		}

		newUser, insertErr := db.Insert(newUserBody.FirstName, newUserBody.LastName, newUserBody.Biography)

		if insertErr != nil {
			sendJSON(w, Response{Error: "There was an error while saving the user to the database"}, http.StatusInternalServerError)
			return
		}

		sendJSON(w, Response{Data: newUser}, http.StatusCreated)
	}
}

func handleGetAllUsers(db methods.UserDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := db.FindAll()

		sendJSON(w, Response{Data: users}, http.StatusOK)
	}
}

func handleGetUser(db methods.UserDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		user, err := db.FindById(id)

		if err != nil {
			sendJSON(w, Response{Error: "The user with the specified ID does not exist"}, http.StatusNotFound)
			return
		}

		sendJSON(w, Response{Data: user}, http.StatusOK)
	}
}

func handleDeleteUser(db methods.UserDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		user, err := db.Delete(id)

		if err != nil {
			sendJSON(w, Response{Error: "The user with the specified ID does not exist"}, http.StatusNotFound)
			return
		}

		sendJSON(w, Response{Data: user}, http.StatusOK)
	}
}

type UpdateUserBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
}

func handleUpdateUser(db methods.UserDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var newUserBody UpdateUserBody

		err := json.NewDecoder(r.Body).Decode(&newUserBody)

		if err != nil {
			sendJSON(w, Response{Error: "invalid body"}, http.StatusInternalServerError)
			return
		}

		if newUserBody.FirstName == "" {
			sendJSON(w, Response{Error: "Please provide first_name for the user"}, http.StatusBadRequest)
			return
		}

		if newUserBody.LastName == "" {
			sendJSON(w, Response{Error: "Please provide last_name for the user"}, http.StatusBadRequest)
			return
		}

		if newUserBody.Biography == "" {
			sendJSON(w, Response{Error: "Please provide biography for the user"}, http.StatusBadRequest)
			return
		}

		newUser, updateErr := db.Update(id, newUserBody.FirstName, newUserBody.LastName, newUserBody.Biography)

		if updateErr != nil {
			sendJSON(w, Response{Error: "The user with the specified ID does not exist"}, http.StatusNotFound)
			return
		}

		sendJSON(w, Response{Data: newUser}, http.StatusOK)
	}
}
