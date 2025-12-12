package main

import (
	"log/slog"
	"main/api"
	"main/types"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Failed to execute code.", "err", err)
	}
	slog.Info("All systems offline.")
}

func run() error {
	db := make(map[uuid.UUID]types.User)

	handler := api.NewHandler(db)
	
	s := http.Server{
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
		WriteTimeout: 10 * time.Second,
		Addr:         ":8080",
		Handler:      handler,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}