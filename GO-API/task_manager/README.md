# Task Manager API — Documentação

API de gestão de tarefas e utilizadores com **autenticação JWT**, **bcrypt**, **logs estruturados (slog)** e **tracing (OpenTelemetry)**.

Stack: Go 1.26 + PostgreSQL (`pgx`) + `net/http` (ServeMux nativo)

---

## Http Errors e Validação

### Problema: validação e erros são coisas distintas mas relacionadas

1. **Validação** — garantir que os dados que chegam fazem sentido antes de tocares na BD.
2. **Tratamento de erros** — comunicar ao cliente o que correu mal de forma consistente.

### Formato de erro consistente

```json
{"error": "Fild email empty"}
```

### Helper de respostas (`response.go`)

Centraliza a escrita de respostas JSON para evitar repetir `Header → WriteHeader → Encode` em todos os handlers.

```go
type ErrorResponse struct {
	Error string `json:"error"`
}

func RespondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func RespondJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
```

### Porque é que a ordem importa

O `ResponseWriter` é **stateful**: assim que chamas `WriteHeader()`, o status code fica fechado. Se chamares `Write()` ou `Encode()` sem `WriteHeader()`, o Go assume `200 OK`. A ordem é sempre: `Header → WriteHeader → Encode`.

### A lógica dos status codes (de quem foi a culpa?)

| Status | Quando usar |
| --- | --- |
| **400 Bad Request** | O cliente mandou JSON inválido ou campo obrigatório vazio |
| **401 Unauthorized** | O cliente não provou quem é (credenciais erradas) |
| **404 Not Found** | Recurso não existe (RowsAffected == 0) |
| **409 Conflict** | Pedido conflitua com estado atual (email duplicado) |
| **500 Internal Server Error** | A culpa foi do servidor (falha de BD, etc.) |

### `errors.Is` vs `errors.As`

- **`errors.Is(err, sentinela)`** — o erro é exatamente esta sentinela? (`pgx.ErrNoRows`)
- **`errors.As(err, &variavelTipo)`** — o erro é deste tipo, dando acesso aos campos? (`pgconn.PgError` para ver `.Code == "23505"`)

### Padrão "guarda + return"

Cada `if err != nil { RespondError(...); return }` corta o caminho assim que algo está errado. O código de sucesso fica no fim sem indentação extra.

---

## Http Routing

### `http.Handler` — a interface fundamental

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

`http.HandlerFunc` é um wrapper que adapta funções normais:

```go
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}
```

### `http.ServeMux` — o router/multiplexer

Desde Go 1.22, o `ServeMux` suporta **métodos HTTP** e **path parameters** nativos:

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /tasks/{id}", handler)
```

Extrair o parâmetro:

```go
id := r.PathValue("id")
```

### Tabela de rotas

| Método | Rota | Handler | Descrição |
| --- | --- | --- | --- |
| `GET` | `/tasks` | `ListTask` | Lista todas as tarefas |
| `GET` | `/tasks/{id}` | `GetTask` | Obtém uma tarefa por ID |
| `POST` | `/tasks` | `CreateTask` | Cria uma nova tarefa |
| `PUT` | `/tasks/{id}` | `UpdateTask` | Atualiza uma tarefa |
| `DELETE` | `/tasks/{id}` | `DeleteTask` | Apaga uma tarefa |
| `GET` | `/users` | `ListUsers` | Lista todos os utilizadores |
| `POST` | `/users` | `CreateUser` | Cria um utilizador |
| `PATCH` | `/users/{id}` | `UpdateUsers` | Atualiza um utilizador |
| `POST` | `/login` | `Login` | Autentica e devolve JWT |

### Detalhes importantes

- **Especificidade vence**: `/tasks/{id}` é mais específico que `/tasks/`
- **Conflitos dão panic**: registar `/tasks/{id}` e `/tasks/{name}` juntos causa panic
- **405 automático**: com `"GET /tasks"`, um `POST` devolve `405 Method Not Allowed`
- **Sem regex**: patterns só suportam literais, `{param}` e `{param...}` (wildcard)

---

## Middlewares

### O que é um middleware

Um **middleware** é uma função que executa **antes e/ou depois** do handler, adicionando comportamentos comuns sem repetir código.

```
Cliente → Middleware (Logs) → Middleware (Auth) → Handler → Resposta
```

### Padrão em Go

```go
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
```

Cada middleware pode:
- Continuar para o próximo (`next.ServeHTTP`)
- Modificar o pedido/resposta
- Interromper a cadeia devolvendo um erro

### LoggingMiddleware (`middleware.go`) — exemplo completo

```go
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type contextKey string
const requestIDKey contextKey = "request_id"

