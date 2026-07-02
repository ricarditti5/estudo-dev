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
	//List
	mux.HandleFunc("GET /tasks", ListTask(db))
	mux.HandleFunc("GET /users", ListUsers(db))
	//Create
	mux.HandleFunc("POST /tasks", CreateTask(db))
	mux.HandleFunc("POST /users", CreateUser(db))

	//Get task by id
	mux.HandleFunc("GET /tasks/{id}", GetTask(db))
	//Update
	mux.HandleFunc("PUT /tasks/{id}", UpdateTask(db))
	mux.HandleFunc("PATCH /users/{id}", UpdateUsers(db))

	//Delete by id
	mux.HandleFunc("DELETE /tasks/{id}", DeleteTask(db))

	//Login
	mux.HandleFunc("POST /login", Login(db))

	fmt.Println("Server runnig...")
	http.ListenAndServe(":8080", mux)
}
