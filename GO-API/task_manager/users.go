package main

import (
	"encoding/json"
	"errors"
	"log/slog"

	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type Role string

const (
	Admin Role = "admin"
	User  Role = "user"
)

var userTracer = otel.Tracer("user-handler")

type Users struct {
	ID           string `json:"id"`
	Nome         string `json:"nome"`
	Role         Role   `json:"role"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type UserService struct {
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Nome  string `json:"nome"`
	Email string `json:"email"`
	Role  Role   `json:"role"`
}

type MessageJSON struct {
	Message string
}

func CreateUser(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := userTracer.Start(r.Context(), "CreateUsers.QueryUsers")
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
			id  string
			req UserService
		)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			RespondError(w, http.StatusBadRequest, "invalid request body")
			logger.Error("Erro to request", "error", err, "request_id", requestID)
			span.SetStatus(codes.Error, "Error decoding request body")
			span.RecordError(err)
			return
		}

		if strings.TrimSpace(req.Nome) == "" ||
			strings.TrimSpace(req.Email) == "" ||
			strings.TrimSpace(req.Password) == "" {
			RespondError(w, http.StatusBadRequest, "All fields are required")
			span.SetStatus(codes.Error, "Missing required fields")
			return
		}

		PasswordHashed, err := HashPassword(req.Password)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Error")
			logger.Error("Erro to hash password", "error", err, "request_id", requestID)
			span.SetStatus(codes.Error, "Error hashing password")
			span.RecordError(err)
			return
		}

		err = db.QueryRow(ctx, "INSERT INTO users(nome, email, password_hash) VALUES($1,$2,$3) RETURNING id", &req.Nome, &req.Email, PasswordHashed).Scan(&id)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				RespondError(w, http.StatusConflict, "email already exists")
				span.SetStatus(codes.Error, "Email already exists")
				span.RecordError(err)
				return
			}
			RespondError(w, http.StatusInternalServerError, "Error to create user")
			logger.Error("Erro to create user", "error", err, "request_id", requestID)
			span.SetStatus(codes.Error, "Error creating user")
			span.RecordError(err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func ListUsers(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := userTracer.Start(r.Context(), "ListUsers.QueryUsers")
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.path", r.URL.Path),
		)

		requestID, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			requestID = "unknown"
		}
		rows, err := db.Query(ctx, "SELECT nome, email, role FROM users")
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			logger.Error("Erro to find users(don't try hack me XD)", "error", err, "request_id", requestID)
			span.SetStatus(codes.Error, "Error querying users")
			span.RecordError(err)
			return
		}
		defer rows.Close()

		var user []UserResponse

		for rows.Next() {
			var u UserResponse
			err := rows.Scan(&u.Nome, &u.Email, &u.Role)
			if err != nil {
				RespondError(w, http.StatusInternalServerError, "Internal error-error to get task")
				span.SetStatus(codes.Error, "Error scanning user")
				span.RecordError(err)
				return
			}
			user = append(user, u)
		}
		if err := RespondJSON(w, http.StatusOK, user); err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			logger.Error("Error to encode data", "error", err, "request_id", requestID)
			span.SetStatus(codes.Error, "Error encoding response")
			span.RecordError(err)
			return
		}
	}
}

func UpdateUsers(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := userTracer.Start(r.Context(), "UpdateUsers.QueryUsers")
		defer span.End()
		requestID, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			requestID = "unknown"
		}
		var req UserService

		//precisa de validação de conversao do id para caso der erro no futuro
		id := r.PathValue("id")
		span.SetAttributes(attribute.String("user.id", id))

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			RespondError(w, http.StatusBadRequest, "Invalid data")
			logger.Error("Internal Error", "error", err, "request_id", requestID)
			span.SetStatus(codes.Error, "Error decoding request body")
			span.RecordError(err)
			return
		}

		if strings.TrimSpace(req.Nome) == "" &&
			strings.TrimSpace(req.Email) == "" &&
			strings.TrimSpace(req.Password) == "" {
			RespondError(w, http.StatusBadRequest, "at least one field must be provided")
			span.SetStatus(codes.Error, "No fields provided for update")
			return
		}
		var (
			passwordHash string
			err          error
		)

		if strings.TrimSpace(req.Password) != "" {
			passwordHash, err = HashPassword(req.Password)
			if err != nil {
				RespondError(w, http.StatusInternalServerError, "Unable to update password")
				logger.Error("Error to update password", "error", err, "request_id", requestID)
				span.SetStatus(codes.Error, "Error hashing password")
				span.RecordError(err)
				return
			}
		}
		result, err := db.Exec(ctx, "UPDATE users SET nome = COALESCE(NULLIF($1,''), nome), email = COALESCE(NULLIF($2,''), email), password_hash = COALESCE(NULLIF($3,''), password_hash) WHERE id = $4", req.Nome, req.Email, passwordHash, id)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Error to update user")
			logger.Error("Error to execute query", "error", err, "request_id", requestID)
			span.SetStatus(codes.Error, "Error updating user")
			span.RecordError(err)
			return
		}

		rows := result.RowsAffected()
		if rows == 0 {
			RespondError(w, http.StatusNotFound, "Error to update user")
			logger.Error("no user founded with id:", "id", id, "request_id", requestID)
			span.SetStatus(codes.Error, "User not found for update")
			return
		}

		if err := RespondJSON(w, http.StatusOK, MessageJSON{
			Message: "Succes to update task",
		}); err != nil {
			logger.Error("Error to encode data", "error", err, "request_id", requestID)
			span.SetStatus(codes.Error, "Error encoding response")
			span.RecordError(err)
			return
		}
	}
}
