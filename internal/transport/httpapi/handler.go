package httpapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/idakhno/weather-service/internal/service"
)

// Handler holds references to application services used by HTTP endpoints.
type Handler struct {
	svc *service.WeatherService
}

func NewHandler(svc *service.WeatherService) *Handler {
	return &Handler{svc: svc}
}

// GetWeather handles GET /{city} requests and returns the latest weather reading.
func (h *Handler) GetWeather(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")
	if city == "" {
		http.Error(w, "city is required", http.StatusBadRequest)
		return
	}

	weather, err := h.svc.GetLatestWeather(r.Context(), city)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(weather); err != nil {
		log.Println("encode error:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
