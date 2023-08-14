package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func (app *Config) routes() http.Handler {
	log.Printf("inside Routes : %s\n", app)
	mux := chi.NewRouter()

	mux.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://swarm.rodo.co.in"},
		AllowCredentials: true,
		MaxAge:           300,
		ExposedHeaders:   []string{"Link"},
		AllowedMethods:   []string{http.MethodPost, http.MethodGet, http.MethodOptions, http.MethodPut},
		AllowedHeaders:   []string{"Accept", "Accept-Language", "Content-Type", "Authorization", "Origin", "X-Requested-With"},
		AllowOriginFunc:  func(origin string) bool { return true },
	}).Handler)
	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/", app.Broker)
	mux.Post("/handle", app.HandleSubmission)
	mux.Post("/log-grpc", app.LogViaGRPC)
	log.Printf("app.Broker : %s\n", app.Broker)
	return mux
}

// mux.Use(cors.Handler(cors.Options{
// 	AllowedOrigins:   []string{"https://*", "http://swarm.rodo.co.in"},
// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
// 	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
// 	ExposedHeaders:   []string{"Link"},
// 	AllowCredentials: true,
// 	MaxAge:           300,
// }))

// Add a middleware that sets the 'Access-Control-Allow-Origin' header
// mux.Use(func(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "http://swarm.rodo.co.in")
// 		next.ServeHTTP(w, r)
// 	})
// })
