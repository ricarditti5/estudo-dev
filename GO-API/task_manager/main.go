package main

import (
	"fmt"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", ListTask)
	fmt.Println("Server runnig...")
	http.ListenAndServe(":8080", mux)
}
