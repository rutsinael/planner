package main

import (
	"log"
	"net/http"
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

	if err = http.ListenAndServe(":7540", nil); err != nil {
		logger.Fatal("Error while server start: ", err)
	}

}
