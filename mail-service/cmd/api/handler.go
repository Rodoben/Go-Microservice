package main

import (
	"log"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {

	log.Println("Inside mail handler")

	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}
	log.Println("REQUEST:", r)
	var requestPayload mailMessage

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	log.Println("ERRJSON:", err)

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	log.Println("ERRSMTP:", err)
	payload := jsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}
	app.writeJson(w, http.StatusAccepted, payload)
}
