package utils

import (
	"encoding/json"
	"net/http"
)

type Envelope struct {
	Error string `json:"error"`
}

func WriteJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(&data)
}

func ReadJson(w http.ResponseWriter, r *http.Request, data any) error {
	if r.Body == nil {
		return http.ErrBodyNotAllowed
	}
	defer r.Body.Close()

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func WriteJsonError(w http.ResponseWriter, status int, message string) error {
	msg := Envelope{Error: message}
	return WriteJson(w, status, msg)
}
