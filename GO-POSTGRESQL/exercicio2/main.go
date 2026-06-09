package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	content, err := os.ReadFile("schema.sql")
	if err != nil {
		fmt.Printf("File not found\n DESCRIPTION OF ERROR: %v", err)
		return
	}

	db, err := sql.Open("postgres", "host=localhost user=postgres password=King dbname=task_tester sslmode=disable")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer db.Close()
	_, err = db.Exec(string(content))
	if err != nil {
		fmt.Println("SQL EXECUTION FAILED\n DESCRIPTION OF ERROR: ", err)
		return
	}
	fmt.Println("Schema executed succefully...")
}
