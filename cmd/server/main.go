package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Init new router
	r := chi.NewRouter()

	// Init middleware
	r.Use(middleware.Logger)

	// Init handlers
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("welcome")); err != nil {
			log.Print(err)
		}
	})

	// Start listening at desired port
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Print(err)
	}
}
