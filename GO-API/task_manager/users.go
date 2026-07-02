package main

import (
	"encoding/json"
	//"errors"
	"fmt"
	"net/http"
	"strings"

	//"github.com/jackc/pgx/v5"

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

func CreateUser(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			id  string
			req UserService
			ver Users
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

		//fazer select para comaprar duplicate email do user e email da db
		if req.Email == ver.Email {
			RespondError(w, http.StatusBadRequest, "Email already exists")
			return
		}

		PasswordHashed, err := HashPassword(req.Password)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Error")
			fmt.Printf("Erro to hash password: %v", err)
			return
		}
		//status code errado no updateusers
		//sem validação de campos copiar o do primeiro issue aqui
		err = db.QueryRow(r.Context(), "INSERT INTO users(nome, email, password_hash) VALUES($1,$2,$3) RETURNING id", &req.Nome, &req.Email, PasswordHashed).Scan(&id)
		if err != nil {
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

		if strings.TrimSpace(req.Nome) == "" ||
			strings.TrimSpace(req.Email) == "" ||
			strings.TrimSpace(req.Password) == "" {
			RespondError(w, http.StatusBadRequest, "All fields are required")
			return
		}

		result, err := db.Exec(r.Context(), "UPDATE users SET nome = $1 ,email = $2 WHERE id = $3", &req.Nome, &req.Email, id)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Error to update user")
			fmt.Printf("Error to execute query: %v", err)
			return
		}

		rows := result.RowsAffected()
		if rows == 0 {
			RespondError(w, http.StatusInternalServerError, "Error to update user")
			fmt.Printf("no task founded with id: %s", id)
			return
		}

		if err := RespondJSON(w, http.StatusOK, req); err != nil {
			fmt.Printf("Error to encode data: %v", err)
			return
		}
	}
}
