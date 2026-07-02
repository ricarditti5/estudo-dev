package main

import (
	"encoding/json"
	"errors"

	"fmt"
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

func CreateUser(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			id  string
			req UserService
		)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			RespondError(w, http.StatusBadRequest, "invalid request body")
			fmt.Printf("Erro to request: %v", err)
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
			fmt.Printf("Erro to hash password: %v", err)
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
			fmt.Printf("Erro to create user: %v", err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func ListUsers(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(r.Context(), "SELECT nome, email, role FROM users")
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			fmt.Printf("Erro to find users(don't try hack me XD): %v", err)
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
			fmt.Printf("Error to encode data: %v", err)
			return
		}
	}
}

func UpdateUsers(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UserService
		//pega o id passado no body ao chamar a rota
		id := r.PathValue("id")

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			RespondError(w, http.StatusBadRequest, "Invalid data")
			fmt.Printf("Internal Error: %v", err)
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
				fmt.Printf("Error to update password: %v", err)
				return
			}
		}
		result, err := db.Exec(r.Context(), "UPDATE users SET nome = COALESCE(NULLIF($1,''), nome), email = COALESCE(NULLIF($2,''), email), password_hash = COALESCE(NULLIF($3,''), password_hash) WHERE id = $4", req.Nome, req.Email, passwordHash, id)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Error to update user")
			fmt.Printf("Error to execute query: %v", err)
			return
		}

		rows := result.RowsAffected()
		if rows == 0 {
			RespondError(w, http.StatusNotFound, "Error to update user")
			fmt.Printf("no user founded with id: %s", id)
			return
		}

		if err := RespondJSON(w, http.StatusOK, MessageJSON{
			Message: "Succes to update task",
		}); err != nil {
			fmt.Printf("Error to encode data: %v", err)
			return
		}
	}
}
