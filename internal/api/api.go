package api

import (
	"bytes"
	"encoding/json"
	"net/http"
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
	var task *db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
