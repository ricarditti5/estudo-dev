package main

import (
	"fmt"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	db, err := ConnectDB()
	if err != nil {
		fmt.Printf("Error to connect to db: %v", err)
	}
	mux.HandleFunc("GET /tasks", ListTask(db))
	mux.HandleFunc("POST /tasks", CreateTask(db))
	mux.HandleFunc("GET /tasks/{id}", GetTask(db))
	mux.HandleFunc("PUT /tasks/{id}", UpdateTask(db))
	mux.HandleFunc("DELETE /tasks/{id}", DeleteTask(db))
	fmt.Println("Server runnig...")
	http.ListenAndServe(":8080", mux)
}
