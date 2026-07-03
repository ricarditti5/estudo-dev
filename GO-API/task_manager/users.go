package main

import (
	"encoding/json"
	"errors"
	"log/slog"

	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Role string

const (
	Admin Role = "admin"
	User  Role = "user"
)

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
			return
		}

		if strings.TrimSpace(req.Nome) == "" ||
			strings.TrimSpace(req.Email) == "" ||
			strings.TrimSpace(req.Password) == "" {
			RespondError(w, http.StatusBadRequest, "All fields are required")
			return
		}

		PasswordHashed, err := HashPassword(req.Password)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Error")
			logger.Error("Erro to hash password", "error", err, "request_id", requestID)
			return
		}

		err = db.QueryRow(r.Context(), "INSERT INTO users(nome, email, password_hash) VALUES($1,$2,$3) RETURNING id", &req.Nome, &req.Email, PasswordHashed).Scan(&id)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				RespondError(w, http.StatusConflict, "email already exists")
				return
			}
			RespondError(w, http.StatusInternalServerError, "Error to create user")
			logger.Error("Erro to create user", "error", err, "request_id", requestID)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func ListUsers(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			requestID = "unknown"
		}
		rows, err := db.Query(r.Context(), "SELECT nome, email, role FROM users")
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			logger.Error("Erro to find users(don't try hack me XD)", "error", err, "request_id", requestID)
			return
		}
		defer rows.Close()

		var user []UserResponse

		for rows.Next() {
			var u UserResponse
			err := rows.Scan(&u.Nome, &u.Email, &u.Role)
			if err != nil {
				RespondError(w, http.StatusInternalServerError, "Internal error-error to get task")
				return
			}
			user = append(user, u)
		}
		//o encoder envia os dados diretamento em json para o w
		if err := RespondJSON(w, http.StatusOK, user); err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			logger.Error("Error to encode data", "error", err, "request_id", requestID)
			return
		}
	}
}

func UpdateUsers(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			requestID = "unknown"
		}
		var req UserService
		//pega o id passado no body ao chamar a rota
		id := r.PathValue("id")

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			RespondError(w, http.StatusBadRequest, "Invalid data")
			logger.Error("Internal Error", "error", err, "request_id", requestID)
			return
		}

		if strings.TrimSpace(req.Nome) == "" &&
			strings.TrimSpace(req.Email) == "" &&
			strings.TrimSpace(req.Password) == "" {
			RespondError(w, http.StatusBadRequest, "at least one field must be provided")
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
				return
			}
		}
		result, err := db.Exec(r.Context(), "UPDATE users SET nome = COALESCE(NULLIF($1,''), nome), email = COALESCE(NULLIF($2,''), email), password_hash = COALESCE(NULLIF($3,''), password_hash) WHERE id = $4", req.Nome, req.Email, passwordHash, id)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Error to update user")
			logger.Error("Error to execute query", "error", err, "request_id", requestID)
			return
		}

		rows := result.RowsAffected()
		if rows == 0 {
			RespondError(w, http.StatusNotFound, "Error to update user")
			logger.Error("no user founded with id:", "id", id, "request_id", requestID)
			return
		}

		if err := RespondJSON(w, http.StatusOK, MessageJSON{
			Message: "Succes to update task",
		}); err != nil {
			logger.Error("Error to encode data", "error", err, "request_id", requestID)
			return
		}
	}
}
