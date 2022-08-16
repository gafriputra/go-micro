package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) (err error) {
	maxBytes := 1048576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err = dec.Decode(data)
	if err != nil {
		return
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		err = errors.New("body must have only a single JSON value")
		return
	}

	err = nil
	return
}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) (err error) {
	out, err := json.Marshal(data)
	if err != nil {
		return
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	return
}

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	return app.writeJSON(w, statusCode, jsonResponse{
		Error:   true,
		Message: err.Error(),
	})
}
