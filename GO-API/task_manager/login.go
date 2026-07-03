package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type TokenKey struct {
	Token string `json:"token"`
}

var loginTracer = otel.Tracer("login-handler")

type ClaimsToken struct {
	UserID string `json:"user_id"`
	Role   Role   `json:"role"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GenerateKey(cfg *Config, userID string, role Role) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("config cannot be nil")
	}
	claims := ClaimsToken{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	convKey := []byte(cfg.JWTSecret)

	FinalSecretKey, err := token.SignedString(convKey)
	if err != nil {
		return "", fmt.Errorf("Error to generate token: %v", err)
	}
	return FinalSecretKey, nil
}

func Login(cfg *Config, db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := loginTracer.Start(r.Context(), "Login.Login")
		defer span.End()
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
			span.SetStatus(codes.Error, "Error decoding login request")
			span.RecordError(err)
			return
		}

		if login.Email == "" {
			RespondError(w, http.StatusBadRequest, "Fild email empty")
			logger.Error("Fild email is empty", "request_id", requestID)
			span.SetStatus(codes.Error, "Email field is empty")
			return
		} else if login.Password == "" {
			RespondError(w, http.StatusBadRequest, "Fild password empty")
			logger.Error("Fild password is empty", "request_id", requestID)
			span.SetStatus(codes.Error, "Password field is empty")
			return
		}
		span.SetAttributes(attribute.String("user.email", login.Email))

		err := db.QueryRow(ctx, "SELECT id, role, password_hash FROM users WHERE email = $1", login.Email).Scan(&reqData.ID, &reqData.Role, &reqData.PasswordHash)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				RespondError(w, http.StatusUnauthorized, "Invalid credentials")
				span.SetStatus(codes.Error, "User not found")
				span.RecordError(err)
				return
			}
			RespondError(w, http.StatusInternalServerError, "Internal error")
			logger.Error("Wrong credentials", "errors", err, "request_id", requestID)
			span.SetStatus(codes.Error, "Error querying user")
			span.RecordError(err)
			return
		}

		pass := CheckPasswordHash(login.Password, reqData.PasswordHash)
		if pass == false {
			RespondError(w, http.StatusUnauthorized, "Wrong credentials")
			span.SetStatus(codes.Error, "Invalid password")
			return
		}

		SecretKey, err := GenerateKey(cfg, reqData.ID, reqData.Role)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error")
			logger.Error("Error to with", "errors", err, "request_id", requestID)
			span.SetStatus(codes.Error, "Error generating token")
			span.RecordError(err)
			return
		}
		if err := RespondJSON(w, http.StatusOK, TokenKey{
			Token: SecretKey,
		}); err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error")
			logger.Error("Error to encode data", "errors", err, "request_id", requestID)
			span.SetStatus(codes.Error, "Error encoding response")
			span.RecordError(err)
			return
		}
	}
}
