package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"planner/db"
	"planner/repeater"
	"time"
)

func Init(r *chi.Mux) {
	r.Get("/api/nextdate", nextDayHandler)

	r.Get("/api/task", getTaskHandler)
	r.Post("/api/task", addTaskHandler)
	r.Put("/api/task", updateTaskHandler)
	r.Delete("/api/task", deleteTaskHandler)

	r.Get("/api/tasks", getTasksHandler)
	r.Post("/api/task/done", taskDoneHandler)
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	dataJson, err := json.Marshal(data)
	if err == nil {
		w.Write(dataJson)
	}
}

func writeError(w http.ResponseWriter, errMessage string) {
	w.WriteHeader(http.StatusBadRequest)
	writeJson(w, ErrorResponse{Error: &errMessage})
}

func writeEmptySuccessResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	writeJson(w, map[string]any{})
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, err.Error())
		return
	}

	task, err := db.GetTaskByID(id)
	if err != nil {
		writeError(w, "Задача не найдена")
		return
	}

	if task.Repeat != "" {
		nextDate, err := repeater.GetNextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeError(w, err.Error())
			return
		}
		db.UpdateTask(&db.Task{
			ID:      task.ID,
			Date:    nextDate,
			Repeat:  task.Repeat,
			Comment: task.Comment,
			Title:   task.Title,
		})

	} else {
		db.DeleteTask(id)
	}
	writeEmptySuccessResponse(w)
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {

	tasks, err := db.GetTasks(tasksMax)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJson(w, GetTasksResponse{
		Tasks: &tasks,
	})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	task, err := db.GetTaskByID(id)
	if err != nil {
		writeError(w, "Задача не найдена")
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJson(w, task)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := addTask(r)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJson(w, PostTaskResponse{ID: &id})
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	err := updateTask(r)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	writeEmptySuccessResponse(w)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	err = db.DeleteTask(id)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	writeEmptySuccessResponse(w)
}

func nextDayHandler(w http.ResponseWriter, request *http.Request) {
	now := request.FormValue("now")
	date := request.FormValue("date")
	repeat := request.FormValue("repeat")

	var timeNow time.Time

	if len(now) > 0 {
		timeNow, _ = time.Parse(repeater.DateFormat, now)
	} else {
		timeNow = time.Now()
	}

	nextDate, err := repeater.GetNextDate(timeNow, date, repeat)
	if err != nil {
		writeError(w, err.Error())
		return
	}

	fmt.Fprintf(w, nextDate)
}
