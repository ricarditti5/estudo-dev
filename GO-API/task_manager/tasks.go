package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
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

func ListTask(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			requestID = "unknown"
		}

		rows, err := db.Query(r.Context(), "SELECT id, title, description, status FROM tasks")
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			logger.Error("Error to find tasks(don't try hack me XD)", "errors", err, "request_id", requestID)
			return
		}
		defer rows.Close()

		var task []Task

		for rows.Next() {
			var t Task
			err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status)
			if err != nil {
				RespondError(w, http.StatusInternalServerError, "Internal error-error to get task")
				logger.Error("Internal error-error to get task", "errors", err, "request_id", requestID)
				return
			}
			task = append(task, t)
		}
		//o encoder envia os dados diretamento em json para o w
		if err := RespondJSON(w, http.StatusOK, task); err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			logger.Error("Error to encode data", "errors", err, "request_id", requestID)
			return
		}
	}
}

func CreateTask(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			requestID = "unknown"
		}
		var (
			id int
			t  Task
		)

		json.NewDecoder(r.Body).Decode(&t)

		err := db.QueryRow(r.Context(), "INSERT INTO tasks (title, description) VALUES($1,$2) RETURNING id", &t.Title, &t.Description).Scan(&id)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Error to create tasks")
			logger.Error("Error to create tasks", "errors", err, "request_id", requestID)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func GetTask(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			requestID = "unknown"
		}
		var t Task

		//pega o id passado no body ao chamar a rota
		id := r.PathValue("id")

		idConv, err := strconv.Atoi(id)
		if err != nil {
			logger.Error("Error to convert the variable", "errors", err, "request_id", requestID)
			return
		}
		err = db.QueryRow(r.Context(), "SELECT id, title, description, status FROM tasks WHERE id = $1", idConv).Scan(&t.ID, &t.Title, &t.Description, &t.Status)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				RespondError(w, http.StatusInternalServerError, "Error to get Tasks")
				logger.Error("Error to execute query", "errors", err, "request_id", requestID)
				return
			}
		}
		if err := RespondJSON(w, http.StatusOK, t); err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			logger.Error("Error to encode data", "errors", err, "request_id", requestID)
			return
		}
	}
}

func UpdateTask(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			requestID = "unknown"
		}
		var t Task
		//pega o id passado no body ao chamar a rota
		id := r.PathValue("id")

		idConv, err := strconv.Atoi(id)
		if err != nil {
			logger.Error("Error to covert the variable", "errors", err, "request_id", requestID)
			return
		}

		json.NewDecoder(r.Body).Decode(&t)

		result, err := db.Exec(r.Context(), "UPDATE tasks SET title = $1 ,description = $2, status = $3 WHERE id = $4", t.Title, t.Description, t.Status, idConv)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Task doesn't exist")
			logger.Error("Error to execute query", "errors", err, "request_id", requestID)
			return
		}

		rows := result.RowsAffected()
		if rows == 0 {
			RespondError(w, http.StatusInternalServerError, "Error to update task")
			logger.Error("No task founded with id", "errors", err, "request_id", requestID)
			return
		}

		if err := RespondJSON(w, http.StatusOK, t); err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			logger.Error("Error to encode data", "errors", err, "request_id", requestID)
			return
		}
	}
}

func DeleteTask(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			requestID = "unknown"
		}
		var t Task
		//pega o id passado no body ao chamar a rota
		id := r.PathValue("id")

		idConv, err := strconv.Atoi(id)
		if err != nil {
			logger.Error("Error to convert the variable ", "errors", err, "request_id", requestID)
			return
		}

		json.NewDecoder(r.Body).Decode(&t)

		result, err := db.Exec(r.Context(), "DELETE FROM tasks WHERE id = $1", idConv)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "Error to delete task")
			logger.Error("Invalid task: ", "errors", err, "request_id", requestID)
			return
		}

		rows := result.RowsAffected()
		if rows == 0 {
			RespondError(w, http.StatusInternalServerError, "Error to delete task")
			logger.Error("no task founded with id: ", "errors", idConv, "request_id", requestID)
			return
		}
		if err := RespondJSON(w, http.StatusOK, t); err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			logger.Error("Error to encode data: ", "errors", err, "request_id", requestID)
			return
		}
	}
}
