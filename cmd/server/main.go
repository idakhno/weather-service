package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron/v2"
)

// Server constants
const httpPort = ":3000"

// Handler constants
const handlerCity = "city"

func main() {
	// Init new router
	r := chi.NewRouter()

	// Init middleware
	r.Use(middleware.Logger)

	// Init handlers
	r.Get("/{city}", func(w http.ResponseWriter, r *http.Request) {
		city := chi.URLParam(r, handlerCity)
		fmt.Printf("Requested city: %s\n", city)

		if _, err := w.Write([]byte("welcome")); err != nil {
			log.Print(err)
		}
	})

	// Init cron jobs
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Print(err)
	}

	jobs, err := initJobs(s)
	if err != nil {
		log.Print(err)
	}

	// Init wait groups for multithreading jobs
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Start listening at desired port, non-blocking
	go func() {
		defer wg.Done()
		fmt.Println("starting server on port " + httpPort)
		if err := http.ListenAndServe(httpPort, r); err != nil {
			log.Print(err)
		}
	}()

	// Start cron jobs
	go func() {
		defer wg.Done()
		fmt.Println("starting job:", jobs[0].ID())
		s.Start()
	}()

	// Waiting for all multithreading jobs to be completed
	wg.Wait()
}

// Initialize cron jobs
func initJobs(scheduler gocron.Scheduler) ([]gocron.Job, error) {
	j, err := scheduler.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func() {
				fmt.Println("hello world")
			},
		),
	)
	if err != nil {
		return nil, err
	}

	return []gocron.Job{j}, nil
}
