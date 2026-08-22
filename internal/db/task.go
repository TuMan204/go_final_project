package db

import (
	"database/sql"
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

type TasksStore struct {
	db *sql.DB
}

func NewTasksStore(db *sql.DB) TasksStore {
	return TasksStore{db: db}
}

func (s *TasksStore) AddTask(task *Task) (int64, error) {
	var id int64

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)"
	res, err := s.db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func (s *TasksStore) GetTasks(search string, limit int) ([]*Task, error) {
	tasks := make([]*Task, 0)

	query := "SELECT * FROM scheduler ORDER BY date LIMIT :limit"
	if len(search) > 0 {
		date, err := time.Parse("02.01.2006", search)
		if err != nil {
			query = "SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit"
			search = "%" + search + "%"
		} else {
			query = "SELECT * FROM scheduler WHERE date = :search ORDER BY title LIMIT :limit"
			search = date.Format("20060102")
		}
	}
	rows, err := s.db.Query(query, sql.Named("limit", limit), sql.Named("search", search))
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

func (s *TasksStore) GetTask(id string) (*Task, error) {
	var task Task

	query := "SELECT * FROM scheduler WHERE id = :id"
	err := s.db.QueryRow(query, sql.Named("id", id)).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TasksStore) UpdateTask(task *Task) error {
	query := "UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id"
	res, err := s.db.Exec(query,
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

func (s *TasksStore) UpdateDate(newDate string, id string) error {
	query := "UPDATE scheduler SET date = :date WHERE id = :id"
	res, err := s.db.Exec(query,
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

func (s *TasksStore) DeleteTask(id string) error {
	query := "DELETE FROM scheduler WHERE id = :id"
	res, err := s.db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for delete task`)
	}

	return nil
}
