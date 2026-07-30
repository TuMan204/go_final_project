package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Server struct {
	logger *log.Logger
	Serv   *http.Server
}

func NewServer(logger *log.Logger) Server {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web")))

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
