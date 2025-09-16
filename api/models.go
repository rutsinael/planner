package api

import "planner/db"

type PostTaskResponse struct {
	ID *int64 `json:"id"`
}

type GetTasksResponse struct {
	Tasks *[]db.Task `json:"tasks"`
	Error *string    `json:"error"`
}

type ErrorResponse struct {
	Error *string `json:"error"`
}
