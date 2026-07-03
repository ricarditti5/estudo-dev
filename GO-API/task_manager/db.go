package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func ConnectDB() (*pgxpool.Pool, error) {
	var logger *slog.Logger
	err := godotenv.Load()
	if err != nil {
		logger.Error(".env doesn't exist.", "errors", err)
		return nil, fmt.Errorf("Error to load .env: %v", err)
	}
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATA_BASE_URL"))
	if err != nil {
		panic(err)
	}
	return pool, nil
}
