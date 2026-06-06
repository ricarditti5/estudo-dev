package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// The driver name for pgx/v5/stdlib is "pgx"
	db, err := sql.Open("pgx", "postgres://postgres:King@localhost:5432/task_manager")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connection succefull")
}