var middlewareTracer = otel.Tracer("middleware")

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
			ctx, span := middlewareTracer.Start(r.Context(), "http.request")
			defer span.End()
			r = r.WithContext(ctx)

			requestID := uuid.New().String()
			ctx = context.WithValue(r.Context(), requestIDKey, requestID)
			r = r.WithContext(ctx)

			rw := &ResponseWriterWrapper{ResponseWriter: w, StatusCode: 200}
			next.ServeHTTP(rw, r)

			span.SetAttributes(
				semconv.HTTPMethodKey.String(r.Method),
				semconv.HTTPTargetKey.String(r.URL.Path),
				semconv.HTTPStatusCodeKey.Int(rw.StatusCode),
			)
			if rw.StatusCode >= 400 {
				span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d %s", rw.StatusCode, http.StatusText(rw.StatusCode)))
			}

			duration := time.Since(start)
			logger.Info("processing request",
				"method", r.Method,
				"path", r.URL.Path,
				"request_id", r.Context().Value(requestIDKey),
				"status", rw.StatusCode,
				"duration_ms", duration.Milliseconds(),
			)
		})
	}
}
```

### No `main.go`

```go
loggerMux := LoggingMiddleware(logger)(mux)
http.ListenAndServe(":"+cfg.PORT, loggerMux)
```

### O que o middleware faz

1. Gera um `requestID` único (UUID) e guarda-o no contexto
2. Cria um span do OpenTelemetry para o pedido
3. Usa um `ResponseWriterWrapper` para capturar o status code
4. No fim, regista no log: método, path, request_id, status, duração
5. Marca o span como erro se status >= 400

---

## Auth (JWT) e Hash

### JWT — o problema que resolve

HTTP é **stateless**. JWT resolve a autenticação sem guardar sessões no servidor:

- O servidor assina um "bilhete" com info do utilizador
- O cliente guarda e envia em cada pedido
- O servidor só verifica a assinatura — não precisa de ir à BD

### Anatomia de um JWT

```
eyJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoxMjN9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
   \_______________/    \_______________/    \____________________________________/
        HEADER              PAYLOAD                   SIGNATURE
```

**Importante**: o payload não é encriptado — é só Base64. Qualquer um consegue ler. A assinatura só garante que **não foi alterado**.

### Estrutura dos claims

```go
type ClaimsToken struct {
	UserID string `json:"user_id"`
	Role   Role   `json:"role"`
	jwt.RegisteredClaims
}
```

### Geração do token (`login.go`)

```go
func GenerateKey(cfg *Config, userID string, role Role) (string, error) {
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
```

### Porque bcrypt e não SHA-256

SHA-256 é **rápido** — foi desenhado para performance. Um atacante testa milhões de passwords por segundo. bcrypt é **lento de propósito**, ajustável via **cost factor**.

### Output do bcrypt

```
$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
```

- `$2a$` → versão do algoritmo
- `10` → cost factor
- Salt está **embutido** no hash (não precisas de o guardar à parte)

### Fluxo completo de password

**Registo**: password → `bcrypt.GenerateFromPassword` → guardas o hash
**Login**: password → `bcrypt.CompareHashAndPassword` → `nil` se bater

### `password.go` — exemplo completo

```go
package main

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bcryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bcryptedPass), nil
}

