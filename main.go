package main

import (
	"log"

	"github.com/TuMan204/go_final_project/internal/server"
	"github.com/TuMan204/go_final_project/pkg/db"
)

func main() {
	var logger log.Logger

	dbFile := "scheduler.db"
	if err := db.Init(dbFile); err != nil {
		logger.Fatal()
	}

	server := server.NewServer(&logger)

	if err := server.Serv.ListenAndServe(); err != nil {
		logger.Fatal()
	}
}
