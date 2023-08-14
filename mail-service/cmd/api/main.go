package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	webPort = "80"
)

type Config struct {
	Mailer Mail
}

func main() {
	app := Config{
		Mailer: createMail(),
	}

	log.Println("Starting the server Mail-Service...")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}
	log.Println("I am Listening at :", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}

	return m
}
