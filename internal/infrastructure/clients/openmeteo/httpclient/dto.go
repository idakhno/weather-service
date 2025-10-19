package httpclient

// geocodingResponse models the JSON structure returned by the
// Open-Meteo geocoding endpoint. It is used internally to unmarshal
// coordinates for a given city.
type geocodingResponse struct {
	Results []struct {
		Name      string  `json:"name"`
		Country   string  `json:"country"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

// weatherResponse models the JSON structure returned by the
// Open-Meteo forecast endpoint for current weather data. It is used
// internally to decode the "current" section of the API response.
type weatherResponse struct {
	Current struct {
		Time          string  `json:"time"`
		Temperature2m float64 `json:"temperature_2m"`
	} `json:"current"`
}
