package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"planner/db"
	"strconv"
	"time"
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", getTasksHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	dataJson, err := json.Marshal(data)
	if err == nil {
		w.Write(dataJson)
	}
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			errMessage := "Не указан идентификатор"
			writeJson(w, PostTaskResponse{Error: &errMessage})
		} else {
			i, err := strconv.Atoi(id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				errMessage := "Некорректный формат id"
				writeJson(w, PostTaskResponse{Error: &errMessage})
			} else {
				task, err := db.GetTaskByID(i)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					errMessage := "Задача не найдена"
					writeJson(w, PostTaskResponse{Error: &errMessage})
				} else {
					if task.Repeat != "" {
						nextDate, err := db.GetNextDate(time.Now(), task.Date, task.Repeat)
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							errMessage := err.Error()
							writeJson(w, PostTaskResponse{Error: &errMessage})
						} else {
							newTask := db.Task{
								ID:      task.ID,
								Date:    nextDate,
								Repeat:  task.Repeat,
								Comment: task.Comment,
								Title:   task.Title,
							}
							db.UpdateTask(&newTask)
							w.WriteHeader(http.StatusOK)
							writeJson(w, map[string]any{})
						}

					} else {
						db.DeleteTask(i)
						w.WriteHeader(http.StatusOK)
						writeJson(w, map[string]any{})
					}
				}
			}
		}
	}
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tasks, err := db.GetTasks(tasksMax)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errMessage := err.Error()
			writeJson(w, GetTasksResponse{
				Error: &errMessage,
			})
		} else {
			w.WriteHeader(http.StatusOK)
			writeJson(w, GetTasksResponse{
				Tasks: &tasks,
			})
		}
	}
}

func taskHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		id, err := addTask(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errMessage := err.Error()
			writeJson(w, PostTaskResponse{Error: &errMessage})
		} else {
			w.WriteHeader(http.StatusOK)
			writeJson(w, PostTaskResponse{ID: &id})
		}

	case http.MethodGet:
		id := r.FormValue("id")

		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			errMessage := "Не указан идентификатор"
			writeJson(w, PostTaskResponse{Error: &errMessage})
		} else {
			i, err := strconv.Atoi(id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				errMessage := "Некорректный формат id"
				writeJson(w, PostTaskResponse{Error: &errMessage})
			} else {
				task, err := db.GetTaskByID(i)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					errMessage := "Задача не найдена"
					writeJson(w, PostTaskResponse{Error: &errMessage})
				} else {
					w.WriteHeader(http.StatusOK)
					writeJson(w, task)
				}
			}
		}
	case http.MethodPut:
		err := updateTask(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errMessage := err.Error()
			writeJson(w, PostTaskResponse{Error: &errMessage})
		} else {
			w.WriteHeader(http.StatusOK)
			writeJson(w, map[string]any{})
		}
	case http.MethodDelete:
		id := r.FormValue("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			errMessage := "Не указан идентификатор"
			writeJson(w, PostTaskResponse{Error: &errMessage})
		} else {
			i, err := strconv.Atoi(id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				errMessage := "Некорректный формат id"
				writeJson(w, PostTaskResponse{Error: &errMessage})
			} else {
				err := db.DeleteTask(i)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					errMessage := "Задача не найдена"
					writeJson(w, PostTaskResponse{Error: &errMessage})
				} else {
					w.WriteHeader(http.StatusOK)
					writeJson(w, map[string]any{})
				}
			}
		}
	}
}

func nextDayHandler(writer http.ResponseWriter, request *http.Request) {

	if request.Method == http.MethodGet {
		now := request.FormValue("now")
		date := request.FormValue("date")
		repeat := request.FormValue("repeat")

		timeNow, _ := time.Parse("20060102", now)

		nextDate, err := db.GetNextDate(timeNow, date, repeat)
		if err != nil {
			return
		}

		fmt.Fprintf(writer, nextDate)

	} else if request.Method == http.MethodPost {

	}
}
