package api

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"main/types"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

func NewHandler(db map[uuid.UUID]types.User) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Get("/users", handleListUsers(db))
	r.Get("/user/{id}", handleGetUserById(db))
	r.Delete("/user/{id}", handleDeleteUserById(db))
	r.Put("/user/{id}", handlePutUser(db))
	r.Post("/user", handlePostUser(db))

	return r
}

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func sendJSON(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("failed to marshal json data", "error", err)
		sendJSON(
			w,
			Response{Error: "something went wrong"},
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write json data", "error", err)
		return
	}
}

func handlePostUser(db map[uuid.UUID]types.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 10000)
		data, err := io.ReadAll(r.Body)
		if err != nil {
			var maxErr *http.MaxBytesError
			if errors.As(err, &maxErr) {
				sendJSON(w, Response{Error: "Body too large."}, http.StatusRequestEntityTooLarge)
				return 
			}

			slog.Error("Falha ao ler o JSON do usuário.", "error", err)
			sendJSON(w, Response{Error: "Something went wrong."}, http.StatusInternalServerError)
			return
		}

		var user types.User
		if err := json.Unmarshal(data, &user); err != nil {
			sendJSON(w, Response{Error: "Invalid body."}, http.StatusUnprocessableEntity)
			return
		}

		if len(user.Biography) < 20 || len(user.Biography) > 450 {
			sendJSON(w, Response{Error: "Biography must have between 20 to 450 carachters."}, http.StatusBadRequest)
			return
		}

		if len(user.FirstName) < 2 || len(user.FirstName) > 20 {
			sendJSON(w, Response{Error: "First name must have between 2 to 20 carachters."}, http.StatusBadRequest)
			return
		}

		if len(user.LastName) < 2 || len(user.LastName) > 20 {
			sendJSON(w, Response{Error: "Last name must have between 2 to 20 carachters."}, http.StatusBadRequest)
			return
		}

		id, err := uuid.NewRandom()
		if err != nil {
			sendJSON(w, Response{Error: "Could not generate uuid."}, http.StatusInternalServerError)
			return
		}

		user.ID = id.String()
		db[id] = user
		sendJSON(w, Response{Data: user}, http.StatusCreated)
	}
}

func handleGetUserById(db map[uuid.UUID]types.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			sendJSON(w, Response{Error: "Invalid UUID."}, http.StatusBadRequest)
			return
		}

		user, ok := db[id]
		if !ok {
			sendJSON(w, Response{Error: "User not found."}, http.StatusNotFound)
			return 
		}

		sendJSON(w, Response{Data: user}, http.StatusOK)
	}
}

func handleDeleteUserById(db map[uuid.UUID]types.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			sendJSON(w, Response{Error: "Invalid UUID."}, http.StatusBadRequest)
			return 
		}

		user, ok := db[id]
		if !ok {
			sendJSON(w, Response{Error: "User not found."}, http.StatusNotFound)
			return 
		}

		delete(db, id)
		sendJSON(w, Response{Data: user}, http.StatusOK)
	}
}

func handleListUsers(db map[uuid.UUID]types.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := make([]types.User, 0, len(db))

		for _, user := range db {
			users = append(users, user)
		}
		sendJSON(w, Response{Data: users}, http.StatusOK)
	}
}

func handlePutUser(db map[uuid.UUID]types.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 10000)
		data, err := io.ReadAll(r.Body)
		if err != nil {
			var maxErr *http.MaxBytesError
			if errors.As(err, &maxErr) {
				sendJSON(w, Response{Error: "Body too large."}, http.StatusRequestEntityTooLarge)
				return 
			}

			slog.Error("Falha ao ler o JSON do usuário.", "error", err)
			sendJSON(w, Response{Error: "Something went wrong."}, http.StatusInternalServerError)
			return
		}

		var userUpdated types.User
		if err := json.Unmarshal(data, &userUpdated); err != nil {
			sendJSON(w, Response{Error: "Invalid body."}, http.StatusUnprocessableEntity)
			return
		}

		if len(userUpdated.Biography) < 20 || len(userUpdated.Biography) > 450 {
			sendJSON(w, Response{Error: "Biography must have between 20 to 450 carachters."}, http.StatusBadRequest)
			return
		}

		if len(userUpdated.FirstName) < 2 || len(userUpdated.FirstName) > 20 {
			sendJSON(w, Response{Error: "First name must have between 2 to 20 carachters."}, http.StatusBadRequest)
			return
		}

		if len(userUpdated.LastName) < 2 || len(userUpdated.LastName) > 20 {
			sendJSON(w, Response{Error: "Last name must have between 2 to 20 carachters."}, http.StatusBadRequest)
			return
		}

		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			sendJSON(w, Response{Error: "Invalid UUID."}, http.StatusBadRequest)
			return
		}

		_, ok := db[id]
		if !ok {
			sendJSON(w, Response{Error: "User not found."}, http.StatusNotFound)
			return 
		}

		userUpdated.ID = id.String()
		db[id] = userUpdated
		sendJSON(w, Response{Data: userUpdated}, http.StatusOK)
	}
}