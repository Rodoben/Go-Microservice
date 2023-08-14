package main

import (
	"broker/event"
	"broker/logs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/rpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayLoad `json:"auth,omitempty"`
	Log    LogPayLoad  `json:"log,omitempty"`
	Mail   Mailpayload `json:"mail,omitempty"`
}

type AuthPayLoad struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LogPayLoad struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type Mailpayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
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

	log.Println("REQUEST:", r)
	var requestPayload RequestPayload
	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	log.Println("requestPayload:", requestPayload)
	switch requestPayload.Action {
	case "auth":
		//app.authenticate(w, requestPayload.Auth)
		app.authEventViaRabbit(w, requestPayload.Auth)
	case "log":
		log.Println("Logging through RPC")
		app.logItemViaRPC(w, requestPayload.Log)
		//app.logitem(w, requestPayload.Log)
		//app.LogViaGRPC(w,r)
		//app.logEventViaRabbit(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unknown action"))
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
	log.Println("RES BODY:", res.Body)
	log.Println("RES:", res)
	if res.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if res.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return

	}

	var jsonFromService jsonResponse1

	err = json.NewDecoder(res.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse1
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) logitem(w http.ResponseWriter, entry LogPayLoad) {
	log.Println("logging item from broker service")
	log.Println("Entry_____________:", entry)
	jsonData, err := json.Marshal(entry)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	log.Println("11111111111111:", jsonData)
	logServiceUrl := "http://logger-service/log"

	req, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	log.Println("222222222222222222222222222")

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
	}
	defer res.Body.Close()
	log.Println("SSSSSSSSSSSSSS:", res.StatusCode)
	if res.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse1
	payload.Error = false
	payload.Message = "logged!"
	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) sendMail(w http.ResponseWriter, msg Mailpayload) {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		app.errorJSON(w, err)
	}
	log.Println("MSG:", msg)
	mailServicerURL := "http://mail-service/send"

	req, err := http.NewRequest("POST", mailServicerURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	log.Println("RESPONSE:", res)
	defer res.Body.Close()
	log.Println("RESPONSESTATUS:", res.StatusCode)
	if res.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mail service"))
		return
	}

	var payload jsonResponse1
	payload.Error = false
	payload.Message = "Message sent to " + msg.To
	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayLoad) {
	log.Println("I am logging event form Rabbit")
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse1
	payload.Error = false
	payload.Message = "logged via RabitMQ"

	app.writeJson(w, http.StatusAccepted, payload)

}

func (app *Config) authEventViaRabbit(w http.ResponseWriter, l AuthPayLoad) {
	log.Println("I am Authenticating event from RabbitEvent")
	err := app.pushToQueue(l.Email, l.Password)
	if err != nil {
		app.errorJSON(w, err)
		log.Println("I am returniing------------", err)
		return
	}

	var payload jsonResponse1
	payload.Error = false
	payload.Message = "Authenticated via RabbitMQ"
	app.writeJson(w, http.StatusAccepted, payload)

}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayLoad{
		Name: name,
		Data: msg,
	}
	log.Println("payload", payload)
	payload1 := AuthPayLoad{
		Email:    name,
		Password: msg,
	}
	log.Println("_______________PAYLOAD______________")
	log.Println("payload1", payload1)

	j, _ := json.Marshal(&payload)
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}

	return nil

}

type RPCpayload struct {
	Name string
	Data string
}

func (app *Config) logItemViaRPC(w http.ResponseWriter, l LogPayLoad) {

	log.Println("I am logging through RPC:", l)

	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {

		log.Println("ERROR:", err)
		app.errorJSON(w, err)
		return
	}

	log.Println("CLIENT:", client)

	rpcPayload := RPCpayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse1{
		Error:   false,
		Message: result,
	}

	app.writeJson(w, http.StatusAccepted, payload)

}

func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {

	var requestPayload RequestPayload

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	log.Println("CONN:", conn)
	defer conn.Close()

	c := logs.NewLogServiceClient(conn)
	log.Println("C:", c)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	log.Println("Ctx:", ctx)
	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})

	if err != nil {
		app.errorJSON(w, err)
	}
	var payload jsonResponse1
	payload.Error = false
	payload.Message = "logged!"
	app.writeJson(w, http.StatusAccepted, payload)

}
