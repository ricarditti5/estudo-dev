package main

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func ListTask(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, _ := db.Query(r.Context(), "SELECT id, title, description FROM tasks")

		//o encoder envia os dados diretamento em json para o w
		json.NewEncoder(w).Encode(tasks)
	}
}
