package main

import (
	"log"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    []any  `json:"data,omitempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	log.Printf("inside Handler : %s\n", app)
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
		Data:    []any{"ronald", "benjamin", 1, 2, 3, "bjkm", 909.9, jsonResponse{Error: true, Message: "I am 420", Data: []any{4, 2, 0}}},
	}

	_ = app.writeJson(w, http.StatusOK, payload)
}
