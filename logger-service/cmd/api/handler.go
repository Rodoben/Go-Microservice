package main

import (
	"log"
	"loggerservice/data"
	"net/http"
)

type Jsonpayload struct {
	Name string `json:"name"`
	Data string `json:"data" `
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestpayload Jsonpayload
	_ = app.readJson(w, r, &requestpayload)

	//insert data

	event := data.LogEntry{
		Name: requestpayload.Name,
		Data: requestpayload.Data,
	}

	log.Println("EVENT:", event)

	err := app.Model.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJson(w, http.StatusAccepted, resp)
}
