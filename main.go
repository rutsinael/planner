package main

import (
	"log"
	"os"
	"planner/db"
	"planner/server"
)

func main() {

	logger := log.New(os.Stdout, "INFO: ", log.LstdFlags)

	err := db.Init("scheduler.db")
	if err != nil {
		logger.Fatal(err)
	}

	server.Run()
}
