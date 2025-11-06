// Package httpclient provides an HTTP implementation of the Open-Meteo API client.
// It fetches geocoding and weather data and maps responses into domain models.
package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/idakhno/weather-service/internal/domain"
)

// Client defines operations supported by the Open-Meteo API client.
type Client interface {
	GetCoords(ctx context.Context, city string) (domain.Coords, error)
	GetTemperature(ctx context.Context, lat, lon float64) (domain.Weather, error)
}

type client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) Client {
	return &client{httpClient: httpClient}
}

// GetCoords calls Open-Meteo geocoding API and returns city coordinates.
func (c *client) GetCoords(ctx context.Context, city string) (domain.Coords, error) {
	// QueryEscape ensures special characters in city names are properly encoded
	apiURL := fmt.Sprintf(
		"https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=ru&format=json",
		url.QueryEscape(city),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return domain.Coords{}, fmt.Errorf("build request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return domain.Coords{}, fmt.Errorf("send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return domain.Coords{}, fmt.Errorf("geocoding API returned %d", res.StatusCode)
	}

	var geoResp geocodingResponse
	if err := json.NewDecoder(res.Body).Decode(&geoResp); err != nil {
		return domain.Coords{}, fmt.Errorf("decode response: %w", err)
	}

	if len(geoResp.Results) == 0 {
		return domain.Coords{}, fmt.Errorf("no coordinates found for city %q", city)
	}

	// Use the first result from geocoding API (most relevant match)
	r := geoResp.Results[0]
	return domain.Coords{
		Name:      r.Name,
		Country:   r.Country,
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
	}, nil
}

// GetTemperature calls Open-Meteo forecast API and returns current temperature.
func (c *client) GetTemperature(ctx context.Context, lat, lon float64) (domain.Weather, error) {
	apiURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m",
		lat, lon,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return domain.Weather{}, fmt.Errorf("build request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return domain.Weather{}, fmt.Errorf("send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return domain.Weather{}, fmt.Errorf("weather API returned %d", res.StatusCode)
	}

	var weatherResp weatherResponse
	if err := json.NewDecoder(res.Body).Decode(&weatherResp); err != nil {
		return domain.Weather{}, fmt.Errorf("decode response: %w", err)
	}

	// Open-Meteo may return time in various ISO 8601 formats depending on the API version.
	// Try multiple formats to handle different response variations.
	var t time.Time
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04Z",
		"2006-01-02T15:04",
	}
	var parseErr error
	for _, format := range formats {
		t, parseErr = time.Parse(format, weatherResp.Current.Time)
		if parseErr == nil {
			break
		}
	}
	if parseErr != nil {
		return domain.Weather{}, fmt.Errorf("parse time %q: %w", weatherResp.Current.Time, parseErr)
	}

	return domain.Weather{
		Time:          t,
		Temperature2m: weatherResp.Current.Temperature2m,
	}, nil
}
