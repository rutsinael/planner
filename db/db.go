package db

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
)

var db *sql.DB

const schema = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR,
    comment TEXT,
    repeat VARCHAR);
CREATE INDEX IX_DATE
ON scheduler (date)`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)

	if install {
		_, err = db.Exec(schema)
	}

	return err
}

func AddTask(task *Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) values ($1, $2, $3, $4)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func GetTaskByID(id int) (Task, error) {
	var task Task
	row := db.QueryRow("SELECT * FROM scheduler where id=$1", id)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func GetTasks(max int) ([]Task, error) {
	var tasks = make([]Task, 0)

	rows, err := db.Query("SELECT * FROM scheduler ORDER BY date ASC LIMIT $1", max)
	if err != nil {
		return tasks, err
	}

	defer rows.Close()
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date=$1, title=$2, comment=$3, repeat=$4 WHERE id=$5`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id int) error {
	query := `DELETE FROM scheduler WHERE id=$1`
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
