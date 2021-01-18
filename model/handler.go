package model

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (q *Queue) UploadFile(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go q.PushVal(handler.Filename, content)

	return
}

func (q *Queue) CheckStatus(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")

	taskStatus := q.Status[filename]
	if err := json.NewEncoder(w).Encode(map[string]string{"status": taskStatus}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