func CheckPasswordHash(password, hash string) bool {
	checkPassword := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return checkPassword == nil
}
```

### Login — `login.go` (versão final)

```go
func Login(cfg *Config, db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := loginTracer.Start(r.Context(), "Login.Login")
		defer span.End()
		requestID := extractRequestID(r)

		var login LoginRequest
		var reqData Users

		if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
			RespondError(w, http.StatusBadRequest, "Error to login")
			span.SetStatus(codes.Error, "Error decoding login request")
			return
		}

		// validacao de campos vazios
		if login.Email == "" {
			RespondError(w, http.StatusBadRequest, "Fild email empty")
			return
		} else if login.Password == "" {
			RespondError(w, http.StatusBadRequest, "Fild password empty")
			return
		}

		// buscar user na BD
		err := db.QueryRow(ctx, "SELECT id, role, password_hash FROM users WHERE email = $1",
			login.Email).Scan(&reqData.ID, &reqData.Role, &reqData.PasswordHash)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				RespondError(w, http.StatusUnauthorized, "Invalid credentials")
				return
			}
			RespondError(w, http.StatusInternalServerError, "Internal error")
			return
		}

		// comparar password
		if !CheckPasswordHash(login.Password, reqData.PasswordHash) {
			RespondError(w, http.StatusUnauthorized, "Wrong credentials")
			return
		}

		// gerar token
		SecretKey, err := GenerateKey(cfg, reqData.ID, reqData.Role)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error")
			return
		}

		RespondJSON(w, http.StatusOK, TokenKey{Token: SecretKey})
	}
}
```

### Segurança: não vazar informação

No login, "email não existe" e "password errada" devolvem a **mesma mensagem** (`"Invalid credentials"`) para não revelares que emails estão registados.

---

## Config & Secrets

### Problema: `os.Getenv` espalhado

1. **Sem validação centralizada** — descobres que falta uma env var a meio da execução
2. **Sem tipagem** — `os.Getenv` devolve sempre `string`
3. **Código espalhado** — não há um sítio único que diga "isto é a config da app"

### Solução: struct `Config` + Viper

```go
type Config struct {
	PORT      string `mapstructure:"PORT"`
	DB_URL    string `mapstructure:"DATABASE_URL"`
	JWTSecret string `mapstructure:"SECRET_KEY"`
}
```

### `config.go` — exemplo completo

```go
package main

import (
	"fmt"
	"log/slog"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	PORT      string `mapstructure:"PORT"`
	DB_URL    string `mapstructure:"DATABASE_URL"`
	JWTSecret string `mapstructure:"SECRET_KEY"`
}

