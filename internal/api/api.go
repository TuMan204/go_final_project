package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/TuMan204/go_final_project/internal/api/nextdate"
	"github.com/TuMan204/go_final_project/internal/db"
)

func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Query()

	now, err := time.Parse("20060102", request.Get("now"))
	if err != nil {
		now = time.Now().UTC()
	}
	date := request.Get("date")
	repeat := request.Get("repeat")

	nextDate, err := nextdate.NextDate(now, date, repeat)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte(nextDate))
}

func HandleAddTask(w http.ResponseWriter, r *http.Request) {
	var (
		task db.Task
		buf  bytes.Buffer
	)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJson(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(task.Title) == 0 {
		writeJson(w, "the title field is not filled in", http.StatusBadRequest)
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJson(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, err.Error(), http.StatusBadRequest)
		return
	}
	idStr := strconv.Itoa(int(id))
	writeJson(w, idStr, http.StatusOK)
}

func checkDate(task *db.Task) error {
	now := time.Now().UTC()

	if len(task.Date) == 0 {
		task.Date = now.Format("20060102")
		return nil
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

	if nextdate.AfterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format("20060102")
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

func writeJson(w http.ResponseWriter, data any, status int) {
	msg := make(map[string]any)
	if status != http.StatusOK {
		msg["error"] = data
	} else {
		msg["id"] = data
	}

	resp, err := json.Marshal(msg)
	if err != nil {
		writeJson(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	w.Write(resp)
}
