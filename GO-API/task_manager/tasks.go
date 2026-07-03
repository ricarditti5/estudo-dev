package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

var taskTracer = otel.Tracer("task-handler")

func ListTask(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := taskTracer.Start(r.Context(), "ListTask.QueryTasks")
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.path", r.URL.Path),
		)

		requestID, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			requestID = "unknown"
		}

		rows, err := db.Query(ctx, "SELECT id, title, description, status FROM tasks")
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to query tasks")
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
				span.RecordError(err)
				span.SetStatus(codes.Error, "failed to scan task")
				RespondError(w, http.StatusInternalServerError, "Internal error-error to get task")
				logger.Error("Internal error-error to get task", "errors", err, "request_id", requestID)
				return
			}
			task = append(task, t)
		}

		span.SetAttributes(attribute.Int("tasks.count", len(task)))

		if err := RespondJSON(w, http.StatusOK, task); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to encode response")
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			logger.Error("Error to encode data", "errors", err, "request_id", requestID)
			return
		}
	}
}

func CreateTask(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := taskTracer.Start(r.Context(), "CreateTask.QueryTasks")
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.path", r.URL.Path),
		)

		requestID, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			requestID = "unknown"
		}
		var (
			id int
			t  Task
		)

		json.NewDecoder(r.Body).Decode(&t)

		err := db.QueryRow(ctx, "INSERT INTO tasks (title, description) VALUES($1,$2) RETURNING id", &t.Title, &t.Description).Scan(&id)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to create task")
			RespondError(w, http.StatusInternalServerError, "Error to create tasks")
			logger.Error("Error to create tasks", "errors", err, "request_id", requestID)
			return
		}
		span.SetAttributes(attribute.Int("task.id", id))
		w.WriteHeader(http.StatusCreated)
	}
}

func GetTask(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := taskTracer.Start(r.Context(), "GetTask.QueryTasks")
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.path", r.URL.Path),
		)

		requestID, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			requestID = "unknown"
		}
		var t Task

		id := r.PathValue("id")

		idConv, err := strconv.Atoi(id)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "invalid task id")
			RespondError(w, http.StatusBadRequest, "Invalid task id")
			logger.Error("Error to convert the variable", "errors", err, "request_id", requestID)
			return
		}
		span.SetAttributes(attribute.Int("task.id", idConv))

		err = db.QueryRow(ctx, "SELECT id, title, description, status FROM tasks WHERE id = $1", idConv).Scan(&t.ID, &t.Title, &t.Description, &t.Status)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				span.SetStatus(codes.Error, "task not found")
				span.RecordError(err)
				RespondError(w, http.StatusInternalServerError, "Error to get Tasks")
				logger.Error("Error to execute query", "errors", err, "request_id", requestID)
				return
			}
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to query task")
			RespondError(w, http.StatusInternalServerError, "Error to get Tasks")
			logger.Error("Error to execute query", "errors", err, "request_id", requestID)
			return
		}

		if err := RespondJSON(w, http.StatusOK, t); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to encode response")
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			logger.Error("Error to encode data", "errors", err, "request_id", requestID)
			return
		}
	}
}

func UpdateTask(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := taskTracer.Start(r.Context(), "UpdateTask.QueryTasks")
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.path", r.URL.Path),
		)

		requestID, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			requestID = "unknown"
		}
		var t Task
		id := r.PathValue("id")

		idConv, err := strconv.Atoi(id)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "invalid task id")
			RespondError(w, http.StatusBadRequest, "Invalid task id")
			logger.Error("Error to covert the variable", "errors", err, "request_id", requestID)
			return
		}
		span.SetAttributes(attribute.Int("task.id", idConv))

		json.NewDecoder(r.Body).Decode(&t)

		result, err := db.Exec(ctx, "UPDATE tasks SET title = $1 ,description = $2, status = $3 WHERE id = $4", t.Title, t.Description, t.Status, idConv)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to update task")
			RespondError(w, http.StatusInternalServerError, "Task doesn't exist")
			logger.Error("Error to execute query", "errors", err, "request_id", requestID)
			return
		}

		rows := result.RowsAffected()
		if rows == 0 {
			span.SetStatus(codes.Error, "task not found for update")
			RespondError(w, http.StatusInternalServerError, "Error to update task")
			logger.Error("No task founded with id", "errors", idConv, "request_id", requestID)
			return
		}

		if err := RespondJSON(w, http.StatusOK, t); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to encode response")
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			logger.Error("Error to encode data", "errors", err, "request_id", requestID)
			return
		}
	}
}

func DeleteTask(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := taskTracer.Start(r.Context(), "DeleteTask.QueryTasks")
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.path", r.URL.Path),
		)

		requestID, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			requestID = "unknown"
		}
		var t Task
		id := r.PathValue("id")

		idConv, err := strconv.Atoi(id)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "invalid task id")
			RespondError(w, http.StatusBadRequest, "Invalid task id")
			logger.Error("Error to convert the variable ", "errors", err, "request_id", requestID)
			return
		}
		span.SetAttributes(attribute.Int("task.id", idConv))

		json.NewDecoder(r.Body).Decode(&t)

		result, err := db.Exec(ctx, "DELETE FROM tasks WHERE id = $1", idConv)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to delete task")
			RespondError(w, http.StatusBadRequest, "Error to delete task")
			logger.Error("Invalid task: ", "errors", err, "request_id", requestID)
			return
		}

		rows := result.RowsAffected()
		if rows == 0 {
			span.SetStatus(codes.Error, "task not found for deletion")
			RespondError(w, http.StatusInternalServerError, "Error to delete task")
			logger.Error("no task founded with id: ", "errors", idConv, "request_id", requestID)
			return
		}
		if err := RespondJSON(w, http.StatusOK, t); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to encode response")
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			logger.Error("Error to encode data: ", "errors", err, "request_id", requestID)
			return
		}
	}
}
