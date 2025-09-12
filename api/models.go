package api

import "planner/db"

type PostTaskResponse struct {
	ID    *int64  `json:"id"`
	Error *string `json:"error"`
}

type GetTasksResponse struct {
	Tasks *[]db.Task `json:"tasks"`
	Error *string    `json:"error"`
}
