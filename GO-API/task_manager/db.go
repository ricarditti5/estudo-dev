package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func ConnectDB() *pgxpool.Pool {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env doesn't exist.")
		return nil
	}
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATA_BASE_URL"))
	if err != nil {
		panic(err)
	}
	return pool
}
