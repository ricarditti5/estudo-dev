package main

import (
	"encoding/json"
	"net/http"

	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func ListTask(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, _ := db.Query(r.Context(), "SELECT id, title, description, status FROM tasks")
		defer rows.Close()

		var task []Task

		for rows.Next() {
			var t Task
			rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status)
			task = append(task, t)
		}
		//o encoder envia os dados diretamento em json para o w
		json.NewEncoder(w).Encode(task)
	}
}

func CreateTask(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			id int
			t  Task
		)

		json.NewDecoder(r.Body).Decode(&t)

		err := db.QueryRow(r.Context(), "INSERT INTO tasks (title, description) VALUES($1,$2) RETURNING id", &t.Title, &t.Description).Scan(&id)
		if err != nil {
			fmt.Printf("could not insert task: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func GetTask(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t Task

		//pega o id passado no body ao chamar a rota
		id := r.PathValue("id")

		idConv, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("Error to convert the variable: %v", err)
			return
		}
		err = db.QueryRow(r.Context(), "SELECT id, title, description, status FROM tasks WHERE id = $1", idConv).Scan(&t.ID, &t.Title, &t.Description, &t.Status)
		if err != nil {
			fmt.Printf("Error to execute query: %v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(t); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func UpdateTask(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t Task
		//pega o id passado no body ao chamar a rota
		id := r.PathValue("id")

		idConv, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("Error to convert the variable: %v", err)
			return
		}

		json.NewDecoder(r.Body).Decode(&t)

		result, err := db.Exec(r.Context(), "UPDATE tasks SET title = $1 ,description = $2, status = $3 WHERE id = $4", t.Title, t.Description, t.Status, idConv)
		if err != nil {
			fmt.Printf("Error to execute query: %v", err)
			return
		}

		rows := result.RowsAffected()
		if rows == 0 {
			fmt.Printf("no task founded with id: %d", idConv)
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(t); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func DeleteTask(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t Task
		//pega o id passado no body ao chamar a rota
		id := r.PathValue("id")

		idConv, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("Error to convert the variable: %v", err)
			return
		}

		json.NewDecoder(r.Body).Decode(&t)

		result, err := db.Exec(r.Context(), "DELETE FROM tasks WHERE id = $1", idConv)
		if err != nil {
			fmt.Printf("Error to execute query: %v", err)
			return
		}

		rows := result.RowsAffected()
		if rows == 0 {
			fmt.Printf("no task founded with id: %d", idConv)
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(t); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
