package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/TuMan204/go_final_project/internal/api"
)

type Server struct {
	logger *log.Logger
	Serv   *http.Server
}

func NewServer(logger *log.Logger, webPath string) Server {
	mux := http.NewServeMux()
	//mux.Handle("/", http.FileServer(http.Dir(webPath)))	// Основное API
	mux.Handle("/", http.FileServer(http.Dir("./web"))) // Для работы тестов workflow при команде go run main.go
	mux.HandleFunc("/api/nextdate", api.HandleNextDate)
	mux.HandleFunc("GET  /api/tasks", api.HandleGetTasks)
	mux.HandleFunc("POST /api/task", api.HandleAddTask)

	addr := 7540
	envPort := os.Getenv("TODO_PORT")
	if len(envPort) > 0 {
		port, err := strconv.Atoi(envPort)
		if err == nil {
			addr = port
		}
	}

	serv := http.Server{
		Addr:     fmt.Sprintf(":%d", addr),
		Handler:  mux,
		ErrorLog: logger,
		//ReadTimeout:  10 * time.Second,
		//WriteTimeout: 10 * time.Second,
		//IdleTimeout:  15 * time.Second,
	}

	return Server{logger: logger, Serv: &serv}
}
