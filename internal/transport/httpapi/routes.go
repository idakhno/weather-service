package httpapi

import (
	"github.com/go-chi/chi/v5"
	"github.com/idakhno/weather-service/internal/service"
)

// RegisterRoutes registers all HTTP routes for the weather API.
func RegisterRoutes(r *chi.Mux, svc *service.WeatherService) {
	h := NewHandler(svc)

	r.Route("/", func(r chi.Router) {
		r.Get("/{city}", h.GetWeather)
	})
}
