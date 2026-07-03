package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type TokenKey struct {
	Token string `json:"token"`
}

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
		return "", fmt.Errorf("Error to generate token: %v", err)
	}
	return FinalSecretKey, nil
}

func Login(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			requestID = "unknown"
		}

		var (
			login   LoginRequest
			reqData Users
		)

		if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
			RespondError(w, http.StatusBadRequest, "Error to login")
			logger.Error("Error to login", "errors", err, "request_id", requestID)
			return
		}

		if login.Email == "" {
			RespondError(w, http.StatusBadRequest, "Fild email empty")
			logger.Error("Fild email is empty", "request_id", requestID)
			return
		} else if login.Password == "" {
			RespondError(w, http.StatusBadRequest, "Fild password empty")
			logger.Error("Fild password is empty", "request_id", requestID)
			return
		}

		err := db.QueryRow(r.Context(), "SELECT id, role, password_hash FROM users WHERE email = $1", login.Email).Scan(&reqData.ID, &reqData.Role, &reqData.PasswordHash)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) { //o errors.Is(compara o err que recebi da db com o erro sentinela)
				RespondError(w, http.StatusUnauthorized, "Invalid credentials")
				return
			}
			RespondError(w, http.StatusInternalServerError, "Internal error")
			logger.Error("Wrong credentials", "errors", err, "request_id", requestID)
			return
		}

		pass := CheckPasswordHash(login.Password, reqData.PasswordHash)
		if pass == false {
			RespondError(w, http.StatusUnauthorized, "Wrong credentials")
			return
		}

		SecretKey, err := GenerateKey(reqData.ID, reqData.Role)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error")
			logger.Error("Error to with", "errors", err, "request_id", requestID)
			return
		}
		if err := RespondJSON(w, http.StatusOK, TokenKey{
			Token: SecretKey,
		}); err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error")
			logger.Error("Error to encode data", "errors", err, "request_id", requestID)
			return
		}
	}
}