func LoadConfig(logger *slog.Logger) (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Error to load .env %v", err)
	}

	viper.SetDefault("PORT", "8080")
	viper.BindEnv("DATABASE_URL")
	viper.BindEnv("SECRET_KEY")
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("Error to get variables")
	}

	if cfg.DB_URL == "" || cfg.JWTSecret == "" {
		return nil, fmt.Errorf("Empty variables")
	}

	return &cfg, nil
}
```

### Fluxo

```
.env (ficheiro) → godotenv.Load() → env vars do processo → viper.AutomaticEnv() → struct Config
```

### Diferença entre config e secrets

- **Config não-sensível** (porta, log level) — pode ir para um ficheiro versionado
- **Secrets** (DB password, JWT secret) — **nunca** vão para o git

### Nota sobre `AutomaticEnv()` + `Unmarshal()`

`viper.AutomaticEnv()` sozinho não funciona com `Unmarshal()` para chaves que o Viper não conhece. É necessário registá-las com `viper.BindEnv("DATABASE_URL")` ou dar um `viper.SetDefault(...)` para cada variável que queres que o Unmarshal encontre.

---

## Database

### `db.go` — conexão PostgreSQL com pgxpool

```go
package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(cfg *Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DB_URL)
	if err != nil {
		return nil, fmt.Errorf("error to connect in database: %v", err)
	}
	return pool, nil
}
```

### Migrations

**UP** (`db/migrations/000001_alter_tasks_users.up.sql`):

```sql
ALTER TABLE users ADD COLUMN email varchar(255) NOT NULL;
ALTER TABLE users ADD CONSTRAINT up_users_email UNIQUE(email);
ALTER TABLE users ADD COLUMN password_hash varchar(255) NOT NULL;
```

**DOWN** (`db/migrations/000001_alter_tasks_users.down.sql`):

```sql
ALTER TABLE users DROP COLUMN password_hash;
ALTER TABLE users DROP CONSTRAINT up_users_email;
ALTER TABLE users DROP COLUMN email;
```

---

## Observability (Logging e Tracing)

### `slog` — logging estruturado

Em vez de `log.Println("user logged in:", userID)` (texto solto), usamos pares chave-valor:

```go
logger.Info("user logged in", "user_id", userID) // estruturado
```

Criação do logger:

```go
logger := slog.New(slog.NewTextHandler(os.Stdout, nil)) // dev: texto legível
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)) // prod: JSON
```

### Níveis de log

| Nível | O que significa | A aplicação continua? |
| --- | --- | --- |
| **Debug** | Informação detalhada para programadores | ✅ |
| **Info** | Evento normal da aplicação | ✅ |
| **Warn** | Algo inesperado mas recuperável | ✅ |
| **Error** | Uma operação falhou | ✅ |
| **Fatal** | Erro crítico — aplicação termina | ❌ |
| **Panic** | Situação excecional — gera stack trace | ❌ (salvo recover) |

### Logger como dependência (mesmo padrão do `db`)

```go
func ListTask(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
```

Isto é uma **closure**: o handler "lembra-se" das dependências passadas.

### OpenTelemetry — tracing distribuído

**Logs** respondem a "o que aconteceu". **Traces** respondem a "onde foi gasto o tempo e em que ordem".

Conceitos-chave:
- **Trace** — a viagem completa de um pedido
- **Span** — um passo individual dentro da viagem
- **Trace ID / Span ID** — identificadores únicos

### Inicialização (`tracing.go`)

```go
package main

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func InitTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("task-manager"),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}
```

### Como usar spans nos handlers

```go
var taskTracer = otel.Tracer("task-handler") // uma vez por ficheiro

func ListTask(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := taskTracer.Start(r.Context(), "ListTask.QueryTasks")
		defer span.End()

		// usar ctx (nao r.Context()) na query
		rows, err := db.Query(ctx, "SELECT ...")
		// ...
	}
}
```

### Hierarquia

```
TracerProvider (motor, um por aplicação)
    └── Tracer (obtido do provider, um por pacote)
            └── Span (criado pelo tracer, um por operação)
