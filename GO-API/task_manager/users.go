package main

import (
	"encoding/json"
	//"errors"
	"fmt"
	"net/http"

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
		)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			fmt.Printf("Erro to request: %v", err)
			return
		}

		PasswordHashed, err := HashPassword(req.Password)
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			fmt.Printf("Erro to hash password: %v", err)
			return
		}
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
			rows.Scan(&u.Nome, &u.Email, &u.Role)
			user = append(user, u)
		}
		//o encoder envia os dados diretamento em json para o w
		if err := RespondJSON(w, http.StatusOK, user); err != nil {
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
			RespondError(w, http.StatusInternalServerError, "Internal error")
			fmt.Printf("Internal Error: %v", err)
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
