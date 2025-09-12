package server

import (
	"net/http"
	"planner/api"
)

func Run() {
	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	api.Init()
}
