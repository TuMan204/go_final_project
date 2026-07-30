package main

import (
	"log"

	"github.com/TuMan204/go_final_project/internal/server"
)

func main() {
	var logger log.Logger
	server := server.NewServer(&logger)

	err := server.Serv.ListenAndServe()
	if err != nil {
		logger.Fatal()
	}
}
