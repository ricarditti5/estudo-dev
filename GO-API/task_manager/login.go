package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type ClaimsToken struct {
	UserID string `json:"user_id"`
	Role   Role   `json:"role"` //get by users.go
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GenerateKey(userID string, role Role) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf(".env doesn't exist.")
	}
	SecretKey := os.Getenv("SECRET_KEY")

	claims := ClaimsToken{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	convKey := []byte(SecretKey)

	FinalSecretKey, err := token.SignedString(convKey)
	if err != nil {
		fmt.Printf("Error to generate token: %v", err)
		return "", nil
	}
	return FinalSecretKey, nil
}

func Login(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			login   LoginRequest
			reqData Users
		)

		if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
			fmt.Printf("Error to login: %v", err)
			http.Error(w, "Error to login", http.StatusBadRequest)
			return
		}

		if login.Email == "" {
			http.Error(w, "Fild email empty", http.StatusBadRequest)
			fmt.Printf("Fild email empty")
			return
		} else if login.Password == "" {
			http.Error(w, "Fild password empty", http.StatusBadRequest)
			fmt.Printf("Fild password empty")
			return
		}

		err := db.QueryRow(r.Context(), "SELECT password_hash FROM users WHERE email = $1", login.Email).Scan(&reqData.PasswordHash)
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			fmt.Printf("Internal Error: %v", err)
			return
		}

		pass := CheckPasswordHash(login.Password, reqData.PasswordHash)
		if pass == false {
			http.Error(w, "Wrong password", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
