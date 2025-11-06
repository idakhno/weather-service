package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/idakhno/weather-service/internal/infrastructure/clients/openmeteo/httpclient"
	"github.com/idakhno/weather-service/internal/infrastructure/repository/postgres"
	"github.com/idakhno/weather-service/internal/infrastructure/scheduler"
	"github.com/idakhno/weather-service/internal/service"
	"github.com/idakhno/weather-service/internal/transport/httpapi"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Background context is used for long-running operations that should not be cancelled
	ctx := context.Background()

	db, err := pgxpool.New(ctx, "postgresql://postgres:postgres@localhost:5432/weather")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	httpClient := &http.Client{Timeout: 10 * time.Second}
	openMeteo := httpclient.NewClient(httpClient)
	repo := postgres.NewWeatherRepo(db)
	serviceInit := service.NewWeatherService(openMeteo, repo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	httpapi.RegisterRoutes(r, serviceInit)

	schedulerCron, err := scheduler.New()
	if err != nil {
		log.Fatalf("failed to create scheduler: %v", err)
	}

	err = schedulerCron.Every(10*time.Second, func() {
		if err := serviceInit.UpdateWeather(ctx, "moscow"); err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Fatalf("failed to schedule weather update: %v", err)
	}

	// Start HTTP server in a separate goroutine to allow scheduler to run concurrently
	go func() {
		err := http.ListenAndServe(":3000", r)
		if err != nil {
			log.Println(err)
		}
	}()
	schedulerCron.Start()
	// Block forever to keep the application running
	select {}
}
