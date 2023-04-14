package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type Config struct{}

func main() {

	app := Config{}
	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Starting broker server : %s\n", srv)
	// start the server
	log.Printf("I am listening: : %s\n")
	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}
