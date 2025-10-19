package postgres

import (
	"context"

	"github.com/idakhno/weather-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WeatherRepository interface {
	Save(ctx context.Context, city string, w domain.Weather) error
	GetLatest(ctx context.Context, city string) (domain.Weather, error)
}

type weatherRepo struct {
	db *pgxpool.Pool
}

func NewWeatherRepo(db *pgxpool.Pool) WeatherRepository {
	return &weatherRepo{db: db}
}

func (r *weatherRepo) Save(ctx context.Context, city string, w domain.Weather) error {
	_, err := r.db.Exec(ctx,
		`insert into reading (name, temperature, timestamp) values ($1, $2, $3)`,
		city, w.Temperature2m, w.Time,
	)
	return err
}

func (r *weatherRepo) GetLatest(ctx context.Context, city string) (domain.Weather, error) {
	var w domain.Weather
	err := r.db.QueryRow(ctx,
		`select timestamp, temperature from reading where name = $1 order by timestamp desc limit 1`,
		city,
	).Scan(&w.Time, &w.Temperature2m)
	return w, err
}
