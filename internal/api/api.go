package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/TuMan204/go_final_project/internal/api/auth"
	"github.com/TuMan204/go_final_project/internal/api/nextdate"
	"github.com/TuMan204/go_final_project/internal/db"
)

const (
	dateFormat = "20060102"
	limit      = 10
)

var (
	errIdentifierNotSpecified = errors.New("identifier is not specified")
	errTitleNotFilled         = errors.New("the title field is not filled in")
)

func checkDate(task *db.Task) error {
	now := time.Now()

	if len(task.Date) == 0 {
		task.Date = now.Format(dateFormat)
		return nil
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return err
	}

	if nextdate.AfterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(dateFormat)
		} else {
			next, err := nextdate.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	}
	return nil
}

func writeErrorJSON(w http.ResponseWriter, data any, status int) {
	msg := make(map[string]any)
	msg["error"] = data

	resp, err := json.Marshal(msg)
	if err != nil {
		writeErrorJSON(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	w.Write(resp)
}

func writeJSON(w http.ResponseWriter, data any, status int) {
	resp, err := json.Marshal(data)
	if err != nil {
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	w.Write(resp)
}

type TaskService struct {
	store  *db.TasksStore
	logger *log.Logger
}

func NewTaskService(store *db.TasksStore, logger *log.Logger) TaskService {
	return TaskService{store: store, logger: logger}
}

func (s *TaskService) HandleNextDate(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Query()

	nowStr := request.Get("now")
	now, err := time.Parse(dateFormat, nowStr)
	if err != nil {
		now = time.Now().UTC()
	}
	date := request.Get("date")
	repeat := request.Get("repeat")

	nextDate, err := nextdate.NextDate(now, date, repeat)
	if err != nil {
		s.logger.Printf("HandleNextDate: ERR: %s\n", err.Error())
		w.Write([]byte(err.Error()))
		return
	}

	s.logger.Println("HandleNextDate: OK")
	w.Write([]byte(nextDate))
}

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func (s *TaskService) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Query()
	searchStr := request.Get("search")

	tasks, err := s.store.GetTasks(searchStr, limit)
	if err != nil {
		s.logger.Printf("HandleGetTasks: ERR: %s\n", err.Error())
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.logger.Println("HandleGetTasks: OK")
	writeJSON(w, TasksResp{Tasks: tasks}, http.StatusOK)
}

func (s *TaskService) HandleAddTask(w http.ResponseWriter, r *http.Request) {
	var (
		task db.Task
		buf  bytes.Buffer
	)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		s.logger.Printf("HandleAddTask: ERR: %s\n", err.Error())
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		s.logger.Printf("HandleAddTask: ERR: %s\n", err.Error())
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(task.Title) == 0 {
		s.logger.Printf("HandleAddTask: ERR: %s\n", errTitleNotFilled.Error())
		writeErrorJSON(w, errTitleNotFilled.Error(), http.StatusBadRequest)
		return
	}

	err = checkDate(&task)
	if err != nil {
		s.logger.Printf("HandleAddTask: ERR: %s\n", err.Error())
		writeErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := s.store.AddTask(&task)
	if err != nil {
		s.logger.Printf("HandleAddTask: ERR: %s\n", err.Error())
		writeErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	idStr := strconv.Itoa(int(id))
	s.logger.Println("HandleAddTask: OK")
	writeJSON(w, map[string]string{"id": idStr}, http.StatusOK)
}

func (s *TaskService) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		s.logger.Printf("HandleGetTask: ERR: %s\n", errIdentifierNotSpecified.Error())
		writeErrorJSON(w, errIdentifierNotSpecified.Error(), http.StatusBadRequest)
		return
	}

	task, err := s.store.GetTask(id)
	if err != nil {
		s.logger.Printf("HandleGetTask: ERR: %s\n", err.Error())
		writeErrorJSON(w, err.Error(), http.StatusNotFound)
		return
	}
	s.logger.Println("HandleGetTask: OK")
	writeJSON(w, task, http.StatusOK)
}

func (s *TaskService) HandleEditTask(w http.ResponseWriter, r *http.Request) {
	var (
		task db.Task
		buf  bytes.Buffer
	)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		s.logger.Printf("HandleEditTask: ERR: %s\n", err.Error())
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		s.logger.Printf("HandleEditTask: ERR: %s\n", err.Error())
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(task.Title) == 0 {
		s.logger.Printf("HandleEditTask: ERR: %s\n", errTitleNotFilled.Error())
		writeErrorJSON(w, errTitleNotFilled.Error(), http.StatusBadRequest)
		return
	}

	err = checkDate(&task)
	if err != nil {
		s.logger.Printf("HandleEditTask: ERR: %s\n", err.Error())
		writeErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.store.UpdateTask(&task)
	if err != nil {
		s.logger.Printf("HandleEditTask: ERR: %s\n", err.Error())
		writeErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.logger.Println("HandleEditTask: OK")
	writeJSON(w, map[string]string{}, http.StatusOK)
}

func (s *TaskService) HandleTaskDone(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		s.logger.Printf("HandleTaskDone: ERR: %s\n", errIdentifierNotSpecified.Error())
		writeErrorJSON(w, errIdentifierNotSpecified.Error(), http.StatusBadRequest)
		return
	}

	task, err := s.store.GetTask(id)
	if err != nil {
		s.logger.Printf("HandleTaskDone: ERR: %s\n", err.Error())
		writeErrorJSON(w, err.Error(), http.StatusNotFound)
		return
	}

	if len(task.Repeat) == 0 {
		err = s.store.DeleteTask(id)
		if err != nil {
			s.logger.Printf("HandleTaskDone: ERR: %s\n", err.Error())
			writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		now := time.Now()

		newDate, err := nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			s.logger.Printf("HandleTaskDone: ERR: %s\n", err.Error())
			writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = s.store.UpdateDate(newDate, id)
		if err != nil {
			s.logger.Printf("HandleTaskDone: ERR: %s\n", err.Error())
			writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	s.logger.Println("HandleTaskDone: OK")
	writeJSON(w, map[string]string{}, http.StatusOK)
}

func (s *TaskService) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		s.logger.Printf("HandleDeleteTask: ERR: %s\n", errIdentifierNotSpecified.Error())
		writeErrorJSON(w, errIdentifierNotSpecified.Error(), http.StatusBadRequest)
		return
	}

	err := s.store.DeleteTask(id)
	if err != nil {
		s.logger.Printf("HandleDeleteTask: ERR: %s\n", err.Error())
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.logger.Println("HandleDeleteTask: OK")
	writeJSON(w, map[string]string{}, http.StatusOK)
}

func (s *TaskService) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	var token struct {
		Token string `json:"token"`
	}

	var pass struct {
		Pass string `json:"password"`
	}
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		s.logger.Printf("HandleSignIn: ERR: %s\n", err.Error())
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &pass); err != nil {
		s.logger.Printf("HandleSignIn: ERR: %s\n", err.Error())
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	envPass := os.Getenv("TODO_PASSWORD")
	if pass.Pass != envPass {
		s.logger.Printf("HandleSignIn: ERR: %s\n", "Invalid password")
		writeErrorJSON(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token.Token, err = auth.GenerateJWT(pass.Pass)
	if err != nil {
		s.logger.Printf("HandleSignIn: ERR: %s\n", err.Error())
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.logger.Println("HandleSignIn: OK")
	writeJSON(w, token, http.StatusOK)
}
