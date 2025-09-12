package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"planner/db"
	"planner/repeater"
	"strconv"
	"time"
)

func getIDFromRequest(r *http.Request) (int, error) {
	id := r.FormValue("id")

	if id == "" {
		return 0, errors.New("Не указан идентификатор")
	} else {
		i, err := strconv.Atoi(id)
		if err != nil {
			return 0, errors.New("Некорректный формат id")
		} else {
			return i, nil
		}
	}
}

func validateTask(r *http.Request) (*db.Task, error) {
	var task db.Task

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		return nil, errors.New("error reading body")
	}

	err = json.Unmarshal(body, &task)

	if err != nil {
		return nil, errors.New("error reading body")
	}

	if len(task.Title) == 0 {
		return nil, errors.New("task title is empty")
	}

	err = checkDate(&task)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func updateTask(r *http.Request) error {
	task, err := validateTask(r)

	if err != nil {
		return err
	}
	return db.UpdateTask(task)
}

func addTask(r *http.Request) (int64, error) {

	task, err := validateTask(r)

	if err != nil {
		return -1, err
	}

	return db.AddTask(task)
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if len(task.Date) == 0 {
		task.Date = time.Now().Format(repeater.DateFormat)
	}

	date, err := time.Parse(repeater.DateFormat, task.Date)
	if err != nil {
		return errors.New("couldn't parse date")
	}

	var nextDate string
	if len(task.Repeat) != 0 {
		nextDate, err = repeater.GetNextDate(now, task.Date, task.Repeat)
		if err != nil {
			return errors.New("error getting next date")
		}
	}

	if repeater.AfterNow(now, date) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(repeater.DateFormat)
		} else {
			task.Date = nextDate
		}
	}

	return nil
}
