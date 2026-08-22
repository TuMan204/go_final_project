package main

import (
	"log"
	"os"

	"github.com/TuMan204/go_final_project/internal/db"
	"github.com/TuMan204/go_final_project/internal/server"
)

func main() {
	flog, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer flog.Close()

	logger := log.New(flog, "", log.LstdFlags)

	db, err := db.Init()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	server := server.NewServer(logger, db)

	logger.Println("Start server")
	defer logger.Println("Stop server")
	if err := server.Serv.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