```

### Correlacionar logs com traces

```go
span.SetAttributes(attribute.Int("task.id", idConv))
span.SetStatus(codes.Error, "task not found")
span.RecordError(err)
```

---

## Tasks API — CRUD

### Model

```go
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
```

### `ListTask` — GET /tasks

```go
func ListTask(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := taskTracer.Start(r.Context(), "ListTask.QueryTasks")
		defer span.End()
		requestID := extractRequestID(r)

		rows, err := db.Query(ctx, "SELECT id, title, description, status FROM tasks")
		if err != nil {
			span.RecordError(err)
			RespondError(w, http.StatusInternalServerError, "Internal error, try again refreshing")
			logger.Error("Error to find tasks", "errors", err, "request_id", requestID)
			return
		}
		defer rows.Close()

		var tasks []Task
		for rows.Next() {
			var t Task
			if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status); err != nil {
				RespondError(w, http.StatusInternalServerError, "Internal error")
				return
			}
			tasks = append(tasks, t)
		}
		RespondJSON(w, http.StatusOK, tasks)
	}
}
```

### `CreateTask` — POST /tasks

Body: `{"title": "...", "description": "..."}` (status é gerado pela BD com default)

```go
func CreateTask(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := taskTracer.Start(r.Context(), "CreateTask.QueryTasks")
		defer span.End()

		var t Task
		json.NewDecoder(r.Body).Decode(&t)

		var id int
		err := db.QueryRow(ctx,
			"INSERT INTO tasks (title, description) VALUES($1,$2) RETURNING id",
			&t.Title, &t.Description).Scan(&id)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Error to create tasks")
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
```

### `GetTask` — GET /tasks/{id}

```go
func GetTask(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := taskTracer.Start(r.Context(), "GetTask.QueryTasks")
		defer span.End()

		idConv, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			RespondError(w, http.StatusBadRequest, "Invalid task id")
			return
		}

		var t Task
		err = db.QueryRow(ctx, "SELECT id, title, description, status FROM tasks WHERE id = $1",
			idConv).Scan(&t.ID, &t.Title, &t.Description, &t.Status)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				RespondError(w, http.StatusNotFound, "Task not found")
				return
			}
			RespondError(w, http.StatusInternalServerError, "Error to get Tasks")
			return
		}
		RespondJSON(w, http.StatusOK, t)
	}
}
```

### `UpdateTask` — PUT /tasks/{id}

Body: `{"title": "...", "description": "...", "status": "..."}`

```go
func UpdateTask(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := taskTracer.Start(r.Context(), "UpdateTask.QueryTasks")
		defer span.End()

		idConv, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			RespondError(w, http.StatusBadRequest, "Invalid task id")
			return
		}

		var t Task
		json.NewDecoder(r.Body).Decode(&t)

		result, err := db.Exec(ctx,
			"UPDATE tasks SET title=$1, description=$2, status=$3 WHERE id=$4",
			t.Title, t.Description, t.Status, idConv)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Task doesn't exist")
			return
		}
		if result.RowsAffected() == 0 {
			RespondError(w, http.StatusNotFound, "Error to update task")
			return
		}
		RespondJSON(w, http.StatusOK, t)
	}
}
```

### `DeleteTask` — DELETE /tasks/{id}

```go
func DeleteTask(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := taskTracer.Start(r.Context(), "DeleteTask.QueryTasks")
		defer span.End()

		idConv, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			RespondError(w, http.StatusBadRequest, "Invalid task id")
			return
		}

		result, err := db.Exec(ctx, "DELETE FROM tasks WHERE id = $1", idConv)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "Error to delete task")
			return
		}
		if result.RowsAffected() == 0 {
			RespondError(w, http.StatusNotFound, "Error to delete task")
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
```

---

## Users API — CRUD

### Model

```go
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
```

### `CreateUser` — POST /users

**Validação**: Todos os campos obrigatórios. Email tem constraint UNIQUE — se duplicado, devolve 409.

```go
func CreateUser(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := userTracer.Start(r.Context(), "CreateUsers.QueryUsers")
		defer span.End()

		var req UserService
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			RespondError(w, http.StatusBadRequest, "invalid request body")
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
			return
		}

		var id string
		err = db.QueryRow(ctx,
			"INSERT INTO users(nome, email, password_hash) VALUES($1,$2,$3) RETURNING id",
			&req.Nome, &req.Email, PasswordHashed).Scan(&id)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				RespondError(w, http.StatusConflict, "email already exists")
				return
			}
			RespondError(w, http.StatusInternalServerError, "Error to create user")
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
```

### `ListUsers` — GET /users

Devolve apenas `nome`, `email` e `role` (nunca expõe `password_hash`).

```go
func ListUsers(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := userTracer.Start(r.Context(), "ListUsers.QueryUsers")
		defer span.End()

		rows, err := db.Query(ctx, "SELECT nome, email, role FROM users")
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Internal error")
			return
		}
		defer rows.Close()

		var users []UserResponse
		for rows.Next() {
			var u UserResponse
			rows.Scan(&u.Nome, &u.Email, &u.Role)
			users = append(users, u)
		}
		RespondJSON(w, http.StatusOK, users)
	}
}
```

### `UpdateUsers` — PATCH /users/{id}

**Importante**: as validações trocam de `||` (PUT completo) para `&&` (PATCH parcial) — "pelo menos um campo fornecido".

Usa `COALESCE(NULLIF(...))` para só atualizar campos preenchidos:

```sql
UPDATE users SET
  nome = COALESCE(NULLIF($1,''), nome),
  email = COALESCE(NULLIF($2,''), email),
  password_hash = COALESCE(NULLIF($3,''), password_hash)
