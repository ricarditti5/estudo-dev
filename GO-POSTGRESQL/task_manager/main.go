package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type TaskService struct {
	DB *sql.DB
}

type Task struct {
	ID          int
	Title       string
	Description string
	Status      *string
}

type TaskFilters struct {
	Status *string
}

type ErrTaskNotFound struct {
	ID int
}

func (e *ErrTaskNotFound) Error() string {
	return fmt.Sprintf("task with id %d not found", e.ID)
}

func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{DB: db}
}

func (s *TaskService) GetTask(ctx context.Context, id int) (*Task, error) {
	var t Task
	err := s.DB.QueryRowContext(ctx, "SELECT id, title, description, status FROM tasks WHERE id = $1", id).Scan(&t.ID, &t.Title, &t.Description, &t.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &ErrTaskNotFound{ID: id}
		}
		return nil, err
	}
	return &t, nil
}

func (s *TaskService) CreateTask(ctx context.Context, title string, description string) (int, error) {
	var id int
	err := s.DB.QueryRowContext(ctx, "INSERT INTO tasks (title, description) VALUES ($1, $2) RETURNING id", title, description).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("could not insert task: %w", err)
	}
	return id, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id int, title string) error {
	result, err := s.DB.ExecContext(ctx, "UPDATE tasks SET title = $1 WHERE id = $2", title, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no task found with id %d", id)
	}
	return nil
}

func (s *TaskService) ListTasks(ctx context.Context) ([]Task, error) {
	rows, err := s.DB.QueryContext(ctx, "SELECT id, title, description, status FROM tasks ORDER BY status = 'done' ASC")
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status)
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

func (s *TaskService) FetchTasks(ctx context.Context, filters TaskFilters, page int, pageSize int) ([]*Task, error) {
	query := "SELECT id, title, description, status FROM tasks WHERE 1=1"
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

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status)
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

func main() {
	connStr := "user=postgres password=King dbname=task_manager sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	// contexto só para o ping
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()

	if err := db.PingContext(pingCtx); err != nil {
		fmt.Printf("unable to access database: %v\n", err)
		return
	}

	fmt.Println("connected to database")

	service := NewTaskService(db)

	// página 1
	tasks1, err := service.FetchTasks(context.Background(), TaskFilters{}, 1, 10)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Println("--- Página 1 ---")
	for _, t := range tasks1 {
		fmt.Printf("ID: %d | Title: %s | Status: %v\n", t.ID, t.Title, *t.Status)
	}

	// página 2
	tasks2, err := service.FetchTasks(context.Background(), TaskFilters{}, 2, 10)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Println("--- Página 2 ---")
	for _, t := range tasks2 {
		fmt.Printf("ID: %d | Title: %s | Status: %v\n", t.ID, t.Title, *t.Status)
	}

	// página 1 filtrado por status "done"
	status := "done"
	tasksStatus, err := service.FetchTasks(context.Background(), TaskFilters{Status: &status}, 1, 10)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Println("--- Página 1 com status done ---")
	for _, t := range tasksStatus {
		fmt.Printf("ID: %d | Title: %s | Status: %v\n", t.ID, t.Title, *t.Status)
	}

	task, err := service.ListTasks(pingCtx)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	for _, t := range task {
		fmt.Printf("ID: %d | Title: %s | Status: %v\n", t.ID, t.Title, *t.Status)
	}
}
