package main

import (
	"context"
	"fmt"
	"log/slog"
	"main/api"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)	


func main() {
	if err := run(); err != nil {
		slog.Error("Failed to execute code.", "err", err)
	}
	slog.Info("All systems offline.")
}

func run() error {
	urlExample := "postgres://pg:password@localhost:8541/tests"
	db, err := pgxpool.New(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		panic(err)
	}

	query := `
		CREATE EXTENSION IF NOT EXISTS pgcrypto;

		CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		first_name TEXT NOT NULL,
		last_name  TEXT NOT NULL,
		biography  TEXT NOT NULL
		);
	`
	
	_, err = db.Exec(context.Background(), query); 
	if err != nil {
		panic(err)
	}

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
