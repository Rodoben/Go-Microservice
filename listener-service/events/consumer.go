package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for i, s := range topics {
		log.Println("TOPICS:", i, s)
		ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}
	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {

			log.Println("MESSAGESBODY_____________", string(d.Body))
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlepayload(payload)
		}
	}()
	fmt.Printf("Waiting for message [Exchange,Queue] [logs_topic, %s]\n", q.Name)
	<-forever
	return nil
}

func handlepayload(payload Payload) {
	log.Println("I am under handlepayload switch case listener service", payload)
	switch payload.Name {
	case "log", "event":
		log.Println("I am under log switch case listener service")
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		log.Println("I am under auth switch case listener service")
		err := authEvent(payload)
		if err != nil {
			log.Println(err)
		}
	default:
		log.Println("defualt")
		err := authEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	log.Println("I am under LogEVENT Listeber service")
	jsonData, _ := json.Marshal(entry)

	logServiceurl := "http://logger-service/log"

	req, err := http.NewRequest("POST", logServiceurl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusAccepted {
		return err
	}
	return nil
}

func authEvent(entry Payload) error {
	log.Println("I m under authEVENT Listener service", entry)
	var AuthPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	AuthPayload.Email = entry.Name
	AuthPayload.Password = entry.Data
	jsonData, err := json.Marshal(AuthPayload)
	log.Println(string(jsonData))
	if err != nil {
		log.Println("JSONDATAERROR:", err)
		return err
	}
	authServiceUrl := "http://authentication-service/authenticate"

	req, err := http.NewRequest("POST", authServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	log.Println(res.Body)
	log.Println(res.StatusCode)
	if res.StatusCode != http.StatusAccepted {
		return err
	}
	return nil

}
