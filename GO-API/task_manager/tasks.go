package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
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
		rows, err := db.Query(r.Context(), "SELECT id, title, description, status FROM tasks")
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			fmt.Printf("Erro to find tasks(don't try hack me XD): %v", err)
			return
		}
		defer rows.Close()

		var task []Task

		for rows.Next() {
			var t Task
			err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status)
			if err != nil {
				RespondError(w, http.StatusInternalServerError, "Internal error-error to get task")
			}
			task = append(task, t)
		}
		//o encoder envia os dados diretamento em json para o w
		if err := RespondJSON(w, http.StatusOK, task); err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			fmt.Printf("Error to encode data: %v", err)
			return
		}
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
			RespondError(w, http.StatusInternalServerError, "Error to create tasks")
			fmt.Printf("Erro to create tasks: %v", err)
			return
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
			if errors.Is(err, pgx.ErrNoRows) {
				RespondError(w, http.StatusInternalServerError, "Error to get Tasks")
				fmt.Printf("Error to execute query: %v", err)
				return
			}
		}
		if err := RespondJSON(w, http.StatusOK, t); err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			fmt.Printf("Error to encode data: %v", err)
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
			RespondError(w, http.StatusInternalServerError, "Task doesn't exist")
			fmt.Printf("Error to execute query: %v", err)
			return
		}

		rows := result.RowsAffected()
		if rows == 0 {
			RespondError(w, http.StatusInternalServerError, "Error to update task")
			fmt.Printf("no task founded with id: %d", idConv)
			return
		}

		if err := RespondJSON(w, http.StatusOK, t); err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			fmt.Printf("Error to encode data: %v", err)
			return
		}
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
			RespondError(w, http.StatusBadRequest, "Error to delete task")
			fmt.Printf("Invalid task: %v", err)
			return
		}

		rows := result.RowsAffected()
		if rows == 0 {
			RespondError(w, http.StatusInternalServerError, "Error to delete task")
			fmt.Printf("no task founded with id: %d", idConv)
		}
		if err := RespondJSON(w, http.StatusOK, t); err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			fmt.Printf("Error to encode data: %v", err)
			return
		}
	}
}
