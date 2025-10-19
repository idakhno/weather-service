package service

import (
	"context"

	"github.com/idakhno/weather-service/internal/domain"
	"github.com/idakhno/weather-service/internal/infrastructure/clients/openmeteo/httpclient"
	"github.com/idakhno/weather-service/internal/infrastructure/repository/postgres"
)

type WeatherService struct {
	client httpclient.Client
	repo   postgres.WeatherRepository
}

func NewWeatherService(client httpclient.Client, repo postgres.WeatherRepository) *WeatherService {
	return &WeatherService{client: client, repo: repo}
}

func (s *WeatherService) UpdateWeather(ctx context.Context, city string) error {
	coords, err := s.client.GetCoords(ctx, city)
	if err != nil {
		return err
	}

	weather, err := s.client.GetTemperature(ctx, coords.Latitude, coords.Longitude)
	if err != nil {
		return err
	}

	return s.repo.Save(ctx, city, weather)
}

func (s *WeatherService) GetLatestWeather(ctx context.Context, city string) (domain.Weather, error) {
	return s.repo.GetLatest(ctx, city)
}
