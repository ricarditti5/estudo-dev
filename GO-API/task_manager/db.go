package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(cfg *Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DB_URL)
	if err != nil {
		return nil, fmt.Errorf("error to connect in database: %v", err)
	}
	return pool, nil
}
