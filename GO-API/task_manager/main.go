package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
)

func main() {

	port := 8080
	mux := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db, err := ConnectDB(logger)
	if err != nil {
		logger.Error("error to connect database", "errors", err)
		os.Exit(1)
	}

	tp, err := InitTracer() // chama o tracer que acabei de criar, e devolve o TracerProvider já configurado e registado globalmente.
	if err != nil {
		logger.Error("error to initialize tracer", "error", err)
		os.Exit(1)
	}
	defer tp.Shutdown(context.Background()) // usamos aqui um contexto "vazio"/raiz, porque este Shutdown não está associado a nenhum pedido HTTP específico — é uma operação de ciclo de vida da aplicação, não de um pedido.
	//---------------------------
	/* isto é crucial. O defer garante que, quando main() terminar (por qualquer motivo), o Shutdown é chamado, o que força o envio/impressão de quaisquer spans que ainda estejam em buffer no exporter. Sem isto, podias perder os últimos traces gerados antes do programa fechar.*/

	//List
	mux.HandleFunc("GET /tasks", ListTask(db, logger))
	mux.HandleFunc("GET /users", ListUsers(db, logger))
	//Create
	mux.HandleFunc("POST /tasks", CreateTask(db, logger))
	mux.HandleFunc("POST /users", CreateUser(db, logger))

	//Get task by id
	mux.HandleFunc("GET /tasks/{id}", GetTask(db, logger))
	//Update
	mux.HandleFunc("PUT /tasks/{id}", UpdateTask(db, logger))
	mux.HandleFunc("PATCH /users/{id}", UpdateUsers(db, logger))

	//Delete by id
	mux.HandleFunc("DELETE /tasks/{id}", DeleteTask(db, logger))

	//Login
	mux.HandleFunc("POST /login", Login(db, logger))

	loggerMux := LoggingMiddleware(logger)(mux)
	logger.Info("-Server runnig...", "port", port)
	if err := http.ListenAndServe(":8080", loggerMux); err != nil {
		logger.Error("error to innicialize server", "errors", err)
		os.Exit(1)
	}
}
