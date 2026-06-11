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
	err := s.DB.QueryRowContext(ctx, "SELECT id, title, description FROM tasks WHERE id = $1", id).Scan(&t.ID, &t.Title, &t.Description)
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
	rows, err := s.DB.QueryContext(ctx, "SELECT id, title, description FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("Error to execute the Query %w", err)
	}
	defer rows.Close()
	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description)
		if err != nil {
			return nil, fmt.Errorf("Error to scan row %w", err)
		}
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error to iterate rows %w", err)
	}
	return tasks, nil
}

func main() {
	connStr := "user=postgres password=King dbname=task_manager sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		fmt.Printf("Unable to acces database %v\n", err)
		return
	}
	fmt.Println("Hello task manager")
}
