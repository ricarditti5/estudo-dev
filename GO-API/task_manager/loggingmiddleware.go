package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const requestIDKey contextKey = "request_id"

type ResponseWriterWrapper struct {
	http.ResponseWriter
	StatusCode int
}

func (rw *ResponseWriterWrapper) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			//isto serve para gerar um identificador unico e guarda no contexto, para depois poder ler lido dentro dos handlers quando forem chamados nas requisições
			requestID := uuid.New().String()
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			r = r.WithContext(ctx)
			//---------------------------------------------------

			//wrapper para poder passar o status code após a chamada do middleware, com 200 como default
			rw := &ResponseWriterWrapper{ResponseWriter: w, StatusCode: 200}

			next.ServeHTTP(rw, r)

			duration := time.Since(start)
			logger.Info(
				"processing request",
				"method", r.Method,
				"path", r.URL.Path,
				//aqui passo no contextt.Value--> a chave q criamos para ser o identificador no caso o const requestIDKey
				"request_id", r.Context().Value(requestIDKey),
				"status", rw.StatusCode,
				"duration_ms", duration.Milliseconds(),
			)
		})
	}
}
