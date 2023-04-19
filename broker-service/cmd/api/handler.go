package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayLoad `json:"auth,omitempty"`
}

type AuthPayLoad struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	log.Printf("inside Broker Handler : %s\n", app)
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
		Data:    []any{"ronald", "benjamin", 1, 2, 3, "bjkm", 909.9, jsonResponse{Error: true, Message: "I am 420", Data: []any{4, 2, 0}}},
	}

	_ = app.writeJson(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		app.errorJSON(w, errors.New("unknown acrion"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayLoad) {
	log.Println("inside authentication call")
	jsonData, _ := json.Marshal(a)
	req, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer res.Body.Close()
	// make sure we get back the correct status code
	log.Println("_________________PPPPPPPPPPPPPPP_____________")
	log.Println("STATUS CODE:", res.StatusCode)
	if res.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if res.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return

	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(res.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJson(w, http.StatusAccepted, payload)
}
