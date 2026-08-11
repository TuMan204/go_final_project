package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/TuMan204/go_final_project/internal/db"
	"github.com/TuMan204/go_final_project/internal/server"
)

func main() {
	var logger log.Logger

	execPath, err := os.Executable()
	if err != nil {
		logger.Fatal()
	}
	rootPath := filepath.Dir(execPath)
	webPath := filepath.Join(rootPath, "web")

	dbFile := "scheduler.db"
	envDBFile := os.Getenv("TODO_DBFILE")
	if len(envDBFile) > 0 {
		dbFile = envDBFile
	}
	dbPath := filepath.Join(rootPath, dbFile)

	if err := db.Init(dbPath); err != nil {
		logger.Fatal()
	}

	server := server.NewServer(&logger, webPath)

	if err := server.Serv.ListenAndServe(); err != nil {
		logger.Fatal()
	}
}
