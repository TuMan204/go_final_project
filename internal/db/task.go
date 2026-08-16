package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return id, err
	}
	defer db.Close()

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)"
	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func GetTasks(search string, limit int) ([]*Task, error) {
	tasks := make([]*Task, 0)

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT * FROM scheduler ORDER BY date DESC LIMIT :limit"
	if len(search) > 0 {
		date, err := time.Parse("02.01.2006", search)
		if err != nil {
			query = "SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date DESC LIMIT :limit"
			search = "%" + search + "%"
		} else {
			query = "SELECT * FROM scheduler WHERE date = :search LIMIT :limit"
			search = date.Format("20060102")
		}
	}
	rows, err := db.Query(query, sql.Named("limit", limit), sql.Named("search", search))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var task Task

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT * FROM scheduler WHERE id = :id"
	err = db.QueryRow(query, sql.Named("id", id)).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, errors.New("task not found")
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	query := "UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id"
	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
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

func UpdateDate(newDate string, id string) error {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	query := "UPDATE scheduler SET date = :date WHERE id = :id"
	res, err := db.Exec(query,
		sql.Named("date", newDate),
		sql.Named("id", id))
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

func DeleteTask(id string) error {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = GetTask(id)
	if err != nil {
		return err
	}

	query := "DELETE FROM scheduler WHERE id = :id"
	_, err = db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}

	return nil
}