WHERE id = $4
```

---

## Main — ponto de entrada

```go
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := LoadConfig(logger)
	if err != nil {
		logger.Error("Error to load .env", "error", err)
		os.Exit(1)
	}

	db, err := ConnectDB(cfg)
	if err != nil {
		logger.Error("error to connect database", "errors", err)
		os.Exit(1)
	}

	tp, err := InitTracer()
	if err != nil {
		logger.Error("error to initialize tracer", "error", err)
		os.Exit(1)
	}
	defer tp.Shutdown(context.Background())

	// List
	mux.HandleFunc("GET /tasks", ListTask(db, logger))
	mux.HandleFunc("GET /users", ListUsers(db, logger))
	// Create
	mux.HandleFunc("POST /tasks", CreateTask(db, logger))
	mux.HandleFunc("POST /users", CreateUser(db, logger))
	// Get by id
	mux.HandleFunc("GET /tasks/{id}", GetTask(db, logger))
	// Update
	mux.HandleFunc("PUT /tasks/{id}", UpdateTask(db, logger))
	mux.HandleFunc("PATCH /users/{id}", UpdateUsers(db, logger))
	// Delete
	mux.HandleFunc("DELETE /tasks/{id}", DeleteTask(db, logger))
	// Login
	mux.HandleFunc("POST /login", Login(cfg, db, logger))

	loggerMux := LoggingMiddleware(logger)(mux)
	logger.Info("-Server runnig...", "port", cfg.PORT)
	if err := http.ListenAndServe(":"+cfg.PORT, loggerMux); err != nil {
		logger.Error("error to innicialize server", "errors", err)
		os.Exit(1)
	}
}
```

### O padrão "logger como dependência"

```go
func X(db *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
```

A função recebe `db` e `logger` "de fora" e devolve o handler — uma **closure** que se "lembra" das dependências. Nunca variáveis globais, porque:
- Dificultam testes (não injetas logger diferente)
- Escondem dependências reais

---

## Testes

### `password_test.go`

```go
package main

import (
	"testing"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("123456")
	if err != nil {
		t.Fatal(err)
	}

	if hash == "" {
		t.Fatal("hash vazio")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("123456")); err != nil {
		t.Fatalf("hash não corresponde à senha original: %v", err)
	}
}
```

Para correr:

```bash
go test -v
```

---

## Como correr

### Pré-requisitos

- PostgreSQL
- Go 1.26+

### Setup

```bash
# clonar
git clone <repo>
cd task_manager

# configurar variaveis
cp .env.example .env
# editar .env com os teus dados

# correr migrations (manualmente ou com ferramenta de migracoes)

# iniciar servidor
go run .
```

### Variáveis de ambiente (`.env`)

```env
DATABASE_URL=postgres://user:password@localhost:5432/task_manager?sslmode=disable
SECRET_KEY=He_-OmLgdRqggyF0e8dPWrVO7IGZP4CCG5CAyhwaLd8
```

---

## Dependências

| Biblioteca | Uso |
| --- | --- |
| `github.com/jackc/pgx/v5` | Driver PostgreSQL com pool |
| `github.com/golang-jwt/jwt/v5` | Geração de tokens JWT |
| `golang.org/x/crypto` | bcrypt para passwords |
| `github.com/spf13/viper` | Config management |
| `github.com/joho/godotenv` | Load de `.env` |
| `github.com/google/uuid` | Geração de request IDs |
| `go.opentelemetry.io/otel` | Tracing (OpenTelemetry) |
