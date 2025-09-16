package server

import (
	"fmt"
	"net/http"
	"planner/api"

	"github.com/go-chi/chi/v5"
)

func Run() {
	webDir := "./web"

	r := chi.NewRouter()
	//http.Handle("/", http.FileServer(http.Dir(webDir)))

	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	api.Init(r)

	if err := http.ListenAndServe(":7540", r); err != nil {
		fmt.Printf("Error while server start: %s", err)
	}
}
