package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/TuMan204/go_final_project/internal/api"
	"github.com/TuMan204/go_final_project/internal/api/auth"
	"github.com/TuMan204/go_final_project/internal/db"
)

type Server struct {
	logger *log.Logger
	Serv   *http.Server
}

func NewServer(logger *log.Logger, database *sql.DB) Server {
	store := db.NewTasksStore(database)
	service := api.NewTaskService(&store, logger)

	pass := os.Getenv("TODO_PASSWORD")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web")))
	mux.HandleFunc("GET /api/nextdate", service.HandleNextDate)
	mux.HandleFunc("GET /api/tasks", auth.Auth(service.HandleGetTasks, pass))
	mux.HandleFunc("POST /api/task", auth.Auth(service.HandleAddTask, pass))
	mux.HandleFunc("GET /api/task", auth.Auth(service.HandleGetTask, pass))
	mux.HandleFunc("PUT /api/task", auth.Auth(service.HandleEditTask, pass))
	mux.HandleFunc("DELETE /api/task", auth.Auth(service.HandleDeleteTask, pass))
	mux.HandleFunc("POST /api/task/done", auth.Auth(service.HandleTaskDone, pass))
	mux.HandleFunc("POST /api/signin", service.HandleSignIn)

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
	}

	return Server{logger: logger, Serv: &serv}
}
