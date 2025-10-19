package domain

import "time"

// Coords represents a geographical location resolved by a geocoding API.
type Coords struct {
	Name      string
	Country   string // ISO 3166-1 alpha-2 country code
	Latitude  float64
	Longitude float64
}

// Weather describes a current weather measurement in UTC.
type Weather struct {
	Time          time.Time
	Temperature2m float64 // Temperature in Celsius at 2 meters above ground
}
