package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type RoleUser string

const (
	UserAdmin             RoleUser = "admin"
	UserUser              RoleUser = "user"
	maxActiveTasksPerUser          = 3
)

type Task struct {
	ID          int
	Title       string
	Description string
	Status      *string
	Tags        []string
	user_id     *string
}

type User struct {
	ID   string
	Nome string
	Role RoleUser
}

type TaskFilters struct {
	Status *string
}

type ErrUserNotFound struct {
	ID string
}

type ErrTaskNotFound struct {
	ID int
}

// If task not found
func (e *ErrTaskNotFound) Error() string {
	return fmt.Sprintf("task with id %d not found", e.ID)
}

// If user not found
func (u *ErrUserNotFound) Error() string {
	return fmt.Sprintf("user with id %s not found", u.ID)
}

type TaskService struct {
	DB *pgxpool.Pool
}

func NewTaskService(db *pgxpool.Pool) *TaskService {
	return &TaskService{DB: db}
}

type UserService struct {
	DB *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{DB: db}
}

// Pegar as TASKS
func (s *TaskService) GetTask(ctx context.Context, id int) (*Task, error) {
	var t Task
	err := s.DB.QueryRow(ctx, "SELECT id, title, description, status, tags, user_id FROM tasks WHERE id = $1", id).Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Tags, &t.user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &ErrTaskNotFound{ID: id}
		}
		return nil, err
	}
	return &t, nil
}

// Pegar os USERS
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
	var u User
	err := s.DB.QueryRow(ctx, "SELECT id, nome, role FROM users WHERE id = $1", id).Scan(&u.ID, &u.Nome, &u.Role)
	if err != nil {
		return nil, &ErrUserNotFound{ID: id} // ErrUserNotFound também precisa de ID string agora
	}
	return &u, nil
}

// Pegar as TASKS
func (s *TaskService) CreateTask(ctx context.Context, title string, description string, tags []string, user_id string) (int, error) {
	var id int

	//o context ja vem implementado no struct pelo pgxPool.Pool
	err := s.DB.QueryRow(ctx, "INSERT INTO tasks (title, description, tags, user_id) VALUES ($1, $2, $3, $4) RETURNING id", title, description, tags, user_id).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("could not insert task: %w", err)
	}
	return id, nil
}

func (s *UserService) CreateUser(ctx context.Context, nome string) (string, error) {
	var id string
	err := s.DB.QueryRow(ctx, "INSERT INTO users (nome) VALUES ($1) RETURNING id").Scan(&id, nome)
	if err != nil {
		return "invalid id", fmt.Errorf("could not insert user:%w", err)
	}
	return id, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id int, title string) error {
	//o context ja vem implementado no struct pelo pgxPool.Pool
	result, err := s.DB.Exec(ctx, "UPDATE tasks SET title = $1 WHERE id = $2", title, id)
	if err != nil {
		return err
	}

	rows := result.RowsAffected() //so retorna o int64 sem o error
	if rows == 0 {
		return fmt.Errorf("no task found with id %d", id)
	}
	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, nome string) error {
	result, err := s.DB.Exec(ctx, "UPDATE users SET nome = $1 WHERE id = $2", nome, id)
	if err != nil {
		return err
	}
	rows := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no user found with id %d", id)
	}
	return nil
}

func (s *TaskService) ListTasks(ctx context.Context) ([]Task, error) {
	//o context ja vem implementado no struct pelo pgxPool.Pool
	rows, err := s.DB.Query(ctx, "SELECT id, title, description, status, tags, user_id FROM tasks ORDER BY id ASC")
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Tags, &t.user_id)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return tasks, nil
}

func (s *UserService) ListUser(ctx context.Context) ([]User, error) {
	rows, err := s.DB.Query(ctx, "SELECT id, nome, role FROM users ORDER by id ASC")
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var t User
		err := rows.Scan(&t.ID, &t.Nome, &t.Role)
		if err != nil {
			return nil, fmt.Errorf("error scaning row: %w", err)
		}
		users = append(users, t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating row: %w", err)
	}
	return users, nil
}

func (s *TaskService) FetchTasks(ctx context.Context, filters TaskFilters, page int, pageSize int) ([]*Task, error) {
	query := "SELECT id, title, description, status, tags, user_id FROM tasks WHERE 1=1"
	args := []any{}
	argPos := 1

	if filters.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, *filters.Status)
		argPos++
	}

	// pagination
	offset := (page - 1) * pageSize
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, pageSize, offset)

	rows, err := s.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Tags, &t.user_id)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		tasks = append(tasks, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return tasks, nil
}

func (s *TaskService) CreateTaskWithQuota(ctx context.Context, userId string, title string, description string) (int, error) {

	tx, err := s.DB.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return 0, fmt.Errorf("Error to conect: %w", err)
	}
	defer tx.Rollback(ctx)

	var count int
	err = tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM tasks WHERE user_id = $1 AND status != 'done'`,
		userId,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count tasks: %w", err)
	}

	if count >= maxActiveTasksPerUser {
		return 0, fmt.Errorf("usuário %s já atingiu o limite de %d tasks ativas", userId, maxActiveTasksPerUser)
	}

	var tskId int
	err = tx.QueryRow(ctx, "INSERT INTO tasks (title, description, user_id) VALUES($1, $2, $3) RETURNING id", title, description, userId).Scan(&tskId)
	if err != nil {
		return 0, fmt.Errorf("Error to create task: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit tx: %w", err)
	}

	return tskId, nil
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env doesn't exist.")
		return
	}
	db, err := pgxpool.New(context.Background(), os.Getenv("CONNECTION_STRING"))
	if err != nil {
		panic(err)
	}

	// contexto só para o ping
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()

	if err := db.Ping(pingCtx); err != nil {
		fmt.Printf("unable to access database: %v\n", err)
		return
	}

	fmt.Println("connected to database")

	userservice := NewUserService(db)
	service := NewTaskService(db)

	//alexandreID := "692920e0-8c85-4feb-a8c8-1cdd5b2a38fb"
	ricardoID := "15d1cefa-ef63-4646-a451-af37e5fd7df7"

	for i := 1; i < 4; i++ {
		task, err := service.CreateTaskWithQuota(pingCtx, ricardoID, fmt.Sprintf("Teste quota nº %d", i), "teste quota")
		if err != nil {
			fmt.Printf("tentativa %d de criar a task %d: ERRO -> %v\n", i, task, err)
		}
	}
	task, err := service.ListTasks(pingCtx)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	for _, t := range task {
		nomeUser := "Empty"
		if t.user_id != nil {
			u, err := userservice.GetUser(pingCtx, *t.user_id)
			if err == nil {
				nomeUser = u.Nome
			}
		}
		fmt.Printf("ID: %d | Title: %s | Status: %v | Tags: %v | Created by: %v\n", t.ID, t.Title, *t.Status, t.Tags, nomeUser)
	}
}
